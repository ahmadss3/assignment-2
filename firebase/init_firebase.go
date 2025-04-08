// file assignment-2/firebase/init_firebase.go
package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp is a global reference to the Firebase application
var FirebaseApp *firebase.App

// FirestoreClient is used for Firestore operations
var FirestoreClient *firestore.Client

// InitFirebase initializes Firestore.
func InitFirebase() error {
	opt := option.WithCredentialsFile("assignment-2-firebasekey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize firebase app: %v", err)
	}
	FirebaseApp = app

	client, err := app.Firestore(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize firestore client: %v", err)
	}
	FirestoreClient = client

	return nil
}

// ensureClient ensures that FirestoreClient is not nil
func ensureClient() error {
	if FirestoreClient == nil {
		return fmt.Errorf("firestore client is not initialized")
	}
	return nil
}
