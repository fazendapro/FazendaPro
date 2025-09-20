import { BaseHttpResponse } from '../../../../components/services/axios/base-https-response';

export interface RefreshTokenParams {
  refresh_token: string;
}

export interface RefreshTokenResponse extends BaseHttpResponse<{
  access_token: string;
}> {
  success: boolean;
  message: string;
  status: number;
  access_token: string;
}

export interface LogoutParams {
  refresh_token: string;
}

export interface LogoutResponse extends BaseHttpResponse<null> {
  success: boolean;
  message: string;
  status: number;
}

export interface AuthDomain {
  refreshToken(params: RefreshTokenParams): Promise<RefreshTokenResponse>;
  logout(params: LogoutParams): Promise<LogoutResponse>;
}
