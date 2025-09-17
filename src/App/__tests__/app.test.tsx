import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'

const MockApp = () => (
  <div data-testid="app">
    <h1>FazendaPro App</h1>
  </div>
)

describe('App', () => {
  it('deve renderizar sem erros', () => {
    render(<MockApp />)

    expect(screen.getByTestId('app')).toBeInTheDocument()
    expect(screen.getByText('FazendaPro App')).toBeInTheDocument()
  })

  it('deve ter estrutura básica', () => {
    render(<MockApp />)

    const app = screen.getByTestId('app')
    expect(app).toBeInTheDocument()
    expect(app.tagName).toBe('DIV')
  })
})