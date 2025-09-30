package weaviatego

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {
	// test invalid host
	t.Run("InvalidHost", func(t *testing.T) {
		_, err := InitClient(context.Background(), "invalid-host", "test")
		assert.Error(t, err)
	})
}

type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (d *TestData) ClassName() string {
	return "Test"
}

func (d *TestData) ID() uuid.UUID {
	return NewUUIDFromString(d.Name)
}

func TestGenerateUUIDFromString(t *testing.T) {
	t.Run("GenerateUUIDFromString", func(t *testing.T) {
		uuid := NewUUIDFromString("test")
		uuid2 := NewUUIDFromString("test")
		uuid3 := NewUUIDFromString("test2")
		assert.NotEqual(t, uuid, uuid3)
		assert.Equal(t, uuid, uuid2)
	})
}
