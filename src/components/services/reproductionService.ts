import { api } from '../../config/api';
import { Reproduction } from '../../types/reproduction';
import i18n from '../../locale/i18n';
import { ReproductionPhase } from '../../pages/contents/AnimalTable/domain/model/reproduction';

interface ReproductionResponse {
  success: boolean;
  data: any;
  message?: string;
}

const getPhaseTranslation = (phaseNumber: number): string => {
  const phaseMap: { [key: number]: string } = {
    [ReproductionPhase.LACTACAO]: 'animalTable.reproduction.phases.lactacao',
    [ReproductionPhase.SECANDO]: 'animalTable.reproduction.phases.secando',
    [ReproductionPhase.VAZIAS]: 'animalTable.reproduction.phases.vazias',
    [ReproductionPhase.PRENHAS]: 'animalTable.reproduction.phases.prenhas'
  };

  const translationKey = phaseMap[phaseNumber];
  if (translationKey) {
    return i18n.t(translationKey);
  }

  return i18n.t('common.notInformed', { defaultValue: 'Não informado' });
};

export const reproductionService = {
  async getReproductionsByAnimal(animalId: number): Promise<Reproduction[]> {
    try {
      const response = await api.get<ReproductionResponse>(`/reproductions/animal`, {
        params: { animalId }
      });
      
      if (response.data.success && response.data.data) {
        const data = response.data.data;
        
        const reproductions = Array.isArray(data) ? data : [data];
        
        return reproductions.map((rep: any) => {
          const date = rep.pregnancy_date || rep.insemination_date || rep.date || rep.Date;
          
          const phaseNumber = rep.current_phase !== undefined ? rep.current_phase : (rep.phase || 0);
          const phase = typeof phaseNumber === 'number' 
            ? getPhaseTranslation(phaseNumber)
            : (rep.phase || rep.Phase || i18n.t('common.notInformed', { defaultValue: 'Não informado' }));
          
          return {
            id: rep.id || rep.ID,
            animal_id: rep.animal_id || rep.AnimalID,
            farm_id: rep.farm_id || rep.FarmID,
            date: date,
            phase: phase,
            notes: rep.notes || rep.Notes || rep.observations || rep.Observations,
            created_at: rep.created_at || rep.createdAt || rep.CreatedAt,
            updated_at: rep.updated_at || rep.updatedAt || rep.UpdatedAt,
            animal: rep.animal || rep.Animal || (rep.animal_name ? { 
              id: rep.animal_id || rep.AnimalID,
              animal_name: rep.animal_name || rep.AnimalName,
              ear_tag_number_local: rep.ear_tag || rep.EarTag
            } : undefined)
          };
        });
      }
      
      return [];
    } catch (error: any) {
      if (error.response?.status === 404) {
        return [];
      }
      if (error.response?.status !== 404) {
        console.error('Error fetching reproductions for animal:', error);
      }
      return [];
    }
  }
};
