import { api } from '../../../../../components/services/axios/api';
import { Reproduction } from '../../domain/model/reproduction';

export const remoteGetReproductionsByFarm = async (farmId: number): Promise<Reproduction[]> => {
  const response = await api().get(`/api/v1/reproductions/farm?farmId=${farmId}`);
  return response.data.data;
};
