import { useState } from 'react';
import { LayoutTemplate, Mic, Video, Filter, Plus, Download } from 'lucide-react';
import Card from '../components/Card';
import ProgressBar from '../components/ProgressBar';
import { materiais, macrorregioes, type TipoMaterial } from '../data/Data';

const tipoIcon: Record<TipoMaterial, React.ReactNode> = {
  card: <LayoutTemplate size={14} />,
  podcast: <Mic size={14} />,
  video: <Video size={14} />,
  cartilha: <LayoutTemplate size={14} />,
  infografico: <LayoutTemplate size={14} />,
  radio: <Mic size={14} />,
};

const tipoLabel: Record<TipoMaterial, string> = {
  card: 'Card', podcast: 'Podcast', video: 'Vídeo curto',
  cartilha: 'Cartilha', infografico: 'Infográfico', radio: 'Programa de rádio',
};

export default function Comunicacao() {
  const [filtroRegiao, setFiltroRegiao] = useState<string>('todas');
  const kitInicial = materiais.filter(m => ['card', 'podcast', 'video'].includes(m.tipo));
  const lista = filtroRegiao === 'todas' ? kitInicial : kitInicial.filter(m => m.macrorregiao === filtroRegiao);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card title="Plano Transversal de Comunicação em Saúde e Clima">
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6, marginBottom: 16 }}>
          Estratégia de mídia, linguagem e identidade visual do projeto, adaptada às 4 macrorregiões
          do Piauí, com foco no enfrentamento da desinformação (fake news) e na promoção do
          letramento em saúde.
        </p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
          {macrorregioes.map(m => (
            <div key={m.nome} style={{ padding: 12, border: '1px solid var(--border)', borderRadius: 10 }}>
              <div style={{ fontSize: 12.5, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 6 }}>{m.nome}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 6 }}>{m.municipiosAtendidos} / {m.meta} municípios</div>
              <ProgressBar value={m.municipiosAtendidos} max={m.meta} />
            </div>
          ))}
        </div>
      </Card>

      <Card
        title="Kit de Comunicação Regional"
        icon={<LayoutTemplate size={16} />}
        action={
          <button style={{
            display: 'flex', alignItems: 'center', gap: 6, fontSize: 12, fontWeight: 600,
            color: '#fff', background: 'var(--primary)', border: 'none', borderRadius: 999,
            padding: '6px 14px', cursor: 'pointer',
          }}>
            <Plus size={13} /> Novo material
          </button>
        }
      >
        <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 4 }}>
          Meta do Semestre 1: 10 materiais iniciais (cards, podcasts e vídeos curtos) sobre prevenção de riscos climáticos.
        </div>
        <div style={{ marginBottom: 14 }}>
          <ProgressBar value={kitInicial.length} max={10} />
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 14 }}>
          <Filter size={13} color="var(--text-muted)" />
          <select
            value={filtroRegiao}
            onChange={e => setFiltroRegiao(e.target.value)}
            style={{
              fontSize: 12, padding: '5px 10px', borderRadius: 8, border: '1px solid var(--border)',
              background: 'var(--bg-input)', color: 'var(--text-primary)',
            }}
          >
            <option value="todas">Todas as macrorregiões</option>
            {macrorregioes.map(m => <option key={m.nome} value={m.nome}>{m.nome}</option>)}
          </select>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 12 }}>
          {lista.map(m => (
            <div key={m.id} style={{ border: '1px solid var(--border)', borderRadius: 10, padding: 12 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 6, color: 'var(--primary)', marginBottom: 8 }}>
                {tipoIcon[m.tipo]}
                <span style={{ fontSize: 11, fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.04em' }}>{tipoLabel[m.tipo]}</span>
              </div>
              <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)', marginBottom: 6 }}>{m.titulo}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 10 }}>{m.macrorregiao} · {m.tema}</div>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: 11, color: 'var(--text-muted)' }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Download size={12} /> {m.downloads}</span>
                <span>{new Date(m.data).toLocaleDateString('pt-BR')}</span>
              </div>
            </div>
          ))}
          {lista.length === 0 && (
            <div style={{ gridColumn: '1 / -1', fontSize: 12, color: 'var(--text-muted)', padding: 20, textAlign: 'center' }}>
              Nenhum material para essa macrorregião ainda.
            </div>
          )}
        </div>
      </Card>
    </div>
  );
}
