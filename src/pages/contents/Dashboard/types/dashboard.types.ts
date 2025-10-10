export interface NextToCalveAnimal {
  id: number;
  animal_name: string;
  ear_tag_number_local: number;
  photo?: string;
  pregnancy_date: string;
  expected_birth_date: string;
  days_until_birth: number;
  status: 'Alto' | 'Médio' | 'Baixo';
}

export interface GetNextToCalveParams {
  farm_id: number;
}

export interface GetNextToCalveResponse {
  data: NextToCalveAnimal[];
  success: boolean;
  message: string;
  status: number;
}
