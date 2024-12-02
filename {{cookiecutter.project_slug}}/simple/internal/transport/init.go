package transport

import (
	"{{cookiecutter.project_slug}}/internal/pkg/logging"
)

// AuthService defines authentication-related methods.
// Extend this interface with the actual methods you need for authentication.
type AuthService interface {
	// Example method: AuthenticateUser(username, password string) (bool, error)
}

// UseCase aggregates all the use case interfaces needed for the transport layer.
type UseCase interface {
	AuthService // Example: Include the authentication use case
}

// Handler is responsible for managing transport-level interactions.
// It depends on the UseCase layer for business logic.
type Handler struct {
	useCase UseCase        // UseCase aggregates business logic interfaces.
	logger  logging.Logger // Logger for logging transport events and errors.
}

// NewHandler creates a new Handler instance.
func NewHandler(useCase UseCase, logger logging.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		logger:  logger,
	}
}

// Run initializes and runs the transport layer.
// Extend this method to start HTTP, gRPC, or other types of servers.
func (h *Handler) Run() {
	h.logger.InfoWithoutContext("Transport layer is running...")
	// TODO: Implement transport logic, such as starting an HTTP or gRPC server or some broker.
}
