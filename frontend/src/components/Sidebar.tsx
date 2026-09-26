import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard, Megaphone, GraduationCap, BookOpenCheck,
  MapPinned, FileBarChart, Archive, CalendarClock, Settings,
} from 'lucide-react';

const navPrincipal = [
  { to: '/dashboard', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/cronograma', icon: CalendarClock, label: 'Cronograma' },
];

const navConteudo = [
  { to: '/comunicacao', icon: Megaphone, label: 'Comunicação' },
  { to: '/formacao', icon: GraduationCap, label: 'Educação Permanente' },
  { to: '/ava', icon: BookOpenCheck, label: 'AVA · Cursos' },
];

const navTerritorio = [
  { to: '/mapa', icon: MapPinned, label: 'Mapa de Alcance' },
  { to: '/relatorios', icon: FileBarChart, label: 'Relatórios' },
  { to: '/repositorio', icon: Archive, label: 'Repositório' },
];

const itemStyle = (active: boolean): React.CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  padding: '9px 12px',
  borderRadius: '10px',
  fontSize: '13px',
  fontWeight: 500,
  color: active ? 'var(--primary)' : 'var(--text-secondary)',
  background: active ? 'var(--bg-chip-active)' : 'transparent',
  marginBottom: '2px',
  textDecoration: 'none',
});

function NavGroup({ label, items }: { label: string; items: typeof navPrincipal }) {
  return (
    <div>
      <div style={{
        fontSize: '10px', fontWeight: 700, color: 'var(--text-muted)',
        padding: '12px 10px 6px', letterSpacing: '0.08em', textTransform: 'uppercase',
      }}>
        {label}
      </div>
      {items.map(({ to, icon: Icon, label }) => (
        <NavLink key={to} to={to} style={({ isActive }) => itemStyle(isActive)}>
          <Icon size={17} strokeWidth={2} />
          <span>{label}</span>
        </NavLink>
      ))}
    </div>
  );
}

export default function Sidebar() {
  return (
    <div style={{
      width: 'var(--sidebar-width, 228px)',
      background: 'var(--bg-sidebar)',
      borderRight: '1px solid var(--border)',
      display: 'flex',
      flexDirection: 'column',
      flexShrink: 0,
      height: '100vh',
      position: 'sticky',
      top: 0,
      overflow: 'hidden',
    }}>
      <div style={{ padding: '24px 20px 20px', borderBottom: '1px solid var(--border)', flexShrink: 0 }}>
        <NavLink to="/dashboard" style={{ display: 'flex', alignItems: 'center', gap: 12, textDecoration: 'none' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <img src="petsaudeclima_icon.png" alt="PET-Saúde Clima Logo" style={{ height: 32 }} />
            <span style={{ fontSize: 13, fontWeight: 700, color: 'var(--primary)' }}>PET-Saúde Clima</span>
          </div>
        </NavLink>
      </div>

      <div style={{ padding: '12px 10px', flex: 1, overflowY: 'auto' }}>
        <NavGroup label="Visão Geral" items={navPrincipal} />
        <NavGroup label="Comunicação e Formação" items={navConteudo} />
        <NavGroup label="Território" items={navTerritorio} />
      </div>

      <div style={{ padding: '10px', borderTop: '1px solid var(--border)' }}>
        <NavLink to="/configuracoes" style={({ isActive }) => itemStyle(isActive)}>
          <Settings size={17} />
          <span>Configurações</span>
        </NavLink>
      </div>
    </div>
  );
}
