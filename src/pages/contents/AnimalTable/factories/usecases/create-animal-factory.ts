import { RemoteCreateAnimal } from '../../data/usecases/remote-create-animal';

// TODO: add csrf token

export const CreateAnimalFactory = () => {
  return new RemoteCreateAnimal();
};