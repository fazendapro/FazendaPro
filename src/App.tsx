import { Layout, Grid } from 'antd';
import { useAuth, Login, Dashboard, Animals } from './pages';
import { Sidebar, Spinner } from './components';
import { Routes, Route, Navigate } from 'react-router'

const { useBreakpoint } = Grid;
const isAuthenticated = true; // TODO: remove this and remove useAuth

const ProtectedLayout = ({ children }: { children: React.ReactNode }) => {
  const { isLoading } = useAuth()
  const screens = useBreakpoint();

  if (isLoading) return <Spinner />

  if (!isAuthenticated) return <Navigate to="/login" replace />

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout style={{
        marginLeft: screens.xs ? 0 : 280,
        transition: 'all 0.2s'
      }}>
        <Layout.Content style={{ 
          height: '100vh',
          overflow: 'auto',
          background: '#f5f5f5'
        }}>
          <div style={{
            padding: '24px',
            minHeight: 'calc(100vh - 48px)'
          }}>
            {children}
          </div>
        </Layout.Content>
      </Layout>
    </Layout>
  )
}

export const App = () => {
  const { isLoading } = useAuth();

  if (isLoading) return <Spinner />

  return (
    <Routes>
      <Route
        path="/login"
        element={isAuthenticated ? <Navigate to="/" replace /> : <Login />}
      />
      <Route
        path="/"
        element={
          <ProtectedLayout>
            <Dashboard />
          </ProtectedLayout>
        }
      />
      <Route
        path="/vacas"
        element={
          <ProtectedLayout>
            <Animals />
          </ProtectedLayout>
        }
      />
      <Route
        path="/relatorios"
        element={
          <ProtectedLayout>
            <h1>Página de Relatórios</h1>
          </ProtectedLayout>
        }
      />
      <Route
        path="/fornecedores"
        element={
          <ProtectedLayout>
            <h1>Página de Fornecedores</h1>
          </ProtectedLayout>
        }
      />
      <Route
        path="/vendas"
        element={
          <ProtectedLayout>
            <h1>Página de Vendas</h1>
          </ProtectedLayout>
        }
      />
      <Route
        path="/estoque"
        element={
          <ProtectedLayout>
            <h1>Página de Estoque</h1>
          </ProtectedLayout>
        }
      />
      <Route
        path="/configuracoes"
        element={
          <ProtectedLayout>
            <h1>Página de Configurações</h1>
          </ProtectedLayout>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};