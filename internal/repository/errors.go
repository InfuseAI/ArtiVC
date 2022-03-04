package repository

import "errors"

var ErrUnsupportedRepository = errors.New("Unsupported repository")

type UnsupportedRepositoryError struct {
	Message string
}

func (err UnsupportedRepositoryError) Error() string {
	return err.Message
}
