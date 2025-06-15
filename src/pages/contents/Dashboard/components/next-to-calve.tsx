import React from 'react';
import { Card } from 'antd';
import { useTranslation } from 'react-i18next';

interface CalveData {
  name: string;
  days: number;
  image: string;
  status: string;
}

const NextToCalve: React.FC = () => {
  const { t } = useTranslation();
  const data: CalveData[] = [
    { name: 'Tata Salt', days: 15, image: 'cow1.jpg', status: 'Alto' },
    { name: 'Lays', days: 15, image: 'cow2.jpg', status: 'Alto' },
  ];

  return (
    <Card title={t('dashboard.nextToCalve')} extra={<a href="#">Ver todos</a>}>
      {data.map((item) => (
        <div key={item.name} style={{ display: 'flex', alignItems: 'center', marginBottom: 16 }}>
          <img src={item.image} alt={item.name} style={{ width: 50, marginRight: 16 }} />
          <div>
            <p>{item.name}</p>
            <p>Última vez: {item.days} dias</p>
            <span style={{ color: 'red' }}>{item.status}</span>
          </div>
        </div>
      ))}
    </Card>
  );
};

export {NextToCalve};