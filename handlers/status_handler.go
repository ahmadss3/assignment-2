// File: assignment-2/handlers/status_handler.go
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/tools"
)

var serviceStartTime time.Time

// AssignStartTime is called from main.go to share the service start time
func AssignStartTime(t time.Time) {
	serviceStartTime = t
	fmt.Printf("StatusHandler assigned startTime: %v\n", t)
}

// StatusHandler shows the status of external services
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Only GET is allowed on status")
		return
	}

	countriesStatus := checkCountriesAPI()
	meteoStatus := checkMeteoAPI()
	currencyStatus := checkCurrencyAPI()

	notifStatus, notifCount := checkNotificationsDB()
	overallStatus := http.StatusOK
	if countriesStatus != http.StatusOK || meteoStatus != http.StatusOK ||
		currencyStatus != http.StatusOK || notifStatus != http.StatusOK {
		overallStatus = http.StatusServiceUnavailable
	}

	uptimeSec := int(time.Since(serviceStartTime).Seconds())

	resp := map[string]interface{}{
		"countries_api":   countriesStatus,
		"meteo_api":       meteoStatus,
		"currency_api":    currencyStatus,
		"notification_db": notifStatus,
		"webhooks":        notifCount,
		"version":         constants.ServiceVersion,
		"uptime":          uptimeSec,
	}

	tools.WriteJsonResponse(w, overallStatus, resp)
}

func checkCountriesAPI() int {
	client := &http.Client{Timeout: 3 * time.Second}
	// minimal call
	url := constants.REST_COUNTRIES_ALPHA + "NO?fields=name"
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error checking countries API: %v\n", err)
		return http.StatusServiceUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

func checkMeteoAPI() int {
	client := &http.Client{Timeout: 3 * time.Second}
	url := constants.OPEN_METEO_API + "?latitude=10&longitude=10&hourly=temperature_2m"
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error checking open-meteo: %v\n", err)
		return http.StatusServiceUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

func checkCurrencyAPI() int {
	client := &http.Client{Timeout: 3 * time.Second}
	url := constants.CURRENCY_API + "NOK"
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error checking currency API: %v\n", err)
		return http.StatusServiceUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

func checkNotificationsDB() (int, int) {
	ctx := context.Background()
	notifs, err := firebase.GetAllNotifications(ctx)
	if err != nil {
		log.Printf("Error checking notifications DB: %v\n", err)
		return http.StatusServiceUnavailable, 0
	}
	return http.StatusOK, len(notifs)
}
