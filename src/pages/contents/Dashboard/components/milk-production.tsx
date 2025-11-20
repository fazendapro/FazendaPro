import React, { useMemo } from 'react';
import { Table, Card, Spin, Tooltip } from 'antd';
import type { Breakpoint } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useTopMilkProducers } from '../hooks/useTopMilkProducers';

const DashboardMilkProduction: React.FC = () => {
  const { t } = useTranslation();
  const { topProducers, loading, error } = useTopMilkProducers(undefined, 10, 30);

  const columns = useMemo(() => [
    { 
      title: t('dashboard.name'), 
      dataIndex: 'animal_name', 
      key: 'animal_name',
      width: '40%',
      responsive: ['xs', 'sm', 'md', 'lg', 'xl'] as Breakpoint[],
      render: (text: string, record: { ear_tag_number_local: number }) => (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <span style={{ fontWeight: 'bold', fontSize: '14px' }}>{text}</span>
          <span style={{ fontSize: '12px', color: '#666' }}>
            #{record.ear_tag_number_local}
          </span>
        </div>
      )
    },
    { 
      title: t('dashboard.production'), 
      dataIndex: 'average_daily_production', 
      key: 'average_daily_production',
      width: '20%',
      responsive: ['sm', 'md', 'lg', 'xl'] as Breakpoint[],
      render: (value: number) => (
        <span style={{ fontWeight: 'bold', color: '#1890ff' }}>
          {value.toFixed(1)}L/dia
        </span>
      )
    },
    { 
      title: t('dashboard.fat'), 
      dataIndex: 'fat_content', 
      key: 'fat_content',
      width: '20%',
      responsive: ['md', 'lg', 'xl'] as Breakpoint[],
      render: (value: number) => (
        <span style={{ color: '#52c41a' }}>
          {value.toFixed(1)}%
        </span>
      )
    },
    { 
      title: t('dashboard.daysInLactation'), 
      dataIndex: 'days_in_lactation', 
      key: 'days_in_lactation',
      width: '20%',
      responsive: ['lg', 'xl'] as Breakpoint[],
      render: (value: number) => (
        <span style={{ color: '#722ed1' }}>
          {value} dias
        </span>
      )
    },
  ], [t]);

  const dataSource = useMemo(() => topProducers.map((producer) => ({
    ...producer,
    key: producer.id
  })), [topProducers]);

  if (loading) {
    return (
      <Card title={t('dashboard.topMilkProducers')} style={{ height: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card title={t('dashboard.topMilkProducers')} style={{ height: '100%' }}>
        <div style={{ textAlign: 'center', padding: '20px', color: 'red' }}>
          {error}
        </div>
      </Card>
    );
  }

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          {t('dashboard.topMilkProducers')}
          <Tooltip 
            title="Ranking das vacas com maior produção de leite nos últimos 30 dias, baseado na produção média diária."
            placement="top"
          >
            <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
          </Tooltip>
        </div>
      }
      style={{ height: '100%' }}
    >
      <Table 
        columns={columns} 
        dataSource={dataSource} 
        pagination={false}
        scroll={{ x: 600 }}
        size="small"
        bordered={false}
        style={{ 
          fontSize: '14px'
        }}
      />
    </Card>
  );
};

export { DashboardMilkProduction };