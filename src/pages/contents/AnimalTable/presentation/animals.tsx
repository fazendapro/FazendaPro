import React, { useRef, useState } from 'react';
import { AnimalDashboard } from "./animal-dashboard";
import { AnimalTable } from "./animal-table";
import { useAnimalColumnBuilder } from "./column-builder";

const Animals = () => {
  const tableRef = useRef<{ refetch: () => void }>(null);
  const { getDefaultColumnKeys } = useAnimalColumnBuilder();

  const [selectedColumns, setSelectedColumns] = useState<string[]>(getDefaultColumnKeys());

  const handleAnimalCreated = () => {
    if (tableRef.current) {
      tableRef.current.refetch();
    }
  };

  const handleColumnsChanged = (columns: string[]) => {
    setSelectedColumns(columns);
  };

  return (
    <div>
      <AnimalDashboard 
        onAnimalCreated={handleAnimalCreated}
        onColumnsChanged={handleColumnsChanged}
        selectedColumns={selectedColumns}
      />
      <AnimalTable 
        ref={tableRef} 
        selectedColumns={selectedColumns}
      />
    </div>
  );
};

export { Animals };