import { RemoteGetVaccineApplications } from '../../data/usecases/remote-get-vaccine-applications'

export const getVaccineApplicationsFactory = () => {
  return new RemoteGetVaccineApplications(
    undefined,
  )
}

