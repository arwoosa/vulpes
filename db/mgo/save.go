package mgo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

func Save[T DocInter](ctx context.Context, doc T) (T, error) {
	if v := reflect.ValueOf(doc); v.Kind() == reflect.Ptr && v.IsNil() {
		var zero T // 宣告一個 T 型別的零值
		return zero, errors.Join(ErrInvalidDocument, errors.New("document cannot be nil"))
	}
	if err := doc.Validate(); err != nil {
		return doc, fmt.Errorf("%w: %w", ErrInvalidDocument, err)
	}
	collection := GetCollection(doc.C())
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return doc, fmt.Errorf("%w: %w", ErrWriteFailed, err)
		// return doc, errors.Join(ErrWriteFailed, err)
	}
	doc.SetId(result.InsertedID)
	return doc, nil
}
