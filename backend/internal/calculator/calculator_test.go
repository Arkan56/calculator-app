package calculator_test

import (
	"errors"
	"math"
	"testing"

	"calculator-api/internal/calculator"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

func approxEqual(a, b float64) bool {
	if math.IsInf(a, 0) || math.IsInf(b, 0) || math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	const eps = 1e-9
	return math.Abs(a-b) < eps
}

// ─── Add ─────────────────────────────────────────────────────────────────────

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 3, 7, 10, nil},
		{"negative numbers", -4, -6, -10, nil},
		{"mixed signs", -3, 5, 2, nil},
		{"zeros", 0, 0, 0, nil},
		{"decimals", 1.1, 2.2, 3.3, nil},
		{"large numbers", 1e308, 0, 1e308, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Add(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Add(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Add(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Subtract ────────────────────────────────────────────────────────────────

func TestSubtract(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 10, 4, 6, nil},
		{"negative numbers", -3, -7, 4, nil},
		{"result is zero", 5, 5, 0, nil},
		{"result is negative", 2, 8, -6, nil},
		{"decimals", 5.5, 2.2, 3.3, nil},
		{"zero minus zero", 0, 0, 0, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Subtract(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Subtract(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Multiply ────────────────────────────────────────────────────────────────

func TestMultiply(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 3, 4, 12, nil},
		{"negative × positive", -3, 4, -12, nil},
		{"negative × negative", -3, -4, 12, nil},
		{"multiply by zero", 99, 0, 0, nil},
		{"decimals", 2.5, 4, 10, nil},
		{"identity", 1, 42, 42, nil},
		{"overflow to Inf", math.MaxFloat64, 2, 0, calculator.ErrInvalidInput},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Multiply(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Multiply(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Divide ──────────────────────────────────────────────────────────────────

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 10, 2, 5, nil},
		{"negative dividend", -10, 2, -5, nil},
		{"negative divisor", 10, -2, -5, nil},
		{"both negative", -10, -2, 5, nil},
		{"decimal result", 7, 2, 3.5, nil},
		{"zero dividend", 0, 5, 0, nil},
		{"divide by zero", 5, 0, 0, calculator.ErrDivisionByZero},
		{"negative divide by zero", -5, 0, 0, calculator.ErrDivisionByZero},
		{"zero divide by zero", 0, 0, 0, calculator.ErrDivisionByZero},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Divide(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Divide(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Divide(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Power ───────────────────────────────────────────────────────────────────

func TestPower(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive base and exp", 2, 10, 1024, nil},
		{"zero exponent", 99, 0, 1, nil},
		{"zero base", 0, 5, 0, nil},
		{"negative base even exp", -2, 2, 4, nil},
		{"negative base odd exp", -2, 3, -8, nil},
		{"fractional exponent", 9, 0.5, 3, nil},
		{"negative exponent", 2, -1, 0.5, nil},
		// 0^(-1) → Inf → ErrInvalidInput
		{"zero to negative power", 0, -1, 0, calculator.ErrInvalidInput},
		// (-1)^0.5 → NaN → ErrInvalidInput
		{"negative base fractional exp", -1, 0.5, 0, calculator.ErrInvalidInput},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Power(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Power(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Power(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Percentage ──────────────────────────────────────────────────────────────

func TestPercentage(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"50% of 200", 50, 200, 100, nil},
		{"10% of 150", 10, 150, 15, nil},
		{"100% of 75", 100, 75, 75, nil},
		{"0% of anything", 0, 999, 0, nil},
		{"percentage of zero", 50, 0, 0, nil},
		{"negative percentage", -10, 100, -10, nil},
		{"decimal percentage", 12.5, 80, 10, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Percentage(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Percentage(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ─── Sqrt ────────────────────────────────────────────────────────────────────

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{"perfect square", 16, 4, nil},
		{"zero", 0, 0, nil},
		{"decimal", 2.25, 1.5, nil},
		{"one", 1, 1, nil},
		{"large number", 1e10, 1e5, nil},
		{"negative number", -4, 0, calculator.ErrNegativeSquareRoot},
		{"negative one", -1, 0, calculator.ErrNegativeSquareRoot},
		{"large negative", -1e10, 0, calculator.ErrNegativeSquareRoot},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Sqrt(tc.a)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Sqrt(%v) error = %v, want %v", tc.a, err, tc.wantErr)
			}
			if tc.wantErr == nil && !approxEqual(got, tc.want) {
				t.Errorf("Sqrt(%v) = %v, want %v", tc.a, got, tc.want)
			}
		})
	}
}
