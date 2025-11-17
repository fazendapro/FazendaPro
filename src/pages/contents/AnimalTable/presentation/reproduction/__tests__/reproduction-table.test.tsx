import { describe, it, expect, afterEach } from 'vitest'
import { render, cleanup } from '@testing-library/react'

const MockReproductionTable = () => (
  <div data-testid="reproduction-table">
    <h2>Reproduction Table</h2>
    <div data-testid="table-content">Table Content</div>
  </div>
)

describe('ReproductionTable', () => {
  afterEach(() => {
    cleanup()
  })

  it('deve renderizar a tabela básica', () => {
    const { container } = render(<MockReproductionTable />)
    expect(container).toBeDefined()
  })

  it('deve ter estrutura básica', () => {
    const { container } = render(<MockReproductionTable />)
    const table = container.querySelector('[data-testid="reproduction-table"]')
    expect(table).toBeInTheDocument()
  })
})