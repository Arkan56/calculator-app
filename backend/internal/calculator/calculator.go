package calculator

import (
	"errors"
	"math"
)

// Sentinel errors returned by the domain layer.
var (
	ErrDivisionByZero      = errors.New("division by zero is not allowed")
	ErrNegativeSquareRoot  = errors.New("square root of a negative number is not allowed")
	ErrInvalidInput        = errors.New("invalid input: result is NaN or Inf")
)

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	result := a + b
	return validate(result)
}

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	result := a - b
	return validate(result)
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	result := a * b
	return validate(result)
}

// Divide returns a / b. Returns ErrDivisionByZero when b == 0.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	result := a / b
	return validate(result)
}

// Power returns a raised to the power of b (a^b).
func Power(a, b float64) (float64, error) {
	result := math.Pow(a, b)
	return validate(result)
}

// Percentage returns a% of b, i.e. (a / 100) * b.
func Percentage(a, b float64) (float64, error) {
	result := (a / 100) * b
	return validate(result)
}

// Sqrt returns the square root of a. Returns ErrNegativeSquareRoot when a < 0.
func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	result := math.Sqrt(a)
	return validate(result)
}

// validate rejects NaN and ±Inf results that can arise from edge-case inputs
// (e.g. 0^0 edge cases, very large numbers, etc.).
func validate(v float64) (float64, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, ErrInvalidInput
	}
	return v, nil
}
