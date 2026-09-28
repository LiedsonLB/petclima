import { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Mail, ArrowRight, ArrowLeft, MailCheck, AlertCircle, KeyRound,
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { ApiError } from '../lib/api';

export default function EsqueciSenha() {
  const { esqueciSenha } = useAuth();
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [enviado, setEnviado] = useState(false);
  const [erro, setErro] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErro('');
    setIsLoading(true);
    try {
      await esqueciSenha(email);
      setEnviado(true);
    } catch (err) {
      if (err instanceof ApiError) {
        setErro(err.message);
      } else {
        setErro('Não foi possível conectar ao servidor. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden hero-gradient"
      style={{ background: 'var(--bg)', color: 'var(--text-primary)' }}
    >
      {/* Background decorativo */}
      <div className="absolute inset-0 z-0">
        <div className="absolute -top-32 -left-32 w-96 h-96 rounded-full bg-[#004e47]/20 blur-3xl" />
        <div className="absolute -bottom-32 -right-32 w-96 h-96 rounded-full bg-[#a1f1e5]/20 blur-3xl" />
      </div>

      <div className="relative z-10 w-full max-w-md">
        <div className="glass-panel-login rounded-3xl p-8 lg:p-10 space-y-8">
          {!enviado ? (
            <>
              {/* Cabeçalho */}
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 mb-3">
                  <img
                    src="petsaudeclima_icon.png"
                    alt="PET-Saúde Clima Logo"
                    className="h-12 w-auto"
                  />
                </div>
                <h1 className="text-[26px] font-bold text-[#181c1c] mb-2">
                  Recuperar senha
                </h1>
                <p className="text-sm text-[#3e4947] leading-relaxed">
                  Informe o e-mail institucional usado no cadastro. Enviaremos
                  um link seguro para você redefinir sua senha.
                </p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="space-y-1.5">
                  <label
                    htmlFor="email"
                    className="block text-xs font-semibold text-[#3e4947] mb-1.5"
                  >
                    E-mail institucional
                  </label>
                  <div className="relative">
                    <Mail
                      size={18}
                      className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2"
                    />
                    <input
                      id="email"
                      className="auth-input w-full pl-10"
                      type="email"
                      required
                      placeholder="nome@instituicao.org"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      autoComplete="email"
                    />
                  </div>
                </div>

                {erro && (
                  <div
                    className="flex items-center gap-2 text-xs px-3 py-2.5 rounded-lg"
                    style={{
                      background: 'var(--danger-bg)',
                      color: 'var(--danger-text)',
                    }}
                  >
                    <AlertCircle size={14} className="shrink-0" />
                    {erro}
                  </div>
                )}

                <button
                  type="submit"
                  disabled={isLoading}
                  className="btn-primary w-full flex items-center justify-center gap-2 text-[15px]"
                >
                  {isLoading ? (
                    <>
                      <span className="inline-block w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      Enviando...
                    </>
                  ) : (
                    <>
                      Enviar link de recuperação
                      <ArrowRight size={18} />
                    </>
                  )}
                </button>
              </form>
            </>
          ) : (
            /* Estado de sucesso */
            <div className="text-center py-4">
              <div
                className="w-16 h-16 rounded-full mx-auto mb-5 flex items-center justify-center"
                style={{
                  background: 'var(--success-bg)',
                  color: 'var(--success-text)',
                }}
              >
                <MailCheck size={30} />
              </div>
              <h2 className="text-[22px] font-bold text-[#181c1c] mb-3">
                Verifique seu e-mail
              </h2>
              <p className="text-sm text-[#3e4947] leading-relaxed">
                Se <strong className="text-[#004e47]">{email}</strong> estiver
                cadastrado, você receberá em instantes um link para redefinir
                sua senha.
              </p>
              <p className="text-xs text-[#9aaba7] mt-4">
                Não recebeu? Verifique a caixa de spam ou tente novamente em
                alguns minutos.
              </p>
            </div>
          )}

          <Link
            to="/login"
            className="flex items-center justify-center gap-1.5 text-sm font-semibold text-[#004e47] hover:underline transition-all pt-2"
          >
            <ArrowLeft size={14} /> Voltar para o login
          </Link>
        </div>

        <p className="text-center text-xs text-white/60 mt-6">
          © 2026 PET-Saúde Clima — Territórios do Piauí
        </p>
      </div>
    </div>
  );
}