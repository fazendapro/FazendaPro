import { VaccineApplication, UpdateVaccineApplicationRequest } from '../model/vaccine-application'

export interface UpdateVaccineApplicationUseCase {
  updateVaccineApplication: (id: number, data: UpdateVaccineApplicationRequest) => Promise<VaccineApplication>
}

