import { forwardRef, useImperativeHandle, useState } from 'react';
import { Table, Spin, Alert } from 'antd';
import { useAnimals } from '../../hooks/useAnimals';
import { useFarm } from '../../../../../hooks/useFarm';
import { useAnimalColumnBuilder } from './column-builder.tsx';
import { useResponsive } from '../../../../../hooks';
import { CustomPagination } from '../../../../../components/lib/Pagination/custom-pagination';

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
  const { isMobile, isTablet } = useResponsive();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

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

  const handlePageChange = (page: number, size: number) => {
    setCurrentPage(page);
    setPageSize(size);
  };

  const handleShowSizeChange = (_: number, size: number) => {
    setCurrentPage(1);
    setPageSize(size);
  };

  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = startIndex + pageSize;
  const paginatedData = filteredAnimals.slice(startIndex, endIndex);

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
    <>
      <Table 
        columns={filteredColumns} 
        dataSource={paginatedData} 
        pagination={false}
        rowKey="id"
        scroll={{ 
          x: isMobile ? 800 : isTablet ? 1200 : 1500,
          y: isMobile ? 400 : undefined
        }}
        size={isMobile ? 'small' : 'middle'}
        style={{
          fontSize: isMobile ? '12px' : '14px'
        }}
      />
      
      <CustomPagination
        current={currentPage}
        total={filteredAnimals.length}
        pageSize={pageSize}
        onChange={handlePageChange}
        onShowSizeChange={handleShowSizeChange}
        showSizeChanger={!isMobile}
        showTotal={!isMobile}
      />
    </>
  );
});

export { AnimalTable };