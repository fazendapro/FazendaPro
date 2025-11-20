export interface Animal {
  id: number;
  farm_id: number;
  ear_tag_number_local: number;
  ear_tag_number_register?: number;
  animal_name: string;
  sex: number;
  breed: string;
  type: string;
  birth_date?: string;
  photo?: string;
  father_id?: number;
  father?: Animal;
  mother_id?: number;
  mother?: Animal;
  confinement: boolean;
  animal_type: number;
  status: number;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  created_at: string;
  updated_at: string;
}
