import { api } from '../../../../components';
import { AuthDomain, RefreshTokenParams, RefreshTokenResponse, LogoutParams, LogoutResponse } from '../../domain/usecases/auth-domain';
import { AxiosError } from 'axios';

const axiosInstance = api('');

export class RemoteAuth implements AuthDomain {
  constructor(
    private readonly csrfToken?: string
  ) {}

  async refreshToken({ refresh_token }: RefreshTokenParams): Promise<RefreshTokenResponse> {
    try {
      const { data, status } = await axiosInstance.post(
        '/api/v1/auth/refresh',
        { refresh_token },
        {
          headers: {
            'X-CSRF-Token': this.csrfToken,
          },
        }
      );

      const { message, access_token } = data;
      return {
        data: { access_token },
        status,
        message: message || 'Token renovado com sucesso',
        success: true,
        access_token: access_token
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao renovar token');
      }
      throw new Error('Erro desconhecido durante a renovação do token');
    }
  }

  async logout({ refresh_token }: LogoutParams): Promise<LogoutResponse> {
    try {
      const { data, status } = await axiosInstance.post(
        '/api/v1/auth/logout',
        { refresh_token },
        {
          headers: {
            'X-CSRF-Token': this.csrfToken,
          },
        }
      );

      const { message } = data;
      return { 
        data: null, 
        status, 
        message: message || 'Logout realizado com sucesso', 
        success: true 
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro no logout');
      }
      throw new Error('Erro desconhecido durante o logout');
    }
  }
}
