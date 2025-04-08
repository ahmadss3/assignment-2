// File: assignment-2/cmd/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/handlers"
	"assignment-2/tools"
)

var startTime time.Time

func main() {
	// Record the application start time, used for uptime reporting
	startTime = time.Now()

	// Assign the startTime to status handler so it can calculate uptime
	handlers.AssignStartTime(startTime)

	// Initialize Firebase and Firestore
	if err := firebase.InitFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer func() {
		if firebase.FirestoreClient != nil {
			closeErr := firebase.FirestoreClient.Close()
			if closeErr != nil {
				log.Printf("Failed to close Firestore client: %v", closeErr)
			}
		}
	}()

	// Start a goroutine that periodically purges old cache data every hour
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			ctx := context.Background()
			err := firebase.PurgeOldCache(ctx, 24*time.Hour)
			if err != nil {
				log.Printf("Periodic cache purge failed: %v\n", err)
			} else {
				log.Println("Periodic cache purge successful")
			}
		}
	}()

	// Registrations
	http.HandleFunc(constants.REGISTRATIONS_PATH, handlers.RegistrationRouter)
	// Dashboards
	http.HandleFunc(constants.DASHBOARDS_PATH, handlers.DashboardsRouter)
	// Notifications
	http.HandleFunc(constants.NOTIFICATIONS_PATH, handlers.NotificationsRouter)
	// Status
	http.HandleFunc(constants.STATUS_PATH, handlers.StatusHandler)

	// Determine the port from environment or use default
	port := tools.GetServerPort(constants.DefaultPort)

	// Start the HTTP server
	fmt.Printf("Server running on port %s (version: %s)\n", port, constants.ServiceVersion)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
