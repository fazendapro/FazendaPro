import { RemoteUpdateVaccineApplication } from '../../data/usecases/remote-update-vaccine-application'

export const updateVaccineApplicationFactory = () => {
  return new RemoteUpdateVaccineApplication(
    undefined,
    undefined
  )
}

