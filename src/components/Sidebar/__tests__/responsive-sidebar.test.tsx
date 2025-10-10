import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ResponsiveSidebar } from '../responsive-sidebar'
import { Grid } from 'antd'

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd')
  return {
    ...actual,
    Grid: {
      useBreakpoint: vi.fn()
    }
  }
})

vi.mock('../sidebar', () => ({
  Sidebar: () => <div data-testid="desktop-sidebar">Desktop Sidebar</div>
}))

vi.mock('../mobile-sidebar', () => ({
  MobileSidebar: () => <div data-testid="mobile-sidebar">Mobile Sidebar</div>
}))

const mockUseBreakpoint = vi.mocked(Grid.useBreakpoint)

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('ResponsiveSidebar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar MobileSidebar quando a tela for extra pequena (xs)', () => {
    mockUseBreakpoint.mockReturnValue({
      xs: true,
      sm: false,
      md: false,
      lg: false,
      xl: false,
      xxl: false,
    })

    renderWithRouter(<ResponsiveSidebar />)
    
    expect(screen.getByTestId('mobile-sidebar')).toBeInTheDocument()
    expect(screen.queryByTestId('desktop-sidebar')).not.toBeInTheDocument()
  })

  it('deve renderizar Sidebar quando a tela for maior que xs', () => {
    mockUseBreakpoint.mockReturnValue({
      xs: false,
      sm: true,
      md: false,
      lg: false,
      xl: false,
      xxl: false,
    })

    renderWithRouter(<ResponsiveSidebar />)
    
    expect(screen.getByTestId('desktop-sidebar')).toBeInTheDocument()
    expect(screen.queryByTestId('mobile-sidebar')).not.toBeInTheDocument()
  })

  it('deve renderizar Sidebar para telas médias (md)', () => {
    mockUseBreakpoint.mockReturnValue({
      xs: false,
      sm: false,
      md: true,
      lg: false,
      xl: false,
      xxl: false,
    })

    renderWithRouter(<ResponsiveSidebar />)
    
    expect(screen.getByTestId('desktop-sidebar')).toBeInTheDocument()
    expect(screen.queryByTestId('mobile-sidebar')).not.toBeInTheDocument()
  })

  it('deve renderizar Sidebar para telas grandes (lg)', () => {
    mockUseBreakpoint.mockReturnValue({
      xs: false,
      sm: false,
      md: false,
      lg: true,
      xl: false,
      xxl: false,
    })

    renderWithRouter(<ResponsiveSidebar />)
    
    expect(screen.getByTestId('desktop-sidebar')).toBeInTheDocument()
    expect(screen.queryByTestId('mobile-sidebar')).not.toBeInTheDocument()
  })

  it('deve renderizar Sidebar para telas extra grandes (xl)', () => {
    mockUseBreakpoint.mockReturnValue({
      xs: false,
      sm: false,
      md: false,
      lg: false,
      xl: true,
      xxl: false,
    })

    renderWithRouter(<ResponsiveSidebar />)
    
    expect(screen.getByTestId('desktop-sidebar')).toBeInTheDocument()
    expect(screen.queryByTestId('mobile-sidebar')).not.toBeInTheDocument()
  })
})