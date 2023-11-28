package invostatus

type InvoStatus string

const (
	CLOSED InvoStatus = "CLOSED"
	OPEN   InvoStatus = "OPEN"
	FUTURE InvoStatus = "FUTURE"
)
