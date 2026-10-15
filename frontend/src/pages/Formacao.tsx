import { useState } from 'react';
import { CalendarDays, Users, ShieldAlert, Plus } from 'lucide-react';
import Card from '../components/Card';
import KpiCard from '../components/KpiCard';
import Modal from '../components/Modal';
import ProgressBar from '../components/ProgressBar';
import SalasAoVivo from '../components/SalasAoVivo';
import { useAuth } from '../contexts/AuthContext';
import {
  cancelarInscricao, criarAtividade, ehStaff, inscreverAtividade, listarAtividades,
  listarMunicipios, mensagemDeErro, painelResumo, type Atividade,
} from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro, botaoContorno, botaoPrimario, campo, dataHoraBR, rotulo } from '../lib/ui';

const TIPOS = 'oficina,capacitacao';

export default function Formacao() {
  const { usuario } = useAuth();
  const [novo, setNovo] = useState(false);
  const [erro, setErro] = useState('');
  const [ocupado, setOcupado] = useState<number | null>(null);
  const atividades = useRecurso(listarAtividades, { tipo: TIPOS });
  const resumo = useRecurso(painelResumo);

  const lista = atividades.data ?? [];
  const inscritos = lista.reduce((acc, o) => acc + (o.inscritos ?? 0), 0);
  const comVagas = lista.filter((o) => o.vagas);
  const vagasTotais = comVagas.reduce((acc, o) => acc + (o.vagas ?? 0), 0);
  const inscritosComVagas = comVagas.reduce((acc, o) => acc + (o.inscritos ?? 0), 0);
  const capacitados = resumo.data?.profissionais_capacitados ?? 0;

  async function alternar(a: Atividade) {
    setErro(''); setOcupado(a.id);
    try {
      if (a.minha_inscricao) await cancelarInscricao(a.id);
      else await inscreverAtividade(a.id);
      atividades.recarregar();
    } catch (e) {
      setErro(mensagemDeErro(e));
    } finally {
      setOcupado(null);
    }
  }

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
        <KpiCard label="Inscritos até agora" value={`${inscritos}`} sub="meta: 200 profissionais/ACS" icon={<Users size={16} />} />
        <KpiCard label="Oficinas programadas" value={`${lista.length}`} sub="presenciais e remotas" icon={<CalendarDays size={16} />} />
        <KpiCard label="Ocupação média" value={vagasTotais ? `${Math.round((inscritosComVagas / vagasTotais) * 100)}%` : '—'} sub="das vagas oferecidas" icon={<ShieldAlert size={16} />} />
      </div>

      <Card>
        <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 6 }}>
          Progresso rumo à meta de 200 profissionais capacitados ({capacitados} com presença confirmada)
        </div>
        <ProgressBar value={capacitados} max={200} />
      </Card>

      <SalasAoVivo />

      <Card
        title="Oficinas e capacitações"
        action={ehStaff(usuario?.papel) && (
          <button style={botaoPrimario} onClick={() => setNovo(true)}><Plus size={13} /> Nova oficina</button>
        )}
      >
        {(erro || atividades.erro) && <div style={{ ...alertaErro, marginBottom: 12 }}>{erro || atividades.erro}</div>}
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          {lista.map((o, i) => {
            const cheia = !!o.vagas && (o.inscritos ?? 0) >= o.vagas;
            return (
              <div key={o.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 16, padding: '12px 0', borderTop: i === 0 ? 'none' : '1px solid var(--border)' }}>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>{o.titulo}</div>
                  <div style={{ fontSize: 11, color: 'var(--text-muted)', marginTop: 2 }}>
                    {o.municipio_nome ?? 'Remota'} · {dataHoraBR(o.data_realizacao)} · {o.modalidade} · {o.carga_horaria}h
                  </div>
                  {!!o.vagas && (
                    <div style={{ marginTop: 6, maxWidth: 260 }}>
                      <ProgressBar value={o.inscritos ?? 0} max={o.vagas} />
                      <div style={{ fontSize: 10.5, color: 'var(--text-muted)', marginTop: 3 }}>{o.inscritos ?? 0} / {o.vagas} vagas</div>
                    </div>
                  )}
                </div>
                <button
                  disabled={ocupado === o.id || (cheia && !o.minha_inscricao)}
                  onClick={() => alternar(o)}
                  style={o.minha_inscricao ? botaoContorno : botaoPrimario}
                >
                  {o.minha_inscricao ? 'Cancelar inscrição' : cheia ? 'Vagas esgotadas' : 'Inscrever-me'}
                </button>
              </div>
            );
          })}
          {!atividades.carregando && lista.length === 0 && (
            <div style={{ fontSize: 12, color: 'var(--text-muted)', padding: 16, textAlign: 'center' }}>Nenhuma oficina programada.</div>
          )}
        </div>
      </Card>

      {novo && <NovaOficina onClose={() => setNovo(false)} onSalvo={() => { setNovo(false); atividades.recarregar(); }} />}
    </div>
  );
}

