import { useRef, useState } from 'react';
import { AnimalDashboard } from "./animal-dashboard";
import { AnimalTable } from "./animal-table";
import { useAnimalColumnBuilder } from "./column-builder";

const Animals = () => {
  const tableRef = useRef<{ refetch: () => void }>(null);
  const { getDefaultColumnKeys } = useAnimalColumnBuilder();

  const [selectedColumns, setSelectedColumns] = useState<string[]>(getDefaultColumnKeys());
  const [searchTerm, setSearchTerm] = useState<string>('');

  const handleAnimalCreated = () => {
    if (tableRef.current) {
      tableRef.current.refetch();
    }
  };

  const handleColumnsChanged = (columns: string[]) => {
    setSelectedColumns(columns);
  };

  const handleSearchChange = (searchTerm: string) => {
    setSearchTerm(searchTerm);
  };

  return (
    <div id="animals-list">
      <AnimalDashboard
        onAnimalCreated={handleAnimalCreated}
        onColumnsChanged={handleColumnsChanged}
        onSearchChange={handleSearchChange}
        selectedColumns={selectedColumns}
      />
      <AnimalTable
        ref={tableRef}
        selectedColumns={selectedColumns}
        searchTerm={searchTerm}
      />
    </div>
  );
};

export { Animals };