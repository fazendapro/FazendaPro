import { useState, useEffect, useCallback, useMemo, useRef } from 'react'
import { VaccineApplication, CreateVaccineApplicationRequest, VaccineApplicationFilters, UpdateVaccineApplicationRequest } from '../domain/model/vaccine-application'
import { getVaccineApplicationsFactory, createVaccineApplicationFactory, updateVaccineApplicationFactory, deleteVaccineApplicationFactory } from '../factories'
import { GetVaccineApplicationsResponse } from '../data/usecases/remote-get-vaccine-applications'

export const useVaccineApplication = (farmId: number, filters?: VaccineApplicationFilters) => {
  const [vaccineApplications, setVaccineApplications] = useState<VaccineApplication[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const filtersRef = useRef(filters)

  useEffect(() => {
    filtersRef.current = filters
  }, [filters])

  const filterKey = useMemo(() => {
    if (!filters) return ''
    return JSON.stringify({
      startDate: filters.startDate,
      endDate: filters.endDate,
      animalId: filters.animalId,
      vaccineId: filters.vaccineId,
      page: filters.page,
      limit: filters.limit
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters?.startDate, filters?.endDate, filters?.animalId, filters?.vaccineId, filters?.page, filters?.limit])

  const fetchVaccineApplications = useCallback(async () => {
    if (!farmId || farmId <= 0) {
      setVaccineApplications([])
      setTotal(0)
      return
    }

    setLoading(true)
    setError(null)
    
    try {
      const getVaccineApplicationsUseCase = getVaccineApplicationsFactory()
      const data = await getVaccineApplicationsUseCase.getVaccineApplications(farmId, filtersRef.current)
      
      if (data && typeof data === 'object' && 'vaccine_applications' in data) {
        const paginatedData = data as GetVaccineApplicationsResponse
        setVaccineApplications(paginatedData.vaccine_applications || [])
        setTotal(paginatedData.total || 0)
      } else {
        const applications = Array.isArray(data) ? data : []
        setVaccineApplications(applications)
        setTotal(applications.length)
      }
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || 'Erro ao carregar aplicações de vacinas')
      setVaccineApplications([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }, [farmId])

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
  }, [farmId, filterKey, fetchVaccineApplications])

  return {
    vaccineApplications,
    total,
    loading,
    error,
    refetch: fetchVaccineApplications,
    createVaccineApplication,
    updateVaccineApplication,
    deleteVaccineApplication
  }
}

