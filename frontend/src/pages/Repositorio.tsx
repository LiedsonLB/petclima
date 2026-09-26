import { Archive, Landmark, CheckCircle2 } from 'lucide-react';
import Card from '../components/Card';
import { StatusBadge } from '../components/ProgressBar';

const itensRepositorio = [
  { nome: 'Portal PET-Saúde Clima Piauí — Versão 2.0', status: 'planejado' as const },
  { nome: 'Kit de Comunicação Regional completo', status: 'andamento' as const },
  { nome: 'Pacote de Mídias Educativas (30 materiais)', status: 'planejado' as const },
  { nome: 'Cursos autoinstrucionais do AVA', status: 'planejado' as const },
  { nome: 'Relatório Final de Educação Permanente', status: 'planejado' as const },
];

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
    </div>
  );
}
