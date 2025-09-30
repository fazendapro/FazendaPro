import { api } from '../../../../components';
import { LoginDomain, LoginParams, LoginResponse } from '../../domain/usecases/login-domain';
import { AxiosError } from 'axios';

const axiosInstance = api('');

export class RemoteLogin implements LoginDomain {
  constructor(
    private readonly csrfToken?: string
  ) {}

  async authenticate({ email, password }: LoginParams): Promise<LoginResponse> {
    try {
      const { data, status } = await axiosInstance.post(
        '/api/v1/auth/login',
        { email, password },
        {
          headers: {
            'X-CSRF-Token': this.csrfToken,
          },
        }
      );

      const { message, access_token, refresh_token } = data;
      return { 
        data: { access_token, refresh_token }, 
        status, 
        message: message || 'Login realizado com sucesso', 
        success: true, 
        access_token: access_token,
        refresh_token: refresh_token
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro no login');
      }
      throw new Error('Erro desconhecido durante o login');
    }
  }
}