package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Find[T DocInter](ctx context.Context, doc T, filter any, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	collection := GetCollection(doc.C())
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	var result []T
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return result, nil
}

func FindOne[T DocInter](ctx context.Context, doc T, filter any, opts ...options.Lister[options.FindOneOptions]) error {
	collection := GetCollection(doc.C())
	err := collection.FindOne(ctx, filter, opts...).Decode(&doc)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return nil
}

func FindById[T DocInter](ctx context.Context, doc T) error {
	collection := GetCollection(doc.C())
	err := collection.FindOne(ctx, bson.M{"_id": doc.GetId()}).Decode(&doc)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return nil
}
