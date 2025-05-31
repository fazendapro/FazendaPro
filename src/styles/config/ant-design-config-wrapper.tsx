import { ConfigProvider, theme as antTheme } from 'antd';
import { useTheme } from '../useTheme';

export const AntConfigWrapper = ({ children }: { children: React.ReactNode }) => {
  const { isDarkMode } = useTheme();
  
  const lightColors = {
    colorPrimary: '#1890ff',
    colorText: '#000',
    colorTextBase: '#000',
    colorTextLightSolid: '#fff',
    colorTextSecondary: '#666',
    colorTextTertiary: '#999',
    colorTextQuaternary: '#bfbfbf',
    colorBgBase: '#fff',
  };
  
  const darkColors = {
    colorPrimary: '#177ddc',
    colorText: '#fff',
    colorTextBase: '#fff',
    colorTextLightSolid: '#000',
    colorTextSecondary: '#a6a6a6',
    colorTextTertiary: '#8c8c8c',
    colorTextQuaternary: '#595959',
    colorBgBase: '#141414',
  };
  
  return (
    <ConfigProvider
      theme={{
        algorithm: isDarkMode ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm,
        token: isDarkMode ? darkColors : lightColors,
        components: {
          Button: {
            colorPrimary: isDarkMode ? '#177ddc' : '#1890ff',
          },
        },
      }}
    >
      {children}
    </ConfigProvider>
  );
};