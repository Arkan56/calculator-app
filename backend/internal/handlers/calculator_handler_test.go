package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"calculator-api/internal/handlers"
)

// ─── test router ─────────────────────────────────────────────────────────────

func init() {
	gin.SetMode(gin.TestMode)
}

// newRouter builds a minimal Gin engine with all calculator routes registered,
// mirroring exactly what main.go does.
func newRouter() *gin.Engine {
	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/add",        handlers.Add)
		v1.POST("/subtract",   handlers.Subtract)
		v1.POST("/multiply",   handlers.Multiply)
		v1.POST("/divide",     handlers.Divide)
		v1.POST("/power",      handlers.Power)
		v1.POST("/percentage", handlers.Percentage)
		v1.POST("/sqrt",       handlers.Sqrt)
	}
	return r
}

// post is a helper that fires a POST to the given path with a JSON body and
// returns the recorded response.
func post(t *testing.T, r *gin.Engine, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// assertStatus checks the HTTP status code.
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Errorf("status = %d, want %d  (body: %s)", w.Code, want, w.Body.String())
	}
}

// assertResult decodes the response JSON and checks the "result" field.
func assertResult(t *testing.T, w *httptest.ResponseRecorder, want float64) {
	t.Helper()
	var resp map[string]float64
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v  (body: %s)", err, w.Body.String())
	}
	got, ok := resp["result"]
	if !ok {
		t.Fatalf("response has no 'result' field: %s", w.Body.String())
	}
	if got != want {
		t.Errorf("result = %v, want %v", got, want)
	}
}

// assertError decodes the response JSON and checks the "error" field.
func assertError(t *testing.T, w *httptest.ResponseRecorder, wantMsg string) {
	t.Helper()
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v  (body: %s)", err, w.Body.String())
	}
	got, ok := resp["error"]
	if !ok {
		t.Fatalf("response has no 'error' field: %s", w.Body.String())
	}
	if got != wantMsg {
		t.Errorf("error = %q, want %q", got, wantMsg)
	}
}

// ─── Add ─────────────────────────────────────────────────────────────────────

func TestAddHandler(t *testing.T) {
	r := newRouter()

	t.Run("valid positive numbers", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":3,"b":7}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 10)
	})
	t.Run("zero values accepted", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":0,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("a is zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":0,"b":5}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 5)
	})
	t.Run("negative numbers", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":-3,"b":-7}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, -10)
	})
	t.Run("missing field b", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":5}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("missing field a", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"b":5}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("empty body", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("invalid json", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `not-json`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("wrong type for a", func(t *testing.T) {
		w := post(t, r, "/api/v1/add", `{"a":"abc","b":2}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Subtract ────────────────────────────────────────────────────────────────

func TestSubtractHandler(t *testing.T) {
	r := newRouter()

	t.Run("valid subtraction", func(t *testing.T) {
		w := post(t, r, "/api/v1/subtract", `{"a":10,"b":4}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 6)
	})
	t.Run("b is zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/subtract", `{"a":5,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 5)
	})
	t.Run("both zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/subtract", `{"a":0,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("missing field b", func(t *testing.T) {
		w := post(t, r, "/api/v1/subtract", `{"a":5}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("wrong type for b", func(t *testing.T) {
		w := post(t, r, "/api/v1/subtract", `{"a":5,"b":"x"}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Multiply ────────────────────────────────────────────────────────────────

func TestMultiplyHandler(t *testing.T) {
	r := newRouter()

	t.Run("valid multiplication", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{"a":6,"b":7}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 42)
	})
	t.Run("multiply by zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{"a":99,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("a is zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{"a":0,"b":99}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("negative numbers", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{"a":-3,"b":-4}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 12)
	})
	t.Run("overflow produces ErrInvalidInput", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{"a":1.7976931348623157e+308,"b":2}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "invalid input: result is not a finite number")
	})
	t.Run("missing both fields", func(t *testing.T) {
		w := post(t, r, "/api/v1/multiply", `{}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Divide ──────────────────────────────────────────────────────────────────

func TestDivideHandler(t *testing.T) {
	r := newRouter()

	t.Run("valid division", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":10,"b":2}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 5)
	})
	t.Run("decimal result", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":7,"b":2}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 3.5)
	})
	t.Run("dividend is zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":0,"b":5}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("divide by zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":5,"b":0}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "division by zero is not allowed")
	})
	t.Run("zero divide by zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":0,"b":0}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "division by zero is not allowed")
	})
	t.Run("missing field b", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":5}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("invalid json", func(t *testing.T) {
		w := post(t, r, "/api/v1/divide", `{"a":5,"b":}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Power ───────────────────────────────────────────────────────────────────

func TestPowerHandler(t *testing.T) {
	r := newRouter()

	t.Run("positive base and exponent", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":2,"b":10}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 1024)
	})
	t.Run("zero exponent", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":99,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 1)
	})
	t.Run("base is zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":0,"b":5}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("zero to negative power -> Inf -> 400", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":0,"b":-1}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "invalid input: result is not a finite number")
	})
	t.Run("negative base fractional exp -> NaN -> 400", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":-1,"b":0.5}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "invalid input: result is not a finite number")
	})
	t.Run("missing exponent", func(t *testing.T) {
		w := post(t, r, "/api/v1/power", `{"a":2}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Percentage ──────────────────────────────────────────────────────────────

func TestPercentageHandler(t *testing.T) {
	r := newRouter()

	t.Run("50% of 200", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"a":50,"b":200}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 100)
	})
	t.Run("0% of anything", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"a":0,"b":500}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("percentage of zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"a":50,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("both zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"a":0,"b":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("missing field a", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"b":200}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("wrong type", func(t *testing.T) {
		w := post(t, r, "/api/v1/percentage", `{"a":"fifty","b":200}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ─── Sqrt ────────────────────────────────────────────────────────────────────

func TestSqrtHandler(t *testing.T) {
	r := newRouter()

	t.Run("perfect square", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{"a":16}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 4)
	})
	t.Run("zero", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{"a":0}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 0)
	})
	t.Run("decimal input", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{"a":2.25}`)
		assertStatus(t, w, http.StatusOK)
		assertResult(t, w, 1.5)
	})
	t.Run("negative number -> 400", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{"a":-4}`)
		assertStatus(t, w, http.StatusBadRequest)
		assertError(t, w, "square root of a negative number is not allowed")
	})
	t.Run("missing field a", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("wrong type", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `{"a":"four"}`)
		assertStatus(t, w, http.StatusBadRequest)
	})
	t.Run("invalid json", func(t *testing.T) {
		w := post(t, r, "/api/v1/sqrt", `not-json`)
		assertStatus(t, w, http.StatusBadRequest)
	})
}
