package client

import (
	"context"
)

type Client interface {
	Invoke(ctx context.Context, address, serviceName, methodName string, req []byte) ([]byte, error)
}
