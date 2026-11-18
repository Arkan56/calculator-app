import { useEffect } from "react";
import { Display } from "./components/Display";
import { Keypad } from "./components/Keypad";
import { useCalculator } from "./hooks/useCalculator";
import type { BinaryOperation } from "./types/calculator";

const BINARY_KEYS: Record<string, BinaryOperation> = {
  "+": "add",
  "-": "subtract",
  "*": "multiply",
  "/": "divide",
  "^": "power",
  "%": "percentage",
};

function App() {
  const {
    currentInput,
    previousValue,
    activeOperation,
    isLoading,
    error,
    inputDigit,
    inputDecimal,
    toggleSign,
    clear,
    backspace,
    selectBinaryOperation,
    executeUnaryOperation,
    evaluate,
  } = useCalculator();

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey || e.altKey) return; // keep browser shortcuts working

      const target = e.target as HTMLElement | null;
      if (
        target &&
        (target.tagName === "INPUT" ||
          target.tagName === "TEXTAREA" ||
          target.tagName === "SELECT" ||
          target.isContentEditable)
      ) {
        return;
      }

      const { key } = e;
      const isHandled =
        /^[0-9]$/.test(key) ||
        key === "." ||
        key === "," ||
        key in BINARY_KEYS ||
        key === "Enter" ||
        key === "=" ||
        key === "Backspace" ||
        key === "Escape";

      if (!isHandled) return;

      // Stops "/" opening quick-find, and Enter re-clicking a focused button.
      e.preventDefault();
      if (isLoading || e.repeat) return;

      if (/^[0-9]$/.test(key)) {
        inputDigit(key);
      } else if (key === "." || key === ",") {
        inputDecimal();
      } else if (key in BINARY_KEYS) {
        void selectBinaryOperation(BINARY_KEYS[key]);
      } else if (key === "Enter" || key === "=") {
        void evaluate();
      } else if (key === "Backspace") {
        backspace();
      } else if (key === "Escape") {
        clear();
      }
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [
    isLoading,
    inputDigit,
    inputDecimal,
    selectBinaryOperation,
    evaluate,
    backspace,
    clear,
  ]);

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-200 p-4">
      <h1 className="sr-only">Calculator</h1>
      <div className="w-full max-w-sm space-y-4 rounded-3xl bg-slate-100 p-4 shadow-xl ring-1 ring-slate-300">
        <Display
          currentInput={currentInput}
          previousValue={previousValue}
          activeOperation={activeOperation}
          isLoading={isLoading}
          error={error}
        />
        <Keypad
          activeOperation={activeOperation}
          isLoading={isLoading}
          onDigit={inputDigit}
          onDecimal={inputDecimal}
          onToggleSign={toggleSign}
          onClear={clear}
          onBackspace={backspace}
          onBinaryOperation={selectBinaryOperation}
          onUnaryOperation={executeUnaryOperation}
          onEvaluate={evaluate}
        />
      </div>
    </main>
  );
}

export default App;
