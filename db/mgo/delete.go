package mgo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DeleteMany deletes all documents matching the filter.
// doc: An instance of the document type, used to determine the collection.
func DeleteMany[T DocInter](ctx context.Context, doc T, filter bson.D) (int64, error) {
	collection := GetCollection(doc.C())
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Join(ErrWriteFailed, err)
	}
	return result.DeletedCount, nil
}

// DeleteById deletes a single document identified by the _id of the provided document instance.
// doc: An instance of the document, from which the _id is extracted for the filter.
func DeleteById[T DocInter](ctx context.Context, doc T) (int64, error) {
	collection := GetCollection(doc.C())
	result, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: doc.GetId()}})
	if err != nil {
		return 0, errors.Join(ErrWriteFailed, err)
	}
	return result.DeletedCount, nil
}