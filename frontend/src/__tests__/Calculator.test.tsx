import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from '../App'
import { calculateBinary, calculateUnary } from '../api/calculatorApi'

vi.mock('../api/calculatorApi', () => ({
  calculateBinary: vi.fn(),
  calculateUnary: vi.fn(),
}))

function readout(): HTMLElement {
  const output = document.querySelector('output')
  if (!(output instanceof HTMLElement)) {
    throw new Error('Calculator readout was not rendered')
  }
  return output
}

describe('Calculator', () => {
  beforeEach(() => {
    vi.mocked(calculateBinary).mockReset()
    vi.mocked(calculateUnary).mockReset()
    vi.mocked(calculateBinary).mockResolvedValue(0)
    vi.mocked(calculateUnary).mockResolvedValue(0)
  })

  it('keeps a single decimal point while entering a multi-digit number', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(screen.getByRole('button', { name: '1' }))
    await user.click(screen.getByRole('button', { name: 'Decimal point' }))
    await user.click(screen.getByRole('button', { name: '2' }))
    await user.click(screen.getByRole('button', { name: 'Decimal point' }))
    await user.click(screen.getByRole('button', { name: '3' }))

    expect(readout()).toHaveTextContent('1.23')
  })

  it('keeps a single decimal point when the number is typed', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.keyboard('1.2.3')

    expect(readout()).toHaveTextContent('1.23')
  })

  it('evaluates a binary operation from button clicks', async () => {
    const user = userEvent.setup()
    vi.mocked(calculateBinary).mockResolvedValue(12)
    render(<App />)

    await user.click(screen.getByRole('button', { name: '8' }))
    await user.click(screen.getByRole('button', { name: 'Add' }))
    await user.click(screen.getByRole('button', { name: '4' }))
    await user.click(screen.getByRole('button', { name: 'Equals' }))

    await waitFor(() => {
      expect(readout()).toHaveTextContent('12')
    })
    expect(calculateBinary).toHaveBeenCalledWith('add', 8, 4)
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('evaluates a binary operation from the keyboard', async () => {
    const user = userEvent.setup()
    vi.mocked(calculateBinary).mockResolvedValue(12)
    render(<App />)

    await user.keyboard('8+4{Enter}')

    await waitFor(() => {
      expect(readout()).toHaveTextContent('12')
    })
    expect(calculateBinary).toHaveBeenCalledWith('add', 8, 4)
  })

  it('runs a unary operation immediately, without pressing equals', async () => {
    const user = userEvent.setup()
    vi.mocked(calculateUnary).mockResolvedValue(4)
    render(<App />)

    await user.click(screen.getByRole('button', { name: '1' }))
    await user.click(screen.getByRole('button', { name: '6' }))
    await user.click(screen.getByRole('button', { name: 'Square root' }))

    await waitFor(() => {
      expect(readout()).toHaveTextContent('4')
    })
    expect(calculateUnary).toHaveBeenCalledWith('sqrt', 16)
    expect(calculateBinary).not.toHaveBeenCalled()
  })

  it('shows a rejected API error and clears it when a new digit is typed', async () => {
    const user = userEvent.setup()
    vi.mocked(calculateBinary).mockRejectedValueOnce(
      new Error('cannot divide by zero'),
    )
    render(<App />)

    await user.click(screen.getByRole('button', { name: '8' }))
    await user.click(screen.getByRole('button', { name: 'Divide' }))
    await user.click(screen.getByRole('button', { name: '0' }))
    await user.click(screen.getByRole('button', { name: 'Equals' }))

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'cannot divide by zero',
    )
    expect(readout()).toHaveTextContent('0')

    await user.keyboard('5')

    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(readout()).toHaveTextContent('5')
  })
})
