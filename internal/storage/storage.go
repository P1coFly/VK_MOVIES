package storage

import "errors"

var (
	ErrActorNotFound = errors.New("actor not found")
	ErrMovieNotFound = errors.New("movie not found")
)
