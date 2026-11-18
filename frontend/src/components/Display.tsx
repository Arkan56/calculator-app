import type { Operation } from "../types/calculator";

interface DisplayProps {
  currentInput: string;
  previousValue: number | null;
  activeOperation: Operation | null;
  isLoading: boolean;
  error: string | null;
}

const OPERATION_SYMBOLS: Record<Operation, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
  power: "^",
  percentage: "%",
  sqrt: "√",
};

/**
 * Step the font size down as the number gets longer so it never overflows,
 * with a smaller scale on mobile than on >= sm screens. `break-all` on the
 * readout is the last-resort fallback for extremely long values.
 */
function readoutSize(length: number): string {
  if (length <= 8) return "text-4xl sm:text-5xl";
  if (length <= 11) return "text-3xl sm:text-4xl";
  if (length <= 15) return "text-2xl sm:text-3xl";
  return "text-lg sm:text-2xl";
}

export function Display({
  currentInput,
  previousValue,
  activeOperation,
  isLoading,
  error,
}: DisplayProps) {
  const pendingLabel =
    previousValue !== null && activeOperation !== null
      ? `${previousValue} ${OPERATION_SYMBOLS[activeOperation]}`
      : null;

  return (
    <div
      className="flex min-h-40 flex-col justify-between gap-2 rounded-xl bg-[#c8d6c0] p-4 text-[#1c261c] shadow-inner ring-1 ring-black/10"
      aria-label="Calculator display"
    >
      {/* Secondary line: pending operation + loading indicator */}
      <div className="flex h-6 items-center justify-between gap-2">
        <span className="min-w-0 flex-1 truncate text-right font-mono text-sm text-[#1c261c]/70 sm:text-base">
          {pendingLabel}
        </span>
        {isLoading && (
          <span
            role="status"
            className="flex shrink-0 items-center text-[#1c261c]/60"
          >
            <svg
              className="h-4 w-4 motion-safe:animate-spin"
              viewBox="0 0 24 24"
              fill="none"
              aria-hidden="true"
            >
              <circle
                cx="12"
                cy="12"
                r="9"
                stroke="currentColor"
                strokeWidth="3"
                opacity="0.25"
              />
              <path
                d="M21 12a9 9 0 0 0-9-9"
                stroke="currentColor"
                strokeWidth="3"
                strokeLinecap="round"
              />
            </svg>
            <span className="sr-only">Calculating</span>
          </span>
        )}
      </div>

      {/* Error banner */}
      {error && (
        <p
          role="alert"
          className="rounded-md bg-red-700 px-3 py-2 text-sm font-semibold text-white"
        >
          {error}
        </p>
      )}

      {/* Main readout (<output> is announced politely by screen readers) */}
      <output
        className={`block break-all text-right font-mono font-semibold tabular-nums leading-tight ${readoutSize(
          currentInput.length,
        )} ${isLoading ? "opacity-60" : ""}`}
      >
        {currentInput}
      </output>
    </div>
  );
}
