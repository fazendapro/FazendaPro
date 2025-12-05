import React, { useState, useEffect } from "react";
import { Button, Dropdown } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { HomeOutlined, UserOutlined, LogoutOutlined, SettingOutlined, DollarOutlined, SwapOutlined } from "@ant-design/icons";
import { useAuth } from "../../contexts/AuthContext";
import { useSelectedFarm } from "../../hooks/useSelectedFarm";
import { useFarmSwitcher } from "../../hooks/useFarmSwitcher";

interface MenuItem {
  key: string;
  icon: React.ReactNode;
  label: string;
}

export const MobileSidebar = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, isAuthenticated } = useAuth();
  const { selectedFarm } = useSelectedFarm();
  const { farms, switchFarm, loadFarms } = useFarmSwitcher();
  const [showFarmMenu, setShowFarmMenu] = useState(false);

  useEffect(() => {
    if (farms.length === 0) {
      loadFarms();
    }
  }, [farms.length, loadFarms]);

  const handleMenuClick = (key: string) => {
    if (key === '/sair') {
      logout();
    } else if (key === '/trocar-fazenda') {
      setShowFarmMenu(true);
    } else {
      navigate(key);
    }
  };

  const handleFarmSelect = (farmId: number) => {
    if (selectedFarm?.ID !== farmId) {
      switchFarm(farmId);
    }
    setShowFarmMenu(false);
  };

  const farmMenuItems = farms.length > 1 ? farms.map((farm) => ({
    key: `farm-${farm.ID}`,
    label: farm.Company?.CompanyName || `Fazenda ${farm.ID}`,
    onClick: () => handleFarmSelect(farm.ID),
    disabled: selectedFarm?.ID === farm.ID,
  })) : [];

  const menuItems: MenuItem[] = [
    { key: '/', icon: <HomeOutlined />, label: t('navigation.dashboard') },
    { key: '/vacas', icon: <UserOutlined />, label: t('navigation.cattle') },
    { key: '/vendas', icon: <DollarOutlined />, label: t('navigation.sales') },
    { key: '/configuracoes', icon: <SettingOutlined />, label: t('navigation.settings') },
    ...(farms.length > 1 ? [{ key: '/trocar-fazenda', icon: <SwapOutlined />, label: t('navigation.switchFarm') || 'Trocar Fazenda' }] : []),
    { key: '/sair', icon: <LogoutOutlined />, label: t('navigation.logout') },
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
        const isFarmSwitch = item.key === '/trocar-fazenda';

        if (isFarmSwitch && farms.length > 1) {
          return (
            <Dropdown
              key={item.key}
              menu={{ items: farmMenuItems }}
              trigger={['click']}
              open={showFarmMenu}
              onOpenChange={setShowFarmMenu}
            >
              <Button
                type="text"
                icon={item.icon}
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
            </Dropdown>
          );
        }
        
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
