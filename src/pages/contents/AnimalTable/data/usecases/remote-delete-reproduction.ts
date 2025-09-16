import { api } from '../../../../../components/services/axios/api';

export const remoteDeleteReproduction = async (id: number): Promise<void> => {
  await api().delete(`/api/v1/reproductions?id=${id}`);
};
