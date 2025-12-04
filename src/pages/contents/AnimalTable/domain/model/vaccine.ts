export interface Vaccine {
  id: number
  farm_id: number
  name: string
  description?: string
  manufacturer?: string
  created_at: string
  updated_at: string
}

export interface CreateVaccineRequest {
  farm_id: number
  name: string
  description?: string
  manufacturer?: string
}

export interface UpdateVaccineRequest {
  farm_id: number
  name: string
  description?: string
  manufacturer?: string
}

