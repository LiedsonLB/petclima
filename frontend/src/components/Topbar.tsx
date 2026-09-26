import { Moon, Sun, UserCircle } from 'lucide-react';
import { useLocation } from 'react-router-dom';
import { useTheme } from '../contexts/ThemeContext';

const routeTitles: Record<string, { title: string; subtitle: string }> = {
  '/dashboard': { title: 'Dashboard', subtitle: 'Visão geral das entregas do GAT 5' },
  '/cronograma': { title: 'Cronograma', subtitle: 'Entregas por semestre, Set/2026 a Ago/2028' },
  '/comunicacao': { title: 'Comunicação', subtitle: 'Plano transversal e kit de comunicação regional' },
  '/formacao': { title: 'Educação Permanente', subtitle: 'Oficinas de letramento em saúde e combate às fake news' },
  '/ava': { title: 'AVA · Cursos', subtitle: 'Ambiente Virtual de Aprendizagem — cursos autoinstrucionais' },
  '/mapa': { title: 'Mapa de Alcance', subtitle: 'Municípios e macrorregiões atendidos no Piauí' },
  '/relatorios': { title: 'Relatórios', subtitle: 'Indicadores de tráfego, acessos e engajamento' },
  '/repositorio': { title: 'Repositório', subtitle: 'Tecnologias sociais e digitais entregues à SESAPI' },
  '/configuracoes': { title: 'Configurações', subtitle: 'Preferências do sistema' },
};

export default function Topbar() {
  const { theme, toggle } = useTheme();
  const location = useLocation();
  const { title, subtitle } = routeTitles[location.pathname] ?? { title: 'GAT 5 Comunica', subtitle: '' };

  return (
    <div style={{
      height: 'var(--topbar-height, 56px)',
      borderBottom: '1px solid var(--border)',
      background: 'var(--bg-topbar)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '0 20px',
      flexShrink: 0,
      position: 'sticky',
      top: 0,
      zIndex: 20,
    }}>
      <div>
        <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)' }}>{title}</div>
        {subtitle && <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{subtitle}</div>}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <button
          onClick={toggle}
          aria-label="Alternar tema"
          style={{
            width: 32, height: 32, borderRadius: 8, border: '1px solid var(--border)',
            background: 'var(--bg-card)', display: 'flex', alignItems: 'center', justifyContent: 'center',
          }}
        >
          {theme === 'light' ? <Moon size={15} /> : <Sun size={15} />}
        </button>
        <div style={{ width: 32, height: 32, borderRadius: '50%', background: 'var(--primary-light)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <UserCircle size={20} color="var(--primary)" />
        </div>
      </div>
    </div>
  );
}
