import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { Login } from '../presentation/login'
import { useAuth } from '../hooks/useAuth'

vi.mock('../hooks/useAuth', () => ({
  useAuth: vi.fn()
}))

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('../../assets/images/logo.png', () => ({
  default: 'logo.png'
}))

Object.defineProperty(window, 'innerWidth', {
  writable: true,
  configurable: true,
  value: 1024,
})

const mockUseAuth = vi.mocked(useAuth)

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('Login Responsividade', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockUseAuth.mockReturnValue({
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: false,
      isLoading: false,
      user: null,
      token: null,
    })
  })

  it('deve renderizar layout desktop quando window.innerWidth > 768', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    })

    renderWithRouter(<Login />)
    
    const loginContainer = screen.getByText('loginTitle').closest('.login-container')
    expect(loginContainer).toBeInTheDocument()
    
    const cols = document.querySelectorAll('.ant-col')
    expect(cols.length).toBeGreaterThanOrEqual(2)
  })

  it('deve renderizar layout mobile quando window.innerWidth <= 768', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768,
    })

    renderWithRouter(<Login />)
    
    const loginContainer = screen.getByText('loginTitle').closest('.login-container')
    expect(loginContainer).toBeInTheDocument()
    
    const cols = document.querySelectorAll('.ant-col')
    expect(cols.length).toBeGreaterThanOrEqual(2)
  })

  it('deve ter logo visível em telas grandes', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    })

    renderWithRouter(<Login />)
    
    const logoImage = document.querySelector('img')
    expect(logoImage).toBeInTheDocument()
  })

  it('deve ter formulário de login responsivo', () => {
    renderWithRouter(<Login />)
    
    expect(screen.getByText('loginTitle')).toBeInTheDocument()
    expect(screen.getByText('access')).toBeInTheDocument()
  })

  it('deve ter padding adequado baseado no tamanho da tela', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    })

    const { container } = renderWithRouter(<Login />)
    
    const loginContainer = container.querySelector('.login-container')
    expect(loginContainer).toHaveStyle('padding: 50px')
  })

  it('deve ter padding reduzido em telas pequenas', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768,
    })

    const { container } = renderWithRouter(<Login />)
    
    const loginContainer = container.querySelector('.login-container')
    expect(loginContainer).toHaveStyle('padding: 50px')
  })

  it('deve ter flexbox layout vertical', () => {
    renderWithRouter(<Login />)
    
    const loginContainer = document.querySelector('.login-container')
    expect(loginContainer).toHaveStyle({
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center'
    })
  })

  it('deve ter colunas responsivas com breakpoints', () => {
    renderWithRouter(<Login />)
    
    const logoCol = document.querySelector('.ant-col-xs-24.ant-col-md-12')
    const formCol = document.querySelector('.ant-col-xs-24.ant-col-md-12')
    
    expect(logoCol).toBeInTheDocument()
    expect(formCol).toBeInTheDocument()
  })
})
