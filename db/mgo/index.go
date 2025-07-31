// Package mgo provides a high-level abstraction layer over the official MongoDB Go driver,
// simplifying connection management, document operations, and schema definitions.
package mgo

import (
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Index defines the interface for a collection's schema, including its name and index specifications.
// This allows for the automatic creation and registration of indexes.
type Index interface {
	// C returns the name of the MongoDB collection.
	C() string
	// Indexes returns a slice of mongo.IndexModel, defining the indexes for the collection.
	Indexes() []mongo.IndexModel
}

// indexes holds all registered Index definitions for the application.
var indexes = []Index{}

// collectDef is a concrete implementation of the Index interface.
// It is used internally to store registered index information.
type collectDef struct {
	collectionName string
	indexes        []mongo.IndexModel
}

// C returns the collection name.
func (c *collectDef) C() string {
	return c.collectionName
}

// Indexes returns the defined index models.
func (c *collectDef) Indexes() []mongo.IndexModel {
	return c.indexes
}

// NewCollectDef creates a new collection definition that satisfies the Index interface.
// This is a helper function to simplify the creation of index definitions.
func NewCollectDef(name string, f func() []mongo.IndexModel) Index {
	return &collectDef{
		collectionName: name,
		indexes:        f(),
	}
}

// RegisterIndex adds a new Index definition to the global registry.
// This is typically called from the init() function of a model package.
func RegisterIndex(index Index) {
	indexes = append(indexes, index)
}

// CreateIndexesIfNotExists iterates through all registered indexes and creates them in the database
// if the corresponding collection does not already exist. This ensures that the database schema
// is aligned with the application's models on startup.
func CreateIndexesIfNotExists(ctx context.Context) error {
	if conn == nil {
		return ErrNotConnected
	}
	db := conn.db
	// Retrieve a list of all existing collection names in the database.
	dbNames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return errors.Join(ErrListCollectionFailed, err)
	}

	// Create a map for quick lookup of existing collections.
	existingCollections := make(map[string]bool)
	for _, name := range dbNames {
		existingCollections[name] = true
	}

	// Iterate over all programmatically registered indexes.
	for _, index := range indexes {
		// If the collection does not exist in the database, create it and its indexes.
		if !existingCollections[index.C()] {
			// If there are no specific indexes to create, just create the collection itself.
			if len(index.Indexes()) == 0 {
				err = db.CreateCollection(ctx, index.C())
				if err != nil {
					return errors.Join(ErrCreateIndexFailed, err)
				}
				continue
			}
			// Create the defined indexes for the new collection.
			_, err := db.Collection(index.C()).Indexes().CreateMany(ctx, index.Indexes())
			if err != nil {
				return errors.Join(ErrCreateIndexFailed, err)
			}
		}
	}
	return nil
}

// contains is a utility function to check for the existence of a string in a slice.
// NOTE: This function was replaced by a map lookup in CreateIndexesIfNotExists for better performance.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
