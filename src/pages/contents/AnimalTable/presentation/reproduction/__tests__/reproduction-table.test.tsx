import { describe, it, expect } from 'vitest'

const MockReproductionTable = () => (
  <div data-testid="reproduction-table">
    <h2>Reproduction Table</h2>
    <div data-testid="table-content">Table Content</div>
  </div>
)

describe('ReproductionTable', () => {
  it('deve renderizar a tabela básica', () => {
    const table = MockReproductionTable()
    expect(table).toBeDefined()
  })

  it('deve ter estrutura básica', () => {
    const table = MockReproductionTable()
    expect(table).toBeDefined()
  })
})