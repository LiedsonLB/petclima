/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  darkMode: ['class', '[data-theme="dark"]'],
  theme: {
    extend: {
      // Mapeado 1:1 para as CSS vars já usadas em src/index.css, para que
      // dê pra migrar os componentes de style={{...}} para classes
      // Tailwind (bg-primary, text-muted, border-default etc.) sem perder
      // o tema claro/escuro que já existe hoje.
      colors: {
        primary: {
          DEFAULT: 'var(--primary)',
          hover: 'var(--primary-hover)',
          light: 'var(--primary-light)',
          dark: 'var(--primary-dark)',
        },
        surface: {
          DEFAULT: 'var(--bg)',
          card: 'var(--bg-card)',
          sidebar: 'var(--bg-sidebar)',
          topbar: 'var(--bg-topbar)',
          input: 'var(--bg-input)',
          hover: 'var(--bg-hover)',
          chip: 'var(--bg-chip)',
          'chip-active': 'var(--bg-chip-active)',
        },
        text: {
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted: 'var(--text-muted)',
          placeholder: 'var(--text-placeholder)',
        },
        border: {
          DEFAULT: 'var(--border)',
          strong: 'var(--border-strong)',
        },
        danger: { DEFAULT: 'var(--danger)', bg: 'var(--danger-bg)', text: 'var(--danger-text)' },
        warning: { DEFAULT: 'var(--warning)', bg: 'var(--warning-bg)', text: 'var(--warning-text)' },
        success: { DEFAULT: 'var(--success)', bg: 'var(--success-bg)', text: 'var(--success-text)' },
      },
      borderRadius: {
        DEFAULT: 'var(--radius)',
        lg: 'var(--radius-lg)',
      },
      boxShadow: {
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
      },
    },
  },
  plugins: [],
}
