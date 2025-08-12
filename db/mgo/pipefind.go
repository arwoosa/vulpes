package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MgoAggregate interface {
	GetPipeline(q bson.M) mongo.Pipeline
	Index
}

func PipeFind[T MgoAggregate](ctx context.Context, aggr T, filter bson.M) ([]T, error) {
	collection := GetCollection(aggr.C())
	sortCursor, err := collection.Aggregate(ctx, aggr.GetPipeline(filter))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	var slice []T
	err = sortCursor.All(ctx, &slice)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return slice, nil
}
