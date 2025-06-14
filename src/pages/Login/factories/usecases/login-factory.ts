import { RemoteLogin } from '../../data/usecases/remote-login';

export const LoginFactory = (csrfToken?: string) => {
  return new RemoteLogin(csrfToken);
};