import { api } from '../../../../components/services/axios/api';

export interface FarmUpdateData {
  name: string;
  logo: string;
}

export interface FarmData {
  id: number;
  name: string;
  logo: string;
  company_id: number;
  created_at: string;
  updated_at: string;
}

export const farmService = {
  async getFarm(farmId: number): Promise<FarmData> {
    const response = await api().get(`/farms?id=${farmId}`);
    return response.data.data;
  },

  async updateFarm(farmId: number, data: FarmUpdateData): Promise<FarmData> {
    const response = await api().put(`/farms?id=${farmId}`, data);
    return response.data.data;
  },
};
