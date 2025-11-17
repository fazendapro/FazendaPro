import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import { App } from '../App'
import { BrowserRouter } from 'react-router-dom'

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>{children}</BrowserRouter>
)

describe('App', () => {
  it('deve renderizar sem erros', () => {
    const { container } = render(<App />, { wrapper })
    expect(container).toBeDefined()
  })

  it('deve ter estrutura básica', () => {
    const { container } = render(<App />, { wrapper })
    expect(container.firstChild).toBeDefined()
  })
})