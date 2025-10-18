import { api } from '../../../../../components/services/axios/api';

export const remoteDeleteReproduction = async (id: number): Promise<void> => {
  await api().delete(`/reproductions?id=${id}`);
};
