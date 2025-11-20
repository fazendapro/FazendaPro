import { describe, it, expect, beforeEach, vi } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useSale, useSaleForm, useSaleList } from '../useSale'
import { useSaleContext } from '../../contexts/SaleContext'

vi.mock('../../contexts/SaleContext')

describe('useSale', () => {
  const mockContext = {
    sales: [],
    loading: false,
    error: null,
    createSale: vi.fn(),
    updateSale: vi.fn(),
    getSalesByFarm: vi.fn(),
    getSalesHistory: vi.fn(),
    getSalesByAnimal: vi.fn(),
    getSalesByDateRange: vi.fn(),
    deleteSale: vi.fn(),
    clearError: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useSaleContext).mockReturnValue(mockContext)
  })

  it('deve retornar o contexto de vendas', () => {
    const { result } = renderHook(() => useSale())

    expect(result.current).toEqual(mockContext)
  })
})

describe('useSaleForm', () => {
  const mockContext = {
    sales: [],
    loading: false,
    error: null,
    createSale: vi.fn().mockResolvedValue({ id: 1 }),
    updateSale: vi.fn().mockResolvedValue({ id: 1 }),
    getSalesByFarm: vi.fn(),
    getSalesHistory: vi.fn(),
    getSalesByAnimal: vi.fn(),
    getSalesByDateRange: vi.fn(),
    deleteSale: vi.fn(),
    clearError: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useSaleContext).mockReturnValue(mockContext)
  })

  it('deve retornar createSale, updateSale e loading', () => {
    const { result } = renderHook(() => useSaleForm())

    expect(result.current.createSale).toBeDefined()
    expect(result.current.updateSale).toBeDefined()
    expect(result.current.loading).toBe(false)
  })

  it('deve chamar createSale do contexto quando handleCreateSale é chamado', async () => {
    const { result } = renderHook(() => useSaleForm())

    const saleData = {
      animal_id: 1,
      buyer_name: 'Test Buyer',
      price: 1000,
      sale_date: '2024-01-01',
    }

    await result.current.createSale(saleData)

    expect(mockContext.createSale).toHaveBeenCalledWith(saleData)
  })

  it('deve chamar updateSale do contexto quando handleUpdateSale é chamado', async () => {
    const { result } = renderHook(() => useSaleForm())

    const saleData = {
      buyer_name: 'Updated Buyer',
      price: 1500,
      sale_date: '2024-01-01',
      notes: '',
    }

    await result.current.updateSale(1, saleData)

    expect(mockContext.updateSale).toHaveBeenCalledWith(1, saleData)
  })
})

describe('useSaleList', () => {
  const mockContext = {
    sales: [{
      id: 1,
      animal_id: 1,
      farm_id: 1,
      buyer_name: 'Test Buyer',
      price: 1000,
      sale_date: '2024-01-01',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }],
    loading: false,
    error: null,
    createSale: vi.fn(),
    updateSale: vi.fn(),
    getSalesByFarm: vi.fn(),
    getSalesHistory: vi.fn(),
    getSalesByAnimal: vi.fn(),
    getSalesByDateRange: vi.fn(),
    deleteSale: vi.fn(),
    clearError: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useSaleContext).mockReturnValue(mockContext)
  })

  it('deve retornar todas as propriedades do contexto de vendas', () => {
    const { result } = renderHook(() => useSaleList())

    expect(result.current.sales).toEqual(mockContext.sales)
    expect(result.current.loading).toBe(false)
    expect(result.current.error).toBe(null)
    expect(result.current.getSalesByFarm).toBe(mockContext.getSalesByFarm)
    expect(result.current.getSalesHistory).toBe(mockContext.getSalesHistory)
    expect(result.current.getSalesByAnimal).toBe(mockContext.getSalesByAnimal)
    expect(result.current.getSalesByDateRange).toBe(mockContext.getSalesByDateRange)
    expect(result.current.deleteSale).toBe(mockContext.deleteSale)
    expect(result.current.clearError).toBe(mockContext.clearError)
  })
})




