import { VaccineApplication, VaccineApplicationFilters } from '../model/vaccine-application'
import { GetVaccineApplicationsResponse } from '../../data/usecases/remote-get-vaccine-applications'

export interface GetVaccineApplicationsUseCase {
  getVaccineApplications: (farmId: number, filters?: VaccineApplicationFilters) => Promise<VaccineApplication[] | GetVaccineApplicationsResponse>
}

