import { useState, useEffect, useCallback } from 'react'
import { Vaccine, CreateVaccineRequest } from '../domain/model/vaccine'
import { getVaccinesFactory, createVaccineFactory } from '../factories'

export const useVaccine = (farmId: number) => {
  const [vaccines, setVaccines] = useState<Vaccine[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchVaccines = useCallback(async () => {
    if (!farmId || farmId <= 0) {
      setVaccines([])
      return
    }

    setLoading(true)
    setError(null)
    
    try {
      const getVaccinesUseCase = getVaccinesFactory()
      const data = await getVaccinesUseCase.getVaccines(farmId)
      setVaccines(Array.isArray(data) ? data : [])
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || 'Erro ao carregar vacinas')
    } finally {
      setLoading(false)
    }
  }, [farmId])

  const createVaccine = async (data: CreateVaccineRequest) => {
    try {
      const createVaccineUseCase = createVaccineFactory()
      const newVaccine = await createVaccineUseCase.createVaccine(data)
      setVaccines(prev => {
        const prevArray = Array.isArray(prev) ? prev : []
        return [newVaccine, ...prevArray]
      })
      return newVaccine
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao criar vacina')
    }
  }

  useEffect(() => {
    fetchVaccines()
  }, [fetchVaccines])

  return {
    vaccines,
    loading,
    error,
    refetch: fetchVaccines,
    createVaccine
  }
}

