import { api } from '../../../../../components';
import { GetAnimalsByFarmDomain, GetAnimalsByFarmResponse } from '../../domain/usecases/get-animals-by-farm-domain';
import { GetAnimalsByFarmParams } from '../../types/type';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetAnimalsByFarm implements GetAnimalsByFarmDomain {
  // constructor(
  //   private readonly csrfToken?: string
  // ) {} // TODO: add csrf token

  async getAnimalsByFarm(params: GetAnimalsByFarmParams): Promise<GetAnimalsByFarmResponse> {
    try {
      const { data, status } = await api().get(
        '/api/v1/animals/farm',
        {
          params: {
            farmId: params.farm_id
          }
          // headers: {
          //   'X-CSRF-Token': this.csrfToken || '', // TODO: add csrf token
          // },
        }
      );

      const { message, data: animalsData } = data;
      return {
        data: animalsData || [],
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