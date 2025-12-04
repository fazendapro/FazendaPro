export interface VaccineApplication {
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
  vaccine_id: number
  vaccine: {
    id: number
    farm_id: number
    name: string
    description?: string
    manufacturer?: string
    created_at: string
    updated_at: string
  }
  application_date: string
  batch_number?: string
  veterinarian?: string
  observations?: string
  created_at: string
  updated_at: string
}

export interface CreateVaccineApplicationRequest {
  animal_id: number
  vaccine_id: number
  application_date: string
  batch_number?: string
  veterinarian?: string
  observations?: string
}

export interface UpdateVaccineApplicationRequest {
  animal_id: number
  vaccine_id: number
  application_date: string
  batch_number?: string
  veterinarian?: string
  observations?: string
}

export interface VaccineApplicationFilters {
  startDate?: string
  endDate?: string
  animalId?: number
  vaccineId?: number
}

