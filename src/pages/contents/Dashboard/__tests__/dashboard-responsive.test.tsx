import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { Dashboard } from '../presentation/dashboard'
import { useIsMobile } from '../../../../hooks/use-is-mobile'

vi.mock('../../../../hooks/use-is-mobile', () => ({
  useIsMobile: vi.fn()
}))

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('react-chartjs-2', () => ({
  Bar: ({ data }: { data: unknown }) => <div data-testid="bar-chart">{JSON.stringify(data)}</div>,
  Line: ({ data }: { data: unknown }) => <div data-testid="line-chart">{JSON.stringify(data)}</div>,
}))

vi.mock('../components/overview', () => ({
  Overview: () => <div data-testid="overview">Overview</div>
}))

vi.mock('../components/cattle-quantity', () => ({
  CattleQuantity: () => <div data-testid="cattle-quantity">Cattle Quantity</div>
}))

vi.mock('../components/shopping-overview', () => ({
  ShoppingOverview: () => <div data-testid="shopping-overview">Shopping Overview</div>
}))

vi.mock('../components/rations', () => ({
  Rations: () => <div data-testid="rations">Rations</div>
}))

vi.mock('../components/dashboard-milk-production', () => ({
  DashboardMilkProduction: () => <div data-testid="milk-production">Milk Production</div>
}))

vi.mock('../components/next-to-calve', () => ({
  NextToCalve: () => <div data-testid="next-to-calve">Next To Calve</div>
}))

const mockUseIsMobile = vi.mocked(useIsMobile)

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('Dashboard Responsividade', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar layout desktop quando não for mobile', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Dashboard />)
    
    expect(screen.getByTestId('overview')).toBeInTheDocument()
    expect(screen.getByTestId('cattle-quantity')).toBeInTheDocument()
    expect(screen.getByTestId('shopping-overview')).toBeInTheDocument()
    expect(screen.getByTestId('rations')).toBeInTheDocument()
    expect(screen.getByTestId('next-to-calve')).toBeInTheDocument()
    
    expect(screen.getByTestId('bar-chart')).toBeInTheDocument()
    expect(screen.getByTestId('line-chart')).toBeInTheDocument()
  })

  it('deve renderizar layout mobile quando for mobile', () => {
    mockUseIsMobile.mockReturnValue(true)

    renderWithRouter(<Dashboard />)
    
    expect(screen.getByTestId('overview')).toBeInTheDocument()
    expect(screen.getByTestId('cattle-quantity')).toBeInTheDocument()
    expect(screen.getByTestId('shopping-overview')).toBeInTheDocument()
    expect(screen.getByTestId('rations')).toBeInTheDocument()
    expect(screen.getByTestId('next-to-calve')).toBeInTheDocument()
    
    expect(screen.getByTestId('bar-chart')).toBeInTheDocument()
    expect(screen.getByTestId('line-chart')).toBeInTheDocument()
  })

  it('deve ter estrutura de grid responsiva', () => {
    mockUseIsMobile.mockReturnValue(false)

    const { container } = renderWithRouter(<Dashboard />)
    
    const rows = container.querySelectorAll('.ant-row')
    expect(rows.length).toBeGreaterThan(0)
    
    const cols = container.querySelectorAll('.ant-col')
    expect(cols.length).toBeGreaterThan(0)
  })

  it('deve ter cards com informações do dashboard', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Dashboard />)
    
    expect(screen.getByText('dashboard.salesAndPurchases')).toBeInTheDocument()
    expect(screen.getByText('dashboard.semenation')).toBeInTheDocument()
  })

  it('deve ter gutter adequado entre colunas', () => {
    mockUseIsMobile.mockReturnValue(false)

    const { container } = renderWithRouter(<Dashboard />)
    
    const rows = container.querySelectorAll('.ant-row')
    rows.forEach(row => {
      expect(row).toHaveStyle('margin-left: -8px; margin-right: -8px;')
    })
  })
})
