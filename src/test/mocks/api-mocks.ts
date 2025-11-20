import { vi } from 'vitest'

interface MockAxiosInstance {
  get: ReturnType<typeof vi.fn>
  post: ReturnType<typeof vi.fn>
  put: ReturnType<typeof vi.fn>
  delete: ReturnType<typeof vi.fn>
  patch: ReturnType<typeof vi.fn>
  create: ReturnType<typeof vi.fn>
  interceptors: {
    request: {
      use: ReturnType<typeof vi.fn>
      eject: ReturnType<typeof vi.fn>
    }
    response: {
      use: ReturnType<typeof vi.fn>
      eject: ReturnType<typeof vi.fn>
    }
  }
}

export const mockAxios: MockAxiosInstance = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
  patch: vi.fn(),
  create: vi.fn((): MockAxiosInstance => mockAxios),
  interceptors: {
    request: {
      use: vi.fn(),
      eject: vi.fn(),
    },
    response: {
      use: vi.fn(),
      eject: vi.fn(),
    },
  },
}

export const createMockApiResponse = <T>(data: T, success = true, message = 'Success') => ({
  data,
  success,
  message,
  status: success ? 200 : 400,
})

export const createMockApiError = (message = 'API Error', status = 400) => ({
  response: {
    data: {
      message,
      success: false,
    },
    status,
  },
  message,
})

export const mockAnimalsData = [
  {
    id: 1,
    name: 'Vaca 001',
    ear_tag: '001',
    breed: 'Holandesa',
    birth_date: '2020-01-01',
    farm_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Vaca 002',
    ear_tag: '002',
    breed: 'Gir',
    birth_date: '2019-05-15',
    farm_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
]

export const mockReproductionData = [
  {
    id: 1,
    animal_id: 1,
    animal_name: 'Vaca 001',
    ear_tag: '001',
    current_phase: 1,
    insemination_date: '2024-01-01',
    pregnancy_date: '2024-01-15',
    expected_birth_date: '2024-10-15',
    veterinary_confirmation: true,
    farm_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    animal_id: 2,
    animal_name: 'Vaca 002',
    ear_tag: '002',
    current_phase: 2,
    insemination_date: '2024-02-01',
    pregnancy_date: '2024-02-15',
    expected_birth_date: '2024-11-15',
    veterinary_confirmation: false,
    farm_id: 1,
    created_at: '2024-02-01T00:00:00Z',
    updated_at: '2024-02-01T00:00:00Z',
  },
]

export const mockMilkProductionData = [
  {
    id: 1,
    animal_id: 1,
    animal_name: 'Vaca 001',
    ear_tag: '001',
    date: '2024-01-01',
    morning_quantity: 15.5,
    afternoon_quantity: 12.3,
    total_quantity: 27.8,
    farm_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    animal_id: 2,
    animal_name: 'Vaca 002',
    ear_tag: '002',
    date: '2024-01-01',
    morning_quantity: 18.2,
    afternoon_quantity: 14.7,
    total_quantity: 32.9,
    farm_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
]

export const mockFarmData = {
  id: 1,
  name: 'Fazenda Teste',
  location: 'Rua Teste, 123',
  created_at: '2021-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

export const mockUserData = {
  id: 1,
  name: 'Usuário Teste',
  email: 'usuario@teste.com',
  farm_id: 1,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

export const mockApiHooks = {
  useAnimals: () => ({
    animals: mockAnimalsData,
    loading: false,
    error: null,
    getAnimals: vi.fn().mockResolvedValue(mockAnimalsData),
    createAnimal: vi.fn().mockResolvedValue(mockAnimalsData[0]),
    updateAnimal: vi.fn().mockResolvedValue(mockAnimalsData[0]),
    deleteAnimal: vi.fn().mockResolvedValue(true),
  }),
  
  useReproduction: () => ({
    reproductions: mockReproductionData,
    loading: false,
    error: null,
    getReproductionsByFarm: vi.fn().mockResolvedValue(mockReproductionData),
    createReproduction: vi.fn().mockResolvedValue(mockReproductionData[0]),
    updateReproduction: vi.fn().mockResolvedValue(true),
    deleteReproduction: vi.fn().mockResolvedValue(true),
  }),
  
  useMilkProduction: () => ({
    milkProductions: mockMilkProductionData,
    loading: false,
    error: null,
    getMilkProductionsByFarm: vi.fn().mockResolvedValue(mockMilkProductionData),
    createMilkProduction: vi.fn().mockResolvedValue(mockMilkProductionData[0]),
    updateMilkProduction: vi.fn().mockResolvedValue(true),
    deleteMilkProduction: vi.fn().mockResolvedValue(true),
  }),
  
  useFarm: () => ({
    farm: mockFarmData,
    loading: false,
    error: null,
    getFarm: vi.fn().mockResolvedValue(mockFarmData),
    updateFarm: vi.fn().mockResolvedValue(mockFarmData),
  }),
  
  useAuth: () => ({
    user: mockUserData,
    loading: false,
    error: null,
    login: vi.fn().mockResolvedValue(mockUserData),
    logout: vi.fn().mockResolvedValue(true),
    register: vi.fn().mockResolvedValue(mockUserData),
  }),
}
