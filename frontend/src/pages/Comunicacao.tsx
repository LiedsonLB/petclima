import { useState } from 'react';
import { LayoutTemplate, Mic, Video, Filter, Plus, Download } from 'lucide-react';
import Card from '../components/Card';
import Modal from '../components/Modal';
import ProgressBar from '../components/ProgressBar';
import { useAuth } from '../contexts/AuthContext';
import {
  assetUrl, baixarMaterial, criarMaterial, ehStaff, enviarArquivoMaterial, listarMacrorregioes,
  listarMateriais, mensagemDeErro, type TipoMaterial,
} from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro, botaoPrimario, campo, dataBR, rotulo } from '../lib/ui';

const tipoIcon: Partial<Record<TipoMaterial, React.ReactNode>> = {
  podcast: <Mic size={14} />, radio: <Mic size={14} />, video: <Video size={14} />,
};
const tipoLabel: Record<string, string> = {
  card: 'Card', podcast: 'Podcast', video: 'Vídeo curto',
  cartilha: 'Cartilha', infografico: 'Infográfico', radio: 'Programa de rádio',
};
const TIPOS_KIT = 'card,podcast,video';

export default function Comunicacao() {
  const { usuario } = useAuth();
  const [filtroRegiao, setFiltroRegiao] = useState('todas');
  const [novo, setNovo] = useState(false);
  const [erro, setErro] = useState('');
  const macros = useRecurso(listarMacrorregioes);
  const kit = useRecurso(listarMateriais, { tipo: TIPOS_KIT, limit: 100 });

  const todos = kit.data ?? [];
  const lista = filtroRegiao === 'todas' ? todos : todos.filter((m) => String(m.macrorregiao_id) === filtroRegiao);

  async function baixar(id: number) {
    setErro('');
    try {
      const r = await baixarMaterial(id);
      window.open(assetUrl(r.arquivo_url), '_blank', 'noopener');
      kit.recarregar();
    } catch (e) {
      setErro(mensagemDeErro(e));
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card title="Plano Transversal de Comunicação em Saúde e Clima">
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6, marginBottom: 16 }}>
          Estratégia de mídia, linguagem e identidade visual do projeto, adaptada às 4 macrorregiões
          do Piauí, com foco no enfrentamento da desinformação (fake news) e na promoção do
          letramento em saúde.
        </p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
          {(macros.data ?? []).map((m) => (
            <div key={m.id} style={{ padding: 12, border: '1px solid var(--border)', borderRadius: 10 }}>
              <div style={{ fontSize: 12.5, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 6 }}>{m.nome}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 6 }}>{m.municipios_ativos ?? 0} / {m.meta_municipios} municípios</div>
              <ProgressBar value={m.municipios_ativos ?? 0} max={m.meta_municipios || 1} />
            </div>
          ))}
        </div>
      </Card>

      <Card
        title="Kit de Comunicação Regional"
        icon={<LayoutTemplate size={16} />}
        action={ehStaff(usuario?.papel) && (
          <button style={botaoPrimario} onClick={() => setNovo(true)}><Plus size={13} /> Novo material</button>
        )}
      >
        <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 4 }}>
          Meta do Semestre 1: 10 materiais iniciais (cards, podcasts e vídeos curtos) sobre prevenção de riscos climáticos.
        </div>
        <div style={{ marginBottom: 14 }}><ProgressBar value={todos.length} max={10} /></div>

        {(erro || kit.erro) && <div style={{ ...alertaErro, marginBottom: 12 }}>{erro || kit.erro}</div>}

        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 14 }}>
          <Filter size={13} color="var(--text-muted)" />
          <select value={filtroRegiao} onChange={(e) => setFiltroRegiao(e.target.value)} style={{ ...campo, width: 'auto' }}>
            <option value="todas">Todas as macrorregiões</option>
            {(macros.data ?? []).map((m) => <option key={m.id} value={m.id}>{m.nome}</option>)}
          </select>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 12 }}>
          {lista.map((m) => (
            <div key={m.id} style={{ border: '1px solid var(--border)', borderRadius: 10, padding: 12 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 6, color: 'var(--primary)', marginBottom: 8 }}>
                {tipoIcon[m.tipo] ?? <LayoutTemplate size={14} />}
                <span style={{ fontSize: 11, fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.04em' }}>{tipoLabel[m.tipo] ?? m.tipo}</span>
              </div>
              <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)', marginBottom: 6 }}>{m.titulo}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 10 }}>{m.macrorregiao_nome ?? 'Estadual'} · {m.tema ?? '—'}</div>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: 11, color: 'var(--text-muted)' }}>
                <button onClick={() => baixar(m.id)} style={{ display: 'flex', alignItems: 'center', gap: 4, background: 'none', border: 'none', cursor: 'pointer', color: 'var(--primary)', fontSize: 11, fontWeight: 600, padding: 0 }}>
                  <Download size={12} /> {m.downloads}
                </button>
                <span>{dataBR(m.created_at)}</span>
              </div>
            </div>
          ))}
          {!kit.carregando && lista.length === 0 && (
            <div style={{ gridColumn: '1 / -1', fontSize: 12, color: 'var(--text-muted)', padding: 20, textAlign: 'center' }}>
              Nenhum material para essa macrorregião ainda.
            </div>
          )}
        </div>
      </Card>

      {novo && (
        <NovoMaterial
          macros={macros.data ?? []}
          onClose={() => setNovo(false)}
          onSalvo={() => { setNovo(false); kit.recarregar(); }}
        />
      )}
    </div>
  );
}

