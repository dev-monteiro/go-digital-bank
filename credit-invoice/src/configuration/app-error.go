package configuration

// TODO: maybe use the error interface?
type AppError struct {
	Message    string
	StatusCode int
}
