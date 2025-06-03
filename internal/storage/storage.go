package storage

import (
	"errors"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileExists   = errors.New("file already exists")
	// ErrInvalidFileType = errors.New("invalid file type")
	// ErrStorageNotInitialized = errors.New("storage not initialized")
)
