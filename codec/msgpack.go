// Package codec provides a flexible framework for encoding and decoding data structures.
// It supports multiple encoding formats (GOB, MessagePack) and uses generics for type safety.
package codec

import (
	"encoding/base64"

	"github.com/arwoosa/vulpes/errors"
	"github.com/vmihailenco/msgpack/v5"
)

// msgpackCodec implements the Codec interface using the MessagePack serialization format.
// The binary output is then encoded into a Base64 string for safe transport.
type msgpackCodec[T any] struct{}

// Encode serializes the value `v` first using MessagePack, then encodes the result into a Base64 string.
func (c *msgpackCodec[T]) Encode(v T) (string, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return "", errors.NewWrapperError(ERR_MsgPackEncodeFailed, err.Error())
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode first decodes the Base64 string `s` into bytes, then deserializes the bytes using MessagePack.
func (c *msgpackCodec[T]) Decode(s string) (T, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return *new(T), errors.NewWrapperError(ERR_Base64DecodeFailed, err.Error())
	}
	var v T
	if err := msgpack.Unmarshal(b, &v); err != nil {
		return *new(T), errors.NewWrapperError(ERR_MsgPackDecodeFailed, err.Error())
	}
	return v, nil
}

// Method returns the MessagePack codec method identifier.
func (c *msgpackCodec[T]) Method() CodecMethod {
	return MSGPACK
}

// encodeMsgPack is a package-level helper function for MessagePack encoding.
func encodeMsgPack[T any](v T) (string, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return "", errors.NewWrapperError(ERR_MsgPackEncodeFailed, err.Error())
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// decodeMsgPack is a package-level helper function for MessagePack decoding.
func decodeMsgPack[T any](s string) (T, error) {
	var out T
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return out, errors.NewWrapperError(ERR_Base64DecodeFailed, err.Error())
	}
	err = msgpack.Unmarshal(b, &out)
	if err != nil {
		return out, errors.NewWrapperError(ERR_MsgPackDecodeFailed, err.Error())
	}
	return out, nil
}
