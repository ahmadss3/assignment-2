// File: assignment-2/firebase/notifications_firebase.go
package firebase

import (
	"context"
	"fmt"
	"time"

	"assignment-2/constants"
	"assignment-2/structs"
)

// FUNCTION VARIABLES
//
// These can be overridden in tests. By default, they point to the 'real' Firestore implementations.

var SaveNotification func(ctx context.Context, notif structs.Notification) (string, error) = realSaveNotification

var GetNotificationByID func(ctx context.Context, docID string) (*structs.Notification, error) = realGetNotificationByID

var GetAllNotifications func(ctx context.Context) ([]structs.Notification, error) = realGetAllNotifications

var DeleteNotification func(ctx context.Context, docID string) error = realDeleteNotification

// REAL IMPLEMENTATIONS
//
// The "real*" functions do the actual Firestore calls. The function variables
// above point to these by default, but can be overridden in tests.

// realSaveNotification is the actual Firestore-based logic for storing a new notification.
func realSaveNotification(ctx context.Context, notif structs.Notification) (string, error) {
	if err := ensureClient(); err != nil {
		return "", err
	}
	colRef := FirestoreClient.Collection(constants.NOTIFICATIONS_COLLECTION)
	docRef, _, err := colRef.Add(ctx, map[string]interface{}{
		"url":     notif.URL,
		"country": notif.Country,
		"event":   notif.Event,
		"created": notif.Created,
	})
	if err != nil {
		return "", fmt.Errorf("failed to save notification: %v", err)
	}
	return docRef.ID, nil
}

// realGetNotificationByID retrieves a single webhook doc from Firestore.
func realGetNotificationByID(ctx context.Context, docID string) (*structs.Notification, error) {
	if err := ensureClient(); err != nil {
		return nil, err
	}
	docRef := FirestoreClient.Collection(constants.NOTIFICATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification doc: %v", err)
	}
	if !snap.Exists() {
		return nil, fmt.Errorf("notification not found")
	}

	var data struct {
		URL     string    `firestore:"url"`
		Country string    `firestore:"country"`
		Event   string    `firestore:"event"`
		Created time.Time `firestore:"created"`
	}
	if err := snap.DataTo(&data); err != nil {
		return nil, fmt.Errorf("failed to parse notification data: %v", err)
	}
	return &structs.Notification{
		ID:      snap.Ref.ID,
		URL:     data.URL,
		Country: data.Country,
		Event:   data.Event,
		Created: data.Created,
	}, nil
}

// realGetAllNotifications lists all notifications in the Firestore collection.
func realGetAllNotifications(ctx context.Context) ([]structs.Notification, error) {
	if err := ensureClient(); err != nil {
		return nil, err
	}
	colRef := FirestoreClient.Collection(constants.NOTIFICATIONS_COLLECTION)
	snaps, err := colRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %v", err)
	}
	var results []structs.Notification
	for _, snap := range snaps {
		if !snap.Exists() {
			continue
		}
		var data struct {
			URL     string    `firestore:"url"`
			Country string    `firestore:"country"`
			Event   string    `firestore:"event"`
			Created time.Time `firestore:"created"`
		}
		if err := snap.DataTo(&data); err != nil {
			// Could log a warning, skip
			continue
		}
		results = append(results, structs.Notification{
			ID:      snap.Ref.ID,
			URL:     data.URL,
			Country: data.Country,
			Event:   data.Event,
			Created: data.Created,
		})
	}
	return results, nil
}

// realDeleteNotification deletes a notification doc in Firestore by ID.
func realDeleteNotification(ctx context.Context, docID string) error {
	if err := ensureClient(); err != nil {
		return err
	}
	docRef := FirestoreClient.Collection(constants.NOTIFICATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get notification doc: %v", err)
	}
	if !snap.Exists() {
		return fmt.Errorf("notification not found")
	}
	_, err = docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %v", err)
	}
	return nil
}
