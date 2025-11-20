import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest } from '../../domain/model/reproduction';

export const remoteUpdateReproduction = async (data: CreateReproductionRequest & { id: number }): Promise<void> => {
  const { id, ...updateData } = data;
  await api().put(`/reproductions/${id}`, updateData);
};
