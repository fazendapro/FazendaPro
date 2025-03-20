import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App';
import { ThemeProvider } from './styles/context/theme-provider';
import { AntConfigWrapper } from './styles/config/ant-design-config-wrapper';
import { BrowserRouter } from 'react-router';
import './locale/i18n';
import { ToastContainer } from 'react-toastify';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <AntConfigWrapper>
        <BrowserRouter>
        <ToastContainer />
          <App />
        </BrowserRouter>
      </AntConfigWrapper>
    </ThemeProvider>
  </StrictMode>
);
