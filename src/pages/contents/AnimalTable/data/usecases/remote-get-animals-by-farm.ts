import { api } from '../../../../../components';
import { GetAnimalsByFarmDomain, GetAnimalsByFarmResponse } from '../../domain/usecases/get-animals-by-farm-domain';
import { GetAnimalsByFarmParams } from '../../types/type';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetAnimalsByFarm implements GetAnimalsByFarmDomain {

  async getAnimalsByFarm(params: GetAnimalsByFarmParams): Promise<GetAnimalsByFarmResponse> {
    try {
      const requestParams: Record<string, string | number> = {
        farmId: params.farm_id
      };

      if (params.page !== undefined) {
        requestParams.page = params.page;
      }

      if (params.limit !== undefined) {
        requestParams.limit = params.limit;
      }

      const { data, status } = await api().get(
        '/animals/farm',
        {
          params: requestParams,
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, data: responseData } = data;

      let paginatedData;
      if (Array.isArray(responseData)) {
        paginatedData = {
          animals: responseData,
          total: responseData.length,
          page: 1,
          limit: responseData.length
        };
      } else {
        paginatedData = responseData || {
          animals: [],
          total: 0,
          page: 1,
          limit: 10
        };
      }

      return {
        data: paginatedData,
        status,
        message: message || t('animalTable.animalsRetrievedSuccessfully'), 
        success: true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar animais');
      }
      throw new Error('Erro desconhecido ao buscar animais');
    }
  }
} 