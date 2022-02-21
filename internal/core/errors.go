package core

import (
	"errors"
	"fmt"
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
)

type ReferenceNotFoundError struct {
	Ref string
	Err error
}

func (err ReferenceNotFoundError) Error() string {
	return fmt.Sprintf("reference not found: %s", err.Ref)
}
