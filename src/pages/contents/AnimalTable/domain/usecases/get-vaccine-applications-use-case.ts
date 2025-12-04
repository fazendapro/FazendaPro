import { VaccineApplication, VaccineApplicationFilters } from '../model/vaccine-application'

export interface GetVaccineApplicationsUseCase {
  getVaccineApplications: (farmId: number, filters?: VaccineApplicationFilters) => Promise<VaccineApplication[]>
}

