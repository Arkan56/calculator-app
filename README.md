# Calculator App

A full-stack calculator built with a **Go / Gin REST API** backend and a **React / TypeScript / Vite** frontend.

---

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Running with Docker](#running-with-docker)
- [Frontend](#frontend)
  - [Technologies](#technologies)
  - [Source layout](#source-layout)
  - [Features](#features)
  - [How calculations are handled](#how-calculations-are-handled)
  - [npm scripts](#npm-scripts)
  - [Frontend tests](#frontend-tests)
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
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── calculatorApi.ts       # fetch client for /api/v1
│   │   ├── components/
│   │   │   ├── Display.tsx            # Readout, pending operation, errors
│   │   │   └── Keypad.tsx             # Buttons for digits and operations
│   │   ├── hooks/
│   │   │   └── useCalculator.ts       # Input state and operation chaining
│   │   ├── types/
│   │   │   └── calculator.ts          # Operation names and request/response types
│   │   ├── __tests__/
│   │   │   ├── Calculator.test.tsx    # UI interaction tests
│   │   │   └── calculatorApi.test.ts  # fetch client tests
│   │   ├── test/
│   │   │   └── setup.ts               # Vitest + Testing Library setup
│   │   ├── App.tsx                    # Layout and keyboard handling
│   │   ├── main.tsx                   # React entry
│   │   └── index.css                  # Tailwind directives
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts                 # Vite config and Vitest options
├── Dockerfile                         # Multi-stage build with `backend` and `frontend` targets
├── docker-compose.yml                 # Runs both services together
├── nginx.conf                         # Serves the SPA and proxies /api to the backend
└── .dockerignore                      # Keeps the Docker build context small
```

---

## Getting Started

### Prerequisites

| Tool              | Version                                         |
| ----------------- | ----------------------------------------------- |
| Go                | ≥ 1.26.10                                       |
| Node              | ≥ 20.19 (22 recommended)                        |
| npm               | ≥ 9                                             |
| Docker (optional) | Docker Desktop or Docker Engine with Compose v2 |

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

Leave this process running, then start the frontend in a second terminal. The keypad calls the API for every operation, so calculations fail until the server is up.

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start the dev server (http://localhost:5173)
npm run dev
```

The client posts to `http://localhost:8080/api/v1` unless `VITE_API_URL` is set. That variable replaces the whole base URL (include `/api/v1`, not just the host).

> [!NOTE]
> The backend allows CORS requests from `http://localhost:5173` out of the box, so no proxy configuration is needed during development.

To run everything without installing Go or Node locally, see [Running with Docker](#running-with-docker).

### Running with Docker

Run the whole app (API and UI) with one command from the project root:

```bash
docker compose up --build
```

Then open `http://localhost:3000`.

| Service    | Image                        | Role                                                          |
| ---------- | ---------------------------- | ------------------------------------------------------------- |
| `backend`  | Go binary on distroless      | REST API on port 8080 (internal only)                         |
| `frontend` | nginx serving the Vite build | Serves the UI on port 3000 and proxies `/api/` to the backend |

How it fits together:

- The `Dockerfile` has two targets, `backend` and `frontend`, which `docker-compose.yml` builds and runs.
- The frontend is built with `VITE_API_URL=/api/v1`, so the browser calls the same origin and nginx forwards the request to the backend. No CORS configuration is needed in this setup.
- `VITE_API_URL` is baked in at **build time**. After changing it, rebuild with `docker compose up --build`.
- The backend port is not published by default. To call the API directly (for example with `curl`), uncomment the `ports` entry on the `backend` service in `docker-compose.yml`.
- If you change the backend `PORT`, update the `proxy_pass` line in `nginx.conf` to match.
- The Go version in the `Dockerfile` (`GO_VERSION`) must satisfy the requirement in `backend/go.mod`.

Stop and clean up:

```bash
docker compose down
```

---

## Frontend

The UI is a single calculator screen. It keeps the current entry locally and sends one operation at a time to the API described in the [API Reference](#api-reference). It does not parse full expressions.

### Technologies

Runtime dependencies are React 19 (`react`, `react-dom`). Everything else is a dev dependency. Ranges below match `frontend/package.json`:

| Area  | Packages                                                                                                                     |
| ----- | ---------------------------------------------------------------------------------------------------------------------------- |
| Build | Vite 8, `@vitejs/plugin-react`, TypeScript 6                                                                                 |
| Style | Tailwind CSS 3, PostCSS, Autoprefixer                                                                                        |
| Lint  | ESLint 10, `typescript-eslint`, `eslint-plugin-react-hooks`, `eslint-plugin-react-refresh`                                   |
| Test  | Vitest 5, jsdom, `@testing-library/react`, `@testing-library/user-event`, `@testing-library/jest-dom`, `@vitest/coverage-v8` |

`vite.config.ts` sets the Vitest environment to jsdom and loads `src/test/setup.ts`, which registers jest-dom matchers and unmounts the rendered tree after each test. Coverage uses the v8 provider and writes text plus HTML under `frontend/coverage/`. `main.tsx`, `src/test/**`, and test files are excluded from that report.

### Source layout

| Path                         | Role                                                        |
| ---------------------------- | ----------------------------------------------------------- |
| `src/main.tsx`               | Mounts `<App />` in `#root`                                 |
| `src/App.tsx`                | Centers the calculator and maps keyboard input to the hook  |
| `src/hooks/useCalculator.ts` | Display state, chaining, and in-flight request lock         |
| `src/api/calculatorApi.ts`   | `calculateBinary` and `calculateUnary`                      |
| `src/types/calculator.ts`    | Operation slugs and the JSON shapes shared with the Go DTOs |
| `src/components/Display.tsx` | Pending line, error alert, numeric readout                  |
| `src/components/Keypad.tsx`  | Digit, utility, operator, and equals buttons                |
| `src/index.css`              | `@tailwind` base, components, and utilities                 |

Tailwind scans `index.html` and `src/**/*.{js,ts,jsx,tsx}`.

### Features

The keypad is a 4-column grid:

| Control                       | Action                                                                                |
| ----------------------------- | ------------------------------------------------------------------------------------- |
| `0`–`9`                       | Append a digit. A leading `0` or `-0` is replaced, so `05` becomes `5`.               |
| `.`                           | Append one decimal point. A second `.` is ignored (`1.2.3` stays `1.23`).             |
| `+/-`                         | Toggle the sign. Right after an operator, the next operand starts as `-0`.            |
| `DEL`                         | Delete the last character, unless the screen is showing a result or a stored operand. |
| `AC`                          | Reset the entry, the pending operation, and any error.                                |
| `+`, `−`, `×`, `÷`, `xʸ`, `%` | Binary operations: add, subtract, multiply, divide, power, percentage.                |
| `√`                           | Square root. It runs immediately and does not need `=`.                               |
| `=`                           | Evaluate the pending binary operation.                                                |

While a request is in flight, every button is disabled and the readout dims. A spinner is exposed to assistive tech as "Calculating". Domain and network failures render in an alert above the number; entering a digit, a decimal, toggling the sign, backspacing, or pressing `AC` clears that alert.

Keyboard input is handled on `window` and ignored when Ctrl, Meta, or Alt is held, when the event target is an input, textarea, select, or contenteditable element, and while a request is loading or the key is repeating.

| Key                          | Action                                             |
| ---------------------------- | -------------------------------------------------- |
| `0`–`9`                      | Digit                                              |
| `.` or `,`                   | Decimal point                                      |
| `+`, `-`, `*`, `/`, `^`, `%` | Add, subtract, multiply, divide, power, percentage |
| `Enter` or `=`               | Equals                                             |
| `Backspace`                  | Delete                                             |
| `Escape`                     | Clear                                              |

`√` and `+/-` are pointer-only.

### How calculations are handled

`useCalculator` stores the current entry as a string, plus the previous operand and the pending operation. Pressing a binary operator stores the current number and waits for the next one. Pressing a different operator before that next number only swaps the pending operation. Pressing an operator when a previous operation is already ready evaluates that operation first, then chains the new one. `=` evaluates and clears the pending operation. `√` replaces the current entry with the result and, if a binary operation is pending, that result becomes its next operand.

Each evaluation calls the API once:

| UI             | Client                             | Request                                       |
| -------------- | ---------------------------------- | --------------------------------------------- |
| `+ − × ÷ xʸ %` | `calculateBinary(operation, a, b)` | `POST {base}/{operation}` with `{ "a", "b" }` |
| `√`            | `calculateUnary('sqrt', a)`        | `POST {base}/sqrt` with `{ "a" }`             |

`{base}` is `import.meta.env.VITE_API_URL` or `http://localhost:8080/api/v1`. Both functions send `Content-Type: application/json`. A 2xx body must be `{ "result": number }`. Any other status throws an `Error` whose message is the JSON `error` field, or `An unexpected error occurred.` when that field is missing. If `fetch` itself rejects, the message is `Network error: unable to reach the calculator API.` A ref blocks a second request until the first settles, because `isLoading` updates asynchronously.

Request and response fields, status codes, and domain errors are specified in the [API Reference](#api-reference). Why the API is one operation per call, rather than an expression string, is covered under [Design Decisions](#design-decisions).

### npm scripts

Run these from `frontend/`:

| Script                  | Runs                    | Purpose                                                |
| ----------------------- | ----------------------- | ------------------------------------------------------ |
| `npm run dev`           | `vite`                  | Dev server with HMR at `http://localhost:5173`         |
| `npm run build`         | `tsc -b && vite build`  | Typecheck, then production bundle                      |
| `npm run preview`       | `vite preview`          | Serve the production bundle                            |
| `npm run lint`          | `eslint .`              | Lint the project (`dist/` and `coverage/` are ignored) |
| `npm test`              | `vitest run`            | Run the suite once                                     |
| `npm run test:coverage` | `vitest run --coverage` | Run the suite and write a v8 report                    |

### Frontend tests

```bash
cd frontend
npm test
npm run test:coverage
```

`npm test` runs 14 tests in two files. Latest run: **14 passed**.

`src/__tests__/calculatorApi.test.ts` stubs `global.fetch` and checks that each binary operation and `sqrt` hit the matching URL with the expected JSON body. It also checks that a `400` body of `{ "error": "division by zero is not allowed" }` rejects with that message.

`src/__tests__/Calculator.test.tsx` mocks the API module and drives the screen with Testing Library and `user-event`:

- Clicks and keystrokes both keep a single decimal (`1.2.3` → `1.23`).
- `8`, `+`, `4`, `=` by click and by `8+4` then Enter both render `12` and call `calculateBinary('add', 8, 4)`.
- `1`, `6`, `√` renders `4` without pressing equals.
- A rejected divide shows `cannot divide by zero`, and typing another digit removes the alert.

`npm run test:coverage` on that same suite:

```
File                 Stmts    Branch   Funcs    Lines
All files            68.68%   59.55%   83.33%   73.86%
src/App.tsx          75%      69.04%   100%     82.14%
src/api/calculatorApi.ts  92.3%  83.33%  100%   92.3%
src/components (rollup)   78.26% 77.27%  100%   85%
src/components/Display.tsx  50%  68.75%  100%   57.14%
src/hooks/useCalculator.ts  63.07% 45.45% 68.42% 67.82%
```

Totals: statements 136/198, branches 81/136, functions 30/36, lines 130/176. The text reporter does not print a separate row for `Keypad.tsx`; that file is part of the `src/components` rollup. `src/types/calculator.ts` has no executable code. Uncovered lines are the network-failure branch in the API client, the longer-readout font steps in `Display`, the Backspace/Escape keyboard paths in `App`, and several hook paths (sign toggle, clear, backspace, operator swap, chaining, and invalid input).

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

| Field | Type     | Required | Description    |
| ----- | -------- | -------- | -------------- |
| `a`   | `number` | ✅       | First operand  |
| `b`   | `number` | ✅       | Second operand |

> [!IMPORTANT]
> Both fields are required. Sending `0` is valid — the API uses pointer fields internally so a zero value is never mistaken for a missing field.

**Available operations:**

| Endpoint                  | Formula         |
| ------------------------- | --------------- |
| `POST /api/v1/add`        | `a + b`         |
| `POST /api/v1/subtract`   | `a - b`         |
| `POST /api/v1/multiply`   | `a × b`         |
| `POST /api/v1/divide`     | `a ÷ b`         |
| `POST /api/v1/power`      | `aᵇ`            |
| `POST /api/v1/percentage` | `(a ÷ 100) × b` |

**Response `200 OK`**

```json
{ "result": 42 }
```

---

### Unary Operations

| Endpoint            | Formula | Input Field |
| ------------------- | ------- | ----------- |
| `POST /api/v1/sqrt` | `√a`    | `a` only    |

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

| Scenario                                        | Error message                                       |
| ----------------------------------------------- | --------------------------------------------------- |
| Missing required field                          | `invalid request: Key: 'BinaryRequest.B' Error:...` |
| Wrong data type (e.g. string instead of number) | `invalid request: ...`                              |
| Malformed JSON                                  | `invalid request: ...`                              |
| Division by zero                                | `division by zero is not allowed`                   |
| Square root of a negative number                | `square root of a negative number is not allowed`   |
| Result is NaN or ±Infinity                      | `invalid input: result is not a finite number`      |

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
{
  "error": "invalid request: Key: 'BinaryRequest.B' Error:Field validation for 'B' failed on the 'required' tag"
}
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

### Reverse Proxy Instead of Cross-Origin Requests in Docker

In the Docker setup, nginx serves the built frontend and forwards `/api/` to the backend container. The frontend is built with a relative `VITE_API_URL` (`/api/v1`), so the browser only ever talks to one origin.

**Why:**

- **No CORS in production-style runs** — the backend's CORS allow-list only needs to cover the Vite dev server (`http://localhost:5173`); the containerized app never makes a cross-origin request.
- **Smaller attack surface** — the backend port is not published to the host by default; only nginx is reachable.
- **Separate containers** — each process has its own lifecycle, logs, and restart policy, instead of running both in one container with a process supervisor.

**Trade-off accepted:** `VITE_API_URL` is resolved at build time, so pointing the frontend at a different API requires rebuilding the image.

---

## Test Coverage

Backend results are below. Frontend results are in [Frontend tests](#frontend-tests).

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

| Package               | Coverage   |
| --------------------- | ---------- |
| `internal/calculator` | **100.0%** |
| `internal/handlers`   | **86.1%**  |
| **Total**             | **91.2%**  |

> [!NOTE]
> The two branches below 100% are intentionally unreachable via the HTTP layer:
>
> - **`handleBinary` / `handleUnary` NaN/Inf guard (83.3%)** — standard JSON does not support `NaN` or `Infinity` literals, so this guard can never be triggered through Gin's JSON parser. It exists as a defensive layer if handlers are ever called directly.
> - **`mapDomainError` default branch (80.0%)** — a catch-all for future, currently non-existent error types. No domain function produces an unmapped error today.
