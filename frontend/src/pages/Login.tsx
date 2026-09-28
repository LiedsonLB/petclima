import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Cloud, Mail, Lock, Eye, EyeOff, ArrowRight,
  Thermometer, ClipboardCheck, AlertCircle,
} from 'lucide-react';
import LegalModal, { type LegalDoc } from '../components/LegalModal';
import { useAuth } from '../contexts/AuthContext';
import { ApiError } from '../lib/api';

export default function Login() {
  const navigate = useNavigate();
  const { entrar } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [remember, setRemember] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [erro, setErro] = useState('');
  const [legalOpen, setLegalOpen] = useState<LegalDoc | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErro('');
    setIsLoading(true);
    try {
      await entrar(email, password, remember);
      navigate('/dashboard');
    } catch (err) {
      if (err instanceof ApiError) {
        setErro(err.status === 401 ? 'E-mail ou senha inválidos.' : err.message);
      } else {
        setErro('Não foi possível conectar ao servidor. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const togglePassword = () => setShowPassword((v) => !v);

  return (
    <div className="min-h-screen flex items-stretch overflow-hidden relative" style={{ background: 'var(--bg)', color: 'var(--text-primary)' }}>
      {/* Full Page Background */}
      <div className="absolute inset-0 z-0">
        <img
          className="w-full h-full object-cover bg-blur-image"
          alt="Ambiente de pesquisa com iluminação suave em azul e verde"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuDmGpQ4kwBpIy1nRfLfah44jGZtd3BLl0_ipTfbxcYu55Nv1_iPTd7u9MVG_Tez3_TWWLYfHYU79n182ZS90FRtc_oR4jz28WcDk0wi6CLCD9xG1u8oxKiYJXegojd5fS2AfRdoJM6kXAnAaAdOHxdOF-Tluc8-MbKu-hLq4PXlQMbSPSvXjlRejJKNMG8_sioAwtBjEKFqH-Dy6qqveBRCT7x5BlBIW7yaWJ8w485rM9Hkk_bQFUK-Dw"
        />
        <div className="hero-bg-overlay" />
      </div>

      {/* Main Content Container */}
      <div className="relative z-10 w-full min-h-screen flex flex-col lg:flex-row items-center justify-between p-4 lg:p-12 max-w-[1600px] mx-auto gap-12">
        {/* Left Side: Cinematic Narrative & Branding */}
        <section className="w-full lg:w-1/2 flex flex-col justify-center items-start pt-12 lg:pt-0">
          <div className="max-w-xl">
            <div className="flex items-center gap-4 mb-6">
              <div className="">
                <div className="flex items-center justify-center gap-1">
                  <img
                    src="petsaudeclima_icon.png"
                    alt="PET-Saúde Clima Logo"
                    className="h-16 w-auto"
                  />
                </div>
              </div>
              <h1 className="text-[48px] font-bold text-white tracking-tight drop-shadow-lg">PET-Saúde Clima</h1>
            </div>
            <p className="text-[18px] text-white/90 mb-12 drop-shadow-md leading-relaxed">
              Integrando vigilância em saúde e monitoramento climático para proteger as comunidades do Piauí através da ciência e tecnologia de ponta.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="p-6 glass-card-dark-login rounded-2xl hover:bg-[#004e47]/70 transition-colors">
                <Thermometer size={28} className="text-[#a1f1e5] block mb-1" />
                <span className="text-[12px] uppercase tracking-widest text-white/70 font-medium">Dados de Clima</span>
                <p className="text-[20px] font-semibold text-white mt-1">Monitoramento Real</p>
              </div>
              <div className="p-6 glass-card-dark-login rounded-2xl hover:bg-[#004e47]/70 transition-colors">
                <ClipboardCheck size={28} className="text-[#a1f1e5] block mb-1" />
                <span className="text-[12px] uppercase tracking-widest text-white/70 font-medium">Saúde Pública</span>
                <p className="text-[20px] font-semibold text-white mt-1">Vigilância Piauí</p>
              </div>
            </div>
          </div>
        </section>

        {/* Right Side: Login Interface */}
        <main className="w-full lg:w-[480px] flex items-center justify-center shrink-0">
          <div className="w-full glass-panel-login rounded-3xl p-8 lg:p-10 space-y-8">
            {/* Header & Logo */}
            <div className="text-center">
              <div className="flex items-center justify-center gap-1 mb-2">
                <img
                  src="petsaudeclima_icon.png"
                  alt="PET-Saúde Clima Logo"
                  className="h-10 w-auto"
                />
              </div>
              <h2 className="text-[28px] font-bold text-[#181c1c] mb-2">Bem-vindo de volta</h2>
              <p className="text-[#3e4947] text-sm">Acesse o portal institucional de monitoramento territorial.</p>
            </div>

            {/* Login Form */}
            <form className="space-y-6" onSubmit={handleSubmit}>
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-[#3e4947] mb-1.5" htmlFor="email">
                  E-mail Institucional
                </label>
                <div className="relative">
                  <Mail size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    className="auth-input w-full pl-10"
                    id="email"
                    name="email"
                    type="email"
                    required
                    placeholder="nome@instituicao.org"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <div className="flex justify-between items-center mb-1.5">
                  <label className="block text-xs font-semibold text-[#3e4947]" htmlFor="password">
                    Senha
                  </label>
                  <Link to="/esqueci-senha" className="text-xs font-semibold text-[#004e47] hover:text-[#00685f] hover:underline transition-all">
                    Esqueceu a senha?
                  </Link>
                </div>
                <div className="relative">
                  <Lock size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    className="auth-input w-full pl-10 pr-10"
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    required
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  <button
                    type="button"
                    onClick={togglePassword}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[#3e4947]/60 hover:text-[#181c1c] transition-colors"
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              <div className="flex items-center gap-2">
                <input
                  className="login-checkbox"
                  id="remember"
                  type="checkbox"
                  checked={remember}
                  onChange={(e) => setRemember(e.target.checked)}
                />
                <label className="text-sm text-[#3e4947] cursor-pointer select-none" htmlFor="remember">
                  Lembrar-me por 30 dias
                </label>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="btn-primary w-full flex items-center justify-center gap-2 text-[15px]"
              >
                {isLoading ? (
                  <>
                    <span className="inline-block w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Acessando...
                  </>
                ) : (
                  <>
                    Entrar na Plataforma
                    <ArrowRight size={18} />
                  </>
                )}
              </button>

              {erro && (
                <div className="flex items-center gap-2 text-xs px-3 py-2.5 rounded-lg" style={{ background: 'var(--danger-bg)', color: 'var(--danger-text)' }}>
                  <AlertCircle size={14} className="shrink-0" />
                  {erro}
                </div>
              )}
            </form>

            <p className="text-center text-sm text-[#3e4947] pt-2">
              Ainda não tem acesso? <br className="sm:hidden" />
              <Link to="/cadastro" className="font-bold text-[#004e47] hover:underline">
                Solicitar cadastro institucional
              </Link>
            </p>

            <p className="text-center text-xs text-[#3e4947]">
              Ao continuar, você concorda com os{' '}
              <button type="button" onClick={() => setLegalOpen('termos')} className="underline hover:text-[#004e47]">
                Termos de Uso
              </button>{' '}
              e a{' '}
              <button type="button" onClick={() => setLegalOpen('privacidade')} className="underline hover:text-[#004e47]">
                Política de Privacidade
              </button>.
            </p>
          </div>
        </main>
      </div>

      <LegalModal open={legalOpen} onClose={() => setLegalOpen(null)} />
    </div>
  );
}