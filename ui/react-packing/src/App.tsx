import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import './App.css'

type PackSizesResponse = {
  pack_sizes: number[]
}

type PackBreakdown = {
  size: number
  count: number
}

type ApiErrorResponse = {
  error?: {
    code?: string
    message?: string
  }
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 10_000,
})

function parsePackSizes(values: string[]): { packSizes: number[]; error: string | null } {
  if (values.length === 0) {
    return { packSizes: [], error: 'Add at least one pack size before saving.' }
  }

  const seen = new Set<number>()
  const packSizes: number[] = []

  for (let i = 0; i < values.length; i += 1) {
    const raw = values[i].trim()
    if (raw.length === 0) {
      return { packSizes: [], error: `Pack size ${i + 1} is empty.` }
    }
    if (!/^\d+$/.test(raw)) {
      return { packSizes: [], error: `Pack size ${i + 1} must be a positive integer.` }
    }

    const size = Number(raw)
    if (!Number.isSafeInteger(size) || size <= 0) {
      return { packSizes: [], error: `Pack size ${i + 1} must be greater than zero.` }
    }
    if (seen.has(size)) {
      return { packSizes: [], error: `Pack size ${size} is duplicated.` }
    }

    seen.add(size)
    packSizes.push(size)
  }

  return { packSizes, error: null }
}

function parseAmount(value: string): { amount: number; error: string | null } {
  const trimmed = value.trim()
  if (trimmed.length === 0) {
    return { amount: 0, error: 'Amount is required.' }
  }
  if (!/^\d+$/.test(trimmed)) {
    return { amount: 0, error: 'Amount must be a positive integer.' }
  }

  const amount = Number(trimmed)
  if (!Number.isSafeInteger(amount) || amount <= 0) {
    return { amount: 0, error: 'Amount must be greater than zero.' }
  }

  return { amount, error: null }
}

function toApiErrorMessage(error: unknown, fallback: string): string {
  if (!axios.isAxiosError(error)) {
    return fallback
  }

  const responseData = error.response?.data as ApiErrorResponse | undefined
  if (responseData?.error?.message) {
    return responseData.error.message
  }

  return error.message || fallback
}

