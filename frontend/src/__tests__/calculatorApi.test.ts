import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { calculateBinary, calculateUnary } from '../api/calculatorApi'
import type { BinaryOperation } from '../types/calculator'

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('calculatorApi', () => {
  const fetchMock = vi.fn<typeof fetch>()

  beforeEach(() => {
    fetchMock.mockReset()
    vi.stubGlobal('fetch', fetchMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it.each<[BinaryOperation, number, number, number]>([
    ['add', 3, 7, 10],
    ['subtract', 0, 5, -5],
    ['multiply', -3, -4, 12],
    ['divide', 10, 4, 2.5],
    ['power', 2, 10, 1024],
    ['percentage', 15, 200, 30],
  ])(
    'calculateBinary POSTs /%s with both operands',
    async (operation, a, b, result) => {
      fetchMock.mockResolvedValueOnce(jsonResponse({ result }))

      await expect(calculateBinary(operation, a, b)).resolves.toBe(result)

      expect(fetchMock).toHaveBeenCalledWith(`${BASE_URL}/${operation}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ a, b }),
      })
    },
  )

  it('calculateUnary POSTs /sqrt with the single operand', async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse({ result: 4 }))

    await expect(calculateUnary('sqrt', 16)).resolves.toBe(4)

    expect(fetchMock).toHaveBeenCalledWith(`${BASE_URL}/sqrt`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a: 16 }),
    })
  })

  it('throws an Error containing the backend error message on 400', async () => {
    fetchMock.mockResolvedValueOnce(
      jsonResponse(
        { error: 'division by zero is not allowed' },
        400,
      ),
    )

    await expect(calculateBinary('divide', 5, 0)).rejects.toThrow(
      new Error('division by zero is not allowed'),
    )
  })
})
