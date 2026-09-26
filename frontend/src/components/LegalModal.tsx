import { useEffect } from 'react';
import { X, FileText, ShieldCheck } from 'lucide-react';

export type LegalDoc = 'termos' | 'privacidade';

interface Props {
  open: LegalDoc | null;
  onClose: () => void;
  onAccept?: () => void;
}

const CONTEUDO: Record<LegalDoc, { titulo: string; icon: React.ReactNode; secoes: { h: string; p: string }[] }> = {
  termos: {
    titulo: 'Termos de Uso',
    icon: <FileText size={18} />,
    secoes: [
      { h: '1. Sobre a plataforma', p: 'O PET-Saúde Clima é o portal institucional de monitoramento territorial e comunicação em saúde e clima do Piauí, voltado à divulgação de materiais educativos, oficinas e indicadores do projeto.' },
      { h: '2. Cadastro e acesso', p: 'O acesso à área restrita é destinado à equipe do projeto, orientadores e colaboradores autorizados. Cada usuário é responsável por manter a confidencialidade de sua senha e pelas atividades realizadas em sua conta.' },
      { h: '3. Uso adequado', p: 'É vedado utilizar a plataforma para fins ilícitos, divulgar desinformação em saúde ou clima, ou tentar acessar dados de outros usuários sem autorização.' },
      { h: '4. Conteúdo produzido', p: 'Materiais publicados (cards, cartilhas, vídeos, podcasts) seguem as diretrizes do Plano Transversal de Comunicação e podem ser reutilizados por instituições parceiras mediante crédito ao PET-Saúde Clima.' },
      { h: '5. Alterações', p: 'Estes termos podem ser atualizados conforme a evolução do projeto; alterações relevantes serão comunicadas aos usuários cadastrados.' },
    ],
  },
  privacidade: {
    titulo: 'Política de Privacidade',
    icon: <ShieldCheck size={18} />,
    secoes: [
      { h: '1. Dados coletados', p: 'Coletamos nome, e-mail, instituição e perfil informados no cadastro, além de dados de uso da plataforma (acessos, downloads e participação em oficinas), para fins estatísticos do projeto.' },
      { h: '2. Finalidade', p: 'Os dados são usados para autenticação, comunicação institucional, emissão de certificados e elaboração de relatórios de alcance exigidos pela SESAPI e demais instituições parceiras.' },
      { h: '3. Compartilhamento', p: 'Não vendemos ou compartilhamos dados pessoais com terceiros para fins comerciais. Dados agregados e anonimizados podem ser usados em relatórios públicos do projeto.' },
      { h: '4. Segurança', p: 'Adotamos medidas técnicas razoáveis para proteger as informações armazenadas, incluindo controle de acesso e criptografia de senhas.' },
      { h: '5. Direitos do titular', p: 'Você pode solicitar a atualização ou exclusão dos seus dados a qualquer momento entrando em contato com a coordenação do PET-Saúde Clima.' },
    ],
  },
};

export default function LegalModal({ open, onClose, onAccept }: Props) {
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    if (open) document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [open, onClose]);

  if (!open) return null;
  const doc = CONTEUDO[open];

  return (
    <div
      onClick={onClose}
      style={{
        position: 'fixed', inset: 0, zIndex: 100,
        background: 'rgba(15, 21, 19, 0.55)',
        backdropFilter: 'blur(4px)', WebkitBackdropFilter: 'blur(4px)',
        display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 20,
      }}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        className="glass-panel-login"
        style={{
          width: '100%', maxWidth: 560, maxHeight: '85vh', borderRadius: 20,
          display: 'flex', flexDirection: 'column', overflow: 'hidden',
          animation: 'slideDown 0.2s ease-out', background: 'var(--bg-card)',
        }}
      >
        <div style={{
          display: 'flex', alignItems: 'center', justifyContent: 'space-between',
          padding: '18px 24px', borderBottom: '1px solid var(--border)',
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <span style={{
              width: 34, height: 34, borderRadius: 10, background: 'var(--primary-light)',
              color: 'var(--primary)', display: 'flex', alignItems: 'center', justifyContent: 'center',
            }}>{doc.icon}</span>
            <h2 style={{ fontSize: 17, fontWeight: 700, color: 'var(--text-primary)' }}>{doc.titulo}</h2>
          </div>
          <button
            onClick={onClose}
            aria-label="Fechar"
            style={{
              background: 'transparent', border: 'none', cursor: 'pointer',
              color: 'var(--text-muted)', padding: 6, borderRadius: 8,
            }}
          >
            <X size={20} />
          </button>
        </div>

        <div style={{ padding: '20px 24px', overflowY: 'auto', flex: 1 }}>
          <p style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 16 }}>
            PET-Saúde Clima — Territórios do Piauí. Última atualização: setembro de 2026.
          </p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            {doc.secoes.map((s) => (
              <div key={s.h}>
                <h3 style={{ fontSize: 13, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>{s.h}</h3>
                <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>{s.p}</p>
              </div>
            ))}
          </div>
        </div>

        <div style={{
          padding: '16px 24px', borderTop: '1px solid var(--border)',
          display: 'flex', justifyContent: 'flex-end', gap: 10,
        }}>
          <button onClick={onClose} className="btn-outline" style={{ fontSize: 13, padding: '10px 18px' }}>
            Fechar
          </button>
          {onAccept && (
            <button
              onClick={() => { onAccept(); onClose(); }}
              className="btn-primary"
              style={{ fontSize: 13, padding: '10px 18px' }}
            >
              Li e concordo
            </button>
          )}
        </div>
      </div>
    </div>
  );
}