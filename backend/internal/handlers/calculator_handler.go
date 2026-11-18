package handlers

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"calculator-api/internal/calculator"
	"calculator-api/internal/models"
)

// ─── Individual handlers ──────────────────────────────────────────────────────

// Add handles POST /api/v1/add
func Add(c *gin.Context) {
	handleBinary(c, calculator.Add)
}

// Subtract handles POST /api/v1/subtract
func Subtract(c *gin.Context) {
	handleBinary(c, calculator.Subtract)
}

// Multiply handles POST /api/v1/multiply
func Multiply(c *gin.Context) {
	handleBinary(c, calculator.Multiply)
}

// Divide handles POST /api/v1/divide
func Divide(c *gin.Context) {
	handleBinary(c, calculator.Divide)
}

// Power handles POST /api/v1/power
func Power(c *gin.Context) {
	handleBinary(c, calculator.Power)
}

// Percentage handles POST /api/v1/percentage
func Percentage(c *gin.Context) {
	handleBinary(c, calculator.Percentage)
}

// Sqrt handles POST /api/v1/sqrt
func Sqrt(c *gin.Context) {
	handleUnary(c, calculator.Sqrt)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// handleBinary binds a BinaryRequest, calls op, and writes the result or error.
func handleBinary(c *gin.Context, op func(a, b float64) (float64, error)) {
	var req models.BinaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}

	if math.IsNaN(*req.A) || math.IsInf(*req.A, 0) ||
		math.IsNaN(*req.B) || math.IsInf(*req.B, 0) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid input: values must be finite numbers"})
		return
	}

	result, err := op(*req.A, *req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: mapDomainError(err)})
		return
	}

	c.JSON(http.StatusOK, models.CalculationResponse{Result: result})
}

// handleUnary binds a UnaryRequest, calls op, and writes the result or error.
func handleUnary(c *gin.Context, op func(a float64) (float64, error)) {
	var req models.UnaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}

	if math.IsNaN(*req.A) || math.IsInf(*req.A, 0) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid input: value must be a finite number"})
		return
	}

	result, err := op(*req.A)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: mapDomainError(err)})
		return
	}

	c.JSON(http.StatusOK, models.CalculationResponse{Result: result})
}

// mapDomainError converts known domain errors into user-friendly messages.
func mapDomainError(err error) string {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		return "division by zero is not allowed"
	case errors.Is(err, calculator.ErrNegativeSquareRoot):
		return "square root of a negative number is not allowed"
	case errors.Is(err, calculator.ErrInvalidInput):
		return "invalid input: result is not a finite number"
	default:
		return err.Error()
	}
}
