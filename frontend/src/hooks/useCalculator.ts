import { useCallback, useRef, useState } from "react";
import { calculateBinary, calculateUnary } from "../api/calculatorApi";
import type {
  BinaryOperation,
  Operation,
  UnaryOperation,
} from "../types/calculator";

export interface UseCalculatorReturn {
  currentInput: string;
  previousValue: number | null;
  activeOperation: Operation | null;
  isLoading: boolean;
  error: string | null;
  inputDigit: (digit: string) => void;
  inputDecimal: () => void;
  toggleSign: () => void;
  clear: () => void;
  backspace: () => void;
  selectBinaryOperation: (op: BinaryOperation) => Promise<void>;
  executeUnaryOperation: (op: UnaryOperation) => Promise<void>;
  evaluate: () => Promise<void>;
}

const INITIAL_INPUT = "0";

/** Parses the display string; returns null if it is not a valid finite number. */
const parseInput = (input: string): number | null => {
  const value = parseFloat(input);
  return Number.isFinite(value) ? value : null;
};

export function useCalculator(): UseCalculatorReturn {
  const [currentInput, setCurrentInput] = useState<string>(INITIAL_INPUT);
  const [previousValue, setPreviousValue] = useState<number | null>(null);
  const [activeOperation, setActiveOperation] = useState<Operation | null>(
    null,
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Internal UI-flow flags (not part of the public state contract).
  // overwrite: the next digit/decimal replaces the display instead of appending
  //            (set after an operator, a result, or a unary result).
  // awaitingOperand: an operator was just pressed and the second number
  //            has not been entered yet (prevents premature chaining).
  const [overwrite, setOverwrite] = useState<boolean>(false);
  const [awaitingOperand, setAwaitingOperand] = useState<boolean>(false);

  // Synchronous lock: state updates are async, so a rapid double click could
  // slip past an `isLoading` state check. The ref closes that gap.
  const loadingRef = useRef<boolean>(false);

  /**
   * Runs an API call with a loading lock and error capture.
   * Returns the result, or null if the call was blocked or failed.
   */
  const runRequest = useCallback(
    async (request: () => Promise<number>): Promise<number | null> => {
      if (loadingRef.current) return null; // duplicate submission guard
      loadingRef.current = true;
      setIsLoading(true);
      setError(null);
      try {
        return await request();
      } catch (e) {
        setError(e instanceof Error ? e.message : "Something went wrong");
        return null;
      } finally {
        loadingRef.current = false;
        setIsLoading(false);
      }
    },
    [],
  );

  const showResult = useCallback((result: number) => {
    setCurrentInput(String(result));
    setOverwrite(true);
  }, []);

  const inputDigit = useCallback(
    (digit: string) => {
      if (!/^[0-9]$/.test(digit)) return;
      setError(null);

      if (overwrite) {
        setCurrentInput(
          currentInput === "-0" && awaitingOperand ? `-${digit}` : digit,
        );
        setOverwrite(false);
        setAwaitingOperand(false);
        return;
      }

      setAwaitingOperand(false);
      setCurrentInput((prev) => {
        if (prev === "0") return digit; // '05' -> '5'
        if (prev === "-0") return `-${digit}`; // '-05' -> '-5'
        return prev + digit;
      });
    },
    [overwrite, awaitingOperand, currentInput],
  );

  const inputDecimal = useCallback(() => {
    setError(null);

    if (overwrite) {
      setCurrentInput("0.");
      setOverwrite(false);
      setAwaitingOperand(false);
      return;
    }

    setAwaitingOperand(false);
    setCurrentInput((prev) => (prev.includes(".") ? prev : `${prev}.`));
  }, [overwrite]);

  const toggleSign = useCallback(() => {
    setError(null);

    // Right after an operator, start the second operand as a negative number.
    if (overwrite && awaitingOperand) {
      setCurrentInput("-0");
      setOverwrite(false);
      setAwaitingOperand(false);
      return;
    }

    setCurrentInput((prev) => {
      if (prev === "0") return prev;
      return prev.startsWith("-") ? prev.slice(1) : `-${prev}`;
    });
  }, [overwrite, awaitingOperand]);

  const clear = useCallback(() => {
    setCurrentInput(INITIAL_INPUT);
    setPreviousValue(null);
    setActiveOperation(null);
    setError(null);
    setOverwrite(false);
    setAwaitingOperand(false);
  }, []);

  const backspace = useCallback(() => {
    setError(null);
    // Don't edit a computed result or a stored operand shown on screen.
    if (overwrite) return;

    setCurrentInput((prev) => {
      const next = prev.slice(0, -1);
      return next === "" || next === "-" ? INITIAL_INPUT : next;
    });
  }, [overwrite]);

  const evaluate = useCallback(async () => {
    if (previousValue === null || activeOperation === null) return;
    if (loadingRef.current) return;

    const operand = parseInput(currentInput);
    if (operand === null) {
      setError("Invalid number");
      return;
    }

    const result = await runRequest(() =>
      calculateBinary(
        activeOperation as BinaryOperation,
        previousValue,
        operand,
      ),
    );
    if (result === null) return; // on error, keep state so the user can correct the input

    showResult(result);
    setPreviousValue(null);
    setActiveOperation(null);
    setAwaitingOperand(false);
  }, [previousValue, activeOperation, currentInput, runRequest, showResult]);

  const selectBinaryOperation = useCallback(
    async (op: BinaryOperation) => {
      if (loadingRef.current) return;

      // Operator pressed twice in a row: just swap the pending operation.
      if (activeOperation !== null && awaitingOperand) {
        setActiveOperation(op);
        return;
      }

      const operand = parseInput(currentInput);
      if (operand === null) {
        setError("Invalid number");
        return;
      }

      // Operation chaining: 2 + 3 + ... -> evaluate 2 + 3 first.
      if (previousValue !== null && activeOperation !== null) {
        const result = await runRequest(() =>
          calculateBinary(
            activeOperation as BinaryOperation,
            previousValue,
            operand,
          ),
        );
        if (result === null) return;

        showResult(result);
        setPreviousValue(result);
      } else {
        setPreviousValue(operand);
        setOverwrite(true);
      }

      setActiveOperation(op);
      setAwaitingOperand(true);
    },
    [
      activeOperation,
      awaitingOperand,
      currentInput,
      previousValue,
      runRequest,
      showResult,
    ],
  );

  const executeUnaryOperation = useCallback(
    async (op: UnaryOperation) => {
      if (loadingRef.current) return;

      const operand = parseInput(currentInput);
      if (operand === null) {
        setError("Invalid number");
        return;
      }

      const result = await runRequest(() => calculateUnary(op, operand));
      if (result === null) return; // e.g. backend: "cannot take square root of a negative number"

      showResult(result);
      // The unary result now acts as the current operand of any pending operation.
      setAwaitingOperand(false);
    },
    [currentInput, runRequest, showResult],
  );

  return {
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
  };
}
