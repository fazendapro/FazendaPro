import React from 'react';
import { Button, Input, Select, Card, Row, Col, Statistic, Space } from 'antd';
import { useModal } from '../../../../hooks';
import { useTranslation } from 'react-i18next';
import { CreateAnimalModal } from './create-animal-modal';

const { Search } = Input;
const { Option } = Select;

const AnimalDashboard: React.FC = () => {
  const { isOpen, onOpen, onClose } = useModal();
  const { t } = useTranslation();

  return (
    <div style={{ padding: '16px' }}>
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
          <Select defaultValue="filtro" style={{ width: 120 }}>
            <Option value="filtro">{t('animalTable.filter')}</Option>
          </Select>
          <Search 
            placeholder={t('animalTable.search')} 
            style={{ width: 200 }} 
          />
        </Space>
      </Card>

      <CreateAnimalModal isOpen={isOpen} onClose={onClose} />
    </div>
  );
};

export { AnimalDashboard };