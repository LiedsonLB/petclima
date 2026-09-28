import { useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import {
  Lock, Eye, EyeOff, ArrowRight, ArrowLeft, AlertCircle,
  CheckCircle2, ShieldCheck, KeyRound,
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

export default function RedefinirSenha() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { redefinirSenha } = useAuth();

  const email = searchParams.get('email') ?? '';
  const token = searchParams.get('token') ?? '';

  const [senha, setSenha] = useState('');
  const [confirmarSenha, setConfirmarSenha] = useState('');
  const [showSenha, setShowSenha] = useState(false);
  const [showConfirmar, setShowConfirmar] = useState(false);
  const [erro, setErro] = useState('');
  const [sucesso, setSucesso] = useState('');
  const [carregando, setCarregando] = useState(false);

  // Regras simples de força da senha
  const temMinimo = senha.length >= 6;
  const temMaiuscula = /[A-Z]/.test(senha);
  const temNumero = /\d/.test(senha);
  const senhasIguais = senha && senha === confirmarSenha;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErro('');
    setSucesso('');

    if (!email || !token) {
      setErro('Link de redefinição inválido ou expirado.');
      return;
    }

    if (!temMinimo) {
      setErro('A senha deve ter pelo menos 6 caracteres.');
      return;
    }

    if (senha !== confirmarSenha) {
      setErro('As senhas não coincidem.');
      return;
    }

    try {
      setCarregando(true);
      await redefinirSenha(email, token, senha);
      setSucesso('Senha redefinida com sucesso! Redirecionando...');
      setTimeout(() => navigate('/login'), 1800);
    } catch (err: any) {
      setErro(err?.message || 'Não foi possível redefinir a senha.');
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div
      className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden hero-gradient"
      style={{ background: 'var(--bg)', color: 'var(--text-primary)' }}
    >
      {/* Background decorativo */}
      <div className="absolute inset-0 z-0">
        <div className="absolute -top-32 -right-32 w-96 h-96 rounded-full bg-[#004e47]/20 blur-3xl" />
        <div className="absolute -bottom-32 -left-32 w-96 h-96 rounded-full bg-[#a1f1e5]/20 blur-3xl" />
      </div>

      <div className="relative z-10 w-full max-w-md">
        <div className="glass-panel-login rounded-3xl p-8 lg:p-10 space-y-7">
          {/* Cabeçalho */}
          <div className="text-center">
            <div className="flex items-center justify-center mb-3">
              <div className="w-14 h-14 rounded-2xl bg-[#004e47]/10 flex items-center justify-center">
                <KeyRound size={28} className="text-[#004e47]" />
              </div>
            </div>
            <h1 className="text-[26px] font-bold text-[#181c1c] mb-2">
              Definir nova senha
            </h1>
            <p className="text-sm text-[#3e4947] leading-relaxed">
              Crie uma nova senha segura para acessar sua conta do PET-Saúde
              Clima.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Nova senha */}
            <div className="space-y-1.5">
              <label
                htmlFor="senha"
                className="block text-xs font-semibold text-[#3e4947] mb-1.5"
              >
                Nova senha
              </label>
              <div className="relative">
                <Lock
                  size={18}
                  className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2"
                />
                <input
                  id="senha"
                  className="auth-input w-full pl-10 pr-10"
                  type={showSenha ? 'text' : 'password'}
                  required
                  placeholder="••••••••"
                  value={senha}
                  onChange={(e) => setSenha(e.target.value)}
                  autoComplete="new-password"
                />
                <button
                  type="button"
                  onClick={() => setShowSenha((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[#3e4947]/60 hover:text-[#181c1c] transition-colors"
                  aria-label={showSenha ? 'Ocultar senha' : 'Mostrar senha'}
                >
                  {showSenha ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Confirmar senha */}
            <div className="space-y-1.5">
              <label
                htmlFor="confirmarSenha"
                className="block text-xs font-semibold text-[#3e4947] mb-1.5"
              >
                Confirmar senha
              </label>
              <div className="relative">
                <Lock
                  size={18}
                  className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2"
                />
                <input
                  id="confirmarSenha"
                  className="auth-input w-full pl-10 pr-10"
                  type={showConfirmar ? 'text' : 'password'}
                  required
                  placeholder="••••••••"
                  value={confirmarSenha}
                  onChange={(e) => setConfirmarSenha(e.target.value)}
                  autoComplete="new-password"
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmar((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[#3e4947]/60 hover:text-[#181c1c] transition-colors"
                  aria-label={showConfirmar ? 'Ocultar senha' : 'Mostrar senha'}
                >
                  {showConfirmar ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Checklist de requisitos */}
            {(senha || confirmarSenha) && (
              <div className="grid grid-cols-1 gap-1.5 pt-1">
                <Requisito ok={temMinimo} texto="Pelo menos 6 caracteres" />
                <Requisito ok={temMaiuscula} texto="Uma letra maiúscula" />
                <Requisito ok={temNumero} texto="Um número" />
                <Requisito ok={senhasIguais} texto="As senhas coincidem" />
              </div>
            )}

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

            {sucesso && (
              <div
                className="flex items-center gap-2 text-xs px-3 py-2.5 rounded-lg"
                style={{
                  background: 'var(--success-bg)',
                  color: 'var(--success-text)',
                }}
              >
                <CheckCircle2 size={14} className="shrink-0" />
                {sucesso}
              </div>
            )}

            <button
              type="submit"
              disabled={carregando || !!sucesso}
              className="btn-primary w-full flex items-center justify-center gap-2 text-[15px]"
            >
              {carregando ? (
                <>
                  <span className="inline-block w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Redefinindo...
                </>
              ) : (
                <>
                  Redefinir senha
                  <ArrowRight size={18} />
                </>
              )}
            </button>
          </form>

          {/* Aviso de segurança */}
          <div className="flex items-start gap-2 p-3 rounded-lg bg-[#004e47]/5 border border-[#004e47]/10">
            <ShieldCheck
              size={16}
              className="text-[#004e47] shrink-0 mt-0.5"
            />
            <p className="text-[11px] text-[#3e4947] leading-relaxed">
              Por segurança, este link é de uso único e expira em poucos
              minutos. Nunca compartilhe sua senha com terceiros.
            </p>
          </div>

          <Link
            to="/login"
            className="flex items-center justify-center gap-1.5 text-sm font-semibold text-[#004e47] hover:underline transition-all"
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

/* Sub-componente: item de requisito com check */
function Requisito({ ok, texto }: { ok: boolean; texto: string }) {
  return (
    <div className="flex items-center gap-2 text-[11px]">
      <span
        className="inline-flex items-center justify-center w-4 h-4 rounded-full transition-colors"
        style={{
          background: ok ? 'var(--success-bg)' : 'rgba(0,0,0,0.05)',
          color: ok ? 'var(--success-text)' : 'var(--text-muted)',
        }}
      >
        {ok ? <CheckCircle2 size={10} /> : <span className="w-1 h-1 rounded-full bg-current" />}
      </span>
      <span style={{ color: ok ? 'var(--success-text)' : 'var(--text-muted)' }}>
        {texto}
      </span>
    </div>
  );
}