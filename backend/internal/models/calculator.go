package models

// BinaryRequest holds two operands for binary operations.
// Pointers are used so that 0 is accepted as a valid value by Gin's
// "required" binding (a zero float64 is not treated as "missing").
type BinaryRequest struct {
	A *float64 `json:"a" binding:"required"`
	B *float64 `json:"b" binding:"required"`
}

// UnaryRequest holds a single operand for unary operations (e.g. sqrt).
type UnaryRequest struct {
	A *float64 `json:"a" binding:"required"`
}

// CalculationResponse wraps the numeric result of a successful calculation.
type CalculationResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse wraps a human-readable error message returned on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}
