import React from 'react';
import { Card, List, Avatar, Spin, Tooltip } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useNextToCalve } from '../hooks/useNextToCalve';
import { NextToCalveAnimal } from '../types/dashboard.types';

const NextToCalve: React.FC = () => {
  const { t } = useTranslation();
  const { nextToCalve, loading, error } = useNextToCalve();

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

  if (loading) {
    return (
      <Card title={t('dashboard.nextToCalve')}>
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card title={t('dashboard.nextToCalve')}>
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
          {t('dashboard.nextToCalve')}
          <Tooltip 
            title="Uma vaca fica grávida por aproximadamente 9 meses, ou cerca de 280 a 290 dias, a partir do dia em que a prenhez é confirmada. O período de gestação pode variar ligeiramente dependendo da raça, da saúde do animal e de outros fatores, mas a média é de 283 dias."
            placement="top"
          >
            <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
          </Tooltip>
        </div>
      }
    >
      <List
        itemLayout="horizontal"
        dataSource={nextToCalve}
        renderItem={(item: NextToCalveAnimal) => (
          <List.Item
            style={{
              padding: '12px 0',
              borderBottom: '1px solid #f0f0f0'
            }}
          >
            <List.Item.Meta
              avatar={
                <Avatar 
                  src={item.photo} 
                  shape="square" 
                  size={{ xs: 60, sm: 80, md: 100 }}
                  style={{ 
                    minWidth: '60px',
                    minHeight: '60px'
                  }}
                />
              }
              title={
                <div style={{ 
                  display: 'flex', 
                  flexDirection: 'column',
                  gap: '4px'
                }}>
                  <span style={{ 
                    fontSize: '16px', 
                    fontWeight: 'bold',
                    color: '#262626'
                  }}>
                    {item.animal_name}
                  </span>
                  <span style={{ 
                    fontSize: '12px', 
                    color: '#8c8c8c'
                  }}>
                    #{item.ear_tag_number_local}
                  </span>
                </div>
              }
              description={
                <div style={{ 
                  display: 'flex', 
                  flexDirection: 'column',
                  gap: '8px',
                  marginTop: '8px'
                }}>
                  <div style={{ 
                    display: 'flex', 
                    alignItems: 'center', 
                    gap: '8px',
                    flexWrap: 'wrap'
                  }}>
                    <span style={{ 
                      fontSize: '14px',
                      color: '#595959'
                    }}>
                      Faltam: <strong>{item.days_until_birth} dias</strong>
                    </span>
                    <span style={{ 
                      ...getStatusStyles(item.status),
                      padding: '4px 8px', 
                      borderRadius: '12px',
                      fontSize: '12px',
                      fontWeight: 'bold'
                    }}>
                      {item.status}
                    </span>
                  </div>
                  <div style={{ 
                    fontSize: '12px',
                    color: '#8c8c8c'
                  }}>
                    Data prevista: {item.expected_birth_date}
                  </div>
                </div>
              }
            />
          </List.Item>
        )}
      />
    </Card>
  );
};

export { NextToCalve };