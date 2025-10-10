import { RemoteGetNextToCalve } from '../../data/usecases/remote-get-next-to-calve';

export const GetNextToCalveFactory = () => {
  return new RemoteGetNextToCalve();
};
