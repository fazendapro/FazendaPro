import React, { useEffect } from 'react';
import { Card, Button, Spin, Alert, Typography, Row, Col, Image } from 'antd';
import { toast } from 'react-toastify';
import { useFarmSelection } from './hooks/useFarmSelection';
import { useResponsive } from '../../hooks';

const { Title, Text } = Typography;

export const FarmSelection: React.FC = () => {
  const { farms, loading, error, autoSelect, selectFarm } = useFarmSelection();
  const { isMobile } = useResponsive();

  useEffect(() => {
    if (error && !loading) {
      if (farms.length > 0) {
        toast.error(error, {
          toastId: 'farm-selection-error-toast',
          autoClose: 5000,
        });
      }
    }
  }, [error, loading, farms.length]);

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: '100vh',
        flexDirection: 'column',
        gap: '16px'
      }}>
        <Spin size="large" />
        <Text>Carregando fazendas...</Text>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: '100vh',
        padding: '20px'
      }}>
        <Alert
          message="Erro"
          description={error}
          type="error"
          showIcon
          style={{ maxWidth: '500px' }}
        />
      </div>
    );
  }

  if (autoSelect) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: '100vh',
        flexDirection: 'column',
        gap: '16px'
      }}>
        <Spin size="large" />
        <Text>Redirecionando para sua fazenda...</Text>
      </div>
    );
  }

  return (
    <div style={{ 
      minHeight: '100vh', 
      padding: isMobile ? '20px' : '50px',
      backgroundColor: '#f5f5f5'
    }}>
      <div style={{ 
        maxWidth: '1200px', 
        margin: '0 auto',
        textAlign: 'center'
      }}>
        <Title level={2} style={{ marginBottom: '32px' }}>
          Selecione uma Fazenda
        </Title>
        
        <Text type="secondary" style={{ marginBottom: '32px', display: 'block' }}>
          Escolha a fazenda que você deseja acessar
        </Text>

        <Row gutter={[24, 24]} justify="center">
          {farms.map((farm) => (
            <Col key={farm.ID} xs={24} sm={12} md={8} lg={6}>
              <Card
                hoverable
                style={{ 
                  height: '100%',
                  textAlign: 'center',
                  borderRadius: '12px',
                  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)'
                }}
                styles={{ body: { padding: '24px' } }}
                onClick={() => selectFarm(farm.ID)}
              >
                <div style={{ marginBottom: '16px' }}>
                  {farm.Logo ? (
                    <Image
                      src={farm.Logo}
                      alt={`Fazenda ${farm.ID}`}
                      style={{ 
                        width: '80px', 
                        height: '80px', 
                        objectFit: 'cover',
                        borderRadius: '50%'
                      }}
                      fallback="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAMIAAADDCAYAAADQvc6UAAABRWlDQ1BJQ0MgUHJvZmlsZQAAKJFjYGASSSwoyGFhYGDIzSspCnJ3UoiIjFJgf8LAwSDCIMogwMCcmFxc4BgQ4ANUwgCjUcG3awyMIPqyLsis7PPOq3QdDFcvjV3QOD1XvZ1YLJHAV0pqcTKQ/gP6mBJFRuYJpPICQnYrlJ4OYCoA4B7A6hcg6y7QO8BzoIOA7DkgDqIFzIPiDxCXsCVAsgDEDlDpYHkIG2+APDpFtC4GXDY7HAc4OQj4NnYZALp+UsgHAt6uA7p8gJgSagP4J5ecllYGHTiAKDO472Q4HA7wSDgUJBYlAh3AOM3luI0YyMIm3oGBhYQDAKM/AsW1SAFC1AICQw4Lg2PAbxUwkTicIoaZgf2+BxwX8B2WOBdaFXwC8sTg4MChQLE3EH4F0k5ikVzI4lLr+gf4rM8Tg/f78JwV3fKvHoBC1HOAEAAP//AwD3fwDi9vsPAAAAAElFTkSuQmCC"
                    />
                  ) : (
                    <div style={{
                      width: '80px',
                      height: '80px',
                      backgroundColor: '#f0f0f0',
                      borderRadius: '50%',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      margin: '0 auto',
                      fontSize: '24px',
                      color: '#999'
                    }}>
                      🏡
                    </div>
                  )}
                </div>
                
                <Title level={4} style={{ marginBottom: '8px' }}>
                  {farm.Company?.CompanyName}
                </Title>
                
                <Button 
                  type="primary" 
                  size="large"
                  style={{ 
                    width: '100%',
                    borderRadius: '8px',
                    height: '40px'
                  }}
                >
                  Acessar Fazenda
                </Button>
              </Card>
            </Col>
          ))}
        </Row>

        {farms.length === 0 && (
          <Card style={{ marginTop: '32px', textAlign: 'center' }}>
            <Title level={4}>Nenhuma fazenda encontrada</Title>
            <Text type="secondary">
              Entre em contato com o administrador para ter acesso a uma fazenda.
            </Text>
          </Card>
        )}
      </div>
    </div>
  );
};
