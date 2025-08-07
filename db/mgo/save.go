package mgo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Save(ctx context.Context, doc DocInter) (string, error) {
	collection := GetCollection(doc.C())
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", errors.Join(ErrWriteFailed, err)
	}
	oid, _ := result.InsertedID.(bson.ObjectID)
	return oid.Hex(), nil
}
