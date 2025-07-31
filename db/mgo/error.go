// Package mgo provides a high-level abstraction layer over the official MongoDB Go driver,
// simplifying connection management, document operations, and schema definitions.
package mgo

import "errors"

// Standardized errors for the mgo package, providing consistent error types for common database operations.
var (
	// ErrInvalidDocument is returned when a document fails validation or is otherwise malformed.
	ErrInvalidDocument = errors.New("invalid document")
	// ErrNotConnected is returned when a database operation is attempted before a connection is established.
	ErrNotConnected = errors.New("mongodb not connected")
	// ErrConnectionFailed is returned when the initial connection attempt to the MongoDB server fails.
	ErrConnectionFailed = errors.New("mongodb connection failed")
	// ErrPingFailed is returned when the client fails to ping the MongoDB server to verify the connection.
	ErrPingFailed = errors.New("mongodb ping failed")
	// ErrCreateIndexFailed is returned when an attempt to create one or more indexes fails.
	ErrCreateIndexFailed = errors.New("mongodb create index failed")
	// ErrListCollectionFailed is returned when listing the collections in the database fails.
	ErrListCollectionFailed = errors.New("mongodb list collection failed")
	// ErrWriteFailed is returned when a write operation fails.
	ErrWriteFailed = errors.New("mongodb write failed")
	// ErrReadFailed is returned when a read operation fails.
	ErrReadFailed = errors.New("mongodb read failed")
)
