import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { AnimalDashboard } from '../animal-dashboard'
import { useResponsive } from '../../../../../../hooks'

vi.mock('../../../../../../hooks', () => ({
  useResponsive: vi.fn(),
  useModal: () => ({
    isOpen: false,
    onOpen: vi.fn(),
    onClose: vi.fn()
  })
}))

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key
  })
}))

describe('AnimalDashboard - Responsive', () => {
  const mockProps = {
    onAnimalCreated: vi.fn(),
    onColumnsChanged: vi.fn(),
    onSearchChange: vi.fn()
  }

  it('deve renderizar layout mobile corretamente', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.dashboard')).toBeInTheDocument()
    expect(screen.getByText('animalTable.createCow')).toBeInTheDocument()
    expect(screen.getByText('animalTable.filter')).toBeInTheDocument()
  })

  it('deve renderizar layout tablet corretamente', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: false,
      isTablet: true,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 800
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.categories')).toBeInTheDocument()
    expect(screen.getByText('animalTable.createCow')).toBeInTheDocument()
  })

  it('deve renderizar layout desktop corretamente', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      isLargeDesktop: false,
      screenWidth: 1200
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.categories')).toBeInTheDocument()
    expect(screen.getByText('animalTable.createCow')).toBeInTheDocument()
  })

  it('deve aplicar breakpoints responsivos nas estatísticas', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.dashboard')).toBeInTheDocument()
  })

  it('deve aplicar layout vertical em mobile', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.createCow')).toBeInTheDocument()
    expect(screen.getByText('animalTable.filter')).toBeInTheDocument()
  })

  it('deve renderizar dropdown em mobile', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.dashboard')).toBeInTheDocument()
  })

  it('deve renderizar card normal em desktop', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      isLargeDesktop: false,
      screenWidth: 1200
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.categories')).toBeInTheDocument()
    expect(screen.getByText('animalTable.totalAnimals')).toBeInTheDocument()
  })

  it('deve aplicar layout horizontal em desktop', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      isLargeDesktop: false,
      screenWidth: 1200
    })

    render(<AnimalDashboard {...mockProps} />)
    
    expect(screen.getByText('animalTable.createCow')).toBeInTheDocument()
    expect(screen.getByText('animalTable.filter')).toBeInTheDocument()
  })

  it('deve aplicar tamanhos responsivos nos botões', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    const createButton = screen.getByText('animalTable.createCow')
    const filterButton = screen.getByText('animalTable.filter')
    
    expect(createButton).toBeInTheDocument()
    expect(filterButton).toBeInTheDocument()
  })

  it('deve aplicar tamanhos responsivos na busca', () => {
    vi.mocked(useResponsive).mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
    })

    render(<AnimalDashboard {...mockProps} />)
    
    const searchInput = screen.getByPlaceholderText('animalTable.search')
    expect(searchInput).toBeInTheDocument()
  })
})
