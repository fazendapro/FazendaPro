import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest } from '../../domain/model/reproduction';

export const remoteUpdateReproduction = async (data: CreateReproductionRequest & { id: number }): Promise<void> => {
  const { id, ...updateData } = data;
  await api().put(`/api/v1/reproductions/${id}`, updateData);
};
