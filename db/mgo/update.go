package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateOneWithPartial(ctx context.Context, d DocInter, fields bson.D) (int64, error) {
	collection := GetCollection(d.C())
	result, err := collection.UpdateOne(ctx, bson.M{"_id": d.GetId()},
		bson.D{
			{Key: "$set", Value: fields},
		},
	)
	if result != nil {
		return result.ModifiedCount, err
	}
	return 0, err
}

func UpdateOneWithIncrField(ctx context.Context, d DocInter, query bson.D, incrFields bson.D) (int64, error) {
	collection := GetCollection(d.C())
	result, err := collection.UpdateOne(ctx, query,
		bson.D{
			{Key: "$inc", Value: incrFields},
		},
	)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}
	if result != nil {
		return result.ModifiedCount, err
	}
	return 0, err
}
