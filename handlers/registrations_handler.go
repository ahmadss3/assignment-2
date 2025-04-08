// file assignment-2/handlers/registrations_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/structs"
	"assignment-2/tools"
)

// RegistrationRouter manages /dashboard/v1/registrations/ and optionally an ID
func RegistrationRouter(w http.ResponseWriter, r *http.Request) {
	// Check if the path is exactly /dashboard/v1/registrations/ or includes an ID
	if r.URL.Path == constants.REGISTRATIONS_PATH {
		handleRegistrationsCollection(w, r)
		return
	}
	// Otherwise, we assume there's an ID
	id := r.URL.Path[len(constants.REGISTRATIONS_PATH):]
	handleRegistrationWithID(w, r, id)
}

func handleRegistrationsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostRegistration(w, r)
	case http.MethodGet:
		handleGetAllRegistrations(w, r)
	case http.MethodHead:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	default:
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed on registrations collection")
	}
}

func handleRegistrationWithID(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		handleGetRegistrationByID(w, r, id)
	case http.MethodPut:
		handlePutRegistration(w, r, id)
	case http.MethodPatch:
		handlePatchRegistration(w, r, id)
	case http.MethodDelete:
		handleDeleteRegistration(w, r, id)
	default:
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed for single registration")
	}
}

// handlePostRegistration
func handlePostRegistration(w http.ResponseWriter, r *http.Request) {
	var req structs.Registration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding registration body: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	req.LastChange = time.Now()

	ctx := context.Background()
	newID, err := firebase.SaveRegistration(ctx, req)
	if err != nil {
		log.Printf("Error saving registration: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Could not save registration in the database")
		return
	}

	resp := map[string]interface{}{
		"id":         newID,
		"lastChange": req.LastChange.Format("20060102 15:04"),
	}
	tools.WriteJsonResponse(w, http.StatusCreated, resp)

	countryFilter := req.Country
	if countryFilter == "" {
		countryFilter = req.ISOCode
	}
	TriggerWebhookEventVar("REGISTER", countryFilter)
}

// handleGetAllRegistrations
func handleGetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	regs, err := firebase.GetAllRegistrations(ctx)
	if err != nil {
		log.Printf("Error fetching registrations: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Could not retrieve registrations")
		return
	}
	tools.WriteJsonResponse(w, http.StatusOK, regs)
}

// handleGetRegistrationByID ... GET /.../registrations/{id}
func handleGetRegistrationByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	reg, err := firebase.GetRegistrationByID(ctx, id)
	if err != nil {
		log.Printf("Error getting registration by ID: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Registration not found")
		return
	}
	tools.WriteJsonResponse(w, http.StatusOK, reg)
}

// handlePutRegistration
func handlePutRegistration(w http.ResponseWriter, r *http.Request, id string) {
	var req structs.Registration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding PUT body: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	req.LastChange = time.Now()

	ctx := context.Background()
	err := firebase.UpdateRegistration(ctx, id, req)
	if err != nil {
		log.Printf("Error updating registration %s: %v\n", id, err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Could not update registration")
		return
	}
	w.WriteHeader(http.StatusNoContent)

	countryFilter := req.Country
	if countryFilter == "" {
		countryFilter = req.ISOCode
	}
	TriggerWebhookEventVar("CHANGE", countryFilter)
}

// handlePatchRegistration ... partial update
func handlePatchRegistration(w http.ResponseWriter, r *http.Request, id string) {
	var partial structs.Registration
	if err := json.NewDecoder(r.Body).Decode(&partial); err != nil {
		log.Printf("Error decoding PATCH body: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	ctx := context.Background()
	err := firebase.PatchRegistration(ctx, id, partial)
	if err != nil {
		log.Printf("Error patching registration %s: %v\n", id, err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Could not patch registration")
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// Trigger "CHANGE"
	countryFilter := partial.Country
	if countryFilter == "" {
		// If partial.Country was not set, read from updated doc
		updatedReg, _ := firebase.GetRegistrationByID(ctx, id)
		if updatedReg != nil && updatedReg.Country != "" {
			countryFilter = updatedReg.Country
		} else if updatedReg != nil {
			countryFilter = updatedReg.ISOCode
		}
	}
	TriggerWebhookEventVar("CHANGE", countryFilter)
}

// handleDeleteRegistration
func handleDeleteRegistration(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	existing, err := firebase.GetRegistrationByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching registration for delete: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Registration not found")
		return
	}

	err = firebase.DeleteRegistration(ctx, id)
	if err != nil {
		log.Printf("Error deleting registration %s: %v\n", id, err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Could not delete registration")
		return
	}
	w.WriteHeader(http.StatusNoContent)

	countryFilter := existing.Country
	if countryFilter == "" {
		countryFilter = existing.ISOCode
	}
	TriggerWebhookEventVar("DELETE", countryFilter)
}
