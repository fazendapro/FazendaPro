import { describe, it, expect } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useFarm } from '../useFarm'

describe('useFarm', () => {
  it('deve retornar dados mockados da fazenda', () => {
    const { result } = renderHook(() => useFarm())

    expect(result.current.farm).toEqual({
      id: 1,
      name: 'Fazenda 1',
      location: 'Rua 1, 123',
      created_at: '2021-01-01',
      updated_at: '2021-01-01'
    })
  })

  it('deve retornar um objeto farm válido', () => {
    const { result } = renderHook(() => useFarm())

    expect(result.current.farm).toBeDefined()
    expect(result.current.farm?.id).toBe(1)
    expect(result.current.farm?.name).toBe('Fazenda 1')
    expect(result.current.farm?.location).toBe('Rua 1, 123')
  })

  it('deve manter a mesma referência entre re-renders', () => {
    const { result, rerender } = renderHook(() => useFarm())

    const firstFarm = result.current.farm

    rerender()

    expect(result.current.farm).toStrictEqual(firstFarm)
  })

  it('deve ter estrutura consistente', () => {
    const { result } = renderHook(() => useFarm())

    const farm = result.current.farm

    expect(farm).toHaveProperty('id')
    expect(farm).toHaveProperty('name')
    expect(farm).toHaveProperty('location')
    expect(farm).toHaveProperty('created_at')
    expect(farm).toHaveProperty('updated_at')

    expect(typeof farm?.id).toBe('number')
    expect(typeof farm?.name).toBe('string')
    expect(typeof farm?.location).toBe('string')
    expect(typeof farm?.created_at).toBe('string')
    expect(typeof farm?.updated_at).toBe('string')
  })
})
