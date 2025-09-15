import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest } from '../../domain/model/reproduction';

export const remoteUpdateReproduction = async (data: CreateReproductionRequest): Promise<void> => {
  await api().put('/api/v1/reproductions', data);
};
