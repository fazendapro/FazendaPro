export interface MilkProduction {
  id: number
  animal_id: number
  animal: {
    id: number
    farm_id: number
    ear_tag_number_local: number
    ear_tag_number_register: number
    animal_name: string
    sex: number
    breed: string
    type: string
    birth_date: string
    confinement: boolean
    animal_type: number
    status: number
    fertilization: boolean
    castrated: boolean
    purpose: number
    current_batch: number
  }
  liters: number
  date: string
  created_at: string
  updated_at: string
}

export interface CreateMilkProductionRequest {
  animal_id: number
  liters: number
  date: string
}

export interface MilkProductionFilters {
  period: 'week' | 'month' | 'all'
  startDate?: string
  endDate?: string
  animalId?: number
}
