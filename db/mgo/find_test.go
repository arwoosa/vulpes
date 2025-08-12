package mgo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/arwoosa/vulpes/db/mgo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// testUser is a simple struct used for testing purposes.
type testUser struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string
	Age  int
}

// Implement the DocInter interface for testUser.
func (u *testUser) C() string                   { return "users" }
func (u *testUser) Indexes() []mongo.IndexModel { return nil }
func (u *testUser) Validate() error             { return nil }
func (u *testUser) GetId() any                  { return u.ID }
func (u *testUser) SetId(id any)                { u.ID = id.(bson.ObjectID) }

func TestFind(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedUsers := []any{
			testUser{ID: bson.NewObjectID(), Name: "Peter", Age: 30},
			testUser{ID: bson.NewObjectID(), Name: "Alice", Age: 25},
		}
		mockDB := &mgo.MockDatastore{
			OnFind: mgo.NewOnFindMock(expectedUsers...),
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		var result []*testUser
		result, err := mgo.Find(context.Background(), &testUser{}, nil)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedUsers[0].(testUser).Name, result[0].Name)
		assert.Equal(t, expectedUsers[1].(testUser).Name, result[1].Name)
	})

	t.Run("Error from Datastore", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("datastore find failed")
		mockDB := &mgo.MockDatastore{
			OnFind: mgo.NewErrOnFind(expectedErr),
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		var result []*testUser
		result, err := mgo.Find(context.Background(), &testUser{}, nil)

		// Assert
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestFindOne(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedUser := testUser{ID: bson.NewObjectID(), Name: "Peter"}
		mockDB := &mgo.MockDatastore{
			OnFindOne: mgo.NewOnFindOneMock(expectedUser),
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		var foundUser testUser
		err := mgo.FindOne(context.Background(), &foundUser, nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, foundUser)
	})

	t.Run("Error No Documents", func(t *testing.T) {
		// Arrange
		mockDB := &mgo.MockDatastore{
			OnFindOne: mgo.NewErrOnFindOne(mongo.ErrNoDocuments),
		}
		restore := mgo.SetDatastore(mockDB)
		defer restore()

		// Act
		var foundUser testUser
		err := mgo.FindOne(context.Background(), &foundUser, nil)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
}

func TestFindById(t *testing.T) {
	// Arrange
	userID := bson.NewObjectID()
	expectedUser := testUser{ID: userID, Name: "Peter"}

	mockDB := &mgo.MockDatastore{
		// We mock FindOne because FindById calls it internally.
		OnFindOne: func(ctx context.Context, collection string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
			// Assert that the filter passed by FindById is correct.
			filterMap := filter.(bson.M)
			assert.Equal(t, userID, filterMap["_id"])

			// Return the expected user.
			return mongo.NewSingleResultFromDocument(expectedUser, nil, nil)
		},
	}
	restore := mgo.SetDatastore(mockDB)
	defer restore()

	// Act: Call FindById with a user struct that has the ID we want to find.
	userToFind := &testUser{ID: userID}
	err := mgo.FindById(context.Background(), userToFind)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Name, userToFind.Name)
}
