import { BaseHttpResponse } from '../../../../components/services/axios/base-https-response';

export interface LoginParams {
  email: string;
  password: string;
}

export interface LoginResponse extends BaseHttpResponse<{
  access_token: string;
}> {
  success: boolean;
  message: string;
  status: number;
  access_token: string;
}

export interface LoginDomain {
  authenticate(params: LoginParams): Promise<LoginResponse>;
}