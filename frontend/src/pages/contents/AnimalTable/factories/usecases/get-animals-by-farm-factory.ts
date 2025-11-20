import { RemoteGetAnimalsByFarm } from '../../data/usecases/remote-get-animals-by-farm';

export const GetAnimalsByFarmFactory = () => {
  return new RemoteGetAnimalsByFarm();
}; 