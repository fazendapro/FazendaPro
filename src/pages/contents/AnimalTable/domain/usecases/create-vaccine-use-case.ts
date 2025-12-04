import { Vaccine, CreateVaccineRequest } from '../model/vaccine'

export interface CreateVaccineUseCase {
  createVaccine: (data: CreateVaccineRequest) => Promise<Vaccine>
}

