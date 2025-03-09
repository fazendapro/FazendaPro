import { Layout, Button, Switch, Typography, Space, Card, Divider } from 'antd';
import { BulbOutlined, BulbFilled } from '@ant-design/icons';
import { useTheme } from './styles/index';
import { useTranslation } from 'react-i18next';

const { Header, Content, Footer } = Layout;
const { Title, Paragraph, Text } = Typography;

export const App = () => {
  const { isDarkMode, toggleTheme } = useTheme();
  const { t } = useTranslation();

  const fixedColorCard = {
    backgroundColor: isDarkMode ? '#141414' : '#f0f2f5',
    color: isDarkMode ? '#fff' : '#000',
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '0 20px',
        backgroundColor: fixedColorCard.backgroundColor,
      }}>
        <Title level={4} style={{ color: fixedColorCard.color, margin: 0 }}>
          {t('home.welcome')}
        </Title>
        <Space>
          <Text style={{ color: fixedColorCard.color }}>
            {isDarkMode ? 'Modo Escuro' : 'Modo Claro'}
          </Text>
          <Switch
            checked={isDarkMode}
            onChange={toggleTheme}
            checkedChildren={<BulbFilled />}
            unCheckedChildren={<BulbOutlined />}
          />
        </Space>
      </Header>
      <Content style={{ padding: '50px', backgroundColor: fixedColorCard.backgroundColor }}>
        <Card>
          <Title level={2}>Configuração do Tema</Title>
          <Paragraph>
            {t('home.headerTitle')}
          </Paragraph>
          <Divider />
          <Space direction="vertical" size="large">
            <Card title="Demonstração de componentes">
              <Space wrap>
                <Button type="primary">Botão Primário</Button>
                <Button>Botão Padrão</Button>
                <Button type="dashed">Botão Tracejado</Button>
                <Button type="text">Botão Texto</Button>
                <Button type="link">Botão Link</Button>
              </Space>
            </Card>

            <Button
              type="primary"
              size="large"
              onClick={toggleTheme}
              icon={isDarkMode ? <BulbFilled /> : <BulbOutlined />}
            >
              Alternar para {isDarkMode ? 'Modo Claro' : 'Modo Escuro'}
            </Button>
          </Space>
        </Card>
      </Content>

      <Footer style={{ textAlign: 'center', backgroundColor: fixedColorCard.backgroundColor }}>
        <Text style={{ color: fixedColorCard.color }}>
          FarmPro - Soluções em Agro ©{new Date().getFullYear()}
        </Text>
      </Footer>
    </Layout>
  );
};