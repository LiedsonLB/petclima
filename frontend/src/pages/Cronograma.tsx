import Card from '../components/Card';
import { StatusBadge } from '../components/ProgressBar';
import { cronograma } from '../data/Data';

export default function Cronograma() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20, maxWidth: 900 }}>
      {cronograma.map((sem, idx) => (
        <div key={sem.id} style={{ display: 'flex', gap: 16 }}>
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', paddingTop: 4 }}>
            <div style={{
              width: 32, height: 32, borderRadius: '50%',
              background: sem.status === 'andamento' ? 'var(--primary)' : 'var(--bg-chip)',
              color: sem.status === 'andamento' ? '#fff' : 'var(--text-muted)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 13, fontWeight: 700, flexShrink: 0,
            }}>
              {sem.id}
            </div>
            {idx < cronograma.length - 1 && (
              <div style={{ width: 2, flex: 1, background: 'var(--border)', marginTop: 4 }} />
            )}
          </div>

          <Card style={{ flex: 1, marginBottom: 4 }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 4 }}>
              <div style={{ fontSize: 15, fontWeight: 700, color: 'var(--text-primary)' }}>Semestre {sem.id}: {sem.titulo}</div>
              <StatusBadge status={sem.status} />
            </div>
            <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 14 }}>{sem.periodo}</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              {sem.entregas.map(e => (
                <div key={e.nome} style={{ display: 'flex', gap: 10, alignItems: 'flex-start', padding: '8px 0', borderTop: '1px solid var(--border)' }}>
                  <div style={{ flex: 1 }}>
                    <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>{e.nome}</div>
                    <div style={{ fontSize: 11.5, color: 'var(--text-muted)' }}>{e.detalhe}</div>
                  </div>
                  <StatusBadge status={e.status} />
                </div>
              ))}
            </div>
          </Card>
        </div>
      ))}

      <Card title="Suporte do Orientador de Serviço 3">
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          O GAT 5 conta com o apoio do Orientador de Serviço 3, responsável por dar suporte às
          metodologias de Educação Permanente, validar as tecnologias de informação e comunicação
          em saúde e acompanhar as ações de telessaúde ao longo das quatro etapas do cronograma.
        </p>
      </Card>
    </div>
  );
}
