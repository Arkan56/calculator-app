# Calculator App

A full-stack calculator built with a **Go / Gin REST API** backend and a **React / TypeScript / Vite** frontend.

---

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
- [API Reference](#api-reference)
  - [Health Check](#health-check)
  - [Binary Operations](#binary-operations)
  - [Unary Operations](#unary-operations)
  - [Error Responses](#error-responses)
- [API Examples](#api-examples)
- [Design Decisions](#design-decisions)
- [Test Coverage](#test-coverage)

---

## Project Structure

```
calculator-app/
├── backend/
│   ├── cmd/api/
│   │   └── main.go                  # Server entry point, router wiring, CORS
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go        # Domain logic & sentinel errors
│   │   │   └── calculator_test.go  # 100% covered unit tests
│   │   ├── handlers/
│   │   │   ├── calculator_handler.go      # HTTP handlers (one per operation)
│   │   │   └── calculator_handler_test.go # Integration tests via httptest
│   │   └── models/
│   │       └── calculator.go        # Request / response DTOs
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── src/
    │   ├── App.tsx
    │   └── main.tsx
    ├── index.html
    └── package.json
```

---

## Getting Started

### Prerequisites

| Tool | Version |
|------|---------|
| Go   | ≥ 1.21  |
| Node | ≥ 18    |
| npm  | ≥ 9     |

### Backend Setup

```bash
cd backend

# Download dependencies
go mod download

# Run the server (defaults to :8080)
go run ./cmd/api/main.go

# Or set a custom port
PORT=3000 go run ./cmd/api/main.go
```

The API will be available at `http://localhost:8080`.

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start the dev server (http://localhost:5173)
npm run dev
```

> [!NOTE]
> The backend allows CORS requests from `http://localhost:5173` out of the box, so no proxy configuration is needed during development.

---

## API Reference

**Base URL:** `http://localhost:8080/api/v1`

All endpoints accept and return `application/json`.

### Health Check

```
GET /health
```

**Response `200 OK`**
```json
{ "status": "ok" }
```

---

### Binary Operations

All binary endpoints share the same request shape:

```
POST /api/v1/{operation}
```

| Field | Type     | Required | Description   |
|-------|----------|----------|---------------|
| `a`   | `number` | ✅       | First operand |
| `b`   | `number` | ✅       | Second operand |

> [!IMPORTANT]
> Both fields are required. Sending `0` is valid — the API uses pointer fields internally so a zero value is never mistaken for a missing field.

**Available operations:**

| Endpoint            | Formula          |
|---------------------|------------------|
| `POST /api/v1/add`        | `a + b`          |
| `POST /api/v1/subtract`   | `a - b`          |
| `POST /api/v1/multiply`   | `a × b`          |
| `POST /api/v1/divide`     | `a ÷ b`          |
| `POST /api/v1/power`      | `aᵇ`             |
| `POST /api/v1/percentage` | `(a ÷ 100) × b`  |

**Response `200 OK`**
```json
{ "result": 42 }
```

---

### Unary Operations

| Endpoint          | Formula  | Input Field |
|-------------------|----------|-------------|
| `POST /api/v1/sqrt` | `√a`   | `a` only    |

**Request**
```json
{ "a": 16 }
```

**Response `200 OK`**
```json
{ "result": 4 }
```

---

### Error Responses

All errors return `400 Bad Request` with a JSON body:

```json
{ "error": "<human-readable message>" }
```

| Scenario                              | Error message                                        |
|---------------------------------------|------------------------------------------------------|
| Missing required field                | `invalid request: Key: 'BinaryRequest.B' Error:...` |
| Wrong data type (e.g. string instead of number) | `invalid request: ...`                     |
| Malformed JSON                        | `invalid request: ...`                               |
| Division by zero                      | `division by zero is not allowed`                    |
| Square root of a negative number      | `square root of a negative number is not allowed`    |
| Result is NaN or ±Infinity            | `invalid input: result is not a finite number`       |

---

## API Examples

### Add two numbers

```bash
curl -s -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a": 3, "b": 7}'
```
```json
{ "result": 10 }
```

### Subtract (with zero as an operand)

```bash
curl -s -X POST http://localhost:8080/api/v1/subtract \
  -H "Content-Type: application/json" \
  -d '{"a": 0, "b": 5}'
```
```json
{ "result": -5 }
```

### Multiply negative numbers

```bash
curl -s -X POST http://localhost:8080/api/v1/multiply \
  -H "Content-Type: application/json" \
  -d '{"a": -3, "b": -4}'
```
```json
{ "result": 12 }
```

### Divide — success

```bash
curl -s -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 4}'
```
```json
{ "result": 2.5 }
```

### Divide — domain error (÷ 0)

```bash
curl -s -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 5, "b": 0}'
```
```json
{ "error": "division by zero is not allowed" }
```

### Power

```bash
curl -s -X POST http://localhost:8080/api/v1/power \
  -H "Content-Type: application/json" \
  -d '{"a": 2, "b": 10}'
```
```json
{ "result": 1024 }
```

### Percentage — 15% of 200

```bash
curl -s -X POST http://localhost:8080/api/v1/percentage \
  -H "Content-Type: application/json" \
  -d '{"a": 15, "b": 200}'
```
```json
{ "result": 30 }
```

### Square root — success

```bash
curl -s -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" \
  -d '{"a": 16}'
```
```json
{ "result": 4 }
```

### Square root — domain error (negative input)

```bash
curl -s -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" \
  -d '{"a": -4}'
```
```json
{ "error": "square root of a negative number is not allowed" }
```

### Validation error — missing field

```bash
curl -s -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a": 5}'
```
```json
{ "error": "invalid request: Key: 'BinaryRequest.B' Error:Field validation for 'B' failed on the 'required' tag" }
```

### Validation error — wrong type

```bash
curl -s -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a": "hello", "b": 2}'
```
```json
{ "error": "invalid request: ..." }
```

---

## Design Decisions

### Binary Operations over Expression Parsing

The API processes discrete **binary** (`a op b`) and **unary** (`op a`) operations rather than raw mathematical expression strings (e.g. `"3 + 4 * 2"`).

**Why:**

- **Strongly typed contract** — each endpoint receives two typed `float64` fields. The Go struct binding layer validates presence and type before any business logic runs, eliminating an entire class of injection and parsing bugs.
- **Simpler validation & edge-case testing** — testing a function `Divide(a, b float64)` is straightforward; testing a string parser covering operator precedence, parentheses, unary minus, and whitespace is significantly more complex and error-prone.
- **No PEMDAS/BODMAS parser needed** — building a correct expression parser (lexer → AST → evaluator) is non-trivial. Skipping it keeps the backend focused on correctness over a small, well-defined surface area.
- **Client-side chaining** — multi-step calculations (e.g. `(3 + 4) * 2`) are composed on the frontend by chaining sequential API calls, which is a natural fit for an interactive UI where intermediate results are displayed to the user at each step.

**Trade-off accepted:** clients that need to evaluate a full expression in one round-trip cannot do so. This is an intentional constraint — the API is designed for a step-by-step calculator UI, not a general-purpose math engine.

---

### Pointer Fields for Zero-Value Binding

Request structs use `*float64` instead of `float64` for operand fields:

```go
type BinaryRequest struct {
    A *float64 `json:"a" binding:"required"`
    B *float64 `json:"b" binding:"required"`
}
```

Gin's `binding:"required"` tag fails on the **zero value** of a non-pointer type. A payload of `{"a": 0, "b": 5}` would incorrectly return a validation error if `A` were a plain `float64`, because `0` is indistinguishable from "field not provided" to the validator. Pointer types solve this: `nil` means absent, `*0.0` means explicitly provided.

---

### Sentinel Errors at the Domain Layer

Domain errors (`ErrDivisionByZero`, `ErrNegativeSquareRoot`, `ErrInvalidInput`) are defined as package-level `var` values using `errors.New`. Handlers use `errors.Is` to match them, which makes the error contract explicit and testable without relying on string matching.

---

## Test Coverage

Run tests and generate a coverage report:

```bash
cd backend
go test ./internal/... -coverprofile coverage.out
go tool cover -func coverage.out
```

### Results

```
calculator-api/internal/calculator/calculator.go  Add           100.0%
calculator-api/internal/calculator/calculator.go  Subtract      100.0%
calculator-api/internal/calculator/calculator.go  Multiply      100.0%
calculator-api/internal/calculator/calculator.go  Divide        100.0%
calculator-api/internal/calculator/calculator.go  Power         100.0%
calculator-api/internal/calculator/calculator.go  Percentage    100.0%
calculator-api/internal/calculator/calculator.go  Sqrt          100.0%
calculator-api/internal/calculator/calculator.go  validate      100.0%
calculator-api/internal/handlers/calculator_handler.go  Add         100.0%
calculator-api/internal/handlers/calculator_handler.go  Subtract    100.0%
calculator-api/internal/handlers/calculator_handler.go  Multiply    100.0%
calculator-api/internal/handlers/calculator_handler.go  Divide      100.0%
calculator-api/internal/handlers/calculator_handler.go  Power       100.0%
calculator-api/internal/handlers/calculator_handler.go  Percentage  100.0%
calculator-api/internal/handlers/calculator_handler.go  Sqrt        100.0%
calculator-api/internal/handlers/calculator_handler.go  handleBinary    83.3%
calculator-api/internal/handlers/calculator_handler.go  handleUnary     83.3%
calculator-api/internal/handlers/calculator_handler.go  mapDomainError  80.0%
total:                                                              91.2%
```

| Package                      | Coverage  |
|------------------------------|-----------|
| `internal/calculator`        | **100.0%** |
| `internal/handlers`          | **86.1%** |
| **Total**                    | **91.2%** |

> [!NOTE]
> The two branches below 100% are intentionally unreachable via the HTTP layer:
> - **`handleBinary` / `handleUnary` NaN/Inf guard (83.3%)** — standard JSON does not support `NaN` or `Infinity` literals, so this guard can never be triggered through Gin's JSON parser. It exists as a defensive layer if handlers are ever called directly.
> - **`mapDomainError` default branch (80.0%)** — a catch-all for future, currently non-existent error types. No domain function produces an unmapped error today.
