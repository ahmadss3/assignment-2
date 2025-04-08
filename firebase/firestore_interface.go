// File: assignment-2/firebase/firestore_interface.go
package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
)

// FirestoreClientInterface is main interface for Firestore operations.
type FirestoreClientInterface interface {
	Collection(path string) FirestoreCollectionInterface
	Close() error
}

// FirestoreCollectionInterface abstracts collection-level operations
type FirestoreCollectionInterface interface {
	Doc(id string) FirestoreDocRefInterface
	Add(ctx context.Context, data interface{}) (FirestoreDocRefInterface, interface{}, error)
	Documents(ctx context.Context) FirestoreQueryInterface
}

// FirestoreDocRefInterface abstracts document-level operations
type FirestoreDocRefInterface interface {
	Get(ctx context.Context) (FirestoreDocSnapInterface, error)
	Delete(ctx context.Context) (interface{}, error)
}

// FirestoreQueryInterface abstracts .GetAll() for listing multiple docs
type FirestoreQueryInterface interface {
	GetAll() ([]FirestoreDocSnapInterface, error)
}

// FirestoreDocSnapInterface abstracts the snapshot
type FirestoreDocSnapInterface interface {
	Exists() bool
	DataTo(dst interface{}) error
}

// ---------------------------
// Real (production) adapters
// ---------------------------

// RealFirestoreClient wraps *firestore.Client
type RealFirestoreClient struct {
	Client *firestore.Client
}

// Collection returns a real collection reference wrapper
func (r *RealFirestoreClient) Collection(path string) FirestoreCollectionInterface {
	return &RealCollectionRef{r.Client.Collection(path)}
}

// Close calls the real client's Close()
func (r *RealFirestoreClient) Close() error {
	return r.Client.Close()
}

// RealCollectionRef wraps a *firestore.CollectionRef
type RealCollectionRef struct {
	Coll *firestore.CollectionRef
}

func (r *RealCollectionRef) Doc(id string) FirestoreDocRefInterface {
	return &RealDocRef{r.Coll.Doc(id)}
}

func (r *RealCollectionRef) Add(ctx context.Context, data interface{}) (FirestoreDocRefInterface, interface{}, error) {
	docRef, wr, err := r.Coll.Add(ctx, data)
	if err != nil {
		return nil, nil, err
	}
	return &RealDocRef{docRef}, wr, nil
}

func (r *RealCollectionRef) Documents(ctx context.Context) FirestoreQueryInterface {
	return &RealQuerySnap{r.Coll.Documents(ctx)}
}

// RealDocRef wraps a *firestore.DocumentRef
type RealDocRef struct {
	Ref *firestore.DocumentRef
}

func (d *RealDocRef) Get(ctx context.Context) (FirestoreDocSnapInterface, error) {
	snap, err := d.Ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &RealDocSnap{snap}, nil
}

func (d *RealDocRef) Delete(ctx context.Context) (interface{}, error) {
	return d.Ref.Delete(ctx)
}

// RealQuerySnap wraps a DocumentIterator
type RealQuerySnap struct {
	Iter *firestore.DocumentIterator
}

func (r *RealQuerySnap) GetAll() ([]FirestoreDocSnapInterface, error) {
	docs, err := r.Iter.GetAll()
	if err != nil {
		return nil, err
	}
	var out []FirestoreDocSnapInterface
	for _, d := range docs {
		out = append(out, &RealDocSnap{d})
	}
	return out, nil
}

// RealDocSnap wraps a *firestore.DocumentSnapshot
type RealDocSnap struct {
	Snap *firestore.DocumentSnapshot
}

func (r *RealDocSnap) Exists() bool {
	return r.Snap.Exists()
}

func (r *RealDocSnap) DataTo(dst interface{}) error {
	return r.Snap.DataTo(dst)
}
