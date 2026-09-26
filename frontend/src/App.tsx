import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import AppLayout from './components/AppLayout';

import Landing from './pages/Landing';
import Login from './pages/Login';
import Cadastro from './pages/Cadastro';
import EsqueciSenha from './pages/EsqueciSenha';
import Dashboard from './pages/Dashboard';
import Cronograma from './pages/Cronograma';
import Comunicacao from './pages/Comunicacao';
import Formacao from './pages/Formacao';
import AVA from './pages/AVA';
import Mapa from './pages/Mapa';
import Relatorios from './pages/Relatorios';
import Repositorio from './pages/Repositorio';
import Configuracoes from './pages/Configuracoes';

function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/login" element={<Login />} />
          <Route path="/cadastro" element={<Cadastro />} />
          <Route path="/esqueci-senha" element={<EsqueciSenha />} />

          <Route path="/dashboard" element={<AppLayout><Dashboard /></AppLayout>} />
          <Route path="/cronograma" element={<AppLayout><Cronograma /></AppLayout>} />
          <Route path="/comunicacao" element={<AppLayout><Comunicacao /></AppLayout>} />
          <Route path="/formacao" element={<AppLayout><Formacao /></AppLayout>} />
          <Route path="/ava" element={<AppLayout><AVA /></AppLayout>} />
          <Route path="/mapa" element={<AppLayout><Mapa /></AppLayout>} />
          <Route path="/relatorios" element={<AppLayout><Relatorios /></AppLayout>} />
          <Route path="/repositorio" element={<AppLayout><Repositorio /></AppLayout>} />
          <Route path="/configuracoes" element={<AppLayout><Configuracoes /></AppLayout>} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
