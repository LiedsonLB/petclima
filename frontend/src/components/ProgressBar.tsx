interface Props {
  value: number;
  max: number;
  color?: string;
}

export default function ProgressBar({ value, max, color = 'var(--primary)' }: Props) {
  const pct = Math.min(100, Math.round((value / max) * 100));
  return (
    <div style={{ width: '100%', height: 8, background: 'var(--bg-input)', borderRadius: 999, overflow: 'hidden' }}>
      <div style={{ width: `${pct}%`, height: '100%', background: color, borderRadius: 999, transition: 'width 0.4s ease' }} />
    </div>
  );
}

export function StatusBadge({ status }: { status: 'concluido' | 'andamento' | 'planejado' }) {
  const map = {
    concluido: { bg: 'var(--success-bg)', text: 'var(--success-text)', label: 'Concluído' },
    andamento: { bg: 'var(--warning-bg)', text: 'var(--warning-text)', label: 'Em andamento' },
    planejado: { bg: 'var(--bg-chip)', text: 'var(--text-muted)', label: 'Planejado' },
  } as const;
  const s = map[status];
  return (
    <span style={{
      background: s.bg, color: s.text, fontSize: 11, fontWeight: 600,
      padding: '3px 10px', borderRadius: 999, whiteSpace: 'nowrap',
    }}>
      {s.label}
    </span>
  );
}
