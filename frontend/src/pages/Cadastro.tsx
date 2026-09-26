import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Cloud, Mail, Lock, Eye, EyeOff, ArrowRight,
  User, Building2, GraduationCap, Users, ClipboardCheck, ChevronDown,
} from 'lucide-react';
import LegalModal, { type LegalDoc } from '../components/LegalModal';

export default function Cadastro() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [acceptTerms, setAcceptTerms] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [erro, setErro] = useState('');
  const [legalOpen, setLegalOpen] = useState<LegalDoc | null>(null);
  const [formData, setFormData] = useState({
    nome: '',
    email: '',
    instituicao: '',
    curso: '',
    perfil: '',
    senha: '',
    confirmSenha: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErro('');
    if (!acceptTerms) {
      setErro('Você deve aceitar os Termos de Uso e a Política de Privacidade.');
      return;
    }
    if (formData.senha !== formData.confirmSenha) {
      setErro('As senhas não coincidem.');
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      navigate('/login');
    }, 1500);
  };

  return (
    <div className="min-h-screen flex overflow-hidden relative" style={{ background: 'var(--bg)', color: 'var(--text-primary)' }}>
      {/* Painel esquerdo */}
      <section className="hidden lg:flex lg:w-1/2 xl:w-3/5 relative overflow-hidden p-12 flex-col justify-between hero-gradient">
        <img
          src="petsaudeclima_hero.png"
          alt="Paisagem do Piauí"
          className="auth-hero-image"
        />
        <div className="hero-overlay" />
        <div className="hero-overlay-bottom" />

        <div className="auth-hero-content">
          <div className="flex items-center gap-3">
            <div className="w-11 h-11 rounded-xl bg-white/10 backdrop-blur-sm flex items-center justify-center border border-white/20">
              <Cloud size={28} className="text-white" />
            </div>
            <div>
              <div className="text-[22px] font-bold text-white tracking-tight">PET-Saúde Clima</div>
              <div className="text-[12px] text-white/70 uppercase tracking-[0.1em]">Territórios do Piauí</div>
            </div>
          </div>
        </div>

        <div className="auth-hero-content max-w-md">
          <h2 className="text-[34px] font-bold text-white leading-tight mb-4">
            Dados que informam. Ciência que transforma.
          </h2>
          <p className="text-[16px] text-white/85 leading-relaxed font-light">
            Junte-se à maior rede de monitoramento de saúde ambiental do Piauí. Contribua para pesquisas que salvam vidas através da análise climática.
          </p>
        </div>

        <div className="auth-hero-content flex gap-8">
          <div>
            <div className="text-[24px] font-bold text-primary">224</div>
            <div className="text-[12px] text-primary/70">Municípios Monitorados</div>
          </div>
          <div className="auth-stats-divider" />
          <div>
            <div className="text-[24px] font-bold text-primary">15+</div>
            <div className="text-[12px] text-primary/70">Instituições Parceiras</div>
          </div>
        </div>
      </section>

      {/* Painel direito: formulário */}
      <main className="flex-1 flex flex-col items-center justify-around p-6 mesh-gradient min-h-screen">
        <div className="w-full max-w-lg">
          <div className="flex items-center gap-2 mb-6 lg:hidden">
            <div className="w-9 h-9 rounded-xl bg-[#004e47]/10 flex items-center justify-center">
              <Cloud size={24} className="text-[#004e47]" />
            </div>
            <div>
              <div className="text-[18px] font-bold text-[#004e47]">PET-Saúde Clima</div>
              <div className="text-[10px] text-[#6e7977] uppercase tracking-wider">Monitoramento Socioambiental</div>
            </div>
          </div>

          <div className="flex items-center gap-4 mb-6">
            <img
              src="petsaudeclima_icon.png"
              alt="PET-Saúde Clima Logo"
              className="h-16 w-auto"
            />
            <div>
              <h2 className="text-[28px] font-bold text-[#004e47] mb-1">Criar Conta</h2>
              <p className="text-[#3e4947] text-[15px] font-light">Preencha os dados abaixo para acessar a plataforma.</p>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Nome completo</label>
              <div className="relative">
                <User size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                <input
                  name="nome"
                  type="text"
                  required
                  placeholder="Seu nome completo"
                  value={formData.nome}
                  onChange={handleChange}
                  className="auth-input w-full pl-10"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">E-mail</label>
              <div className="relative">
                <Mail size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                <input
                  name="email"
                  type="email"
                  required
                  placeholder="seu@email.com"
                  value={formData.email}
                  onChange={handleChange}
                  className="auth-input w-full pl-10"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Instituição</label>
                <div className="relative">
                  <Building2 size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    name="instituicao"
                    type="text"
                    required
                    placeholder="Ex: UFPI, UESPI..."
                    value={formData.instituicao}
                    onChange={handleChange}
                    className="auth-input w-full pl-10"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Curso</label>
                <div className="relative">
                  <GraduationCap size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    name="curso"
                    type="text"
                    required
                    placeholder="Ex: Medicina, Geografia..."
                    value={formData.curso}
                    onChange={handleChange}
                    className="auth-input w-full pl-10"
                  />
                </div>
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Perfil</label>
              <div className="relative">
                <Users size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                <select
                  name="perfil"
                  required
                  value={formData.perfil}
                  onChange={handleChange}
                  className="auth-input w-full pl-10 appearance-none"
                >
                  <option value="" disabled>Selecione seu perfil</option>
                  <option value="pesquisador">Pesquisador</option>
                  <option value="estudante">Estudante / Bolsista</option>
                  <option value="gestor">Gestor Público</option>
                  <option value="comunidade">Membro da Comunidade</option>
                </select>
                <ChevronDown size={18} className="text-[#3e4947]/60 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Senha</label>
                <div className="relative">
                  <Lock size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    name="senha"
                    type={showPassword ? 'text' : 'password'}
                    required
                    placeholder="••••••••"
                    value={formData.senha}
                    onChange={handleChange}
                    className="auth-input w-full pl-10 pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[#3e4947]/60 hover:text-[#181c1c] transition-colors"
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-[#3e4947] mb-1.5">Confirmar senha</label>
                <div className="relative">
                  <ClipboardCheck size={18} className="text-[#3e4947]/60 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    name="confirmSenha"
                    type={showConfirmPassword ? 'text' : 'password'}
                    required
                    placeholder="••••••••"
                    value={formData.confirmSenha}
                    onChange={handleChange}
                    className="auth-input w-full pl-10 pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[#3e4947]/60 hover:text-[#181c1c] transition-colors"
                  >
                    {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>
            </div>

            <label className="flex items-start gap-2.5 cursor-pointer">
              <input
                type="checkbox"
                checked={acceptTerms}
                onChange={(e) => setAcceptTerms(e.target.checked)}
                className="accent-primary mt-0.5"
              />
              <span className="text-xs text-[#3e4947]">
                Eu li e concordo com os{' '}
                <button type="button" onClick={() => setLegalOpen('termos')} className="text-[#004e47] font-semibold hover:underline">
                  Termos de Uso
                </button>{' '}
                e a{' '}
                <button type="button" onClick={() => setLegalOpen('privacidade')} className="text-[#004e47] font-semibold hover:underline">
                  Política de Privacidade
                </button>{' '}
                do projeto PET-Saúde Clima.
              </span>
            </label>

            {erro && (
              <div className="text-xs px-3 py-2.5 rounded-lg" style={{ background: 'var(--danger-bg)', color: 'var(--danger-text)' }}>
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
                  Processando...
                </>
              ) : (
                <>
                  Finalizar Cadastro
                  <ArrowRight size={18} />
                </>
              )}
            </button>

            <p className="text-center text-sm text-[#3e4947] mt-4">
              Já possui uma conta?{' '}
              <Link to="/login" className="font-bold text-[#004e47] hover:underline">
                Entrar agora
              </Link>
            </p>
          </form>
        </div>
        <footer className="mt-6 pt-4 border-t border-[#e0e8e5]/30 flex flex-col md:flex-row justify-between items-center gap-4 opacity-60">
          <span className="text-xs text-[#3e4947]">© 2026 PET-Saúde Clima</span>
          <div className="flex gap-4">
            <Link to="#" className="text-xs text-[#3e4947] hover:text-[#004e47] transition-colors">Suporte</Link>
            <Link to="#" className="text-xs text-[#3e4947] hover:text-[#004e47] transition-colors">Sobre o Projeto</Link>
          </div>
        </footer>
      </main>

      <LegalModal open={legalOpen} onClose={() => setLegalOpen(null)} onAccept={() => setAcceptTerms(true)} />
    </div>
  );
}