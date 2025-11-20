export interface Reproduction {
  id: number;
  animal_id: number;
  animal_name?: string;
  ear_tag?: number;
  current_phase: number;
  insemination_date?: string;
  insemination_type?: string;
  pregnancy_date?: string;
  expected_birth_date?: string;
  actual_birth_date?: string;
  lactation_start_date?: string;
  lactation_end_date?: string;
  dry_period_start_date?: string;
  veterinary_confirmation: boolean;
  observations?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateReproductionRequest {
  animal_id: number;
  current_phase: number;
  insemination_date?: string;
  insemination_type?: string;
  pregnancy_date?: string;
  expected_birth_date?: string;
  actual_birth_date?: string;
  lactation_start_date?: string;
  lactation_end_date?: string;
  dry_period_start_date?: string;
  veterinary_confirmation?: boolean;
  observations?: string;
}

export interface UpdateReproductionPhaseRequest {
  animal_id: number;
  new_phase: number;
  additional_data?: {
    insemination_date?: string;
    insemination_type?: string;
    pregnancy_date?: string;
    lactation_start_date?: string;
    lactation_end_date?: string;
    dry_period_start_date?: string;
    actual_birth_date?: string;
    veterinary_confirmation?: boolean;
    observations?: string;
  };
}

export enum ReproductionPhase {
  LACTACAO = 0,
  SECANDO = 1,
  VAZIAS = 2,
  PRENHAS = 3
}

export const ReproductionPhaseLabels = {
  [ReproductionPhase.LACTACAO]: 'Lactação',
  [ReproductionPhase.SECANDO]: 'Secando',
  [ReproductionPhase.VAZIAS]: 'Vazias',
  [ReproductionPhase.PRENHAS]: 'Prenhas'
};

export const ReproductionPhaseColors = {
  [ReproductionPhase.LACTACAO]: '#52c41a',
  [ReproductionPhase.SECANDO]: '#faad14',
  [ReproductionPhase.VAZIAS]: '#1890ff',
  [ReproductionPhase.PRENHAS]: '#f5222d'
};
