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
}

export const farmSelectionService = {
  async getUserFarms(): Promise<GetUserFarmsResponse> {
    try {
      const response = await api().get<GetUserFarmsResponse>('/farms/user');
      return response.data;
    } catch {
      throw new Error('Erro ao buscar fazendas do usuário');
    }
  },

  async selectFarm(farmId: number): Promise<SelectFarmResponse> {
    try {
      const response = await api().post<SelectFarmResponse>('/farms/select', {
        farm_id: farmId,
      });
      return response.data;
    } catch {
      throw new Error('Erro ao selecionar fazenda');
    }
  },
};
