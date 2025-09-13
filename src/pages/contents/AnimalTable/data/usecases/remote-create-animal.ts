import { api } from '../../../../../components';
import { CreateAnimalDomain, CreateAnimalResponse } from '../../domain/usecases/create-animal-domain';
import { CreateAnimalParams } from '../../types/type';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteCreateAnimal implements CreateAnimalDomain {
  // constructor(
  //   private readonly csrfToken?: string
  // ) {} // TODO: add csrf token

  async create(params: CreateAnimalParams): Promise<CreateAnimalResponse> {
    try {
      const { data, status } = await api().post(
        '/api/v1/animals',
        params,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, ...rest } = data;
      return {
        data: rest,
        status,
        message: message || t('animalTable.animalCreatedSuccessfully'), 
        success: true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao criar animal');
      }
      throw new Error('Erro desconhecido ao criar animal');
    }
  }
} 