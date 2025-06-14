import { api } from '../../../../components';
import { LoginDomain, LoginParams, LoginResponse } from '../../domain/usecases/login-domain';
import { AxiosError } from 'axios';

const axiosInstance = api(import.meta.env.VITE_API_URL, 'auth');

export class RemoteLogin implements LoginDomain {
  constructor(
    private readonly csrfToken?: string
  ) {}

  async authenticate({ email, password }: LoginParams): Promise<LoginResponse> {
    try {
      const { data, status } = await axiosInstance.post(
        '/login',
        { email, password },
        {
          headers: {
            'X-CSRF-Token': this.csrfToken,
          },
        }
      );

      const { message, ...rest } = data;
      return { data: rest, status, message: message || 'Login realizado com sucesso', success: true, access_token: rest.access_token };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro no login');
      }
      throw new Error('Erro desconhecido durante o login');
    }
  }
}