package mgo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func DeleteMany(ctx context.Context, cname string, filter bson.D) (int64, error) {
	collection := GetCollection(cname)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Join(ErrWriteFailed, err)
	}
	return result.DeletedCount, nil
}

func DeleteById(ctx context.Context, cname string, id any) (int64, error) {
	collection := GetCollection(cname)
	result, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return 0, errors.Join(ErrWriteFailed, err)
	}
	return result.DeletedCount, nil
}
