import { CalendarDays, Users, ShieldAlert, Plus } from 'lucide-react';
import Card from '../components/Card';
import KpiCard from '../components/KpiCard';
import ProgressBar from '../components/ProgressBar';
import { oficinas } from '../data/Data';

export default function Formacao() {
  const inscritos = oficinas.reduce((acc, o) => acc + o.inscritos, 0);
  const vagasTotais = oficinas.reduce((acc, o) => acc + o.vagas, 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card>
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Programa de Formação em Letramento em Saúde e Combate às Fake News: oficinas presenciais
          e remotas voltadas a profissionais de saúde e Agentes Comunitários de Saúde (ACS), com a
          meta de capacitar no mínimo <strong style={{ color: 'var(--text-primary)' }}>200 pessoas</strong> até o fim do Semestre 2.
        </p>
      </Card>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 14 }}>
        <KpiCard label="Inscritos até agora" value={`${inscritos}`} sub={`meta: 200 profissionais/ACS`} icon={<Users size={16} />} />
        <KpiCard label="Oficinas programadas" value={`${oficinas.length}`} sub="presenciais e remotas" icon={<CalendarDays size={16} />} />
        <KpiCard label="Ocupação média" value={`${Math.round((inscritos / vagasTotais) * 100)}%`} sub="das vagas oferecidas" icon={<ShieldAlert size={16} />} />
      </div>

      <Card>
        <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 6 }}>Progresso rumo à meta de 200 profissionais capacitados</div>
        <ProgressBar value={inscritos} max={200} />
      </Card>

      <Card
        title="Oficinas programadas"
        icon={<CalendarDays size={16} />}
        action={
          <button style={{
            display: 'flex', alignItems: 'center', gap: 6, fontSize: 12, fontWeight: 600,
            color: '#fff', background: 'var(--primary)', border: 'none', borderRadius: 999,
            padding: '6px 14px', cursor: 'pointer',
          }}>
            <Plus size={13} /> Nova oficina
          </button>
        }
        noPad
      >
        <div>
          {oficinas.map((o, i) => (
            <div key={o.id} style={{
              display: 'flex', alignItems: 'center', justifyContent: 'space-between',
              padding: '12px 16px', borderTop: i === 0 ? 'none' : '1px solid var(--border)',
            }}>
              <div>
                <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>{o.titulo}</div>
                <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{o.municipio} · {new Date(o.data).toLocaleDateString('pt-BR')}</div>
              </div>
              <div style={{ width: 160 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 11, color: 'var(--text-muted)', marginBottom: 4 }}>
                  <span>{o.inscritos} inscritos</span>
                  <span>{o.vagas} vagas</span>
                </div>
                <ProgressBar value={o.inscritos} max={o.vagas} />
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
