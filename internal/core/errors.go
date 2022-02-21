package core

import (
	"errors"
	"fmt"
)

var (
	ErrWorkspaceNotFound = errors.New("not an art workspace")
)

type ReferenceNotFoundError struct {
	Ref string
	Err error
}

func (err ReferenceNotFoundError) Error() string {
	return fmt.Sprintf("reference not found: %s", err.Ref)
}
