// ─── Operation types ──────────────────────────────────────────────────────────

export type BinaryOperation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'percentage';

export type UnaryOperation = 'sqrt';

export type Operation = BinaryOperation | UnaryOperation;

// ─── Request DTOs (mirror the Go backend structs) ─────────────────────────────

export interface BinaryRequest {
  a: number;
  b: number;
}

export interface UnaryRequest {
  a: number;
}

// ─── Response DTOs ────────────────────────────────────────────────────────────

export interface CalculationResponse {
  result: number;
}

export interface ErrorResponse {
  error: string;
}
