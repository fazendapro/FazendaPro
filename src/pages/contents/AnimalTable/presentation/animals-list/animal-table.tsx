import { forwardRef, useImperativeHandle } from 'react';
import { Table, Spin, Alert } from 'antd';
import { useAnimals } from '../../hooks/useAnimals';
import { useFarm } from '../../../../../hooks/useFarm';
import { useAnimalColumnBuilder } from './column-builder';

interface AnimalTableRef {
  refetch: () => void;
}

interface AnimalTableProps {
  selectedColumns?: string[];
  searchTerm?: string;
}

const AnimalTable = forwardRef<AnimalTableRef, AnimalTableProps>((props, ref) => {
  const { selectedColumns = [], searchTerm = '' } = props;
  const { farm } = useFarm();
  const { animals, loading, error, refetch } = useAnimals(farm?.id);
  const { buildTableColumns } = useAnimalColumnBuilder();

  useImperativeHandle(ref, () => ({
    refetch
  }));

  const filteredColumns = buildTableColumns(selectedColumns);

  const filteredAnimals = animals.filter(animal => {
    if (!searchTerm.trim()) return true;
    
    const searchLower = searchTerm.toLowerCase();
    return (
      animal.animal_name?.toLowerCase().includes(searchLower) ||
      animal.ear_tag_number_local?.toString().includes(searchLower) ||
      animal.ear_tag_number_register?.toString().includes(searchLower) ||
      animal.breed?.toLowerCase().includes(searchLower) ||
      animal.type?.toLowerCase().includes(searchLower)
    );
  });

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert
        message="Erro"
        description={error}
        type="error"
        showIcon
        style={{ marginBottom: '16px' }}
      />
    );
  }

  return (
    <Table 
      columns={filteredColumns} 
      dataSource={filteredAnimals} 
      pagination={{ showSizeChanger: true }}
      rowKey="id"
      scroll={{ x: 1500 }}
    />
  );
});

export { AnimalTable };