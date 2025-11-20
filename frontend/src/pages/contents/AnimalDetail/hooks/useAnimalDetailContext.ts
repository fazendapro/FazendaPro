import { useContext } from 'react';
import { AnimalDetailContext } from './AnimalDetailProvider';
import { AnimalDetail, AnimalDetailFormData, AnimalParent } from '../types';

interface AnimalDetailContextType {
  animal: AnimalDetail | undefined;
  loading: boolean;
  error: string | undefined;
  isEditing: boolean;
  setIsEditing: (editing: boolean) => void;
  updateAnimal: (data: AnimalDetailFormData) => Promise<void>;
  uploadPhoto: (file: File) => Promise<string | null>;
  fathers: AnimalParent[];
  mothers: AnimalParent[];
  loadingParents: boolean;
  refreshAnimal: () => Promise<void>;
}

export const useAnimalDetailContext = (): AnimalDetailContextType => {
  const context = useContext(AnimalDetailContext);
  if (!context) {
    throw new Error('useAnimalDetailContext must be used within an AnimalDetailProvider');
  }
  return context;
};
