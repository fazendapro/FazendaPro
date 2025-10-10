import { render, screen, fireEvent } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuth } from '../../../contexts/AuthContext'
import { useLocation, useNavigate } from 'react-router-dom'
import { MobileSidebar } from '../mobile-sidebar'

vi.mock('../../../contexts/AuthContext', () => ({
  useAuth: vi.fn()
}))

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useLocation: vi.fn(),
    useNavigate: vi.fn(),
  }
})

const mockUseAuth = vi.mocked(useAuth)
const mockUseLocation = vi.mocked(useLocation)
const mockUseNavigate = vi.mocked(useNavigate)

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('MobileSidebar', () => {
  const mockLogout = vi.fn()
  const mockNavigate = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    
    mockUseLocation.mockReturnValue({
      pathname: '/',
      search: '',
      hash: '',
      state: null,
      key: 'default',
    })
    
    mockUseNavigate.mockReturnValue(mockNavigate)
    
    mockUseAuth.mockReturnValue({
      logout: mockLogout,
      isAuthenticated: true,
      isLoading: false,
      user: null,
      token: null,
      login: vi.fn(),
    })
  })

  it('deve renderizar a sidebar mobile quando autenticado', () => {
    renderWithRouter(<MobileSidebar />)
    
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Vacas')).toBeInTheDocument()
    expect(screen.getByText('Sair')).toBeInTheDocument()
  })

  it('não deve renderizar quando não autenticado', () => {
    mockUseAuth.mockReturnValue({
      logout: mockLogout,
      isAuthenticated: false,
      isLoading: false,
      user: null,
      token: null,
      login: vi.fn(),
    })

    const { container } = renderWithRouter(<MobileSidebar />)
    expect(container.firstChild).toBeNull()
  })

  it('deve destacar o item ativo baseado na rota atual', () => {
    mockUseLocation.mockReturnValue({
      pathname: '/vacas',
      search: '',
      hash: '',
      state: null,
      key: 'default',
    })

    renderWithRouter(<MobileSidebar />)
    
    const vacasButton = screen.getByText('Vacas').closest('button')
    expect(vacasButton).toHaveClass('ant-btn-primary')
  })

  it('deve navegar para a rota correta ao clicar em um item', () => {
    renderWithRouter(<MobileSidebar />)
    
    fireEvent.click(screen.getByText('Dashboard'))
    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('deve fazer logout ao clicar em Sair', () => {
    renderWithRouter(<MobileSidebar />)
    
    fireEvent.click(screen.getByText('Sair'))
    expect(mockLogout).toHaveBeenCalled()
  })

  it('deve ter posicionamento fixo na parte inferior', () => {
    renderWithRouter(<MobileSidebar />)
    
    const sidebar = screen.getByRole('navigation')
    expect(sidebar).toHaveStyle({
      position: 'fixed',
      bottom: '0',
      left: '0',
      right: '0',
      zIndex: '1000'
    })
  })

  it('deve ter ícones para cada item do menu', () => {
    renderWithRouter(<MobileSidebar />)

    const icons = screen.getAllByRole('img', { hidden: true })
    expect(icons.length).toBeGreaterThan(0)
  })

  it('deve ter layout responsivo com flexbox', () => {
    renderWithRouter(<MobileSidebar />)
    
    const sidebar = screen.getByRole('navigation')
    expect(sidebar).toHaveStyle({
      display: 'flex',
      justifyContent: 'space-around',
      alignItems: 'center'
    })
  })

  it('deve ter altura fixa adequada para mobile', () => {
    renderWithRouter(<MobileSidebar />)
    
    const sidebar = screen.getByRole('navigation')
    expect(sidebar).toHaveStyle({
      height: '60px'
    })
  })

  it('deve ter background branco com sombra', () => {
    renderWithRouter(<MobileSidebar />)
    
    const sidebar = screen.getByRole('navigation')
    expect(sidebar).toHaveStyle({
      background: 'white',
      boxShadow: '0 -2px 8px rgba(0, 0, 0, 0.1)'
    })
  })
})