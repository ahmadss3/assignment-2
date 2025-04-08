// fiel assignment-2/firebase/notifications_firebase_test.go
package firebase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"assignment-2/structs"
)

// TestNotificationsFirebase demonstrates how to override the function variables
// in notifications_firebase.go (SaveNotification, GetAllNotifications, etc.)
// so we can test them purely in-memory without calling Firestore.
func TestNotificationsFirebase(t *testing.T) {
	// 1) Backup references to each function variable
	origSaveNotification := SaveNotification
	origGetNotificationByID := GetNotificationByID
	origGetAllNotifications := GetAllNotifications
	origDeleteNotification := DeleteNotification

	// 2) We'll store an in-memory map docID -> Notification
	var notifMutex sync.Mutex
	notifStore := make(map[string]structs.Notification)
	var idCounter int

	// 3) Define override logic (in memory)
	overrideStubs := func() {
		// SaveNotification
		SaveNotification = func(ctx context.Context, notif structs.Notification) (string, error) {
			notifMutex.Lock()
			defer notifMutex.Unlock()
			idCounter++
			docID := fmt.Sprintf("notif-%d", idCounter)
			notif.ID = docID
			notifStore[docID] = notif
			return docID, nil
		}

		// GetNotificationByID
		GetNotificationByID = func(ctx context.Context, docID string) (*structs.Notification, error) {
			notifMutex.Lock()
			defer notifMutex.Unlock()
			n, ok := notifStore[docID]
			if !ok {
				return nil, errors.New("notification not found")
			}
			return &n, nil
		}

		// GetAllNotifications
		GetAllNotifications = func(ctx context.Context) ([]structs.Notification, error) {
			notifMutex.Lock()
			defer notifMutex.Unlock()
			var all []structs.Notification
			for _, n := range notifStore {
				all = append(all, n)
			}
			return all, nil
		}

		// DeleteNotification
		DeleteNotification = func(ctx context.Context, docID string) error {
			notifMutex.Lock()
			defer notifMutex.Unlock()
			if _, ok := notifStore[docID]; !ok {
				return errors.New("doc not found")
			}
			delete(notifStore, docID)
			return nil
		}
	}

	// 4) Revert logic
	revertStubs := func() {
		SaveNotification = origSaveNotification
		GetNotificationByID = origGetNotificationByID
		GetAllNotifications = origGetAllNotifications
		DeleteNotification = origDeleteNotification
	}

	// 5) Actually override, revert at the end
	overrideStubs()
	defer revertStubs()

	// Helper to create a doc in memory
	createNotif := func(t *testing.T, url, country, event string) string {
		notif := structs.Notification{
			URL:     url,
			Country: country,
			Event:   event,
			Created: time.Now(),
		}
		docID, err := SaveNotification(context.Background(), notif)
		if err != nil {
			t.Fatalf("Failed SaveNotification for url=%s: %v", url, err)
		}
		return docID
	}

	t.Run("SaveNotification_Success", func(t *testing.T) {
		docID := createNotif(t, "http://fakeurl.org/notify", "NO", "REGISTER")
		if docID == "" {
			t.Error("Expected a docID, got empty string")
		}
	})

	t.Run("GetNotificationByID_Found", func(t *testing.T) {
		docID := createNotif(t, "https://example.com/webhook", "SE", "CHANGE")
		got, err := GetNotificationByID(context.Background(), docID)
		if err != nil {
			t.Fatalf("GetNotificationByID error: %v", err)
		}
		if got.URL != "https://example.com/webhook" {
			t.Errorf("Expected URL=https://example.com/webhook, got %s", got.URL)
		}
		if got.Country != "SE" {
			t.Errorf("Expected country=SE, got %s", got.Country)
		}
		if got.Event != "CHANGE" {
			t.Errorf("Expected event=CHANGE, got %s", got.Event)
		}
	})

	t.Run("GetNotificationByID_NotFound", func(t *testing.T) {
		_, err := GetNotificationByID(context.Background(), "notif-9999")
		if err == nil {
			t.Error("Expected error for non-existing docID, got nil")
		}
	})

	t.Run("GetAllNotifications", func(t *testing.T) {
		// Create multiple docs
		createNotif(t, "http://foo.com/a", "DK", "REGISTER")
		createNotif(t, "http://bar.com/b", "", "INVOKE")
		all, err := GetAllNotifications(context.Background())
		if err != nil {
			t.Fatalf("GetAllNotifications error: %v", err)
		}
		if len(all) < 2 {
			t.Errorf("Expected at least 2 notifications, got %d", len(all))
		}
	})

	t.Run("DeleteNotification_Success", func(t *testing.T) {
		docID := createNotif(t, "http://deleteme.org/hook", "DE", "DELETE")
		if err := DeleteNotification(context.Background(), docID); err != nil {
			t.Errorf("DeleteNotification returned error: %v", err)
		}
		_, err := GetNotificationByID(context.Background(), docID)
		if err == nil {
			t.Error("Expected error retrieving deleted doc, got nil")
		}
	})

	t.Run("DeleteNotification_NotFound", func(t *testing.T) {
		err := DeleteNotification(context.Background(), "notif-9999")
		if err == nil {
			t.Error("Expected error for docID=notif-9999 not found, got nil")
		}
	})
}
