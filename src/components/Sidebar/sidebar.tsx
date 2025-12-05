import { Menu, Layout, Grid, Button, Avatar, Card, Typography, Select } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { HomeOutlined, UserOutlined, MenuFoldOutlined, MenuUnfoldOutlined, LogoutOutlined, SettingOutlined, DollarOutlined } from "@ant-design/icons";
import { useAuth } from "../../contexts/AuthContext";
import { useSelectedFarm } from "../../hooks/useSelectedFarm";
import { useFarmSwitcher } from "../../hooks/useFarmSwitcher";

const { Sider } = Layout;
const { useBreakpoint } = Grid;
const { Text } = Typography;

export const Sidebar = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const navigate = useNavigate();
  const screens = useBreakpoint();
  const [collapsed, setCollapsed] = useState(screens.xs);
  const { logout, user } = useAuth();
  const { farmName, farmLogo, selectedFarm } = useSelectedFarm();
  const { farms, loading: farmsLoading, switchFarm, loadFarms } = useFarmSwitcher();
  const isAuthenticated = true;

  useEffect(() => {
    if (!collapsed && farms.length === 0) {
      loadFarms();
    }
  }, [collapsed, farms.length, loadFarms]);

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
    { key: '/', icon: <HomeOutlined />, label: t('navigation.dashboard') },
    { key: '/vacas', icon: <UserOutlined />, label: t('navigation.cattle') },
    { key: '/vendas', icon: <DollarOutlined />, label: t('navigation.sales') },
    { key: '/configuracoes', icon: <SettingOutlined />, label: t('navigation.settings') },
    { key: '/sair', icon: <LogoutOutlined />, label: t('navigation.logout') },
  ];

  if (!isAuthenticated) {
    return null;
  }

  return (
    <Sider
      collapsible
      collapsed={collapsed}
      onCollapse={(value: boolean) => setCollapsed(value)}
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
            textAlign: 'center',
            border: 'none'
          }}
          styles={{ body: { padding: '16px' } }}
        >
          {farmLogo && 
           farmLogo.trim() !== '' && 
           (farmLogo.startsWith('data:') || farmLogo.startsWith('http')) ? (
            <>
              <Avatar
                size={100}
                src={farmLogo}
                shape="square"
                style={{ marginBottom: farms.length > 1 ? 12 : 0 }}
              />
              {farms.length > 1 && (
                <div>
                  <Select
                    value={selectedFarm?.ID}
                    onChange={(farmId: number) => switchFarm(farmId)}
                    loading={farmsLoading}
                    style={{ 
                      width: '100%',
                      fontSize: '15px',
                      fontWeight: 'bold'
                    }}
                    placeholder={t('navigation.selectFarm') || 'Selecionar Fazenda'}
                    size="small"
                    showSearch
                    optionFilterProp="children"
                    filterOption={(input: string, option?: { children?: React.ReactNode }) =>
                      (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
                    }
                    styles={{
                      selector: {
                        fontSize: '15px',
                        fontWeight: 'bold',
                      }
                    }}
                  >
                    {farms.map((farm) => (
                      <Select.Option key={farm.ID} value={farm.ID} style={{ fontSize: '15px', fontWeight: 'bold' }}>
                        {farm.Company?.CompanyName || `Fazenda ${farm.ID}`}
                      </Select.Option>
                    ))}
                  </Select>
                </div>
              )}
            </>
          ) : farms.length > 1 ? (
            <Select
              value={selectedFarm?.ID}
              onChange={(farmId: number) => switchFarm(farmId)}
              loading={farmsLoading}
              style={{ 
                width: '100%',
                fontSize: '15px',
                fontWeight: 'bold'
              }}
              placeholder={t('navigation.selectFarm') || 'Selecionar Fazenda'}
              size="small"
              showSearch
              optionFilterProp="children"
              filterOption={(input: string, option?: { children?: React.ReactNode }) =>
                (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
              }
              styles={{
                selector: {
                  fontSize: '15px',
                  fontWeight: 'bold',
                }
              }}
            >
              {farms.map((farm) => (
                <Select.Option key={farm.ID} value={farm.ID}>
                  {farm.Company?.CompanyName || `Fazenda ${farm.ID}`}
                </Select.Option>
              ))}
            </Select>
          ) : (
            <Text
              strong
              style={{ fontSize: '14px' }}
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
        onClick={({ key }: { key: string }) => handleMenuClick(key)}
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