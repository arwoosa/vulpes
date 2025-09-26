package weaviatego

import (
	"github.com/google/uuid"
)

var namespace = uuid.MustParse("a638b704-c36b-4e89-a2e6-77d9c6e3b56a")

func NewUUIDFromString(name string) uuid.UUID {
	return uuid.NewSHA1(namespace, []byte(name))
}

func NewUUID() uuid.UUID {
	return uuid.New()
}
