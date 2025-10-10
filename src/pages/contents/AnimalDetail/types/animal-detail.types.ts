export interface AnimalDetail {
  id: number;
  farm_id: number;
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: string;
  sex: number;
  breed: string;
  birth_date?: string;
  photo?: string;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  father_id?: number;
  mother_id?: number;
  father?: AnimalDetail;
  mother?: AnimalDetail;
  created_at: string;
  updated_at: string;
}

export interface AnimalDetailFormData {
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: string;
  sex: number;
  breed: string;
  birth_date?: string;
  photo?: string;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  father_id?: number;
  mother_id?: number;
}

export interface AnimalDetailResponse {
  success: boolean;
  data?: AnimalDetail;
  message?: string;
}

export interface AnimalDetailUpdateRequest {
  id: number;
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: string;
  sex: number;
  breed: string;
  birth_date?: string;
  photo?: string;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  father_id?: number;
  mother_id?: number;
}

export interface AnimalParent {
  id: number;
  animal_name: string;
  ear_tag_number_local: number;
  sex: number;
}

export interface AnimalParentsResponse {
  success: boolean;
  data?: AnimalParent[];
  message?: string;
}

export const SEX_OPTIONS = [
  { value: 0, label: 'Fêmea' },
  { value: 1, label: 'Macho' }
];

export const ANIMAL_TYPE_OPTIONS = [
  { value: 0, label: 'Bovino' },
  { value: 1, label: 'Suíno' },
  { value: 2, label: 'Ovino' },
  { value: 3, label: 'Caprino' },
  { value: 4, label: 'Equino' },
  { value: 5, label: 'Ave' },
  { value: 6, label: 'Peixe' },
  { value: 7, label: 'Outro' }
];

export const STATUS_OPTIONS = [
  { value: 0, label: 'Vivo' },
  { value: 1, label: 'Morto' }
];

export const PURPOSE_OPTIONS = [
  { value: 0, label: 'Carne' },
  { value: 1, label: 'Leite' },
  { value: 2, label: 'Reprodução' }
];
