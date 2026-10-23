import { MapPinned, Users, FileText, Landmark, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import Card from '../components/Card';
import KpiCard from '../components/KpiCard';
import ProgressBar, { StatusBadge } from '../components/ProgressBar';
import { cronograma } from '../data/Data';
import { listarMacrorregioes, listarMateriais, painelResumo } from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro } from '../lib/ui';

export default function Dashboard() {
  const resumo = useRecurso(painelResumo);
  const macros = useRecurso(listarMacrorregioes);
  const materiais = useRecurso(listarMateriais, { limit: 5 });
  const etapaAtual = cronograma.find((s) => s.status === 'andamento') ?? cronograma[0];
  const k = resumo.data;
  const par = (a?: number, b?: number) => (k ? `${a} / ${b}` : '—');

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      {resumo.erro && <div style={alertaErro}>{resumo.erro}</div>}

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 14 }}>
        <KpiCard label="Municípios atendidos" value={par(k?.municipios_atendidos, k?.municipios_meta)} sub="meta do projeto (4 anos)" icon={<MapPinned size={16} />} />
        <KpiCard label="Profissionais capacitados" value={par(k?.profissionais_capacitados, k?.profissionais_meta)} sub="presença confirmada em atividades" icon={<Users size={16} />} />
        <KpiCard label="Materiais produzidos" value={par(k?.materiais_produzidos, k?.materiais_meta)} sub="Kit de Comunicação Regional" icon={<FileText size={16} />} />
        <KpiCard label="Macrorregiões cobertas" value={par(k?.macrorregioes_cobertas, k?.macrorregioes_total)} sub="Meio-Norte, Semiárido, Cerrado, Entre Rios" icon={<Landmark size={16} />} />
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
            {etapaAtual.entregas.map((e) => (
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
            {(macros.data ?? []).map((m) => (
              <div key={m.id}>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 12, marginBottom: 5 }}>
                  <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>{m.nome}</span>
                  <span style={{ color: 'var(--text-muted)' }}>{m.municipios_ativos ?? 0} / {m.meta_municipios} municípios</span>
                </div>
                <ProgressBar value={m.municipios_ativos ?? 0} max={m.meta_municipios || 1} />
              </div>
            ))}
            {macros.carregando && <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Carregando…</span>}
          </div>
        </Card>
      </div>

      <Card title="Últimos materiais publicados" icon={<FileText size={16} />} noPad>
        <div>
          {(materiais.data ?? []).map((m, i) => (
            <div key={m.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '10px 16px', borderTop: i === 0 ? 'none' : '1px solid var(--border)' }}>
              <div>
                <div style={{ fontSize: 13, color: 'var(--text-primary)', fontWeight: 500 }}>{m.titulo}</div>
                <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{m.macrorregiao_nome ?? 'Estadual'} · {m.tema ?? m.tipo}</div>
              </div>
              <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>{m.downloads} downloads</div>
            </div>
          ))}
          {!materiais.carregando && (materiais.data ?? []).length === 0 && (
            <div style={{ padding: 16, fontSize: 12, color: 'var(--text-muted)' }}>Nenhum material publicado ainda.</div>
          )}
        </div>
      </Card>
    </div>
  );
}
