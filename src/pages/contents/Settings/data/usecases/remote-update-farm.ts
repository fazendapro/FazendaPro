import { api } from '../../../../../components';
import { UpdateFarmDomain, UpdateFarmResponse, GetFarmDomain, GetFarmResponse } from '../../domain/usecases/update-farm-domain';
import { UpdateFarmParams } from '../../types/farm-types';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteUpdateFarm implements UpdateFarmDomain {
  async update(farmId: number, params: UpdateFarmParams): Promise<UpdateFarmResponse> {
    try {
      const { data, status } = await api().put(
        `/api/v1/farm?id=${farmId}`,
        params,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      return {
        data: data.data,
        status,
        message: data.message || t('settings.farmUpdatedSuccessfully'), 
        success: data.success
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao atualizar fazenda');
      }
      throw new Error('Erro desconhecido ao atualizar fazenda');
    }
  }
}

export class RemoteGetFarm implements GetFarmDomain {
  async get(farmId: number): Promise<GetFarmResponse> {
    try {
      const { data, status } = await api().get(
        `/api/v1/farm?id=${farmId}`,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      return {
        data: data.data,
        status,
        message: data.message || t('settings.farmRetrievedSuccessfully'), 
        success: data.success
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar fazenda');
      }
      throw new Error('Erro desconhecido ao buscar fazenda');
    }
  }
}
