// File: assignment-2/handlers/notifications_service.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"assignment-2/firebase"
	"assignment-2/structs"
)

// TriggerWebhookEventVar is a function variable we can override in test code.
// It points by default to realTriggerWebhookEvent.
var TriggerWebhookEventVar func(event, country string) = realTriggerWebhookEvent

// realTriggerWebhookEvent is the actual logic for sending webhooks.
func realTriggerWebhookEvent(event, country string) {
	ctx := context.Background()

	notifs, err := firebase.GetAllNotifications(ctx)
	if err != nil {
		log.Printf("[Webhook] Could not fetch notifications: %v\n", err)
		return
	}

	// Filter relevant notifications
	var relevant []structs.Notification
	for _, n := range notifs {
		if n.Event == event && (n.Country == "" || n.Country == country) {
			relevant = append(relevant, n)
		}
	}

	if len(relevant) == 0 {
		log.Printf("[Webhook] No matching webhooks for event=%s, country=%s\n", event, country)
		return
	}

	now := time.Now().Format("20060102 15:04")
	for _, wh := range relevant {
		body := map[string]string{
			"id":      wh.ID,
			"country": country,
			"event":   event,
			"time":    now,
		}
		bodyBytes, _ := json.Marshal(body)
		resp, postErr := http.Post(wh.URL, "application/json", bytes.NewBuffer(bodyBytes))
		if postErr != nil {
			log.Printf("[Webhook] Failed POST to %s: %v\n", wh.URL, postErr)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Printf("[Webhook] Webhook %s responded %d\n", wh.URL, resp.StatusCode)
		} else {
			log.Printf("[Webhook] Successfully triggered %s, event=%s, country=%s\n",
				wh.URL, event, country)
		}
		resp.Body.Close()
	}
}
