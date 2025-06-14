import { Routes, Route, Navigate } from 'react-router'
import { useAuth, Login, Home } from './pages';
import { Layout } from 'antd';
import { Sidebar } from './components/Sidebar/sidebar';
import { Spinner } from './components/Spinner/spinner';

const ProtectedLayout = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) return <Spinner />

  if (!isAuthenticated) return <Navigate to="/login" replace />

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout>
        <Layout.Content style={{ 
          padding: '24px',
          transition: 'all 0.2s'
        }}>
          {children}
        </Layout.Content>
      </Layout>
    </Layout>
  )
}

export const App = () => {
  const { isAuthenticated, isLoading } = useAuth();

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
            <Home />
          </ProtectedLayout>
        }
      />
      <Route
        path="/vacas"
        element={
          <ProtectedLayout>
            <div>Página de Vacas</div>
          </ProtectedLayout>
        }
      />
      <Route
        path="/relatorios"
        element={
          <ProtectedLayout>
            <div>Página de Relatórios</div>
          </ProtectedLayout>
        }
      />
      <Route
        path="/fornecedores"
        element={
          <ProtectedLayout>
            <div>Página de Fornecedores</div>
          </ProtectedLayout>
        }
      />
      <Route
        path="/vendas"
        element={
          <ProtectedLayout>
            <div>Página de Vendas</div>
          </ProtectedLayout>
        }
      />
      <Route
        path="/estoque"
        element={
          <ProtectedLayout>
            <div>Página de Estoque</div>
          </ProtectedLayout>
        }
      />
      <Route
        path="/configuracoes"
        element={
          <ProtectedLayout>
            <div>Página de Configurações</div>
          </ProtectedLayout>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};