function App() {
  const [packSizeInputs, setPackSizeInputs] = useState<string[]>([])
  const [isLoadingPackSizes, setIsLoadingPackSizes] = useState(true)
  const [isSavingPackSizes, setIsSavingPackSizes] = useState(false)
  const [packError, setPackError] = useState('')
  const [packStatus, setPackStatus] = useState('')

  const [amountInput, setAmountInput] = useState('')
  const [breakdown, setBreakdown] = useState<PackBreakdown[]>([])
  const [isCalculating, setIsCalculating] = useState(false)
  const [calculationError, setCalculationError] = useState('')
  const [calculationStatus, setCalculationStatus] = useState('')

  const totalPacks = useMemo(
    () => breakdown.reduce((sum, item) => sum + item.count, 0),
    [breakdown],
  )
  const totalShipped = useMemo(
    () => breakdown.reduce((sum, item) => sum + item.size * item.count, 0),
    [breakdown],
  )

  useEffect(() => {
    const loadPackSizes = async () => {
      try {
        const response = await api.get<PackSizesResponse>('/pack-sizes')
        setPackSizeInputs(response.data.pack_sizes.map((size) => String(size)))
        setPackStatus(
          response.data.pack_sizes.length > 0
            ? `Loaded ${response.data.pack_sizes.length} pack sizes.`
            : 'No pack sizes configured yet.',
        )
      } catch (error) {
        setPackError(toApiErrorMessage(error, 'Failed to load pack sizes.'))
      } finally {
        setIsLoadingPackSizes(false)
      }
    }

    void loadPackSizes()
  }, [])

  const handlePackSizeChange = (index: number, value: string) => {
    setPackSizeInputs((current) =>
      current.map((item, itemIndex) => (itemIndex === index ? value : item)),
    )
  }

  const handleAddPackSize = () => {
    setPackError('')
    setPackStatus('')
    setPackSizeInputs((current) => [...current, ''])
  }

  const handleDeletePackSize = (index: number) => {
    setPackError('')
    setPackStatus('')
    setPackSizeInputs((current) => current.filter((_, itemIndex) => itemIndex !== index))
  }

  const handleSavePackSizes = async () => {
    setPackError('')
    setPackStatus('')

    const parsed = parsePackSizes(packSizeInputs)
    if (parsed.error) {
      setPackError(parsed.error)
      return
    }

    setIsSavingPackSizes(true)
    try {
      const response = await api.put<PackSizesResponse>('/pack-sizes', {
        pack_sizes: parsed.packSizes,
      })
      setPackSizeInputs(response.data.pack_sizes.map((size) => String(size)))
      setPackStatus(`Saved ${response.data.pack_sizes.length} pack sizes.`)
    } catch (error) {
      setPackError(toApiErrorMessage(error, 'Failed to save pack sizes.'))
    } finally {
      setIsSavingPackSizes(false)
    }
  }

  const handleCalculate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setCalculationError('')
    setCalculationStatus('')

    const parsed = parseAmount(amountInput)
    if (parsed.error) {
      setCalculationError(parsed.error)
      return
    }

    setIsCalculating(true)
    try {
      const response = await api.post<PackBreakdown[]>('/calculate', {
        amount: parsed.amount,
      })
      setBreakdown(response.data)
      setCalculationStatus(`Calculated breakdown for amount ${parsed.amount}.`)
    } catch (error) {
      setBreakdown([])
      setCalculationError(toApiErrorMessage(error, 'Failed to calculate breakdown.'))
    } finally {
      setIsCalculating(false)
    }
  }

  return (
    <main className="app">
      <div className="layout">
        <section className="panel">
          <h1>Pack Configs</h1>
          <p className="panel-subtitle">Manage sizes used by your Go API optimizer.</p>

          {isLoadingPackSizes ? (
            <p className="message">Loading pack sizes...</p>
          ) : (
            <div className="stack">
              <div className="stack">
                {packSizeInputs.length === 0 ? (
                  <p className="empty-state">No pack sizes yet. Add one to start.</p>
                ) : null}

                {packSizeInputs.map((value, index) => (
                  <div className="input-row" key={`pack-size-${index + 1}`}>
                    <input
                      aria-label={`Pack size ${index + 1}`}
                      className="input"
                      inputMode="numeric"
                      placeholder="Pack size"
                      type="number"
                      min="1"
                      step="1"
                      value={value}
                      onChange={(event) => handlePackSizeChange(index, event.target.value)}
                    />
                    <button
                      className="btn btn-danger"
                      type="button"
                      onClick={() => handleDeletePackSize(index)}
                    >
                      Delete
                    </button>
                  </div>
                ))}
              </div>

              <div className="action-row">
                <button className="btn btn-secondary" type="button" onClick={handleAddPackSize}>
                  Add size
                </button>
                <button
                  className="btn btn-primary"
                  disabled={isSavingPackSizes}
                  type="button"
                  onClick={handleSavePackSizes}
                >
                  {isSavingPackSizes ? 'Saving...' : 'Save'}
                </button>
              </div>
            </div>
          )}

          {packError ? <p className="message message-error">{packError}</p> : null}
          {packStatus ? <p className="message message-success">{packStatus}</p> : null}
        </section>

        <section className="panel">
          <h2>Calculation</h2>
          <p className="panel-subtitle">Enter an amount and fetch the optimal breakdown.</p>

          <form className="calc-form" onSubmit={handleCalculate}>
            <input
              aria-label="Order amount"
              className="input"
              inputMode="numeric"
              min="1"
              name="amount"
              placeholder="Amount"
              step="1"
              type="number"
              value={amountInput}
              onChange={(event) => setAmountInput(event.target.value)}
            />
            <button className="btn btn-primary" disabled={isCalculating} type="submit">
              {isCalculating ? 'Calculating...' : 'Calculate'}
            </button>
          </form>

          {calculationError ? <p className="message message-error">{calculationError}</p> : null}
          {calculationStatus ? (
            <p className="message message-success">{calculationStatus}</p>
          ) : null}

          {breakdown.length === 0 ? (
            <p className="empty-state">No breakdown yet.</p>
          ) : (
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Pack Size</th>
                    <th>Count</th>
                    <th>Total Units</th>
                  </tr>
                </thead>
                <tbody>
                  {breakdown.map((item) => (
                    <tr key={`pack-${item.size}`}>
                      <td>{item.size}</td>
                      <td>{item.count}</td>
                      <td>{item.size * item.count}</td>
                    </tr>
                  ))}
                </tbody>
                <tfoot>
                  <tr>
                    <td>Total</td>
                    <td>{totalPacks}</td>
                    <td>{totalShipped}</td>
                  </tr>
                </tfoot>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  )
}

export default App
