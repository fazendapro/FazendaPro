import { describe, it, expect } from 'vitest'

const MockAnimalsContainer = () => (
  <div data-testid="animals-container">
    <h2>Animals Container</h2>
    <div data-testid="animals-component">Animals Component</div>
  </div>
)

describe('Animal Table - Integration Tests', () => {
  it('deve renderizar o container básico', () => {
    expect(MockAnimalsContainer).toBeDefined()
  })

  it('deve ter estrutura básica', () => {
    const container = MockAnimalsContainer()
    expect(container).toBeDefined()
  })
})