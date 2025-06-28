import React from 'react';
import { Card, List, Avatar } from 'antd';
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
    { name: 'Tata Salt', days: 15, image: 'src/assets/images/mocked/cows/tata.png', status: 'Alto' },
    { name: 'Lays', days: 40, image: 'src/assets/images/mocked/cows/lays.png', status: 'Baixo' },
    { name: 'Matilda', days: 20, image: 'src/assets/images/mocked/cows/matilda.png', status: 'Médio' },
  ];

  const getStatusStyles = (status: string) => {
    switch (status) {
      case 'Alto':
        return {
          color: 'red',
          border: '1px solid red',
          backgroundColor: '#ffcccc',
          fontSize: '12px'
        };
      case 'Médio':
        return {
          color: '#ff8c00',
          border: '1px solid #ff8c00',
          backgroundColor: '#ffe4b5',
          fontSize: '12px'
        };
      default:
        return {
           color: 'green',
          border: '1px solid green',
          backgroundColor: '#ccffcc',
          fontSize: '12px'
        };
    }
  };

  return (
    <Card
      title={t('dashboard.nextToCalve')}
    >
      <List
        itemLayout="horizontal"
        dataSource={data}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              avatar={<Avatar src={item.image} shape="square" size={100} />}
              title={<span style={{ fontSize: 16, fontWeight: 'bold' }}>{item.name}</span>}
              description={
                <>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <p>Última vez: {item.days} dias</p>
                    <span style={{ 
                      ...getStatusStyles(item.status),
                      padding: '4px', 
                      borderRadius: '20px' 
                    }}>
                      {item.status}
                    </span>
                  </div>
                </>
              }
            />
          </List.Item>
        )}
      />
    </Card>
  );
};

export { NextToCalve };