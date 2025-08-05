import { RemoteCreateAnimal } from '../../data/usecases/remote-create-animal';

export const CreateAnimalFactory = (csrfToken?: string) => {
  return new RemoteCreateAnimal(csrfToken);
};