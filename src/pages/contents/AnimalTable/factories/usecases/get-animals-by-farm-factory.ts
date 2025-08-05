import { RemoteGetAnimalsByFarm } from '../../data/usecases/remote-get-animals-by-farm';

export const GetAnimalsByFarmFactory = (csrfToken?: string) => {
  return new RemoteGetAnimalsByFarm(csrfToken);
}; 