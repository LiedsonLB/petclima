import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Mail, ArrowRight, ArrowLeft, MailCheck } from 'lucide-react';

export default function EsqueciSenha() {
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [enviado, setEnviado] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      setEnviado(true);
    }, 1200);
  };

  return (
    <div className="hero-gradient" style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 20 }}>
      <div
        className="glass-panel-login"
        style={{ borderRadius: 24, padding: 40, width: '100%', maxWidth: 420, background: 'var(--bg-card)' }}
      >
        {!enviado ? (
          <>
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', marginBottom: 24, textAlign: 'center' }}>
              <img src="petsaudeclima_icon.png" alt="PET-Saúde Clima Logo" style={{ height: 48, marginBottom: 14 }} />
              <div style={{ fontSize: 18, fontWeight: 700, color: 'var(--text-primary)' }}>Recuperar senha</div>
              <p style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 6, lineHeight: 1.5 }}>
                Informe o e-mail institucional usado no cadastro. Enviaremos um link para você redefinir sua senha do PET-Saúde Clima.
              </p>
            </div>

            <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
              <div className="input-group" style={{ position: 'relative' }}>
                <span className="input-icon-wrapper"><Mail size={16} /></span>
                <input
                  className="login-input" type="email" required
                  placeholder="E-mail institucional"
                  value={email} onChange={(e) => setEmail(e.target.value)}
                />
              </div>

              <button type="submit" disabled={isLoading} className="btn-login">
                {isLoading ? (
                  <>
                    <span className="animate-spin" style={{ width: 16, height: 16, border: '2px solid rgba(255,255,255,0.3)', borderTopColor: '#fff', borderRadius: '50%' }} />
                    Enviando...
                  </>
                ) : (
                  <>
                    Enviar link de recuperação <ArrowRight size={16} className="arrow" />
                  </>
                )}
              </button>
            </form>
          </>
        ) : (
          <div style={{ textAlign: 'center' }}>
            <div style={{ width: 56, height: 56, borderRadius: '50%', background: 'var(--success-bg)', color: 'var(--success-text)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
              <MailCheck size={26} />
            </div>
            <h2 style={{ fontSize: 17, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 8 }}>Verifique seu e-mail</h2>
            <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
              Se <strong>{email}</strong> estiver cadastrado, você receberá em instantes um link para redefinir sua senha.
            </p>
          </div>
        )}

        <Link
          to="/login"
          className="link-institucional"
          style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 6, fontSize: 13, marginTop: 24 }}
        >
          <ArrowLeft size={14} /> Voltar para o login
        </Link>
      </div>
    </div>
  );
}