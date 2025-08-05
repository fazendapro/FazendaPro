import { RemoteGetAnimalsByFarm } from '../../data/usecases/remote-get-animals-by-farm';

// TODO: add csrf token
export const GetAnimalsByFarmFactory = () => {
  return new RemoteGetAnimalsByFarm();
}; 