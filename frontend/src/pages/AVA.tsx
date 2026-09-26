import { BookOpenCheck, Clock, Layers } from 'lucide-react';
import Card from '../components/Card';
import { cursosAVA } from '../data/Data';

export default function AVA() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
          <BookOpenCheck size={16} color="var(--primary)" />
          <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>Plataforma Digital Consolidada — Versão 2.0</span>
        </div>
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Previsto para o Semestre 3 (Set/2027 – Fev/2028): integração de um Ambiente Virtual de
          Aprendizagem (AVA) à plataforma, com módulos de cursos autoinstrucionais sobre Clima e
          Saúde para profissionais e estudantes das 4 macrorregiões.
        </p>
      </Card>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
        {cursosAVA.map(c => (
          <Card key={c.id}>
            <div style={{
              display: 'inline-block', fontSize: 10.5, fontWeight: 700, textTransform: 'uppercase',
              letterSpacing: '0.05em', color: 'var(--text-muted)', background: 'var(--bg-chip)',
              padding: '3px 8px', borderRadius: 999, marginBottom: 10,
            }}>
              {c.status}
            </div>
            <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 10, lineHeight: 1.35 }}>
              {c.titulo}
            </div>
            <div style={{ display: 'flex', gap: 14, fontSize: 11.5, color: 'var(--text-muted)' }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Layers size={13} /> {c.modulos} módulos</span>
              <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}><Clock size={13} /> {c.cargaHoraria}h</span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
