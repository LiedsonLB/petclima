import { useTheme } from '../contexts/ThemeContext';
import Card from '../components/Card';

export default function Configuracoes() {
  const { theme, toggle } = useTheme();

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20, maxWidth: 560 }}>
      <Card title="Aparência">
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>Tema</div>
            <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>Alternar entre modo claro e escuro</div>
          </div>
          <button
            onClick={toggle}
            style={{
              fontSize: 12, fontWeight: 600, color: '#fff', background: 'var(--primary)',
              border: 'none', borderRadius: 999, padding: '7px 16px', cursor: 'pointer',
            }}
          >
            {theme === 'light' ? 'Ativar escuro' : 'Ativar claro'}
          </button>
        </div>
      </Card>

      <Card title="Sobre o PET-Saúde Clima">
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Grupo de Atuação Tutorial responsável por Comunicação, Tecnologias Digitais e Educação
          Permanente dentro do PET-Saúde Clima Piauí, com foco em produção de tecnologias sociais e
          digitais, enfrentamento da desinformação e formação contínua de profissionais de saúde.
        </p>
      </Card>
    </div>
  );
}
