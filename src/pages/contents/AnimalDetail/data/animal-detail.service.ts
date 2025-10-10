import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
import { 
  AnimalDetailResponse, 
  AnimalDetailUpdateRequest, 
  AnimalParentsResponse 
} from '../types';

interface ApiError {
  response?: {
    data?: {
      message?: string;
    };
  };
}

export class AnimalDetailService {
  async getAnimalById(id: number): Promise<AnimalDetailResponse> {
    try {
      const response = await api.get(`/animals?id=${id}`);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error: unknown) {
      return {
        success: false,
        message: (error as ApiError).response?.data?.message || 'Erro ao buscar animal'
      };
    }
  }

  async updateAnimal(animalData: AnimalDetailUpdateRequest): Promise<AnimalDetailResponse> {
    try {
      const response = await api.put('/animals', animalData);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error: unknown) {
      return {
        success: false,
        message: (error as ApiError).response?.data?.message || 'Erro ao atualizar animal'
      };
    }
  }

  async getAnimalsBySex(farmId: number, sex: number): Promise<AnimalParentsResponse> {
    try {
      const response = await api.get(`/animals/sex?farmId=${farmId}&sex=${sex}`);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error: unknown) {
      return {
        success: false,
        message: (error as ApiError).response?.data?.message || 'Erro ao buscar animais por sexo'
      };
    }
  }

  async uploadAnimalPhoto(animalId: number, photo: File): Promise<AnimalDetailResponse> {
    try {
      const formData = new FormData();
      formData.append('photo', photo);
      formData.append('animal_id', animalId.toString());

      const response = await api.post('/animals/photo', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });

      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error: unknown) {
      return {
        success: false,
        message: (error as ApiError).response?.data?.message || 'Erro ao fazer upload da foto'
      };
    }
  }
}

export const animalDetailService = new AnimalDetailService();
