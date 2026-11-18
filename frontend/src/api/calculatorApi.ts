import type {
  BinaryOperation,
  BinaryRequest,
  CalculationResponse,
  ErrorResponse,
  UnaryOperation,
  UnaryRequest,
} from '../types/calculator';

const BASE_URL: string =
  import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1';

// ─── shared fetch helper ──────────────────────────────────────────────────────

/**
 * POSTs `body` to `${BASE_URL}/${path}`, resolves with the numeric result, or
 * throws an Error whose message is taken directly from the backend's
 * `ErrorResponse.error` field (or a generic network message when unreachable).
 */
async function post<T>(path: string, body: T): Promise<number> {
  let response: Response;

  try {
    response = await fetch(`${BASE_URL}/${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
  } catch {
    throw new Error('Network error: unable to reach the calculator API.');
  }

  if (!response.ok) {
    const errBody: ErrorResponse = await response.json();
    throw new Error(errBody.error ?? 'An unexpected error occurred.');
  }

  const data: CalculationResponse = await response.json();
  return data.result;
}

// ─── Public API ───────────────────────────────────────────────────────────────

/**
 * Sends a binary operation request (a op b) to the backend.
 *
 * @param operation - One of the binary operation slugs (add, subtract, …)
 * @param a         - First operand
 * @param b         - Second operand
 * @returns         The numeric result
 * @throws          Error with the backend's error message on failure
 */
export async function calculateBinary(
  operation: BinaryOperation,
  a: number,
  b: number
): Promise<number> {
  const body: BinaryRequest = { a, b };
  return post<BinaryRequest>(operation, body);
}

/**
 * Sends a unary operation request to the backend.
 *
 * @param operation - One of the unary operation slugs (sqrt)
 * @param a         - The operand
 * @returns         The numeric result
 * @throws          Error with the backend's error message on failure
 */
export async function calculateUnary(
  operation: UnaryOperation,
  a: number
): Promise<number> {
  const body: UnaryRequest = { a };
  return post<UnaryRequest>(operation, body);
}
