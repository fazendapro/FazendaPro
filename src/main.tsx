import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App';
import { ThemeProvider } from './styles/context/theme/theme-provider';
import { AntConfigWrapper } from './styles/config/ant-design-config-wrapper';
import { BrowserRouter } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';
import './locale/i18n';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <AntConfigWrapper>
        <ToastContainer />
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </AntConfigWrapper>
    </ThemeProvider>
  </StrictMode>
);
