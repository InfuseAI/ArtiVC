package core

import "errors"

var (
	ErrReferenceNotFound = errors.New("reference not found")
	ErrWorkspaceNotFound = errors.New("workspace not found")
)
