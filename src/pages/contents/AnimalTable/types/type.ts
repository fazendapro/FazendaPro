export enum AnimalSex {
  MALE = 0,
  FEMALE = 1
}

export interface AnimalForm {
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_global: number;
  type: 'vaca' | 'bezerro' | 'touro' | 'novilho';
  sex: AnimalSex;
  breed: string;
  birth_date: string;
}

export interface CreateAnimalParams extends AnimalForm {
  farm_id: number;
}

