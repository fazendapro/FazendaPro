import { useState, useEffect, useCallback } from 'react'
import { MilkProduction, CreateMilkProductionRequest, MilkProductionFilters } from '../domain/model/milk-production'
import { UpdateMilkProductionRequest } from '../domain/usecases/update-milk-production-use-case'
import { getMilkProductionsFactory, createMilkProductionFactory, updateMilkProductionFactory } from '../factories'

export const useMilkProduction = (farmId: number, filters?: MilkProductionFilters) => {
  const [milkProductions, setMilkProductions] = useState<MilkProduction[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchMilkProductions = useCallback(async () => {
    setLoading(true)
    setError(null)
    
    try {
      const getMilkProductionsUseCase = getMilkProductionsFactory()
      const data = await getMilkProductionsUseCase.getMilkProductions(farmId, filters)
      setMilkProductions(data)
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || 'Erro ao carregar dados de produção de leite')
    } finally {
      setLoading(false)
    }
  }, [farmId, filters])

  const createMilkProduction = async (data: CreateMilkProductionRequest) => {
    try {
      const createMilkProductionUseCase = createMilkProductionFactory()
      const newProduction = await createMilkProductionUseCase.createMilkProduction(data)
      setMilkProductions(prev => [newProduction, ...prev])
      return newProduction
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao criar registro de produção de leite')
    }
  }

  const updateMilkProduction = async (data: UpdateMilkProductionRequest) => {
    try {
      const updateMilkProductionUseCase = updateMilkProductionFactory()
      const updatedProduction = await updateMilkProductionUseCase.updateMilkProduction(data)
      setMilkProductions(prev => 
        prev.map(production => 
          production.id === data.id ? updatedProduction : production
        )
      )
      return updatedProduction
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao atualizar registro de produção de leite')
    }
  }

  useEffect(() => {
    if (farmId && farmId > 0) {
      fetchMilkProductions()
    } else {
      setMilkProductions([])
    }
  }, [farmId, filters, fetchMilkProductions])

  return {
    milkProductions,
    loading,
    error,
    refetch: fetchMilkProductions,
    createMilkProduction,
    updateMilkProduction
  }
}
