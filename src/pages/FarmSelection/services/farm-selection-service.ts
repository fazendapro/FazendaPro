import { api } from '../../../components/services/axios/api';

export interface Farm {
  ID: number;
  CompanyID: number;
  Logo: string;
  Company?: {
    ID: number;
    CompanyName: string;
    Location: string;
    FarmCNPJ: string;
  };
}

export interface GetUserFarmsResponse {
  success: boolean;
  message: string;
  farms: Farm[];
  auto_select: boolean;
  selected_farm_id?: number;
}

export interface SelectFarmRequest {
  farm_id: number;
}

export interface SelectFarmResponse {
  success: boolean;
  message: string;
  farm_id: number;
  access_token: string;
}

export const farmSelectionService = {
  async getUserFarms(): Promise<GetUserFarmsResponse> {
    try {
      const response = await api().get<GetUserFarmsResponse>('/farms/user');

      if (response.data && !response.data.success) {
        const errorMessage = response.data.message || 'Erro ao buscar fazendas do usuário';
        if (import.meta.env.DEV) {
          console.error('Farm Selection Service Error:', {
            status: response.status,
            message: errorMessage,
            data: response.data
          });
        }
        throw new Error(errorMessage);
      }
      
      return response.data;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      
      const axiosError = error as { response?: { data?: { message?: string }; status?: number } };
      const errorMessage = axiosError?.response?.data?.message || 'Erro ao buscar fazendas do usuário';
      
      if (import.meta.env.DEV) {
        console.error('Farm Selection Service Network Error:', {
          status: axiosError?.response?.status,
          message: errorMessage,
          error
        });
      }
      
      throw new Error(errorMessage);
    }
  },

  async selectFarm(farmId: number): Promise<SelectFarmResponse> {
    try {
      const response = await api().post<SelectFarmResponse>('/farms/select', {
        farm_id: farmId,
      });
      
      if (response.data && !response.data.success) {
        const errorMessage = response.data.message || 'Erro ao selecionar fazenda';
        if (import.meta.env.DEV) {
          console.error('Farm Selection Service Error:', {
            status: response.status,
            message: errorMessage,
            data: response.data
          });
        }
        throw new Error(errorMessage);
      }
      
      return response.data;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }

      const axiosError = error as { response?: { data?: { message?: string }; status?: number } };
      const errorMessage = axiosError?.response?.data?.message || 'Erro ao selecionar fazenda';

      if (import.meta.env.DEV) {
        console.error('Farm Selection Service Network Error:', {
          status: axiosError?.response?.status,
          message: errorMessage,
          error
        });
      }

      throw new Error(errorMessage);
    }
  },
};
