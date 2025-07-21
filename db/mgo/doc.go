// Package mgo provides a high-level abstraction layer over the official MongoDB Go driver,
// simplifying connection management, document operations, and schema definitions.
package mgo

// DocInter defines a standard interface for all database documents.
// By embedding the Index interface, it requires all documents to specify
// their collection name, index definitions, and validation logic.
// This promotes a consistent and self-describing data model.
type DocInter interface {
	Index
	// Validate performs business logic validation on the document's fields.
	Validate() error
}

// DocSlice represents a collection of documents that all conform to the DocInter interface.
type DocSlice []DocInter

// NewDocSlice creates and returns an empty DocSlice, ready to be populated.
func NewDocSlice() DocSlice {
	return DocSlice{}
}

// Append adds a new document to the slice.
func (d *DocSlice) Append(doc DocInter) {
	*d = append(*d, doc)
}

// Len returns the total number of documents in the slice.
func (d DocSlice) Len() int {
	return len(d)
}

// Get retrieves a document from the slice at the specified index.
func (d DocSlice) Get(i int) DocInter {
	return d[i]
}
