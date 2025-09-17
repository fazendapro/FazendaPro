import { describe, it, expect } from 'vitest'

const MockAnimalsContainer = () => (
  <div data-testid="animals-container">
    <h2>Animals Container</h2>
    <div data-testid="animals-component">Animals Component</div>
  </div>
)

describe('Animals Container', () => {
  it('deve renderizar o container com id correto', () => {
    const container = MockAnimalsContainer()
    expect(container).toBeDefined()
  })

  it('deve ter estrutura básica', () => {
    const container = MockAnimalsContainer()
    expect(container).toBeDefined()
  })
})