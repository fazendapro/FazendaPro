import React, { useEffect, useState } from 'react';
import { Card, Alert, Button, Typography, Space } from 'antd';
import { ReloadOutlined, BugOutlined } from '@ant-design/icons';
import { useAuth } from '../../../../contexts/AuthContext';
import { useSelectedFarm } from '../../../../hooks/useSelectedFarm';

const { Title, Text } = Typography;

export const SalesDebug: React.FC = () => {
  const { isAuthenticated, user } = useAuth();
  const { farmName, farmId } = useSelectedFarm();
  const [debugInfo, setDebugInfo] = useState<any>(null);

  const checkAuth = () => {
    const token = localStorage.getItem('token');
    const farmIdFromStorage = localStorage.getItem('selectedFarmId');
    
    setDebugInfo({
      isAuthenticated,
      user: user ? { email: user.email } : null,
      token: token ? `${token.substring(0, 20)}...` : 'Não encontrado',
      farmId,
      farmName,
      farmIdFromStorage,
      timestamp: new Date().toISOString()
    });
  };

  useEffect(() => {
    checkAuth();
  }, []);

  return (
    <div style={{ padding: '24px' }}>
      <Card>
        <Title level={2}>
          <BugOutlined /> Debug - Página de Vendas
        </Title>
        
        <Space direction="vertical" style={{ width: '100%' }}>
          <Alert
            message="Informações de Debug"
            description="Verificando estado da autenticação e fazenda"
            type="info"
            showIcon
          />
          
          <Button 
            type="primary" 
            icon={<ReloadOutlined />}
            onClick={checkAuth}
          >
            Atualizar Debug
          </Button>
          
          {debugInfo && (
            <Card title="Estado Atual" size="small">
              <pre style={{ 
                background: '#f5f5f5', 
                padding: '12px', 
                borderRadius: '4px',
                fontSize: '12px',
                overflow: 'auto'
              }}>
                {JSON.stringify(debugInfo, null, 2)}
              </pre>
            </Card>
          )}
          
          <Alert
            message="Próximos Passos"
            description={
              <div>
                <Text>1. Verifique se o token está presente</Text><br/>
                <Text>2. Verifique se a fazenda está selecionada</Text><br/>
                <Text>3. Se tudo estiver OK, o problema pode estar na API</Text>
              </div>
            }
            type="warning"
          />
        </Space>
      </Card>
    </div>
  );
};
