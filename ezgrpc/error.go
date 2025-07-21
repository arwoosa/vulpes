// Package ezgrpc provides a simplified setup for gRPC services with a grpc-gateway.
// It includes utilities for handling cookies, sessions, and standard interceptors.
package ezgrpc

import (
	"errors"

	wrapperErr "github.com/arwoosa/vulpes/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Pre-defined gRPC status errors for session-related issues.
// These provide consistent error responses for clients.
var (
	// Status_EZgrpc_Session_NotFound indicates that a session could not be found.
	Status_EZgrpc_Session_NotFound = status.New(codes.NotFound, "session not found")
	// Status_EZgrpc_Session_SaveFailed indicates that a session failed to save.
	Status_EZgrpc_Session_SaveFailed = status.New(codes.Internal, "session save failed")

	// Err_SessionNotFound is the underlying error for a missing session.
	Err_SessionNotFound = errors.New("session not found")
	// Err_SessionSaveFailed is the underlying error for a session save failure.
	Err_SessionSaveFailed = errors.New("session save failed")
)

// ToStatus converts a custom wrapper error into a gRPC status.Status.
// This allows for detailed error information to be sent to the client,
// including a descriptive message and structured details.
//
// err: The custom error with a message.
// Returns a gRPC status, or nil if the input error is nil.
func ToStatus(err wrapperErr.ErrorWithMessage) *status.Status {
	if err == nil {
		return nil
	}
	var baseSt *status.Status
	switch err.Err() {
	case Err_SessionNotFound:
		baseSt = Status_EZgrpc_Session_NotFound
	case Err_SessionSaveFailed:
		baseSt = Status_EZgrpc_Session_SaveFailed
	default:
		// For unhandled errors, create a generic internal error status.
		return status.New(codes.Internal, err.Err().Error())
	}

	// Add more details to the status, such as the type of violation and a description.
	st, myErr := baseSt.WithDetails(
		&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "EZGRPE_SESSION",
					Subject:     err.Err().Error(),
					Description: err.Msg(),
				},
			},
		},
	)
	if myErr != nil {
		// If adding details fails, return the original base status.
		return baseSt
	}
	return st
}
