package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SetDatastore replaces the default datastore with a mock for testing.
// It returns a function to restore the original datastore, which should be
// called at the end of the test using defer.
//
// Example:
//
//	restore := SetDatastore(&MockDatastore{...})
//	defer restore()
func SetDatastore(mock Datastore) (restore func()) {
	original := dataStore
	dataStore = mock
	return func() {
		dataStore = original
	}
}

// MockDatastore is a mock implementation of the Datastore interface.
// It allows for setting mock functions for each method, making it easy to
// control the behavior of the datastore in tests.
type MockDatastore struct {
	OnSave             func(ctx context.Context, doc DocInter) (DocInter, error)
	OnFind             func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	OnFindOne          func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	OnUpdateOne        func(ctx context.Context, collection string, filter bson.D, update bson.D) (int64, error)
	OnDeleteOne        func(ctx context.Context, collection string, filter bson.D) (int64, error)
	OnDeleteMany       func(ctx context.Context, collection string, filter bson.D) (int64, error)
	OnPipeFind         func(ctx context.Context, collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	OnNewBulkOperation func(cname string) BulkOperator
	OnGetCollection    func(name string) *mongo.Collection
	OnClose            func(ctx context.Context) error
}

// MockBulkOperator is a mock implementation of the BulkOperator interface.
type MockBulkOperator struct {
	OnInsertOne func(doc DocInter) BulkOperator
	OnUpdateOne func(filter any, update any) BulkOperator
	OnExecute   func(ctx context.Context) (*mongo.BulkWriteResult, error)
}

// Interface implementations for MockDatastore

func (m *MockDatastore) Save(ctx context.Context, doc DocInter) (DocInter, error) {
	return m.OnSave(ctx, doc)
}

func (m *MockDatastore) Find(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	return m.OnFind(ctx, collection, filter, opts...)
}

func (m *MockDatastore) FindOne(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	return m.OnFindOne(ctx, collection, filter, opts...)
}

func (m *MockDatastore) UpdateOne(ctx context.Context, collection string, filter bson.D, update bson.D) (int64, error) {
	return m.OnUpdateOne(ctx, collection, filter, update)
}

func (m *MockDatastore) DeleteOne(ctx context.Context, collection string, filter bson.D) (int64, error) {
	return m.OnDeleteOne(ctx, collection, filter)
}

func (m *MockDatastore) DeleteMany(ctx context.Context, collection string, filter bson.D) (int64, error) {
	return m.OnDeleteMany(ctx, collection, filter)
}

func (m *MockDatastore) PipeFind(ctx context.Context, collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return m.OnPipeFind(ctx, collection, pipeline)
}

func (m *MockDatastore) NewBulkOperation(cname string) BulkOperator {
	return m.OnNewBulkOperation(cname)
}

func (m *MockDatastore) getCollection(name string) *mongo.Collection {
	return m.OnGetCollection(name)
}

func (m *MockDatastore) close(ctx context.Context) error {
	return m.OnClose(ctx)
}

// Interface implementations for MockBulkOperator

func (m *MockBulkOperator) InsertOne(doc DocInter) BulkOperator {
	return m.OnInsertOne(doc)
}

func (m *MockBulkOperator) UpdateOne(filter any, update any) BulkOperator {
	return m.OnUpdateOne(filter, update)
}

func (m *MockBulkOperator) UpdateById(id any, update any) BulkOperator {
	return m.OnUpdateOne(bson.M{"_id": id}, update)
}

func (m *MockBulkOperator) Execute(ctx context.Context) (*mongo.BulkWriteResult, error) {
	return m.OnExecute(ctx)
}

// ===================================================================
// Mock Helper Functions
// ===================================================================

// NewOnFindMock returns an OnFind function that returns a cursor with the given fake data.
func NewOnFindMock(fakeData ...any) func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	return func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
		cursor, err := mongo.NewCursorFromDocuments(fakeData, nil, nil)
		return cursor, err
	}
}

// NewOnFindOneMock returns an OnFindOne function that returns a SingleResult with the given fake data.
func NewOnFindOneMock(fakeData any) func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	return func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(fakeData, nil, nil)
	}
}

// NewErrOnFind returns an OnFind function that always returns the specified error.
func NewErrOnFind(err error) func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	return func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
		return nil, err
	}
}

// NewErrOnFindOne returns an OnFindOne function that returns a SingleResult containing the specified error.
func NewErrOnFindOne(err error) func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	return func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
		// Pass an empty non-nil document to prevent the decoder from returning
		// its own "document is nil" error, ensuring it returns the error we injected.
		return mongo.NewSingleResultFromDocument(bson.D{}, err, nil)
	}
}

// NewOnSaveMock returns an OnSave function that simulates a successful save.
// It assigns a new ObjectID to the document and returns it.
func NewOnSaveMock() func(ctx context.Context, doc DocInter) (DocInter, error) {
	return func(ctx context.Context, doc DocInter) (DocInter, error) {
		doc.SetId(bson.NewObjectID())
		return doc, nil
	}
}

// NewOnPipeFindMock returns an OnPipeFind function that returns a cursor with the given fake data.
func NewOnPipeFindMock(fakeData ...any) func(ctx context.Context, collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return func(ctx context.Context, collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
		cursor, err := mongo.NewCursorFromDocuments(fakeData, nil, nil)
		return cursor, err
	}
}

// NewOnBulkOperationMock returns an OnNewBulkOperation function that creates a mock BulkOperator.
// The mock BulkOperator's chainable methods are pre-configured to return itself,
// and its Execute method is set to return the provided result and error.
func NewOnBulkOperationMock(result *mongo.BulkWriteResult, err error) func(cname string) BulkOperator {
	return func(cname string) BulkOperator {
		// Create a mock operator
		mockOp := &MockBulkOperator{}

		// Make chainable methods return the mock operator itself
		mockOp.OnInsertOne = func(doc DocInter) BulkOperator { return mockOp }
		mockOp.OnUpdateOne = func(filter any, update any) BulkOperator { return mockOp }

		// Set the final return value for the Execute method
		mockOp.OnExecute = func(ctx context.Context) (*mongo.BulkWriteResult, error) {
			return result, err
		}
		return mockOp
	}
}
