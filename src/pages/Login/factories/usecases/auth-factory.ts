import { RemoteAuth } from '../../data/usecases/remote-auth';

export const AuthFactory = (csrfToken?: string) => {
  return new RemoteAuth(csrfToken);
};
