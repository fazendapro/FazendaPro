import React from 'react';
import { Table } from 'antd';
import { useTranslation } from 'react-i18next';

interface MilkData {
  key: string;
  name: string;
  production: number;
  fat: number;
  price: string;
}

const MilkProduction: React.FC = () => {
  const { t } = useTranslation();
  const columns = [
    { title: t('dashboard.name'), dataIndex: 'name', key: 'name' },
    { title: t('dashboard.production'), dataIndex: 'production', key: 'production' },
    { title: t('dashboard.fat'), dataIndex: 'fat', key: 'fat' },
    { title: t('dashboard.price'), dataIndex: 'price', key: 'price' },
  ];

  const data: MilkData[] = [
    { key: '1', name: 'Rin of Excel', production: 30, fat: 12, price: 'R$ 100' },
    { key: '2', name: 'Rin', production: 21, fat: 15, price: 'R$ 207' },
  ];

  return <Table columns={columns} dataSource={data} />;
};

export {MilkProduction};