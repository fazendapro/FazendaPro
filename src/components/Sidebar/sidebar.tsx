import { Menu, Layout, Grid, Button, Avatar, Card, Typography } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useState } from "react";
import { HomeOutlined, UserOutlined, MenuFoldOutlined, MenuUnfoldOutlined, LogoutOutlined, SettingOutlined, DollarOutlined } from "@ant-design/icons";
import { useAuth } from "../../contexts/AuthContext";
import { useSelectedFarm } from "../../hooks/useSelectedFarm";

const { Sider } = Layout;
const { useBreakpoint } = Grid;
const { Text } = Typography;

export const Sidebar = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const screens = useBreakpoint();
  const [collapsed, setCollapsed] = useState(screens.xs);
  const { logout, user } = useAuth();
  const { farmName, farmLogo } = useSelectedFarm();
  const isAuthenticated = true;

  const handleMenuClick = (key: string) => {
    if (key === '/sair') {
      logout();
    } else {
      navigate(key);
      if (screens.xs) {
        setCollapsed(true);
      }
    }
  };

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: 'Dashboard' },
    { key: '/vacas', icon: <UserOutlined />, label: 'Vacas' },
    { key: '/vendas', icon: <DollarOutlined />, label: 'Vendas' },
    { key: '/configuracoes', icon: <SettingOutlined />, label: 'Configurações' },
    { key: '/sair', icon: <LogoutOutlined />, label: 'Sair' },
  ];

  if (!isAuthenticated) {
    return null;
  }

  return (
    <Sider
      collapsible
      collapsed={collapsed}
      onCollapse={(value) => setCollapsed(value)}
      width={280}
      style={{
        background: 'white',
        position: 'fixed',
        height: '100vh',
        zIndex: 1000,
        left: 0,
        top: 0,
        overflow: 'auto',
        transition: 'all 0.2s'
      }}
      trigger={screens.xs ? (
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={() => setCollapsed(!collapsed)}
          style={{
            fontSize: '16px',
            width: 64,
            height: 64,
            position: 'fixed',
            left: collapsed ? 0 : 280,
            top: 0,
            zIndex: 1001,
            transition: 'all 0.2s'
          }}
        />
      ) : null}
    >
      {!collapsed && (
        <Card 
          size="small" 
          style={{ 
            margin: '16px', 
            marginTop: '24px',
            textAlign: 'center'
          }}
          bodyStyle={{ padding: '16px' }}
        >
          {farmLogo && 
           farmLogo.trim() !== '' && 
           (farmLogo.startsWith('data:') || farmLogo.startsWith('http')) ? (
            <Avatar
              size={100}
              src={farmLogo}
              shape="square"
              style={{ marginBottom: 8 }}
            />
          ) : (
            <Text
              strong
              style={{ fontSize: farmLogo ? '12px' : '14px' }}
            >
              {farmName || 'FAZENDA'}
            </Text>
          )}

        </Card>
      )}

      <Menu
        theme="light"
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        style={{ 
          borderRight: 0,
          padding: '0 16px'
        }}
        onClick={({ key }) => handleMenuClick(key)}
      />

      {user && !collapsed && (
        <div style={{
          position: 'absolute',
          bottom: '16px',
          left: '16px',
          right: '16px'
        }}>
          <Card
            size="small"
            style={{
              textAlign: 'center'
            }}
          >
            <Text strong style={{ fontSize: '14px' }}>
              {user.email}
            </Text>
          </Card>
        </div>
      )}
    </Sider>
  );
};