import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import {
  MailCheck, AlertCircle, ArrowRight, ArrowLeft,
  Loader2, CheckCircle2, ShieldCheck,
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

type Estado = 'carregando' | 'sucesso' | 'erro';

export default function ConfirmarEmail() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { confirmarEmail } = useAuth();

  const email = searchParams.get('email') ?? '';
  const token = searchParams.get('token') ?? '';

  const [estado, setEstado] = useState<Estado>('carregando');
  const [mensagem, setMensagem] = useState('');

  useEffect(() => {
    if (!email || !token) {
      setEstado('erro');
      setMensagem('Link de confirmação inválido ou incompleto.');
      return;
    }

    let cancelado = false;

    async function confirmar() {
      try {
        await confirmarEmail(email, token);
        if (cancelado) return;
        setEstado('sucesso');
        setMensagem('E-mail confirmado com sucesso! Redirecionando...');
        setTimeout(() => navigate('/login'), 2200);
      } catch (err: any) {
        if (cancelado) return;
        setEstado('erro');
        setMensagem(
          err?.message ||
            'Não foi possível confirmar o e-mail. O link pode ter expirado.',
        );
      }
    }

    confirmar();
    return () => {
      cancelado = true;
    };
  }, [email, token, confirmarEmail, navigate]);

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
        <div className="glass-panel-login rounded-3xl p-8 lg:p-10 space-y-6">
          {/* Logo */}
          <div className="flex items-center justify-center">
            <img
              src="petsaudeclima_icon.png"
              alt="PET-Saúde Clima Logo"
              className="h-12 w-auto"
            />
          </div>

          {estado === 'carregando' && (
            <div className="text-center py-4">
              <div className="w-16 h-16 rounded-full mx-auto mb-5 flex items-center justify-center bg-[#004e47]/10">
                <Loader2
                  size={30}
                  className="text-[#004e47] animate-spin"
                />
              </div>
              <h1 className="text-[22px] font-bold text-[#181c1c] mb-2">
                Confirmando seu e-mail...
              </h1>
              <p className="text-sm text-[#3e4947] leading-relaxed">
                Estamos validando o token enviado para{' '}
                <strong className="text-[#004e47]">{email || 'seu e-mail'}</strong>.
                Isso leva só um instante.
              </p>
            </div>
          )}

          {estado === 'sucesso' && (
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
              <h1 className="text-[22px] font-bold text-[#181c1c] mb-2">
                E-mail confirmado!
              </h1>
              <p className="text-sm text-[#3e4947] leading-relaxed">
                Sua conta no PET-Saúde Clima está ativa. Você será
                redirecionado para o login em instantes.
              </p>
              <div className="mt-5 flex items-center justify-center gap-2 text-xs text-[#004e47]/70">
                <CheckCircle2 size={14} />
                <span>Conta verificada com sucesso</span>
              </div>
            </div>
          )}

          {estado === 'erro' && (
            <div className="text-center py-4">
              <div
                className="w-16 h-16 rounded-full mx-auto mb-5 flex items-center justify-center"
                style={{
                  background: 'var(--danger-bg)',
                  color: 'var(--danger-text)',
                }}
              >
                <AlertCircle size={30} />
              </div>
              <h1 className="text-[22px] font-bold text-[#181c1c] mb-2">
                Não foi possível confirmar
              </h1>
              <p className="text-sm text-[#3e4947] leading-relaxed">
                {mensagem}
              </p>
              <p className="text-xs text-[#9aaba7] mt-3">
                Links de confirmação expiram por segurança. Se necessário,
                solicite um novo cadastro.
              </p>
            </div>
          )}

          {/* Botões de ação */}
          <div className="space-y-3 pt-2">
            {estado === 'sucesso' && (
              <button
                onClick={() => navigate('/login')}
                className="btn-primary w-full flex items-center justify-center gap-2 text-[15px]"
              >
                Ir para o login
                <ArrowRight size={18} />
              </button>
            )}

            {estado === 'erro' && (
              <>
                <Link
                  to="/cadastro"
                  className="btn-primary w-full flex items-center justify-center gap-2 text-[15px]"
                >
                  Refazer cadastro
                  <ArrowRight size={18} />
                </Link>
                <Link
                  to="/login"
                  className="flex items-center justify-center gap-1.5 text-sm font-semibold text-[#004e47] hover:underline transition-all"
                >
                  <ArrowLeft size={14} /> Voltar para o login
                </Link>
              </>
            )}

            {estado === 'carregando' && (
              <Link
                to="/login"
                className="flex items-center justify-center gap-1.5 text-sm font-semibold text-[#004e47]/70 hover:underline transition-all"
              >
                <ArrowLeft size={14} /> Cancelar e voltar ao login
              </Link>
            )}
          </div>

          {/* Aviso de segurança */}
          {estado !== 'sucesso' && (
            <div className="flex items-start gap-2 p-3 rounded-lg bg-[#004e47]/5 border border-[#004e47]/10">
              <ShieldCheck
                size={16}
                className="text-[#004e47] shrink-0 mt-0.5"
              />
              <p className="text-[11px] text-[#3e4947] leading-relaxed">
                Nunca compartilhe este link com terceiros. Ele é único,
                pessoal e expira automaticamente.
              </p>
            </div>
          )}
        </div>

        <p className="text-center text-xs text-white/60 mt-6">
          © 2026 PET-Saúde Clima — Territórios do Piauí
        </p>
      </div>
    </div>
  );
}