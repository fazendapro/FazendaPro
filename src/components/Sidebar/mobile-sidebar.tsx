import { Button } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { HomeOutlined, UserOutlined, LogoutOutlined, SettingOutlined, DollarOutlined } from "@ant-design/icons";
import { useAuth } from "../../contexts/AuthContext";

interface MenuItem {
  key: string;
  icon: React.ReactNode;
  label: string;
}

export const MobileSidebar = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, isAuthenticated } = useAuth();

  const handleMenuClick = (key: string) => {
    if (key === '/sair') {
      logout();
    } else {
      navigate(key);
    }
  };

  const menuItems: MenuItem[] = [
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
    <nav
      style={{
        position: 'fixed',
        bottom: 0,
        left: 0,
        right: 0,
        height: '60px',
        background: 'white',
        borderTop: '1px solid #f0f0f0',
        boxShadow: '0 -2px 8px rgba(0, 0, 0, 0.1)',
        display: 'flex',
        justifyContent: 'space-around',
        alignItems: 'center',
        zIndex: 1000,
        padding: '0 8px',
      }}
    >
      {menuItems.map((item) => {
        const isActive = location.pathname === item.key;
        
        return (
          <Button
            key={item.key}
            type={isActive ? 'primary' : 'text'}
            icon={item.icon}
            onClick={() => handleMenuClick(item.key)}
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              height: '48px',
              minWidth: '60px',
              fontSize: '10px',
              border: 'none',
              boxShadow: 'none',
            }}
          >
            <span style={{ marginTop: '2px', fontSize: '10px' }}>
              {item.label}
            </span>
          </Button>
        );
      })}
    </nav>
  );
};
