// file assignment-2/firebase/firestore_interface_test.go
package firebase

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// ------------------------
//  Mock Types (Top-Level)
// ------------------------

// mockDocSnap simulates a FirestoreDocSnapInterface.
type mockDocSnap struct {
	exists bool
	data   interface{}
}

func (m *mockDocSnap) Exists() bool {
	return m.exists
}

func (m *mockDocSnap) DataTo(dst interface{}) error {
	switch v := dst.(type) {
	case *map[string]interface{}:
		if m.data == nil {
			return nil
		}
		if mMap, ok := m.data.(map[string]interface{}); ok {
			*v = mMap
			return nil
		}
		return errors.New("mockDocSnap: data is not a map[string]interface{}")
	default:
		return errors.New("mockDocSnap: unsupported DataTo type in this mock")
	}
}

// mockDocRef simulates a FirestoreDocRefInterface.
type mockDocRef struct {
	docID string
	snap  FirestoreDocSnapInterface // The snapshot returned on Get
}

func (m *mockDocRef) Get(ctx context.Context) (FirestoreDocSnapInterface, error) {
	if m.snap == nil {
		return nil, errors.New("mockDocRef: no snapshot set")
	}
	return m.snap, nil
}

func (m *mockDocRef) Delete(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// mockQuerySnap simulates a FirestoreQueryInterface.
type mockQuerySnap struct {
	docs []FirestoreDocSnapInterface
}

func (m *mockQuerySnap) GetAll() ([]FirestoreDocSnapInterface, error) {
	return m.docs, nil
}

// mockCollectionRef simulates a FirestoreCollectionInterface.
type mockCollectionRef struct {
	docRef   FirestoreDocRefInterface
	queryRef FirestoreQueryInterface
}

func (m *mockCollectionRef) Doc(id string) FirestoreDocRefInterface {
	return m.docRef
}

func (m *mockCollectionRef) Add(ctx context.Context, data interface{}) (FirestoreDocRefInterface, interface{}, error) {
	return nil, nil, errors.New("mockCollectionRef.Add not implemented")
}

func (m *mockCollectionRef) Documents(ctx context.Context) FirestoreQueryInterface {
	return m.queryRef
}

// mockClient simulates a FirestoreClientInterface.
type mockClient struct {
	coll FirestoreCollectionInterface
}

func (m *mockClient) Collection(path string) FirestoreCollectionInterface {
	return m.coll
}

func (m *mockClient) Close() error {
	return nil
}

// ------------------------
//  Tests
// ------------------------

// TestRealFirestoreClientIntegration attempts a minimal usage of the real FirestoreClient
// if it's initialized. Otherwise, we skip.
func TestRealFirestoreClientIntegration(t *testing.T) {
	if FirestoreClient == nil {
		t.Skip("Skipping integration test: FirestoreClient is nil (not initialized).")
	}

	realClient := &RealFirestoreClient{Client: FirestoreClient}
	coll := realClient.Collection("testCollectionInterface")
	if coll == nil {
		t.Error("Expected non-nil FirestoreCollectionInterface from realClient.Collection(...)")
	}

	docRef := coll.Doc("testDocID")
	if docRef == nil {
		t.Error("Expected non-nil FirestoreDocRefInterface from coll.Doc(...)")
	}

}

// TestFirestoreInterface_Mock checks how code might behave if we never call real Firestore
// but use our mock types to simulate the interfaces.
func TestFirestoreInterface_Mock(t *testing.T) {
	// Create a doc snapshot as example
	snap := &mockDocSnap{
		exists: true,
		data: map[string]interface{}{
			"hello": "world",
		},
	}
	// Create docRef
	docRef := &mockDocRef{
		docID: "doc-123",
		snap:  snap,
	}
	// Create collectionRef
	collRef := &mockCollectionRef{
		docRef: docRef,
		queryRef: &mockQuerySnap{
			docs: []FirestoreDocSnapInterface{snap},
		},
	}
	// The mock client
	mockCli := &mockClient{
		coll: collRef,
	}

	colIF := mockCli.Collection("somePath")
	if colIF == nil {
		t.Fatal("mockClient.Collection(...) returned nil, expected a mockCollectionRef")
	}
	docIF := colIF.Doc("someDocID")
	if docIF == nil {
		t.Fatal("mockCollectionRef.Doc(...) returned nil")
	}

	snapIF, err := docIF.Get(context.Background())
	if err != nil {
		t.Fatalf("mockDocRef.Get returned error: %v", err)
	}
	if !snapIF.Exists() {
		t.Error("mockDocSnap says it doesn't exist, but we set exists=true")
	}

	// Try DataTo
	var data map[string]interface{}
	if derr := snapIF.DataTo(&data); derr != nil {
		t.Errorf("DataTo returned error: %v", derr)
	} else if val, ok := data["hello"]; !ok || val != "world" {
		t.Errorf("Expected map[hello:world], got %v", data)
	}
}

// TestRealFirestoreTypes_Reflection ensures our real adapters implement the interfaces.
func TestRealFirestoreTypes_Reflection(t *testing.T) {
	rd := &RealDocRef{}
	if _, ok := interface{}(rd).(FirestoreDocRefInterface); !ok {
		t.Error("RealDocRef does not implement FirestoreDocRefInterface")
	}

	ds := &RealDocSnap{}
	if _, ok := interface{}(ds).(FirestoreDocSnapInterface); !ok {
		t.Error("RealDocSnap does not implement FirestoreDocSnapInterface")
	}

	qc := &RealQuerySnap{}
	if _, ok := interface{}(qc).(FirestoreQueryInterface); !ok {
		t.Error("RealQuerySnap does not implement FirestoreQueryInterface")
	}

	rc := &RealCollectionRef{}
	if _, ok := interface{}(rc).(FirestoreCollectionInterface); !ok {
		t.Error("RealCollectionRef does not implement FirestoreCollectionInterface")
	}

	rfc := &RealFirestoreClient{}
	if _, ok := interface{}(rfc).(FirestoreClientInterface); !ok {
		t.Error("RealFirestoreClient does not implement FirestoreClientInterface")
	}

	typ := reflect.TypeOf(&RealFirestoreClient{})
	if typ.Elem().Name() != "RealFirestoreClient" {
		t.Errorf("Expected RealFirestoreClient type name, got %s", typ.Elem().Name())
	}
}
