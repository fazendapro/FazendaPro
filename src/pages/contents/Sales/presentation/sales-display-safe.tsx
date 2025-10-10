import React from 'react';
import { Card, Alert, Button, Spin } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { useSaleList } from '../../../../hooks/useSale';
import { SalesDisplay } from './sales-display';

export const SalesDisplaySafe: React.FC = () => {
  const { loading, error, clearError } = useSaleList();

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '50vh' 
      }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '24px' }}>
        <Card>
          <Alert
            message="Erro ao carregar vendas"
            description={error}
            type="error"
            showIcon
            action={
              <Button 
                size="small" 
                danger 
                icon={<ReloadOutlined />}
                onClick={() => {
                  clearError();
                  window.location.reload();
                }}
              >
                Tentar novamente
              </Button>
            }
          />
        </Card>
      </div>
    );
  }

  return <SalesDisplay />;
};
