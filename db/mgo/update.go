package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// UpdateById updates a single document identified by the _id field of the provided document instance.
// It offers a flexible way to apply any update operator ($set, $inc, etc.).
//
// doc: An instance of the document, from which the _id is extracted for the filter.
//
//	It is also used to determine the target collection.
//
// update: The update document, e.g., bson.D{{"$set", bson.D{{"field", "value"}}}}.
func UpdateById[T DocInter](ctx context.Context, doc T, update bson.D) (int64, error) {
	return UpdateOne(ctx, doc, bson.D{{Key: "_id", Value: doc.GetId()}}, update)
}

// UpdateOne updates the first document that matches a given filter.
// This is a generic and flexible update function.
//
// doc: An instance of the document type, used to determine the collection.
// filter: The filter to select the document to update.
// update: The update document, e.g., bson.D{{"$set", bson.D{{"field", "value"}}}}.
func UpdateOne[T DocInter](ctx context.Context, doc T, filter bson.D, update bson.D) (int64, error) {
	collection := GetCollection(doc.C())
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}
	return result.ModifiedCount, nil
}
