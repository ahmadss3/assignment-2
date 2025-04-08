package firebase

import (
	"context"
	"fmt"
	"time"

	"assignment-2/constants"
	"assignment-2/structs"
)

// FUNCTION VARIABLES
// These can be overridden in tests.
// By default, they point to the real Firestore-based functions below.

var SaveRegistration func(ctx context.Context, reg structs.Registration) (string, error) = realSaveRegistration
var GetRegistrationByID func(ctx context.Context, docID string) (*structs.Registration, error) = realGetRegistrationByID
var GetAllRegistrations func(ctx context.Context) ([]structs.Registration, error) = realGetAllRegistrations
var UpdateRegistration func(ctx context.Context, docID string, reg structs.Registration) error = realUpdateRegistration
var DeleteRegistration func(ctx context.Context, docID string) error = realDeleteRegistration
var PatchRegistration func(ctx context.Context, docID string, partial structs.Registration) error = realPatchRegistration

// REAL IMPLEMENTATIONS
// These are the actual Firestore-based functions we run in production.
// They are called by the function variables above (unless overridden in tests).

func realSaveRegistration(ctx context.Context, reg structs.Registration) (string, error) {
	if err := ensureClient(); err != nil {
		return "", err
	}
	docRef, _, err := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Add(ctx, map[string]interface{}{
		"country":    reg.Country,
		"isoCode":    reg.ISOCode,
		"features":   reg.Features,
		"lastChange": reg.LastChange,
	})
	if err != nil {
		return "", fmt.Errorf("failed to add registration: %v", err)
	}
	return docRef.ID, nil
}

func realGetRegistrationByID(ctx context.Context, docID string) (*structs.Registration, error) {
	if err := ensureClient(); err != nil {
		return nil, err
	}
	docRef := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %v", err)
	}
	if !snap.Exists() {
		return nil, fmt.Errorf("registration not found")
	}

	var data struct {
		Country    string           `firestore:"country"`
		ISOCode    string           `firestore:"isoCode"`
		Features   structs.Features `firestore:"features"`
		LastChange time.Time        `firestore:"lastChange"`
	}
	if err := snap.DataTo(&data); err != nil {
		return nil, fmt.Errorf("failed to parse registration data: %v", err)
	}
	return &structs.Registration{
		ID:         docID,
		Country:    data.Country,
		ISOCode:    data.ISOCode,
		Features:   data.Features,
		LastChange: data.LastChange,
	}, nil
}

func realGetAllRegistrations(ctx context.Context) ([]structs.Registration, error) {
	if err := ensureClient(); err != nil {
		return nil, err
	}
	snaps, err := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registrations: %v", err)
	}

	var regs []structs.Registration
	for _, snap := range snaps {
		if !snap.Exists() {
			continue
		}
		var data struct {
			Country    string           `firestore:"country"`
			ISOCode    string           `firestore:"isoCode"`
			Features   structs.Features `firestore:"features"`
			LastChange time.Time        `firestore:"lastChange"`
		}
		if err := snap.DataTo(&data); err != nil {
			// Could log a warning, but continue
			continue
		}
		regs = append(regs, structs.Registration{
			ID:         snap.Ref.ID,
			Country:    data.Country,
			ISOCode:    data.ISOCode,
			Features:   data.Features,
			LastChange: data.LastChange,
		})
	}
	return regs, nil
}

func realUpdateRegistration(ctx context.Context, docID string, reg structs.Registration) error {
	if err := ensureClient(); err != nil {
		return err
	}
	docRef := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch doc for update: %v", err)
	}
	if !snap.Exists() {
		return fmt.Errorf("document does not exist")
	}

	_, err = docRef.Set(ctx, map[string]interface{}{
		"country":    reg.Country,
		"isoCode":    reg.ISOCode,
		"features":   reg.Features,
		"lastChange": reg.LastChange,
	})
	if err != nil {
		return fmt.Errorf("failed to update registration: %v", err)
	}
	return nil
}

func realDeleteRegistration(ctx context.Context, docID string) error {
	if err := ensureClient(); err != nil {
		return err
	}
	docRef := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get document for deletion: %v", err)
	}
	if !snap.Exists() {
		return fmt.Errorf("registration not found")
	}
	_, err = docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete registration: %v", err)
	}
	return nil
}

func realPatchRegistration(ctx context.Context, docID string, partial structs.Registration) error {
	if err := ensureClient(); err != nil {
		return err
	}
	docRef := FirestoreClient.Collection(constants.REGISTRATIONS_COLLECTION).Doc(docID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch document for patch: %v", err)
	}
	if !snap.Exists() {
		return fmt.Errorf("document does not exist")
	}

	existingReg, err := realGetRegistrationByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to load existing registration: %v", err)
	}
	updateMap := map[string]interface{}{}

	// Only set fields if partial has a non-empty value
	if partial.Country != "" {
		updateMap["country"] = partial.Country
	}
	if partial.ISOCode != "" {
		updateMap["isoCode"] = partial.ISOCode
	}

	newFeatures := existingReg.Features
	changed := false

	if partial.Features.Temperature != existingReg.Features.Temperature {
		newFeatures.Temperature = partial.Features.Temperature
		changed = true
	}
	if partial.Features.Precipitation != existingReg.Features.Precipitation {
		newFeatures.Precipitation = partial.Features.Precipitation
		changed = true
	}
	if partial.Features.Capital != existingReg.Features.Capital {
		newFeatures.Capital = partial.Features.Capital
		changed = true
	}
	if partial.Features.Coordinates != existingReg.Features.Coordinates {
		newFeatures.Coordinates = partial.Features.Coordinates
		changed = true
	}
	if partial.Features.Population != existingReg.Features.Population {
		newFeatures.Population = partial.Features.Population
		changed = true
	}
	if partial.Features.Area != existingReg.Features.Area {
		newFeatures.Area = partial.Features.Area
		changed = true
	}
	if len(partial.Features.TargetCurrencies) > 0 {
		newFeatures.TargetCurrencies = partial.Features.TargetCurrencies
		changed = true
	}

	if changed {
		updateMap["features"] = newFeatures
	}
	if len(updateMap) > 0 {
		updateMap["lastChange"] = time.Now()
	} else {
		// no changes... do nothing
		return nil
	}

	_, err = docRef.Set(ctx, updateMap)
	if err != nil {
		return fmt.Errorf("failed to patch registration: %v", err)
	}
	return nil
}
