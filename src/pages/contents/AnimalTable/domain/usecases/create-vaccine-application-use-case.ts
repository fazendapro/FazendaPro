import { VaccineApplication, CreateVaccineApplicationRequest } from '../model/vaccine-application'

export interface CreateVaccineApplicationUseCase {
  createVaccineApplication: (data: CreateVaccineApplicationRequest) => Promise<VaccineApplication>
}

