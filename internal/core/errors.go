package core

import (
	"errors"
	"fmt"
)

var (
	ErrWorkspaceNotFound = errors.New("not an art workspace")
	ErrEmptyRepository   = errors.New("no commit is found in the repository. please push data to repository first")
)

type ReferenceNotFoundError struct {
	Ref string
	Err error
}

func (err ReferenceNotFoundError) Error() string {
	return fmt.Sprintf("reference not found: %s", err.Ref)
}
