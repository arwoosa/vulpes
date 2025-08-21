package mgo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/arwoosa/vulpes/db/mgo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUpdateMany(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		filter := bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 25}}}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "Name", Value: "Over 25"}}}}
		expectedModifiedCount := int64(2)

		mockDB := &mgo.MockDatastore{
			OnUpdateMany: func(ctx context.Context, collection string, f bson.D, u bson.D) (int64, error) {
				assert.Equal(t, "users", collection)
				assert.Equal(t, filter, f)
				assert.Equal(t, update, u)
				return expectedModifiedCount, nil
			},
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		modifiedCount, err := mgo.UpdateMany(context.Background(), &testUser{}, filter, update)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedModifiedCount, modifiedCount)
	})

	t.Run("Error from Datastore", func(t *testing.T) {
		// Arrange
		filter := bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 25}}}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "Name", Value: "Over 25"}}}}
		expectedErr := errors.New("datastore update many failed")

		mockDB := &mgo.MockDatastore{
			OnUpdateMany: func(ctx context.Context, collection string, f bson.D, u bson.D) (int64, error) {
				return 0, expectedErr
			},
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		modifiedCount, err := mgo.UpdateMany(context.Background(), &testUser{}, filter, update)

		// Assert
		assert.Zero(t, modifiedCount)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}
