import { Archive, Landmark, CheckCircle2 } from 'lucide-react';
import Card from '../components/Card';
import { StatusBadge } from '../components/ProgressBar';
import { useState } from 'react';
import { Download, Search } from 'lucide-react';
import { assetUrl, baixarMaterial, listarMateriais, mensagemDeErro } from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro, campo } from '../lib/ui';

const itensRepositorio = [
  { nome: 'Portal PET-Saúde Clima Piauí — Versão 2.0', status: 'planejado' as const },
  { nome: 'Kit de Comunicação Regional completo', status: 'andamento' as const },
  { nome: 'Pacote de Mídias Educativas (30 materiais)', status: 'planejado' as const },
  { nome: 'Cursos autoinstrucionais do AVA', status: 'planejado' as const },
  { nome: 'Relatório Final de Educação Permanente', status: 'planejado' as const },
];

function Biblioteca() {
  const [q, setQ] = useState('');
  const [tipo, setTipo] = useState('');
  const [erro, setErro] = useState('');
  const { data, carregando, erro: erroLista, recarregar } = useRecurso(listarMateriais, { q, tipo, limit: 100 });

  async function baixar(id: number) {
    setErro('');
    try {
      const r = await baixarMaterial(id);
      window.open(assetUrl(r.arquivo_url), '_blank', 'noopener');
      recarregar();
    } catch (e) {
      setErro(mensagemDeErro(e));
    }
  }

  return (
    <Card title="Biblioteca de materiais" icon={<Archive size={16} />}>
      <div style={{ display: 'flex', gap: 10, marginBottom: 12 }}>
        <div style={{ position: 'relative', flex: 1 }}>
          <Search size={14} style={{ position: 'absolute', left: 10, top: 10, color: 'var(--text-muted)' }} />
          <input value={q} onChange={(e) => setQ(e.target.value)} placeholder="Buscar por título ou tema" style={{ ...campo, paddingLeft: 30 }} />
        </div>
        <select value={tipo} onChange={(e) => setTipo(e.target.value)} style={{ ...campo, width: 'auto' }}>
          <option value="">Todos os tipos</option>
          {['cartilha', 'protocolo', 'guia', 'video', 'podcast', 'card', 'infografico', 'radio', 'pdf'].map((t) => <option key={t} value={t}>{t}</option>)}
        </select>
      </div>
      {(erro || erroLista) && <div style={{ ...alertaErro, marginBottom: 10 }}>{erro || erroLista}</div>}
      {(data ?? []).map((m, i) => (
        <div key={m.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 12, padding: '10px 0', borderTop: i === 0 ? 'none' : '1px solid var(--border)' }}>
          <div>
            <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>{m.titulo}</div>
            <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{m.tipo} · {m.tema ?? '—'} · {m.macrorregiao_nome ?? 'Estadual'}</div>
          </div>
          <button onClick={() => baixar(m.id)} style={{ display: 'flex', alignItems: 'center', gap: 5, background: 'none', border: '1px solid var(--primary)', color: 'var(--primary)', borderRadius: 999, padding: '5px 12px', fontSize: 11.5, fontWeight: 600, cursor: 'pointer' }}>
            <Download size={12} /> {m.downloads}
          </button>
        </div>
      ))}
      {!carregando && (data ?? []).length === 0 && <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>Nenhum material encontrado.</div>}
    </Card>
  );
}

export default function Repositorio() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
          <Archive size={16} color="var(--primary)" />
          <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>
            Repositório Aberto de Tecnologias Sociais e Digitais
          </span>
        </div>
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Entrega prevista para o Semestre 4 (Mar/2028 – Ago/2028): a plataforma finalizada será
          transferida à SESAPI para manutenção e uso institucional continuado, junto ao Relatório
          Final de Educação Permanente, com a meta de atender 12 municípios e as 4 macrorregiões
          do Piauí.
        </p>
      </Card>

      <Card title="Itens do repositório final" icon={<Archive size={16} />} noPad>
        <div>
          {itensRepositorio.map((item, i) => (
            <div key={item.nome} style={{
              display: 'flex', alignItems: 'center', justifyContent: 'space-between',
              padding: '12px 16px', borderTop: i === 0 ? 'none' : '1px solid var(--border)',
            }}>
              <span style={{ fontSize: 13, color: 'var(--text-primary)' }}>{item.nome}</span>
              <StatusBadge status={item.status} />
            </div>
          ))}
        </div>
      </Card>

      <Card title="Transferência institucional">
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <div style={{ width: 44, height: 44, borderRadius: 10, background: 'var(--primary-light)', display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0 }}>
            <Landmark size={22} color="var(--primary)" />
          </div>
          <div>
            <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>SESAPI — Secretaria de Saúde do Estado do Piauí</div>
            <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>Órgão institucional responsável pela manutenção continuada do repositório</div>
          </div>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginTop: 16, fontSize: 12.5, color: 'var(--text-secondary)' }}>
          <CheckCircle2 size={15} color="var(--text-muted)" />
          Checklist de transferência será liberado ao final do Semestre 3, junto ao relatório de alcance.
        </div>
      </Card>
      <Biblioteca />
    </div>
  );
}
