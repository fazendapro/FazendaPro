import { Routes, Route, Navigate } from 'react-router'
import { useAuth } from './hooks/useAuth';
import { Home } from './pages/Home';
import Login from './pages/Login/login';
import { Layout } from 'antd';


const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, token, user } = useAuth()
  console.log('Estado de autenticação:', { isAuthenticated, hasToken: !!token, hasUser: !!user })
  return isAuthenticated ? children : <Navigate to="/login" />
}

export const App = () => {
  return (
    <Layout style={{ minHeight: '100vh', padding: '50px', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center' }}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
        path="/"
        element={
          <ProtectedRoute>
            <Home />
          </ProtectedRoute>
        }
      />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Layout>
  );
};