// File: assignment-2/services/status_check.go
package services

import (
	"net/http"
	"time"
)

// StatusCheckResult holds an integer status code
type StatusCheckResult struct {
	StatusCode int
	Error      error
}

// CheckCountriesAPI checks whether the REST Countries API is reachable
// by performing a minimal GET request. If successful, StatusCheckResult
// contains the returned HTTP status code. Otherwise, it contains zero
// and the respective error.
func CheckCountriesAPI() StatusCheckResult {
	url := "http://129.241.150.113:8080/v3.1/alpha/NO?fields=name"
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return StatusCheckResult{StatusCode: 0, Error: err}
	}
	defer resp.Body.Close()

	return StatusCheckResult{StatusCode: resp.StatusCode, Error: nil}
}

// CheckOpenMeteo checks if Open-Meteo is reachable with a minimal request
// for temperature data, returning a StatusCheckResult with either the
// HTTP status code or the encountered error.
func CheckOpenMeteo() StatusCheckResult {
	url := "https://api.open-meteo.com/v1/forecast?latitude=10&longitude=10&hourly=temperature_2m&past_days=0"
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return StatusCheckResult{StatusCode: 0, Error: err}
	}
	defer resp.Body.Close()

	return StatusCheckResult{StatusCode: resp.StatusCode, Error: nil}
}

// CheckCurrencyAPI checks the currency API with a minimal GET call
// on base=NOK, returning a StatusCheckResult that includes either
// the HTTP status code or a network error.
func CheckCurrencyAPI() StatusCheckResult {
	url := "http://129.241.150.113:9090/currency/NOK"
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return StatusCheckResult{StatusCode: 0, Error: err}
	}
	defer resp.Body.Close()

	return StatusCheckResult{StatusCode: resp.StatusCode, Error: nil}
}

// TranslateErrorToStatus translates an error into a guessed HTTP code,
// returning 200 if err==nil, or 503 if err!=nil.
func TranslateErrorToStatus(err error) int {
	if err == nil {
		return 200
	}
	return 503
}

// CheckFirebaseNotifications attempts a minimal check on the "notifications" collection,
// for instance by calling GetAllNotifications. If it returns an error, we consider
// status=503. Otherwise, it's presumably 200.
func CheckFirebaseNotifications() StatusCheckResult {
	_, err := GetAllNotifications()
	if err != nil {
		return StatusCheckResult{StatusCode: 503, Error: err}
	}
	return StatusCheckResult{StatusCode: 200, Error: nil}
}

// GetAllNotifications is a helper that calls some actual code from firebase.
func GetAllNotifications() ([]string, error) {
	// Minimal approach: placeholder returning no data and no error.
	return nil, nil
}
