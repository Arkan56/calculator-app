import type {
  BinaryOperation,
  Operation,
  UnaryOperation,
} from "../types/calculator";

interface KeypadProps {
  activeOperation: Operation | null;
  isLoading: boolean;
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onToggleSign: () => void;
  onClear: () => void;
  onBackspace: () => void;
  onBinaryOperation: (op: BinaryOperation) => void | Promise<void>;
  onUnaryOperation: (op: UnaryOperation) => void | Promise<void>;
  onEvaluate: () => void | Promise<void>;
}

type Variant = "digit" | "utility" | "operator" | "equals";

interface KeyProps {
  label: string;
  ariaLabel?: string;
  variant: Variant;
  onClick: () => void;
  disabled: boolean;
  selected?: boolean;
  className?: string;
}

const BASE =
  "flex h-14 items-center justify-center rounded-lg text-xl font-semibold select-none sm:h-16 sm:text-2xl " +
  "shadow-[0_2px_0_0_rgb(0_0_0/0.18)] transition-colors " +
  "enabled:active:translate-y-px enabled:active:shadow-none " +
  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-600 focus-visible:ring-offset-2 " +
  "disabled:cursor-not-allowed disabled:opacity-50";

const VARIANTS: Record<Variant, string> = {
  digit:
    "bg-white text-slate-900 ring-1 ring-slate-300 enabled:hover:bg-slate-50",
  utility: "bg-slate-300 text-slate-900 enabled:hover:bg-slate-200",
  operator: "bg-indigo-600 text-white enabled:hover:bg-indigo-500",
  equals: "bg-amber-500 text-slate-900 enabled:hover:bg-amber-400",
};

// Selected operator is inverted so it reads as "pending" without relying on hue alone.
const SELECTED =
  "bg-white text-indigo-700 ring-2 ring-indigo-600 enabled:hover:bg-white";

function Key({
  label,
  ariaLabel,
  variant,
  onClick,
  disabled,
  selected = false,
  className = "",
}: KeyProps) {
  return (
    <button
      type="button"
      aria-label={ariaLabel}
      aria-pressed={variant === "operator" ? selected : undefined}
      disabled={disabled}
      onClick={onClick}
      className={`${BASE} ${selected ? SELECTED : VARIANTS[variant]} ${className}`}
    >
      {label}
    </button>
  );
}

export function Keypad({
  activeOperation,
  isLoading,
  onDigit,
  onDecimal,
  onToggleSign,
  onClear,
  onBackspace,
  onBinaryOperation,
  onUnaryOperation,
  onEvaluate,
}: KeypadProps) {
  const digit = (d: string, className?: string) => (
    <Key
      key={d}
      label={d}
      variant="digit"
      disabled={isLoading}
      onClick={() => onDigit(d)}
      className={className}
    />
  );

  const binary = (op: BinaryOperation, label: string, ariaLabel: string) => (
    <Key
      label={label}
      ariaLabel={ariaLabel}
      variant="operator"
      disabled={isLoading}
      selected={activeOperation === op}
      onClick={() => void onBinaryOperation(op)}
    />
  );

  return (
    <div className="grid grid-cols-4 gap-2">
      {/* Row 1 */}
      <Key
        label="AC"
        ariaLabel="Clear all"
        variant="utility"
        disabled={isLoading}
        onClick={onClear}
      />
      <Key
        label="DEL"
        ariaLabel="Delete last digit"
        variant="utility"
        disabled={isLoading}
        onClick={onBackspace}
      />
      <Key
        label="+/-"
        ariaLabel="Toggle sign"
        variant="utility"
        disabled={isLoading}
        onClick={onToggleSign}
      />
      {binary("divide", "÷", "Divide")}

      {/* Row 2 */}
      {binary("power", "xʸ", "Power")}
      <Key
        label="√"
        ariaLabel="Square root"
        variant="operator"
        disabled={isLoading}
        onClick={() => void onUnaryOperation("sqrt")}
      />
      {binary("percentage", "%", "Percentage")}
      {binary("multiply", "×", "Multiply")}

      {/* Row 3 */}
      {digit("7")}
      {digit("8")}
      {digit("9")}
      {binary("subtract", "−", "Subtract")}

      {/* Row 4 */}
      {digit("4")}
      {digit("5")}
      {digit("6")}
      {binary("add", "+", "Add")}

      {/* Rows 5-6: "=" spans two rows on the right */}
      {digit("1")}
      {digit("2")}
      {digit("3")}
      <Key
        label="="
        ariaLabel="Equals"
        variant="equals"
        disabled={isLoading}
        onClick={() => void onEvaluate()}
        className="row-span-2 h-auto"
      />

      {digit("0", "col-span-2")}
      <Key
        label="."
        ariaLabel="Decimal point"
        variant="digit"
        disabled={isLoading}
        onClick={onDecimal}
      />
    </div>
  );
}
