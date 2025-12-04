import { useState, useEffect, useCallback } from 'react'
import { VaccineApplication, CreateVaccineApplicationRequest, VaccineApplicationFilters, UpdateVaccineApplicationRequest } from '../domain/model/vaccine-application'
import { getVaccineApplicationsFactory, createVaccineApplicationFactory, updateVaccineApplicationFactory, deleteVaccineApplicationFactory } from '../factories'

export const useVaccineApplication = (farmId: number, filters?: VaccineApplicationFilters) => {
  const [vaccineApplications, setVaccineApplications] = useState<VaccineApplication[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchVaccineApplications = useCallback(async () => {
    if (!farmId || farmId <= 0) {
      setVaccineApplications([])
      return
    }

    setLoading(true)
    setError(null)
    
    try {
      const getVaccineApplicationsUseCase = getVaccineApplicationsFactory()
      const data = await getVaccineApplicationsUseCase.getVaccineApplications(farmId, filters)
      setVaccineApplications(Array.isArray(data) ? data : [])
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || 'Erro ao carregar aplicações de vacinas')
    } finally {
      setLoading(false)
    }
  }, [farmId, filters])

  const createVaccineApplication = async (data: CreateVaccineApplicationRequest) => {
    try {
      const createVaccineApplicationUseCase = createVaccineApplicationFactory()
      const newApplication = await createVaccineApplicationUseCase.createVaccineApplication(data)
      setVaccineApplications(prev => {
        const prevArray = Array.isArray(prev) ? prev : []
        return [newApplication, ...prevArray]
      })
      return newApplication
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao criar aplicação de vacina')
    }
  }

  const updateVaccineApplication = async (id: number, data: UpdateVaccineApplicationRequest) => {
    try {
      const updateVaccineApplicationUseCase = updateVaccineApplicationFactory()
      const updatedApplication = await updateVaccineApplicationUseCase.updateVaccineApplication(id, data)
      setVaccineApplications(prev => {
        const prevArray = Array.isArray(prev) ? prev : []
        return prevArray.map(application => 
          application.id === id ? updatedApplication : application
        )
      })
      return updatedApplication
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao atualizar aplicação de vacina')
    }
  }

  const deleteVaccineApplication = async (id: number) => {
    try {
      const deleteVaccineApplicationUseCase = deleteVaccineApplicationFactory()
      await deleteVaccineApplicationUseCase.deleteVaccineApplication(id)
      setVaccineApplications(prev => {
        const prevArray = Array.isArray(prev) ? prev : []
        return prevArray.filter(application => application.id !== id)
      })
    } catch (err: unknown) {
      const error = err as Error
      throw new Error(error.message || 'Erro ao deletar aplicação de vacina')
    }
  }

  useEffect(() => {
    fetchVaccineApplications()
  }, [fetchVaccineApplications])

  return {
    vaccineApplications,
    loading,
    error,
    refetch: fetchVaccineApplications,
    createVaccineApplication,
    updateVaccineApplication,
    deleteVaccineApplication
  }
}

