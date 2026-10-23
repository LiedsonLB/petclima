import { useState } from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft, BookOpenCheck, Clock, Layers, Award } from 'lucide-react';
import Card from '../components/Card';
import ProgressBar from '../components/ProgressBar';
import { atualizarProgresso, listarCursos, matricular, mensagemDeErro, type Curso } from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro, botaoContorno, botaoPrimario } from '../lib/ui';

export default function AVA() {
  const cursos = useRecurso(listarCursos);
  const [erro, setErro] = useState('');
  const [ocupado, setOcupado] = useState<number | null>(null);

  async function agir(c: Curso, acao: () => Promise<unknown>) {
    setErro(''); setOcupado(c.id);
    try {
      await acao();
      cursos.recarregar();
    } catch (e) {
      setErro(mensagemDeErro(e));
    } finally {
      setOcupado(null);
    }
  }

  // Sem controle por aula no backend, cada clique conclui "o próximo módulo".
  function proximoModulo(c: Curso) {
    const atual = Math.round(((c.minha_matricula?.progresso ?? 0) / 100) * c.modulos_total);
    const prox = Math.min(c.modulos_total, atual + 1);
    return Math.round((prox / c.modulos_total) * 100);
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Link to="/formacao" style={{ ...botaoContorno, alignSelf: 'flex-start', textDecoration: 'none' }}>
        <ArrowLeft size={13} /> Voltar para Educação Permanente
      </Link>

      <Card>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
          <BookOpenCheck size={16} color="var(--primary)" />
          <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>Ambiente Virtual de Aprendizagem</span>
        </div>
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Cursos autoinstrucionais sobre Clima e Saúde para profissionais e estudantes das 4
          macrorregiões. Matricule-se, acompanhe seu progresso e emita o certificado ao concluir.
        </p>
      </Card>

      {(erro || cursos.erro) && <div style={alertaErro}>{erro || cursos.erro}</div>}

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
        {(cursos.data ?? []).map((c) => {
          const m = c.minha_matricula;
          const publicado = c.status === 'publicado';
          return (
            <Card key={c.id}>
              <div style={{ display: 'inline-block', fontSize: 10.5, fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', background: 'var(--bg-chip)', padding: '2px 8px', borderRadius: 6, marginBottom: 10 }}>
                {c.categoria ?? 'Curso'} · {publicado ? 'disponível' : c.status}
              </div>
              <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 6 }}>{c.titulo}</div>
              {c.descricao && <div style={{ fontSize: 12, color: 'var(--text-secondary)', lineHeight: 1.5, marginBottom: 10 }}>{c.descricao}</div>}
              <div style={{ display: 'flex', gap: 14, fontSize: 11.5, color: 'var(--text-muted)', marginBottom: 12 }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Clock size={12} /> {c.carga_horaria}h</span>
                <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Layers size={12} /> {c.modulos_total} módulos</span>
              </div>

              {m ? (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                  <ProgressBar value={m.progresso} max={100} />
                  <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{m.progresso}% concluído</div>
                  {m.status === 'concluido' ? (
                    <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12, color: 'var(--success-text)' }}>
                      <Award size={14} /> Certificado: <strong>{m.certificado_codigo}</strong>
                    </div>
                  ) : (
                    <button style={botaoContorno} disabled={ocupado === c.id} onClick={() => agir(c, () => atualizarProgresso(c.id, proximoModulo(c)))}>
                      Concluir próximo módulo
                    </button>
                  )}
                </div>
              ) : (
                <button style={botaoPrimario} disabled={!publicado || ocupado === c.id} onClick={() => agir(c, () => matricular(c.id))}>
                  {publicado ? 'Matricular-me' : 'Em breve'}
                </button>
              )}
            </Card>
          );
        })}
      </div>
    </div>
  );
}
