import { Layout, Grid } from 'antd';
import { Login, Dashboard, Animals, Settings } from './pages';
import { AnimalDetailComponent as AnimalDetail, AnimalDetailProvider } from './pages/contents/AnimalDetail';
import { ResponsiveSidebar, Spinner } from './components';
import { Routes, Route, Navigate, useParams } from 'react-router'
import { FarmSelection } from './pages/FarmSelection';
import { FarmProvider } from './contexts/FarmContext';
import { AuthProvider, useAuth } from './contexts/AuthContext';

const { useBreakpoint } = Grid;

const AnimalDetailWrapper = () => {
  const { id } = useParams<{ id: string }>();
  const animalId = parseInt(id || '0');
  
  return (
    <AnimalDetailProvider animalId={animalId}>
      <AnimalDetail />
    </AnimalDetailProvider>
  );
};

const ProtectedLayout = ({ children }: { children: React.ReactNode }) => {
  const { isLoading, isAuthenticated } = useAuth()
  const screens = useBreakpoint();

  if (isLoading) return <Spinner />

  if (!isAuthenticated) return <Navigate to="/login" replace />

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <ResponsiveSidebar />
      <Layout style={{
        marginLeft: screens.xs ? 0 : 280,
        marginBottom: screens.xs ? '60px' : 0,
        transition: 'all 0.2s'
      }}>
        <Layout.Content style={{ 
          height: screens.xs ? 'calc(100vh - 60px)' : '100vh',
          overflow: 'auto',
          background: '#f5f5f5'
        }}>
          <div style={{
            padding: '24px',
            minHeight: screens.xs ? 'calc(100vh - 84px)' : 'calc(100vh - 48px)'
          }}>
            {children}
          </div>
        </Layout.Content>
      </Layout>
    </Layout>
  )
}

export const App = () => {
  return (
    <FarmProvider>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </FarmProvider>
  );
};

const AppContent = () => {
  const { isLoading, isAuthenticated } = useAuth();

  if (isLoading) return <Spinner />

  return (
    <Routes>
      <Route
        path="/login"
        element={isAuthenticated ? <Navigate to="/" replace /> : <Login />}
      />
      <Route
        path="/farm-selection"
        element={isAuthenticated ? <FarmSelection /> : <Navigate to="/login" replace />}
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
        path="/animal/:id"
        element={
          <ProtectedLayout>
            <AnimalDetailWrapper />
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
            <Settings />
          </ProtectedLayout>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};