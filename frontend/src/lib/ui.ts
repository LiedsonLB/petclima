import type { CSSProperties } from 'react';

// Estilos compartilhados dos formulários internos (modais, Configurações),
// usando as variáveis de tema já definidas em index.css.
export const campo: CSSProperties = {
  width: '100%',
  fontSize: 13,
  padding: '8px 10px',
  borderRadius: 8,
  border: '1px solid var(--border)',
  background: 'var(--bg-input)',
  color: 'var(--text-primary)',
  outline: 'none',
};

export const rotulo: CSSProperties = {
  display: 'block',
  fontSize: 11.5,
  fontWeight: 600,
  color: 'var(--text-secondary)',
  marginBottom: 5,
};

export const botaoPrimario: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  fontSize: 12,
  fontWeight: 600,
  color: '#fff',
  background: 'var(--primary)',
  border: 'none',
  borderRadius: 999,
  padding: '7px 16px',
  cursor: 'pointer',
};

export const botaoContorno: CSSProperties = {
  ...botaoPrimario,
  color: 'var(--primary)',
  background: 'transparent',
  border: '1px solid var(--primary)',
};

export const alertaErro: CSSProperties = {
  fontSize: 12,
  color: 'var(--danger-text)',
  background: 'var(--danger-bg)',
  padding: '8px 12px',
  borderRadius: 10,
};

export const alertaSucesso: CSSProperties = {
  fontSize: 12,
  color: 'var(--success-text)',
  background: 'var(--success-bg)',
  padding: '8px 12px',
  borderRadius: 10,
};

export const dataBR = (iso?: string | null) =>
  iso ? new Date(iso).toLocaleDateString('pt-BR', { timeZone: 'UTC' }) : '—';

export const dataHoraBR = (iso?: string | null) =>
  iso
    ? new Date(iso).toLocaleString('pt-BR', { timeZone: 'UTC', dateStyle: 'short', timeStyle: 'short' })
    : 'data a definir';
