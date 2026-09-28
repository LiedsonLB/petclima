import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { FileDown } from 'lucide-react';
import Card from '../components/Card';
import { engajamentoMensal, alcancePorMacrorregiao, materiais } from '../data/Data';

export default function Relatorios() {
  const totalDownloads = materiais.reduce((acc, m) => acc + m.downloads, 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <Card>
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Relatório de Alcance de Mídia e Engajamento — consolida indicadores de tráfego, acessos
          regionais e downloads de materiais por município, previsto para o Semestre 3 do projeto.
        </p>
      </Card>

      <div style={{ display: 'grid', gridTemplateColumns: '1.3fr 1fr', gap: 20 }}>
        <Card title="Acessos e downloads ao longo do tempo">
          <div style={{ height: 260 }}>
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={engajamentoMensal}>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
                <XAxis dataKey="mes" tick={{ fontSize: 11 }} stroke="var(--text-muted)" />
                <YAxis tick={{ fontSize: 11 }} stroke="var(--text-muted)" />
                <Tooltip contentStyle={{ fontSize: 12, borderRadius: 8 }} />
                <Line type="monotone" dataKey="acessos" stroke="#004e47" strokeWidth={2} dot={{ r: 3 }} />
                <Line type="monotone" dataKey="downloads" stroke="#eda100" strokeWidth={2} dot={{ r: 3 }} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card title="Acessos por macrorregião">
          <div style={{ height: 260 }}>
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={alcancePorMacrorregiao}>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
                <XAxis dataKey="macrorregiao" tick={{ fontSize: 10 }} stroke="var(--text-muted)" />
                <YAxis tick={{ fontSize: 11 }} stroke="var(--text-muted)" />
                <Tooltip contentStyle={{ fontSize: 12, borderRadius: 8 }} />
                <Bar dataKey="acessos" fill="#006c49" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </div>

      <Card
        title="Downloads por material"
        action={
          <button style={{
            display: 'flex', alignItems: 'center', gap: 6, fontSize: 12, fontWeight: 600,
            color: 'var(--primary)', background: 'transparent', border: '1px solid var(--primary)',
            borderRadius: 999, padding: '6px 14px', cursor: 'pointer',
          }}>
            <FileDown size={13} /> Exportar relatório
          </button>
        }
        noPad
      >
        <div>
          {materiais.map((m, i) => (
            <div key={m.id} style={{
              display: 'flex', justifyContent: 'space-between', alignItems: 'center',
              padding: '10px 16px', borderTop: i === 0 ? 'none' : '1px solid var(--border)', fontSize: 13,
            }}>
              <span style={{ color: 'var(--text-primary)' }}>{m.titulo}</span>
              <span style={{ color: 'var(--text-muted)' }}>{m.downloads} downloads</span>
            </div>
          ))}
          <div style={{ padding: '10px 16px', borderTop: '1px solid var(--border)', fontSize: 13, fontWeight: 700, color: 'var(--text-primary)', display: 'flex', justifyContent: 'space-between' }}>
            <span>Total</span><span>{totalDownloads}</span>
          </div>
        </div>
      </Card>
    </div>
  );
}
