import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';

interface Farm {
  ID: number;
  CompanyID: number;
  Logo: string;
  Language?: string;
  Company?: {
    ID: number;
    CompanyName: string;
    Location: string;
    FarmCNPJ: string;
  };
}

interface FarmContextType {
  selectedFarm: Farm | null;
  setSelectedFarm: (farm: Farm | null) => void;
  clearSelectedFarm: () => void;
}

const FarmContext = createContext<FarmContextType | undefined>(undefined);

interface FarmProviderProps {
  children: ReactNode;
}

export const FarmProvider: React.FC<FarmProviderProps> = ({ children }) => {
  const [selectedFarm, setSelectedFarm] = useState<Farm | null>(null);

  useEffect(() => {
    const savedFarm = localStorage.getItem('selectedFarm');
    if (savedFarm) {
      try {
        setSelectedFarm(JSON.parse(savedFarm));
      } catch {
        localStorage.removeItem('selectedFarm');
      }
    }
  }, []);
  const handleSetSelectedFarm = useCallback((farm: Farm | null) => {
    setSelectedFarm(farm);
    if (farm) {
      localStorage.setItem('selectedFarm', JSON.stringify(farm));
    } else {
      localStorage.removeItem('selectedFarm');
    }
  }, []);

  const clearSelectedFarm = useCallback(() => {
    setSelectedFarm(null);
    localStorage.removeItem('selectedFarm');
  }, []);

  return (
    <FarmContext.Provider
      value={{
        selectedFarm,
        setSelectedFarm: handleSetSelectedFarm,
        clearSelectedFarm,
      }}
    >
      {children}
    </FarmContext.Provider>
  );
};

export const useFarm = (): FarmContextType => {
  const context = useContext(FarmContext);
  if (context === undefined) {
    throw new Error('useFarm deve ser usado dentro de um FarmProvider');
  }
  return context;
};
