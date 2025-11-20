export enum AnimalSex {
  MALE = 0,
  FEMALE = 1
}

export interface AnimalForm {
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: 'vaca' | 'bezerro' | 'touro' | 'novilho';
  sex: AnimalSex;
  breed: string;
  birth_date: string;
}

export interface CreateAnimalParams extends AnimalForm {
  farm_id: number;
}

export interface Animal extends AnimalForm {
  id: string;
  farm_id: number;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  createdAt: string;
  updatedAt: string;
  price?: string;
  current_weight?: string;
  ideal_weight?: string;
  last_update?: string;
  milk_production?: string;
  photo?: string;
}

export interface GetAnimalsByFarmParams {
  farm_id: number;
}

