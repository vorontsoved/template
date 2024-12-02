package usecase

import (
	"{{cookiecutter.project_slug}}/internal/pkg/logging"
)

// Storage defines the methods that the storage layer must implement.
type Storage interface {
	// More methods as needed
}

// CoreUseCase is the main use case implementation.
type CoreUseCase struct {
	storage Storage
	logger  logging.Logger
	// More services as needed
}

// New creates a new instance of CoreUseCase.
func New(storage Storage, logger logging.Logger) *CoreUseCase {
	return &CoreUseCase{
		storage: storage,
		logger:  logger,
	}
}
