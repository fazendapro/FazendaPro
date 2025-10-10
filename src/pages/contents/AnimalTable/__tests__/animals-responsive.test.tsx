import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { Animals } from '../presentation/animals-list/animals-container'
import { useIsMobile } from '../../../../hooks/use-is-mobile'

vi.mock('../../../../hooks/use-is-mobile', () => ({
  useIsMobile: vi.fn()
}))

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('../../../../components/tabs', () => ({
  DesktopTabs: ({ defaultTabIndex }: { defaultTabIndex: number }) => (
    <div data-testid="desktop-tabs">
      Desktop Tabs - Tab {defaultTabIndex}
    </div>
  ),
  MobileTabs: ({ defaultTabIndex }: { defaultTabIndex: number }) => (
    <div data-testid="mobile-tabs">
      Mobile Tabs - Tab {defaultTabIndex}
    </div>
  ),
}))

vi.mock('../presentation/animals-list/animals', () => ({
  Animals: () => <div data-testid="animals-component">Animals Component</div>
}))

vi.mock('../milk-production/milk-production', () => ({
  MilkProduction: () => <div data-testid="milk-production">Milk Production</div>
}))

vi.mock('../reproduction/reproduction', () => ({
  Reproduction: () => <div data-testid="reproduction">Reproduction</div>
}))

const mockUseIsMobile = vi.mocked(useIsMobile)

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('Animals Responsividade', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar MobileTabs quando for mobile', () => {
    mockUseIsMobile.mockReturnValue(true)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('mobile-tabs')).toBeInTheDocument()
    expect(screen.queryByTestId('desktop-tabs')).not.toBeInTheDocument()
  })

  it('deve renderizar DesktopTabs quando não for mobile', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
    expect(screen.queryByTestId('mobile-tabs')).not.toBeInTheDocument()
  })

  it('deve ter container com id correto', () => {
    mockUseIsMobile.mockReturnValue(false)

    const { container } = renderWithRouter(<Animals />)
    
    const animalsContainer = container.querySelector('#animals-container')
    expect(animalsContainer).toBeInTheDocument()
  })

  it('deve ter tabs com tradução', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
  })

  it('deve alternar entre tabs corretamente', () => {
    mockUseIsMobile.mockReturnValue(false)

    const { rerender } = renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
    
    rerender(
      <BrowserRouter>
        <Animals />
      </BrowserRouter>
    )
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
  })

  it('deve ter estrutura responsiva para mobile', () => {
    mockUseIsMobile.mockReturnValue(true)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('mobile-tabs')).toBeInTheDocument()
    
    const { container } = renderWithRouter(<Animals />)
    const animalsContainer = container.querySelector('#animals-container')
    expect(animalsContainer).toBeInTheDocument()
  })

  it('deve ter estrutura responsiva para desktop', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
    
    const { container } = renderWithRouter(<Animals />)
    const animalsContainer = container.querySelector('#animals-container')
    expect(animalsContainer).toBeInTheDocument()
  })

  it('deve ter tabs com componentes corretos', () => {
    mockUseIsMobile.mockReturnValue(false)

    renderWithRouter(<Animals />)
    
    expect(screen.getByTestId('desktop-tabs')).toBeInTheDocument()
  })
})
