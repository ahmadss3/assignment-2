// File: assignment-2/firebase/init_firebase_test.go
package firebase

import (
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

// TestInitFirebase attempts to init.
func TestInitFirebase(t *testing.T) {
	if _, err := os.Stat("assignment-2-firebasekey.json"); os.IsNotExist(err) {
		t.Skip("Skipping because assignment-2-firebasekey.json does not exist.")
	}
	err := InitFirebase()
	if err != nil {
		t.Errorf("InitFirebase returned an error: %v", err)
		return
	}
	if FirestoreClient == nil {
		t.Error("Expected FirestoreClient to be non-nil after InitFirebase, got nil.")
	}
}

// TestEnsureClient checks the behavior of ensureClient() with FirestoreClient set or nil.
func TestEnsureClient(t *testing.T) {
	originalClient := FirestoreClient

	// 1) check behavior when FirestoreClient is nil
	FirestoreClient = nil
	if err := ensureClient(); err == nil {
		t.Error("Expected an error when FirestoreClient is nil, got nil.")
	}

	// 2) check behavior when FirestoreClient is non-nil
	FirestoreClient = &firestore.Client{} // Fake client
	if err := ensureClient(); err != nil {
		t.Errorf("Expected no error from ensureClient, got %v", err)
	}

	// revert
	FirestoreClient = originalClient
}

// TestInitFirebase_NoFile
func TestInitFirebase_NoFile(t *testing.T) {
	// If no real file is present, we expect an error or skip.
	if _, err := os.Stat("assignment-2-firebasekey.json"); !os.IsNotExist(err) {
		t.Skip("File actually exists, skipping negative test.")
	}
	// Attempting to init => should fail
	err := InitFirebase()
	if err == nil {
		t.Error("Expected error if key file is missing, got nil.")
	}
}
