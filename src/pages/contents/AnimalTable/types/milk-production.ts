export interface MilkProduction {
  id: number
  animalId: number
  animal: {
    id: number
    name: string
    earTagNumberLocal: string
    earTagNumberRegister: string
  }
  liters: number
  date: string
  createdAt: string
  updatedAt: string
}

export interface CreateMilkProductionRequest {
  animalId: number
  liters: number
  date: string
}

export interface MilkProductionFilters {
  period: 'week' | 'month' | 'all'
  startDate?: string
  endDate?: string
  animalId?: number
}
