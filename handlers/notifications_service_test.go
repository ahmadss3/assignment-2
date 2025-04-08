// File: assignment-2/handlers/notifications_service_test.go
package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"assignment-2/firebase"
	"assignment-2/structs"
)

// We actually use these references for overriding firebase.GetAllNotifications:
var originalGetAllNotifications = firebase.GetAllNotifications

func overrideGetAllNotifications(fn func(ctx context.Context) ([]structs.Notification, error)) {
	firebase.GetAllNotifications = fn
}

func revertGetAllNotifications() {
	firebase.GetAllNotifications = originalGetAllNotifications
}

func TestTriggerWebhookEventVar(t *testing.T) {
	// 1) Override firebase.GetAllNotifications to return a few webhooks
	overrideGetAllNotifications(func(ctx context.Context) ([]structs.Notification, error) {
		return []structs.Notification{
			{ID: "w1", URL: "http://example.org/webhook1", Country: "NO", Event: "REGISTER"},
			{ID: "w2", URL: "http://example.org/webhook2", Country: "", Event: "REGISTER"},
			{ID: "w3", URL: "http://example.org/webhook3", Country: "DE", Event: "CHANGE"},
		}, nil
	})
	defer revertGetAllNotifications()

	// 2) We'll create a test server to capture incoming requests
	var requestBodies []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		requestBodies = append(requestBodies, string(bodyBytes))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 3) Override getAllNotifications again so that w1 = server.URL
	overrideGetAllNotifications(func(ctx context.Context) ([]structs.Notification, error) {
		return []structs.Notification{
			{ID: "w1", URL: server.URL, Country: "NO", Event: "REGISTER"},
			{ID: "w2", URL: "http://example.org/webhook2", Country: "", Event: "REGISTER"},
			{ID: "w3", URL: "http://example.org/webhook3", Country: "DE", Event: "CHANGE"},
		}, nil
	})

	// 4) Actually call TriggerWebhookEventVar ... event="REGISTER", country="NO"
	//    w1 is matched, w2 also matched but fails (not a real server).
	TriggerWebhookEventVar("REGISTER", "NO")

	if len(requestBodies) != 1 {
		t.Errorf("Expected exactly 1 POST to test server, got %d", len(requestBodies))
	} else {
		var posted map[string]string
		if err := json.Unmarshal([]byte(requestBodies[0]), &posted); err != nil {
			t.Errorf("Could not parse posted JSON: %v", err)
		} else {
			if posted["id"] != "w1" {
				t.Errorf("Expected id=w1, got %s", posted["id"])
			}
			if posted["country"] != "NO" {
				t.Errorf("Expected country=NO, got %s", posted["country"])
			}
			if posted["event"] != "REGISTER" {
				t.Errorf("Expected event=REGISTER, got %s", posted["event"])
			}
			if posted["time"] == "" {
				t.Error("Expected time to be set, but it was empty")
			}
		}
	}
}

func TestTriggerWebhookEventVar_NoMatching(t *testing.T) {
	overrideGetAllNotifications(func(ctx context.Context) ([]structs.Notification, error) {
		// No webhooks for event="DELETE", country="NO"
		return []structs.Notification{
			{ID: "wX", URL: "http://fakeurl.org/hookX", Country: "SE", Event: "CHANGE"},
		}, nil
	})
	defer revertGetAllNotifications()

	// We'll set up a small test server to see if any POST occurs
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// No matching webhooks
	TriggerWebhookEventVar("DELETE", "NO")

	if callCount != 0 {
		t.Errorf("Expected 0 calls to test server, got %d", callCount)
	}
}

func TestTriggerWebhookEventVar_GetAllFails(t *testing.T) {
	overrideGetAllNotifications(func(ctx context.Context) ([]structs.Notification, error) {
		return nil, os.ErrNotExist
	})
	defer revertGetAllNotifications()

	// Should just log the error and return, no crash expected
	TriggerWebhookEventVar("REGISTER", "NO")
}
