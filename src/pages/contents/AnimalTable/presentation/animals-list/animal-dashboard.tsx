import React, { useState } from 'react';
import { Button, Input, Card, Row, Col, Statistic, Space, Collapse } from 'antd';
import { DownOutlined, UpOutlined } from '@ant-design/icons';
import { useModal, useResponsive } from '../../../../../hooks';
import { useTranslation } from 'react-i18next';
import { CreateAnimalModal } from './create-animal-modal';
import { FilterModal } from './filter-modal';
import { CreateWeightModal } from './create-weight-modal';

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
  const { isOpen: isWeightModalOpen, onOpen: onWeightModalOpen, onClose: onWeightModalClose } = useModal();
  const { t } = useTranslation();
  const { isMobile } = useResponsive();
  const [isDashboardExpanded, setIsDashboardExpanded] = useState(false);

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

  const handleWeightCreated = () => {
    onWeightModalClose();
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

  const dashboardContent = (
    <Row gutter={[16, 16]}>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.categories')} 
          value={14} 
          valueStyle={{ color: '#1890ff', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.totalAnimals')} 
          value={868} 
          valueStyle={{ color: '#faad14', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.revenue')} 
          value={25000} 
          prefix="R$"
          valueStyle={{ color: '#000', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.bestSales')} 
          value={5} 
          valueStyle={{ color: '#722ed1', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.cost')} 
          value={2500} 
          prefix="R$"
          valueStyle={{ color: '#000', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.lowProduction')} 
          value={12} 
          valueStyle={{ color: '#ff4d4f', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.inseminated')} 
          value={2} 
          valueStyle={{ color: '#000', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
      <Col xs={12} sm={6} md={6} lg={6} xl={6}>
        <Statistic 
          title={t('animalTable.notInseminated')} 
          value={2} 
          valueStyle={{ color: '#000', fontSize: isMobile ? '16px' : '20px' }}
        />
      </Col>
    </Row>
  );

  const isDeveloped = false;

  return (
    <div>
      {isDeveloped ? (<>
      {isMobile ? (
        <Collapse
          activeKey={isDashboardExpanded ? ['dashboard'] : []}
          onChange={(keys: string | string[]) => setIsDashboardExpanded(Array.isArray(keys) && keys.includes('dashboard'))}
          style={{ marginBottom: '16px' }}
          items={[
            {
              key: 'dashboard',
              label: (
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <span style={{ fontWeight: 'bold' }}>{t('animalTable.dashboard')}</span>
                  {isDashboardExpanded ? <UpOutlined /> : <DownOutlined />}
                </div>
              ),
              children: dashboardContent
            }
          ]}
        />
      ) : (
        <Card style={{ marginBottom: '16px' }}>
          {dashboardContent}
        </Card>
      )}
      </>
    ) : (<></>)}

      <Card style={{ marginBottom: '16px' }}>
        <Space 
          direction={isMobile ? 'vertical' : 'horizontal'} 
          style={{ width: isMobile ? '100%' : 'auto' }}
          size={'middle'}
        >
          <Button 
            type="primary" 
            onClick={onOpen}
            size={'middle'}
            block={isMobile}
          >
            {t('animalTable.createCow')}
          </Button>
          <Button 
            onClick={onFilterOpen}
            size={'middle'}
            block={isMobile}
          >
            {t('animalTable.filter')}
          </Button>
          <Button 
            onClick={onWeightModalOpen}
            size={'middle'}
            block={isMobile}
          >
            {t('animalTable.addWeight') || 'Adicionar Peso'}
          </Button>
          <Search 
            placeholder={t('animalTable.search')} 
            style={{ width: isMobile ? '100%' : 'auto' }}
            onSearch={handleSearch}
            allowClear
            size={'middle'}
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
      <CreateWeightModal
        visible={isWeightModalOpen}
        onCancel={onWeightModalClose}
        onSuccess={handleWeightCreated}
      />
    </div>
  );
};

export { AnimalDashboard };