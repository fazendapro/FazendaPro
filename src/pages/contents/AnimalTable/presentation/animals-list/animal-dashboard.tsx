import React, { useState } from 'react';
import { Button, Input, Card, Row, Col, Statistic, Space } from 'antd';
import { useModal } from '../../../../../hooks';
import { useTranslation } from 'react-i18next';
import { CreateAnimalModal } from './create-animal-modal';
import { FilterModal } from './filter-modal';

const { Search } = Input;

interface AnimalDashboardProps {
  onAnimalCreated?: () => void;
  onColumnsChanged?: (columns: string[]) => void;
  onSearchChange?: (searchTerm: string) => void;
  selectedColumns?: string[];
}

const AnimalDashboard: React.FC<AnimalDashboardProps> = ({ 
  onAnimalCreated, 
  onColumnsChanged, 
  onSearchChange,
  selectedColumns: externalSelectedColumns 
}) => {
  const { isOpen, onOpen, onClose } = useModal();
  const { isOpen: isFilterOpen, onOpen: onFilterOpen, onClose: onFilterClose } = useModal();
  const { t } = useTranslation();

  const [selectedColumns, setSelectedColumns] = useState<string[]>(
    externalSelectedColumns || [
      'animal_name',
      'ear_tag_number_local', 
      'ear_tag_number_register',
      'type',
      'sex',
      'breed',
      'birth_date'
    ]
  );

  const handleAnimalCreated = () => {
    onClose();
    if (onAnimalCreated) {
      onAnimalCreated();
    }
  };

  const handleApplyFilters = (columns: string[]) => {
    setSelectedColumns(columns);
    if (onColumnsChanged) {
      onColumnsChanged(columns);
    }
  };

  const handleSearch = (value: string) => {
    if (onSearchChange) {
      onSearchChange(value);
    }
  };

  return (
    <div>
      <Card style={{ marginBottom: '16px' }}>
        <Row gutter={16}>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.categories')} 
              value={14} 
              valueStyle={{ color: '#1890ff' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.totalAnimals')} 
              value={868} 
              valueStyle={{ color: '#faad14' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.revenue')} 
              value={25000} 
              prefix="R$"
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.bestSales')} 
              value={5} 
              valueStyle={{ color: '#722ed1' }}
            />
          </Col>
        </Row>
        <Row gutter={16} style={{ marginTop: '16px' }}>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.cost')} 
              value={2500} 
              prefix="R$"
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.lowProduction')} 
              value={12} 
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.inseminated')} 
              value={2} 
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.notInseminated')} 
              value={2} 
              valueStyle={{ color: '#000' }}
            />
          </Col>
        </Row>
      </Card>

      <Card style={{ marginBottom: '16px' }}>
        <Space>
          <Button type="primary" onClick={onOpen}>
            {t('animalTable.createCow')}
          </Button>
          <Button onClick={onFilterOpen}>
            {t('animalTable.filter')}
          </Button>
          <Search 
            placeholder={t('animalTable.search')} 
            style={{ width: 'auto' }}
            onSearch={handleSearch}
            allowClear
          />
        </Space>
      </Card>

      <CreateAnimalModal isOpen={isOpen} onClose={handleAnimalCreated} />
      <FilterModal 
        isOpen={isFilterOpen} 
        onClose={onFilterClose}
        onApplyFilters={handleApplyFilters}
        currentColumns={selectedColumns}
      />
    </div>
  );
};

export { AnimalDashboard };