import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App';
import { ThemeProvider } from './styles/context/theme-provider';
import { AntConfigWrapper } from './styles/config/ant-design-config-wrapper';
import './locale/i18n';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <AntConfigWrapper>
        <App />
      </AntConfigWrapper>
    </ThemeProvider>
  </StrictMode>
);
