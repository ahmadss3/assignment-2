// File: assignment-2/handlers/notifications_handler_test.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/structs"
)

// In-memory store for notifications
var (
	notifMutex sync.Mutex
	notifStore map[string]structs.Notification
	notifIDSeq int
)

// Original Firebase functions to revert to after tests
var (
	origSaveNotification    = firebase.SaveNotification
	origGetAllNotifications = firebase.GetAllNotifications
	origGetNotificationByID = firebase.GetNotificationByID
	origDeleteNotification  = firebase.DeleteNotification
)

// overrideNotificationStubs replaces Firebase functions with in-memory stubs
func overrideNotificationStubs() {
	notifStore = make(map[string]structs.Notification)
	notifIDSeq = 0

	firebase.SaveNotification = func(ctx context.Context, notif structs.Notification) (string, error) {
		notifMutex.Lock()
		defer notifMutex.Unlock()
		notifIDSeq++
		newID := "notif-" + string(rune(notifIDSeq))
		notif.ID = newID
		notifStore[newID] = notif
		return newID, nil
	}

	firebase.GetAllNotifications = func(ctx context.Context) ([]structs.Notification, error) {
		notifMutex.Lock()
		defer notifMutex.Unlock()
		var all []structs.Notification
		for _, n := range notifStore {
			all = append(all, n)
		}
		return all, nil
	}

	firebase.GetNotificationByID = func(ctx context.Context, docID string) (*structs.Notification, error) {
		notifMutex.Lock()
		defer notifMutex.Unlock()
		n, ok := notifStore[docID]
		if !ok {
			return nil, os.ErrNotExist
		}
		return &n, nil
	}

	firebase.DeleteNotification = func(ctx context.Context, docID string) error {
		notifMutex.Lock()
		defer notifMutex.Unlock()
		_, ok := notifStore[docID]
		if !ok {
			return os.ErrNotExist
		}
		delete(notifStore, docID)
		return nil
	}
}

func revertNotificationStubs() {
	firebase.SaveNotification = origSaveNotification
	firebase.GetAllNotifications = origGetAllNotifications
	firebase.GetNotificationByID = origGetNotificationByID
	firebase.DeleteNotification = origDeleteNotification
}

func TestNotificationsHandler(t *testing.T) {
	overrideNotificationStubs()
	defer revertNotificationStubs()

	t.Run("POST /notifications - success", func(t *testing.T) {
		body := `{"url":"https://example.org/hook","country":"NO","event":"REGISTER"}`
		req := httptest.NewRequest(http.MethodPost, constants.NOTIFICATIONS_PATH, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Expected 201 Created, got %d", rr.Code)
		}

		// Optionally parse response body to verify "id" field
		var resp map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}
		if resp["id"] == "" {
			t.Errorf("Expected an 'id' field in response, got empty")
		}
	})

	t.Run("POST /notifications - invalid JSON", func(t *testing.T) {
		body := `{"url": "bad}`
		req := httptest.NewRequest(http.MethodPost, constants.NOTIFICATIONS_PATH, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 Bad Request, got %d", rr.Code)
		}
	})

	t.Run("GET /notifications - success", func(t *testing.T) {
		// Let's add a second notification first
		_, _ = firebase.SaveNotification(context.Background(), structs.Notification{
			URL:     "https://another.site/hook",
			Country: "",
			Event:   "INVOKE",
			Created: time.Now(),
		})

		req := httptest.NewRequest(http.MethodGet, constants.NOTIFICATIONS_PATH, nil)
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", rr.Code)
		}

		var notifs []structs.Notification
		if err := json.Unmarshal(rr.Body.Bytes(), &notifs); err != nil {
			t.Errorf("Failed to decode notifications array: %v", err)
		}
		if len(notifs) < 2 {
			t.Errorf("Expected at least 2 notifications, got %d", len(notifs))
		}
	})

	t.Run("GET /notifications/{id} - success", func(t *testing.T) {
		// Insert a known notification
		id, _ := firebase.SaveNotification(context.Background(), structs.Notification{
			URL:     "https://single.test/hook",
			Country: "SE",
			Event:   "DELETE",
			Created: time.Now(),
		})

		req := httptest.NewRequest(http.MethodGet, constants.NOTIFICATIONS_PATH+id, nil)
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", rr.Code)
		}

		var notif structs.Notification
		if err := json.Unmarshal(rr.Body.Bytes(), &notif); err != nil {
			t.Errorf("Failed to parse single notif: %v", err)
		}
		if notif.ID != id {
			t.Errorf("Expected ID=%s, got %s", id, notif.ID)
		}
	})

	t.Run("GET /notifications/{id} - not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, constants.NOTIFICATIONS_PATH+"does-not-exist", nil)
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("DELETE /notifications/{id} - success", func(t *testing.T) {
		// Insert a known notification
		id, _ := firebase.SaveNotification(context.Background(), structs.Notification{
			URL:     "https://delete.me/hook",
			Country: "DK",
			Event:   "CHANGE",
		})

		req := httptest.NewRequest(http.MethodDelete, constants.NOTIFICATIONS_PATH+id, nil)
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("Expected 204 No Content, got %d", rr.Code)
		}

		// Double-check it's really gone
		_, err := firebase.GetNotificationByID(context.Background(), id)
		if err == nil {
			t.Errorf("Expected notification to be deleted, but found it in stub")
		}
	})

	t.Run("DELETE /notifications/{id} - not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, constants.NOTIFICATIONS_PATH+"missing-id", nil)
		rr := httptest.NewRecorder()

		NotificationsRouter(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("Method not allowed - single resource", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, constants.NOTIFICATIONS_PATH+"someid", nil)
		rr := httptest.NewRecorder()
		NotificationsRouter(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected 405 Method Not Allowed, got %d", rr.Code)
		}
	})

	t.Run("Method not allowed - collection", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, constants.NOTIFICATIONS_PATH, nil)
		rr := httptest.NewRecorder()
		NotificationsRouter(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected 405 Method Not Allowed, got %d", rr.Code)
		}
	})
}
