import React, { createContext, ReactNode } from 'react';
import { useAnimalDetail } from './useAnimalDetail';
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

export const AnimalDetailContext = createContext<AnimalDetailContextType | undefined>(undefined);

interface AnimalDetailProviderProps {
  children: ReactNode;
  animalId: number;
}

export const AnimalDetailProvider: React.FC<AnimalDetailProviderProps> = ({ 
  children, 
  animalId 
}) => {
  const animalDetailData = useAnimalDetail(animalId);

  return (
    <AnimalDetailContext.Provider value={animalDetailData}>
      {children}
    </AnimalDetailContext.Provider>
  );
};

