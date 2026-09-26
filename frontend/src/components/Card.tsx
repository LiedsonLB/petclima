import type { ReactNode } from 'react';

interface Props {
  title?: string;
  icon?: ReactNode;
  action?: ReactNode;
  children: ReactNode;
  style?: React.CSSProperties;
  noPad?: boolean;
}

export default function Card({ title, icon, action, children, style, noPad }: Props) {
  return (
    <div style={{
      background: 'var(--bg-card)',
      border: '1px solid var(--border)',
      borderRadius: 12,
      overflow: 'hidden',
      ...style,
    }}>
      {title && (
        <div style={{
          padding: '11px 16px',
          borderBottom: '1px solid var(--border)',
          display: 'flex', alignItems: 'center', gap: 8,
        }}>
          {icon && <span style={{ color: 'var(--primary)', display: 'flex' }}>{icon}</span>}
          <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)', flex: 1 }}>{title}</span>
          {action}
        </div>
      )}
      <div style={noPad ? {} : { padding: 16 }}>
        {children}
      </div>
    </div>
  );
}
