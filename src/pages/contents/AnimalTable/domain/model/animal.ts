export interface Animal {
  id: number;
  farm_id: number;
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: string;
  sex: number;
  breed: string;
  birth_date?: string;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  createdAt: string;
  updatedAt: string;
}
