import { MapPinned, Users, FileText, Landmark, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import Card from '../components/Card';
import KpiCard from '../components/KpiCard';
import ProgressBar, { StatusBadge } from '../components/ProgressBar';
import { cronograma, macrorregioes, materiais, oficinas } from '../data/Data';

export default function Dashboard() {
  const municipiosAtendidos = macrorregioes.reduce((acc, m) => acc + m.municipiosAtendidos, 0);
  const inscritosOficinas = oficinas.reduce((acc, o) => acc + o.inscritos, 0);
  const etapaAtual = cronograma.find(s => s.status === 'andamento') ?? cronograma[0];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 14 }}>
        <KpiCard label="Municípios atendidos" value={`${municipiosAtendidos} / 12`} sub="meta do projeto (4 anos)" icon={<MapPinned size={16} />} />
        <KpiCard label="Profissionais capacitados" value={`${inscritosOficinas} / 200`} sub="meta do Semestre 2" icon={<Users size={16} />} />
        <KpiCard label="Materiais produzidos" value={`${materiais.length} / 10`} sub="Kit de Comunicação Regional" icon={<FileText size={16} />} />
        <KpiCard label="Macrorregiões cobertas" value={`${macrorregioes.filter(m => m.municipiosAtendidos > 0).length} / 4`} sub="Meio-Norte, Semiárido, Cerrado, Entre Rios" icon={<Landmark size={16} />} />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 20 }}>
        <Card title="Semestre em andamento" icon={<Landmark size={16} />}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 10 }}>
            <div>
              <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)' }}>{etapaAtual.titulo}</div>
              <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>{etapaAtual.periodo}</div>
            </div>
            <StatusBadge status={etapaAtual.status} />
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {etapaAtual.entregas.map(e => (
              <div key={e.nome} style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 12, padding: '8px 0', borderTop: '1px solid var(--border)' }}>
                <div>
                  <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>{e.nome}</div>
                  <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{e.detalhe}</div>
                </div>
                <StatusBadge status={e.status} />
              </div>
            ))}
          </div>
          <Link to="/cronograma" style={{ display: 'inline-flex', alignItems: 'center', gap: 6, marginTop: 14, fontSize: 12, fontWeight: 600, color: 'var(--primary)', textDecoration: 'none' }}>
            Ver cronograma completo <ArrowRight size={13} />
          </Link>
        </Card>

        <Card title="Alcance por macrorregião" icon={<MapPinned size={16} />}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
            {macrorregioes.map(m => (
              <div key={m.nome}>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 12, marginBottom: 5 }}>
                  <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>{m.nome}</span>
                  <span style={{ color: 'var(--text-muted)' }}>{m.municipiosAtendidos} / {m.meta} municípios</span>
                </div>
                <ProgressBar value={m.municipiosAtendidos} max={m.meta} />
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card title="Últimos materiais publicados" icon={<FileText size={16} />} noPad>
        <div>
          {materiais.slice(0, 5).map((m, i) => (
            <div key={m.id} style={{
              display: 'flex', alignItems: 'center', justifyContent: 'space-between',
              padding: '10px 16px', borderTop: i === 0 ? 'none' : '1px solid var(--border)',
            }}>
              <div>
                <div style={{ fontSize: 13, color: 'var(--text-primary)', fontWeight: 500 }}>{m.titulo}</div>
                <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{m.macrorregiao} · {m.tema}</div>
              </div>
              <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>{m.downloads} downloads</div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