function NovaOficina({ onClose, onSalvo }: { onClose: () => void; onSalvo: () => void }) {
  const municipios = useRecurso(listarMunicipios);
  const [f, setF] = useState({ titulo: '', tipo: 'oficina', modalidade: 'presencial', municipio: '', data: '', carga: '4', vagas: '' });
  const [erro, setErro] = useState('');
  const [salvando, setSalvando] = useState(false);

  async function salvar(e: React.FormEvent) {
    e.preventDefault();
    setErro(''); setSalvando(true);
    try {
      await criarAtividade({
        titulo: f.titulo, tipo: f.tipo, modalidade: f.modalidade,
        municipio_id: f.municipio ? Number(f.municipio) : null,
        data_realizacao: f.data ? new Date(f.data).toISOString() : undefined,
        carga_horaria: Number(f.carga) || 0,
        vagas: f.vagas ? Number(f.vagas) : null,
      });
      onSalvo();
    } catch (err) {
      setErro(mensagemDeErro(err));
    } finally {
      setSalvando(false);
    }
  }

  return (
    <Modal titulo="Nova oficina" onClose={onClose}>
      <form onSubmit={salvar} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        <div><label style={rotulo}>Título</label><input required style={campo} value={f.titulo} onChange={(e) => setF({ ...f, titulo: e.target.value })} /></div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
          <div>
            <label style={rotulo}>Tipo</label>
            <select style={campo} value={f.tipo} onChange={(e) => setF({ ...f, tipo: e.target.value })}>
              <option value="oficina">Oficina</option><option value="capacitacao">Capacitação</option>
            </select>
          </div>
          <div>
            <label style={rotulo}>Modalidade</label>
            <select style={campo} value={f.modalidade} onChange={(e) => setF({ ...f, modalidade: e.target.value })}>
              <option value="presencial">Presencial</option><option value="remota">Remota</option><option value="hibrida">Híbrida</option>
            </select>
          </div>
        </div>
        <div>
          <label style={rotulo}>Município</label>
          <select style={campo} value={f.municipio} onChange={(e) => setF({ ...f, municipio: e.target.value })}>
            <option value="">Remota / sem município</option>
            {(municipios.data ?? []).map((m) => <option key={m.id} value={m.id}>{m.nome}/{m.uf}</option>)}
          </select>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr', gap: 12 }}>
          <div><label style={rotulo}>Data e hora</label><input type="datetime-local" style={campo} value={f.data} onChange={(e) => setF({ ...f, data: e.target.value })} /></div>
          <div><label style={rotulo}>Carga (h)</label><input type="number" min={0} style={campo} value={f.carga} onChange={(e) => setF({ ...f, carga: e.target.value })} /></div>
          <div><label style={rotulo}>Vagas</label><input type="number" min={1} style={campo} value={f.vagas} onChange={(e) => setF({ ...f, vagas: e.target.value })} /></div>
        </div>
        {erro && <div style={alertaErro}>{erro}</div>}
        <button type="submit" disabled={salvando} style={{ ...botaoPrimario, justifyContent: 'center' }}>{salvando ? 'Salvando…' : 'Criar oficina'}</button>
      </form>
    </Modal>
  );
}
