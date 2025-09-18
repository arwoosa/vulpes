package mgo

import "go.mongodb.org/mongo-driver/v2/bson"

func NewObjectID() ObjectID {
	return ObjectID{bson.NewObjectID()}
}

type ObjectID struct {
	bson.ObjectID
}

func (d ObjectID) GetObjectId() bson.ObjectID {
	return d.ObjectID
}

func (d ObjectID) GetId() any {
	if d.IsZero() {
		return nil
	}
	return d
}

func (d *ObjectID) SetId(id any) {
	oid, ok := id.(bson.ObjectID)
	if ok {
		d.ObjectID = oid
		return
	}
}
