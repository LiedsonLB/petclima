import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import AppLayout from './components/AppLayout';
import RequireAuth from './components/RequireAuth';

import Landing from './pages/Landing';
import Login from './pages/Login';
import Cadastro from './pages/Cadastro';
import EsqueciSenha from './pages/EsqueciSenha';
import RedefinirSenha from './pages/RedefinirSenha';
import ConfirmarEmail from './pages/ConfirmarEmail';
import Dashboard from './pages/Dashboard';
import Cronograma from './pages/Cronograma';
import Comunicacao from './pages/Comunicacao';
import Formacao from './pages/Formacao';
import AVA from './pages/AVA';
import Mapa from './pages/Mapa';
import Relatorios from './pages/Relatorios';
import Repositorio from './pages/Repositorio';
import Configuracoes from './pages/Configuracoes';

function Protegida({ children }: { children: React.ReactNode }) {
  return (
    <RequireAuth>
      <AppLayout>{children}</AppLayout>
    </RequireAuth>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Landing />} />
            <Route path="/login" element={<Login />} />
            <Route path="/cadastro" element={<Cadastro />} />
            <Route path="/esqueci-senha" element={<EsqueciSenha />} />
            <Route path="/redefinir-senha" element={<RedefinirSenha />} />
            <Route path="/confirmar-email" element={<ConfirmarEmail />} />

            <Route path="/dashboard" element={<Protegida><Dashboard /></Protegida>} />
            <Route path="/cronograma" element={<Protegida><Cronograma /></Protegida>} />
            <Route path="/comunicacao" element={<Protegida><Comunicacao /></Protegida>} />
            <Route path="/formacao" element={<Protegida><Formacao /></Protegida>} />
            <Route path="/ava" element={<Protegida><AVA /></Protegida>} />
            <Route path="/mapa" element={<Protegida><Mapa /></Protegida>} />
            <Route path="/relatorios" element={<Protegida><Relatorios /></Protegida>} />
            <Route path="/repositorio" element={<Protegida><Repositorio /></Protegida>} />
            <Route path="/configuracoes" element={<Protegida><Configuracoes /></Protegida>} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;