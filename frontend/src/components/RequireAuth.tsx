import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export default function RequireAuth({ children }: { children: React.ReactNode }) {
  const { autenticado, verificando } = useAuth();

  if (verificando) {
    return (
      <div
        style={{
          minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center',
          background: 'var(--bg)', color: 'var(--text-muted)', fontSize: 13,
        }}
      >
        Verificando sessão…
      </div>
    );
  }

  if (!autenticado) return <Navigate to="/login" replace />;
  return <>{children}</>;
}