function NovoMaterial({ macros, onClose, onSalvo }: {
  macros: { id: number; nome: string }[]; onClose: () => void; onSalvo: () => void;
}) {
  const [f, setF] = useState({ titulo: '', tipo: 'card' as TipoMaterial, tema: '', macro: '' });
  const [arquivo, setArquivo] = useState<File | null>(null);
  const [erro, setErro] = useState('');
  const [salvando, setSalvando] = useState(false);

  async function salvar(e: React.FormEvent) {
    e.preventDefault();
    if (!arquivo) { setErro('Selecione o arquivo do material.'); return; }
    setErro(''); setSalvando(true);
    try {
      const up = await enviarArquivoMaterial(arquivo);
      await criarMaterial({
        titulo: f.titulo, tipo: f.tipo, tema: f.tema, arquivo_url: up.url,
        macrorregiao_id: f.macro ? Number(f.macro) : null,
      });
      onSalvo();
    } catch (err) {
      setErro(mensagemDeErro(err));
    } finally {
      setSalvando(false);
    }
  }

  return (
    <Modal titulo="Novo material" onClose={onClose}>
      <form onSubmit={salvar} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        <div><label style={rotulo}>Título</label><input required style={campo} value={f.titulo} onChange={(e) => setF({ ...f, titulo: e.target.value })} /></div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
          <div>
            <label style={rotulo}>Tipo</label>
            <select style={campo} value={f.tipo} onChange={(e) => setF({ ...f, tipo: e.target.value as TipoMaterial })}>
              {Object.entries(tipoLabel).map(([v, l]) => <option key={v} value={v}>{l}</option>)}
            </select>
          </div>
          <div>
            <label style={rotulo}>Macrorregião</label>
            <select style={campo} value={f.macro} onChange={(e) => setF({ ...f, macro: e.target.value })}>
              <option value="">Estadual</option>
              {macros.map((m) => <option key={m.id} value={m.id}>{m.nome}</option>)}
            </select>
          </div>
        </div>
        <div><label style={rotulo}>Tema</label><input style={campo} value={f.tema} onChange={(e) => setF({ ...f, tema: e.target.value })} placeholder="Ex.: Calor extremo" /></div>
        <div><label style={rotulo}>Arquivo (até 50MB)</label><input type="file" onChange={(e) => setArquivo(e.target.files?.[0] ?? null)} /></div>
        {erro && <div style={alertaErro}>{erro}</div>}
        <button type="submit" disabled={salvando} style={{ ...botaoPrimario, justifyContent: 'center' }}>{salvando ? 'Enviando…' : 'Publicar'}</button>
      </form>
    </Modal>
  );
}
