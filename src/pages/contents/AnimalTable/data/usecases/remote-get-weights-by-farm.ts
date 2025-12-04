import { api } from '../../../../../components/services/axios/api';
import { Weight } from '../../domain/model/weight';

export const remoteGetWeightsByFarm = async (farmId: number): Promise<Weight[]> => {
  const response = await api().get(`/weights/farm/${farmId}`);
  return response.data.data || [];
};

