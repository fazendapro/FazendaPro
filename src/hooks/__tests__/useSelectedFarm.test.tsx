import { describe, it, expect, beforeEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { useSelectedFarm } from '../useSelectedFarm'
import { FarmProvider } from '../../contexts/FarmContext'
import { GetFarmFactory } from '../../pages/contents/Settings/factories'
import { ReactNode } from 'react'

vi.mock('../../pages/contents/Settings/factories')

describe('useSelectedFarm', () => {
  const mockSelectedFarm = {
    ID: 1,
    CompanyID: 1,
    Logo: '',
    Company: {
      ID: 1,
      CompanyName: 'Fazenda Teste',
      Location: 'Test Location',
      FarmCNPJ: '12345678901234'
    }
  }

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('deve retornar farmId corretamente', () => {
    localStorage.setItem('selectedFarm', JSON.stringify(mockSelectedFarm))
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.farmId).toBe(1)
  })

  it('deve retornar farmName do Company quando disponível', () => {
    localStorage.setItem('selectedFarm', JSON.stringify(mockSelectedFarm))
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.farmName).toBe('Fazenda Teste')
  })

  it('deve retornar farmName padrão quando Company não está disponível', () => {
    const farmWithoutCompany = {
      ID: 2,
      CompanyID: 2,
      Logo: '',
    }
    localStorage.setItem('selectedFarm', JSON.stringify(farmWithoutCompany))
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.farmName).toBe('Fazenda 2')
  })

  it('deve retornar null quando selectedFarm não está disponível', () => {
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.farmId).toBe(null)
    expect(result.current.farmName).toBe(null)
  })

  it('deve carregar logo da fazenda quando farmId está disponível', async () => {
    const mockGetFarmUseCase = {
      get: vi.fn().mockResolvedValue({
        data: {
          Logo: 'https://example.com/logo.png'
        }
      })
    }

    vi.mocked(GetFarmFactory.create).mockReturnValue(mockGetFarmUseCase)
    localStorage.setItem('selectedFarm', JSON.stringify(mockSelectedFarm))
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    await waitFor(() => {
      expect(result.current.farmLogo).toBe('https://example.com/logo.png')
    }, { timeout: 3000 })

    expect(mockGetFarmUseCase.get).toHaveBeenCalledWith(1)
  })

  it('deve limpar logo quando farmId não está disponível', () => {
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.farmLogo).toBe('')
  })

  it('deve expor setSelectedFarm e clearSelectedFarm', () => {
    localStorage.setItem('selectedFarm', JSON.stringify(mockSelectedFarm))
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <FarmProvider>{children}</FarmProvider>
    )

    const { result } = renderHook(() => useSelectedFarm(), { wrapper })

    expect(result.current.setSelectedFarm).toBeDefined()
    expect(result.current.clearSelectedFarm).toBeDefined()
  })
})

