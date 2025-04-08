// File: assignment-2/handlers/notifications_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/structs"
	"assignment-2/tools"
)

func NotificationsRouter(w http.ResponseWriter, r *http.Request) {
	// Check if path is exactly /dashboard/v1/notifications/ or includes an ID
	if r.URL.Path == constants.NOTIFICATIONS_PATH {
		// No ID => handle collection-level
		handleNotificationsCollection(w, r)
	} else {
		// There's something after /notifications/
		id := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH)
		handleNotificationWithID(w, r, id)
	}
}

func handleNotificationsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostNotification(w, r)
	case http.MethodGet:
		handleGetAllNotifications(w, r)
	default:
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed on notifications collection")
	}
}

func handleNotificationWithID(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		handleGetNotificationByID(w, r, id)
	case http.MethodDelete:
		handleDeleteNotification(w, r, id)
	default:
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed on single notification")
	}
}

func handlePostNotification(w http.ResponseWriter, r *http.Request) {
	var req structs.Notification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding notification body: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	req.Created = time.Now()

	ctx := context.Background()
	newID, err := firebase.SaveNotification(ctx, req)
	if err != nil {
		log.Printf("Error saving notification: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Could not save webhook notification")
		return
	}
	resp := map[string]string{"id": newID}
	tools.WriteJsonResponse(w, http.StatusCreated, resp)
}

func handleGetAllNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	notifs, err := firebase.GetAllNotifications(ctx)
	if err != nil {
		log.Printf("Error fetching notifications: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Could not retrieve notifications")
		return
	}
	tools.WriteJsonResponse(w, http.StatusOK, notifs)
}

func handleGetNotificationByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	notif, err := firebase.GetNotificationByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching notification %s: %v\n", id, err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Notification not found")
		return
	}
	tools.WriteJsonResponse(w, http.StatusOK, notif)
}

func handleDeleteNotification(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	err := firebase.DeleteNotification(ctx, id)
	if err != nil {
		log.Printf("Error deleting notification %s: %v\n", id, err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Notification not found or could not be deleted")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
