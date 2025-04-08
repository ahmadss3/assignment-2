// File: assignment-2/firebase/registrations_firebase_test.go
package firebase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"assignment-2/structs"
)

// TestRegistrationsFirebase tests the CRUD logic in registrations_firebase.go
// by stubbing SaveRegistration, GetRegistrationByID ... as in-memory.
func TestRegistrationsFirebase(t *testing.T) {
	// 1) Backup references to each function variable so we can restore them.
	origSaveRegistration := SaveRegistration
	origGetRegistrationByID := GetRegistrationByID
	origGetAllRegistrations := GetAllRegistrations
	origUpdateRegistration := UpdateRegistration
	origDeleteRegistration := DeleteRegistration
	origPatchRegistration := PatchRegistration

	// 2) We'll keep an in-memory map of docID -> Registration to simulate Firestore
	var storeMutex sync.Mutex
	regStore := make(map[string]structs.Registration)
	var idCounter int

	// 3) Define override logic locally (in memory).
	overrideStubs := func() {
		// SaveRegistration
		SaveRegistration = func(ctx context.Context, reg structs.Registration) (string, error) {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			idCounter++
			docID := fmt.Sprintf("reg-%d", idCounter)
			reg.ID = docID
			regStore[docID] = reg
			return docID, nil
		}

		// GetRegistrationByID
		GetRegistrationByID = func(ctx context.Context, docID string) (*structs.Registration, error) {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			r, ok := regStore[docID]
			if !ok {
				return nil, errors.New("not found")
			}
			return &r, nil
		}

		// GetAllRegistrations
		GetAllRegistrations = func(ctx context.Context) ([]structs.Registration, error) {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			var all []structs.Registration
			for _, r := range regStore {
				all = append(all, r)
			}
			return all, nil
		}

		// UpdateRegistration
		UpdateRegistration = func(ctx context.Context, docID string, reg structs.Registration) error {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			oldReg, ok := regStore[docID]
			if !ok {
				return errors.New("doc not found")
			}
			reg.ID = docID
			regStore[docID] = reg

			// Return an error if to simulate something, or just nil
			if oldReg.Country != "" && reg.Country == "" {
				return errors.New("cannot remove country")
			}
			return nil
		}

		// DeleteRegistration
		DeleteRegistration = func(ctx context.Context, docID string) error {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			if _, ok := regStore[docID]; !ok {
				return errors.New("doc not found")
			}
			delete(regStore, docID)
			return nil
		}

		// PatchRegistration
		PatchRegistration = func(ctx context.Context, docID string, partial structs.Registration) error {
			storeMutex.Lock()
			defer storeMutex.Unlock()
			existing, ok := regStore[docID]
			if !ok {
				return errors.New("doc not found")
			}
			if partial.Country != "" {
				existing.Country = partial.Country
			}
			if partial.ISOCode != "" {
				existing.ISOCode = partial.ISOCode
			}
			// Patch features field by field
			if partial.Features.Temperature != existing.Features.Temperature {
				existing.Features.Temperature = partial.Features.Temperature
			}
			if len(partial.Features.TargetCurrencies) > 0 {
				existing.Features.TargetCurrencies = partial.Features.TargetCurrencies
			}
			existing.LastChange = time.Now()
			regStore[docID] = existing
			return nil
		}
	}

	// 4) Revert logic
	revertStubs := func() {
		SaveRegistration = origSaveRegistration
		GetRegistrationByID = origGetRegistrationByID
		GetAllRegistrations = origGetAllRegistrations
		UpdateRegistration = origUpdateRegistration
		DeleteRegistration = origDeleteRegistration
		PatchRegistration = origPatchRegistration
	}

	// 5) Actually override, revert at the end
	overrideStubs()
	defer revertStubs()

	// Helper function to create doc in memory
	createDoc := func(t *testing.T, countryName string) string {
		reg := structs.Registration{
			Country:    countryName,
			ISOCode:    "FAKE",
			Features:   structs.Features{Temperature: true},
			LastChange: time.Now(),
		}
		docID, err := SaveRegistration(context.Background(), reg)
		if err != nil {
			t.Fatalf("Failed to SaveRegistration for %s: %v", countryName, err)
		}
		return docID
	}

	t.Run("SaveRegistration_Success", func(t *testing.T) {
		docID := createDoc(t, "Norway")
		if docID == "" {
			t.Error("Expected a docID, got empty string")
		}
	})

	t.Run("GetRegistrationByID_Found", func(t *testing.T) {
		docID := createDoc(t, "Sweden")
		got, err := GetRegistrationByID(context.Background(), docID)
		if err != nil {
			t.Fatalf("GetRegistrationByID failed: %v", err)
		}
		if got.Country != "Sweden" {
			t.Errorf("Expected country=Sweden, got %s", got.Country)
		}
	})

	t.Run("GetRegistrationByID_NotFound", func(t *testing.T) {
		_, err := GetRegistrationByID(context.Background(), "not-exist")
		if err == nil {
			t.Error("Expected error for non-existing docID, got nil")
		}
	})

	t.Run("GetAllRegistrations", func(t *testing.T) {
		// create two docs
		createDoc(t, "Denmark")
		createDoc(t, "Finland")
		all, err := GetAllRegistrations(context.Background())
		if err != nil {
			t.Fatalf("GetAllRegistrations error: %v", err)
		}
		if len(all) < 2 {
			t.Errorf("Expected at least 2 docs, got %d", len(all))
		}
	})

	t.Run("UpdateRegistration_Success", func(t *testing.T) {
		docID := createDoc(t, "Germany")
		upd := structs.Registration{
			Country:    "GermanyUpdated",
			ISOCode:    "DEU",
			Features:   structs.Features{Temperature: false},
			LastChange: time.Now(),
		}
		if err := UpdateRegistration(context.Background(), docID, upd); err != nil {
			t.Errorf("UpdateRegistration returned error: %v", err)
		}
		got, _ := GetRegistrationByID(context.Background(), docID)
		if got.Country != "GermanyUpdated" {
			t.Errorf("Expected updated country=GermanyUpdated, got %s", got.Country)
		}
		if got.Features.Temperature != false {
			t.Errorf("Expected updated temperature=false, got %v", got.Features.Temperature)
		}
	})

	t.Run("UpdateRegistration_NotFound", func(t *testing.T) {
		upd := structs.Registration{Country: "Nowhere"}
		err := UpdateRegistration(context.Background(), "doc-9999", upd)
		if err == nil {
			t.Error("Expected error for docID=doc-9999 not found, got nil")
		}
	})

	t.Run("DeleteRegistration_Success", func(t *testing.T) {
		docID := createDoc(t, "DeleteMe")
		if err := DeleteRegistration(context.Background(), docID); err != nil {
			t.Errorf("DeleteRegistration returned error: %v", err)
		}
		_, err := GetRegistrationByID(context.Background(), docID)
		if err == nil {
			t.Error("Expected error retrieving deleted doc, got nil")
		}
	})

	t.Run("DeleteRegistration_NotFound", func(t *testing.T) {
		err := DeleteRegistration(context.Background(), "doc-9999")
		if err == nil {
			t.Error("Expected error for doc-9999 not found, got nil")
		}
	})

	t.Run("PatchRegistration_Success", func(t *testing.T) {
		docID := createDoc(t, "PatchMe")
		partial := structs.Registration{
			ISOCode: "PatchedISO",
			Features: structs.Features{
				TargetCurrencies: []string{"EUR", "USD"},
			},
		}
		err := PatchRegistration(context.Background(), docID, partial)
		if err != nil {
			t.Errorf("PatchRegistration error: %v", err)
		}
		got, _ := GetRegistrationByID(context.Background(), docID)
		if got.ISOCode != "PatchedISO" {
			t.Errorf("Expected ISOCode=PatchedISO, got %s", got.ISOCode)
		}
		if !reflect.DeepEqual(got.Features.TargetCurrencies, []string{"EUR", "USD"}) {
			t.Errorf("Mismatch in TargetCurrencies after patch")
		}
	})

	t.Run("PatchRegistration_NotFound", func(t *testing.T) {
		partial := structs.Registration{Country: "PatchNowhere"}
		err := PatchRegistration(context.Background(), "doc-9999", partial)
		if err == nil {
			t.Error("Expected error for doc-9999 not found, got nil")
		}
	})
}
