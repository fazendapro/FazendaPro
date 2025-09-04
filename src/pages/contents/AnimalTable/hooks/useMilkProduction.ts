import { useState, useEffect } from 'react'
import { MilkProduction, CreateMilkProductionRequest, MilkProductionFilters } from '../types/milk-production'

export const useMilkProduction = (farmId: number, filters?: MilkProductionFilters) => {
  const [milkProductions, setMilkProductions] = useState<MilkProduction[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchMilkProductions = async () => {
    setLoading(true)
    setError(null)
    
    try {
      // TODO: Implementar chamada para API
      // const response = await api.get(`/farms/${farmId}/milk-productions`, { params: filters })
      // setMilkProductions(response.data)
      
      // Dados mockados para desenvolvimento
      const mockData: MilkProduction[] = [
        {
          id: 1,
          animalId: 1,
          animal: {
            id: 1,
            name: "Vaca 001",
            earTagNumberLocal: "001",
            earTagNumberRegister: "BR001"
          },
          liters: 25.5,
          date: "2024-01-15",
          createdAt: "2024-01-15T10:00:00Z",
          updatedAt: "2024-01-15T10:00:00Z"
        },
        {
          id: 2,
          animalId: 2,
          animal: {
            id: 2,
            name: "Vaca 002",
            earTagNumberLocal: "002",
            earTagNumberRegister: "BR002"
          },
          liters: 30.2,
          date: "2024-01-15",
          createdAt: "2024-01-15T10:30:00Z",
          updatedAt: "2024-01-15T10:30:00Z"
        }
      ]
      
      setMilkProductions(mockData)
    } catch (err) {
      setError('Erro ao carregar dados de produção de leite')
    } finally {
      setLoading(false)
    }
  }

  const createMilkProduction = async (data: CreateMilkProductionRequest) => {
    try {
      // TODO: Implementar chamada para API
      // const response = await api.post(`/farms/${farmId}/milk-productions`, data)
      
      // Simular criação
      const newProduction: MilkProduction = {
        id: Date.now(),
        animalId: data.animalId,
        animal: {
          id: data.animalId,
          name: `Vaca ${data.animalId}`,
          earTagNumberLocal: `${data.animalId}`,
          earTagNumberRegister: `BR${data.animalId}`
        },
        liters: data.liters,
        date: data.date,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
      
      setMilkProductions(prev => [...prev, newProduction])
      return newProduction
    } catch (err) {
      throw new Error('Erro ao criar registro de produção de leite')
    }
  }

  useEffect(() => {
    if (farmId) {
      fetchMilkProductions()
    }
  }, [farmId, filters])

  return {
    milkProductions,
    loading,
    error,
    refetch: fetchMilkProductions,
    createMilkProduction
  }
}
