package constants

type contextKey string

const (
	LoggerKey    = contextKey("logger")
	RequestIDKey = contextKey("requestID")
)

const (
	Discriminator = "marketplace"
)
