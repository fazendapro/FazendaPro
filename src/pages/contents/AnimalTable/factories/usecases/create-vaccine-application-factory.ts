import { RemoteCreateVaccineApplication } from '../../data/usecases/remote-create-vaccine-application'

export const createVaccineApplicationFactory = () => {
  return new RemoteCreateVaccineApplication(
    undefined,
    undefined
  )
}

