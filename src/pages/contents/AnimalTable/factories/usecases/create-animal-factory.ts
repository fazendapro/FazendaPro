import { RemoteCreateAnimal } from '../../data/usecases/remote-create-animal';


export const CreateAnimalFactory = () => {
  return new RemoteCreateAnimal();
};