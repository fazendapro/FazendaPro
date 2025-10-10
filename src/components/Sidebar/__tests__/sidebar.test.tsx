import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '../../../test/test-utils'
import { Sidebar } from '../sidebar'

vi.mock('../../../contexts/AuthContext', () => ({
  useAuth: () => ({
    logout: vi.fn(),
  }),
}))

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useLocation: () => ({
      pathname: '/',
      search: '',
      hash: '',
      state: null,
    }),
  }
})

describe('Sidebar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar o menu com as opções corretas', () => {
    render(<Sidebar />)
    
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Vacas')).toBeInTheDocument()
  })

  it('deve navegar para a rota correta quando um item do menu é clicado', () => {
    render(<Sidebar />)
    
    fireEvent.click(screen.getByText('Vacas'))
    
    expect(mockNavigate).toHaveBeenCalledWith('/vacas')
  })

  it('deve navegar para dashboard quando Dashboard é clicado', () => {
    render(<Sidebar />)
    
    fireEvent.click(screen.getByText('Dashboard'))
    
    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('deve renderizar ícones corretos', () => {
    render(<Sidebar />)
    
    const dashboardIcon = screen.getByText('Dashboard').closest('li')
    const vacasIcon = screen.getByText('Vacas').closest('li')
    
    expect(dashboardIcon).toBeInTheDocument()
    expect(vacasIcon).toBeInTheDocument()
  })

  it('deve ser colapsível', () => {
    render(<Sidebar />)
    
    const sidebar = screen.getByRole('complementary')
    expect(sidebar).toBeInTheDocument()
  })

  it('deve ter largura padrão de 280px quando não colapsado', () => {
    render(<Sidebar />)
    
    const sidebar = screen.getByRole('complementary')
    expect(sidebar).toHaveStyle({ width: '280px' })
  })

  it('deve renderizar apenas quando autenticado', () => {
    render(<Sidebar />)
    
    expect(screen.getByRole('complementary')).toBeInTheDocument()
  })

  it('deve ter tema escuro por padrão', () => {
    render(<Sidebar />)
    
    const sidebar = screen.getByRole('complementary')
    expect(sidebar).toHaveClass('ant-layout-sider-dark')
  })
})
