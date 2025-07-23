// Package mgo provides a high-level abstraction layer over the official MongoDB Go driver,
// simplifying connection management, document operations, and schema definitions.
package mgo

import (
	"context"

	"github.com/arwoosa/vulpes/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BatchSave executes a bulk insert operation for a given slice of documents.
// It performs validation on each document and ensures they all belong to the same collection
// before sending them to the database. This provides a safe and efficient way to insert multiple documents.
//
// ctx: The context for the database operation.
// doclist: A DocSlice containing the documents to be inserted.
// Returns the number of successfully inserted documents, or an error if the operation fails.
func BatchSave(ctx context.Context, doclist DocSlice) (int64, error) {
	if len(doclist) == 0 {
		return 0, errors.NewWrapperError(ErrInvalidDocument, "no documents to save")
	}
	var err error
	// All documents in a single batch must belong to the same collection.
	// We determine the target collection from the first document.
	cname := doclist[0].C()
	for _, d := range doclist {
		// Enforce that all documents are for the same collection.
		if d.C() != cname {
			return 0, errors.NewWrapperError(ErrInvalidDocument, "all documents must be in the same collection")
		}
		// Run the document's own validation logic.
		if err = d.Validate(); err != nil {
			return 0, errors.NewWrapperError(ErrInvalidDocument, err.Error())
		}
	}

	// Pre-allocate the slice for write models with the exact capacity needed.
	// This avoids repeated memory allocations inside the loop, improving performance.
	models := make([]mongo.WriteModel, 0, len(doclist))
	for _, d := range doclist {
		models = append(models, mongo.NewInsertOneModel().SetDocument(d))
	}

	// Execute the bulk write operation.
	collection := GetCollection(cname)
	writeResult, err := collection.BulkWrite(ctx, models)
	if err != nil {
		return 0, errors.NewWrapperError(ErrBulkWriteFailed, err.Error())
	}

	return writeResult.InsertedCount, nil
}
