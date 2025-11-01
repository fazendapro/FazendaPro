import React from 'react';
import { Card, Row, Col, Spin } from 'antd';
import { DollarOutlined, RiseOutlined, ManOutlined, WomanOutlined } from '@ant-design/icons';
import { DashboardInfoCard } from '../../../../components';
import { useTranslation } from 'react-i18next';
import { useOverviewStats } from '../hooks/useOverviewStats';

const Overview: React.FC = () => {
  const { t } = useTranslation();
  const { stats, loading, error } = useOverviewStats();

  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  if (loading) {
    return (
      <Card title={t('dashboard.salesOverview')} style={{ marginBottom: 16, borderRadius: 8 }}>
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card title={t('dashboard.salesOverview')} style={{ marginBottom: 16, borderRadius: 8 }}>
        <div style={{ textAlign: 'center', padding: '20px', color: 'red' }}>
          {error}
        </div>
      </Card>
    );
  }

  return (
    <Card title={t('dashboard.salesOverview')} style={{ marginBottom: 16, borderRadius: 8 }}>
      <Row>
        <Col span={6}>
          <DashboardInfoCard 
            title={t('dashboard.machos')} 
            value={stats.males_count.toString()} 
            icon={<ManOutlined style={{ fontSize: 24, color: '#faad14' }} />} 
            isLast={false} 
          />
        </Col>
        <Col span={6}>
          <DashboardInfoCard 
            title={t('dashboard.femeas')} 
            value={stats.females_count.toString()} 
            icon={<WomanOutlined style={{ fontSize: 24, color: '#d3adf7' }} />} 
            isLast={false} 
          />
        </Col>
        <Col span={6}>
          <DashboardInfoCard 
            title={t('dashboard.sales')} 
            value={stats.total_sold.toString()} 
            icon={<DollarOutlined style={{ fontSize: 24, color: '#1890ff' }} />} 
            isLast={false} 
          />
        </Col>
        <Col span={6}>
          <DashboardInfoCard 
            title={t('dashboard.revenue')} 
            value={formatCurrency(stats.total_revenue)} 
            icon={<RiseOutlined style={{ fontSize: 24, color: '#52c41a' }} />} 
            isLast={true} 
          />
        </Col>
      </Row>
    </Card>
  );
};

export { Overview };