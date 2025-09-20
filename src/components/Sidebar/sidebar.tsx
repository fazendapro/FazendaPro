import { Menu, Layout, Grid, Button } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useState } from "react";
import { HomeOutlined, UserOutlined, MenuFoldOutlined, MenuUnfoldOutlined, LogoutOutlined } from "@ant-design/icons";
import { useAuth } from "../../pages/Login/hooks/useAuth";

const { Sider } = Layout;
const { useBreakpoint } = Grid;

export const Sidebar = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const screens = useBreakpoint();
  const [collapsed, setCollapsed] = useState(screens.xs);
  const { logout } = useAuth();
  const isAuthenticated = true; // TODO: remove this and remove useAuth

  const handleMenuClick = (key: string) => {
    if (key === '/sair') {
      logout();
      navigate('/login');
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
    // { key: '/relatorios', icon: <FileTextOutlined />, label: 'Relatórios' },
    // { key: '/fornecedores', icon: <ShoppingCartOutlined />, label: 'Fornecedores' },
    // { key: '/vendas', icon: <ShoppingCartOutlined />, label: 'Vendas' },
    // { key: '/estoque', icon: <FileTextOutlined />, label: 'Estoque' },
    // { key: '/configuracoes', icon: <SettingOutlined />, label: 'Configurações' },
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
      <div style={{
        height: 64,
        margin: 16,
        marginTop: 24,
        textAlign: 'center',
        background: 'rgba(255, 255, 255, 0.2)',
        display: collapsed ? 'none' : 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontWeight: 'bold',
        fontSize: '14px',
        color: '#333'
      }}>
        FAZENDA BOM JARDIM
      </div>
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
    </Sider>
  );
};