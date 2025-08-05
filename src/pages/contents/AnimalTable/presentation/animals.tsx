import React, { useRef } from 'react';
import { AnimalDashboard } from "./animal-dashboard";
import { AnimalTable } from "./animal-table";

const Animals = () => {
  const tableRef = useRef<{ refetch: () => void }>(null);

  const handleAnimalCreated = () => {
    if (tableRef.current) {
      tableRef.current.refetch();
    }
  };

  return (
    <div>
      <AnimalDashboard onAnimalCreated={handleAnimalCreated} />
      <AnimalTable ref={tableRef} />
    </div>
  );
};

export { Animals };