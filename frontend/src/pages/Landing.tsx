import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  ArrowRight, Menu, X, Radio, GraduationCap,
  MapPinned, BarChart3, ShieldCheck, BookOpenCheck, Megaphone,
  Download,
  Monitor,
  HardDriveDownload,
} from 'lucide-react';
import { cronograma, macrorregioes } from '../data/Data';
import LegalModal, { type LegalDoc } from '../components/LegalModal';

/* ==================================================================== */
/* DOWNLOADS — mude só aqui quando lançar versão nova                    */
/* ==================================================================== */
/*
 * `url` precisa ser link DIRETO pro arquivo (abrir no navegador já baixa,
 * sem tela de confirmação). GitHub Releases funciona bem pra isso.
 * Não funciona: Google Drive, OneDrive, Dropbox, MEGA.
 */
const VERSAO = '1.0.0';

const DOWNLOADS = [
  {
    id: 'windows',
    label: 'Windows',
    detalhe: 'Windows 10 ou superior (64 bits)',
    formato: 'Instalador .exe',
    size: '72 MB',
    icon: Monitor,
    url: `https://github.com/LiedsonLB/petsaude/releases/download/v${VERSAO}/Resenha.${VERSAO}.exe`,
    recomendado: 'windows' as const,
  },
  {
    id: 'linux-deb',
    label: 'Linux — Debian / Ubuntu',
    detalhe: 'Pacote nativo para distros baseadas em Debian',
    formato: 'Pacote .deb (amd64)',
    size: '73 MB',
    icon: HardDriveDownload,
    url: `https://github.com/LiedsonLB/petsaude/releases/download/v${VERSAO}/Resenha_${VERSAO}_amd64.deb`,
    recomendado: 'linux' as const,
  },
  {
    id: 'linux-appimage',
    label: 'Linux — AppImage',
    detalhe: 'Portátil, roda em qualquer distribuição',
    formato: 'AppImage',
    size: '105 MB',
    icon: HardDriveDownload,
    url: `https://github.com/LiedsonLB/petsaude/releases/download/v${VERSAO}/Resenha-${VERSAO}.AppImage`,
    recomendado: 'linux' as const,
  },
] as const;

type Plataforma = 'windows' | 'linux' | 'outro';

function detectarPlataforma(): Plataforma {
  if (typeof navigator === 'undefined') return 'outro';
  const ua = navigator.userAgent;
  if (/Windows|Win32|Win64|WOW64/i.test(ua)) return 'windows';
  if (/Linux/i.test(ua) && !/Android/i.test(ua)) return 'linux';
  return 'outro';
}

/* ------------------------------------------------------------------ */
/* Modal de download — abre ao clicar em qualquer botão de download    */
/* ------------------------------------------------------------------ */
function ModalDownload({
  aberto,
  onFechar,
  plataforma,
}: {
  aberto: boolean;
  onFechar: () => void;
  plataforma: Plataforma;
}) {
  const [baixando, setBaixando] = useState<string | null>(null);

  // Fecha com ESC e trava scroll do body
  useEffect(() => {
    if (!aberto) return;

    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onFechar();
    };
    document.addEventListener('keydown', onKey);

    const overflowAnterior = document.body.style.overflow;
    document.body.style.overflow = 'hidden';

    return () => {
      document.removeEventListener('keydown', onKey);
      document.body.style.overflow = overflowAnterior;
    };
  }, [aberto, onFechar]);

  function handleBaixar(id: string) {
    setBaixando(id);
    window.setTimeout(() => setBaixando(null), 2500);
  }

  if (!aberto) return null;

  return (
    <div
      className="fixed inset-0 z-[999] flex items-end justify-center bg-black/70 p-0 backdrop-blur-sm sm:items-center sm:p-4"
      onClick={onFechar}
      role="dialog"
      aria-modal="true"
      aria-labelledby="titulo-download"
    >
      <div
        className="relative w-full max-w-lg overflow-hidden rounded-t-3xl border border-outline-variant/30 bg-surface-container-low shadow-2xl sm:rounded-3xl"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Cabeçalho */}
        <div className="flex items-start justify-between gap-4 border-b border-outline-variant/20 p-6 pb-5">
          <div>
            <h2
              id="titulo-download"
              className="text-headline-md text-on-surface"
            >
              Baixar o Resenha
            </h2>
            <p className="mt-1 text-body-md text-on-surface-variant">
              Escolha a versão para o seu sistema. Versão {VERSAO}.
            </p>
          </div>

          <button
            type="button"
            onClick={onFechar}
            aria-label="Fechar"
            className="shrink-0 rounded-full p-2 text-on-surface-variant transition-colors hover:bg-surface-container-highest hover:text-on-surface"
          >
            <X size={20} />
          </button>
        </div>

        {/* Lista de plataformas */}
        <div className="space-y-3 p-6 pt-5">
          {DOWNLOADS.map(
            ({ id, label, detalhe, formato, size, icon: Icon, url, recomendado }) => {
              const ehRecomendado = plataforma === recomendado;
              const estaBaixando = baixando === id;

              return (
                <a
                  key={id}
                  href={url}
                  onClick={() => handleBaixar(id)}
                  className="group flex items-center gap-4 rounded-2xl border border-outline-variant/20 bg-surface-container p-4 transition-all hover:border-primary/40 hover:bg-surface-container-high"
                >
                  <span className="inline-flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-primary-container/20 text-primary">
                    <Icon size={22} />
                  </span>

                  <span className="min-w-0 flex-1">
                    <span className="flex items-center gap-2">
                      <span className="text-label-md text-on-surface">
                        {label}
                      </span>
                      {ehRecomendado && (
                        <span className="rounded-full bg-tertiary/20 px-2 py-0.5 text-label-sm uppercase tracking-wide text-tertiary">
                          Recomendado
                        </span>
                      )}
                    </span>
                    <span className="mt-0.5 block text-label-sm text-on-surface-variant">
                      {detalhe}
                    </span>
                    <span className="mt-1 block text-label-sm text-on-surface-variant/60">
                      {formato} · {size}
                    </span>
                  </span>

                  <span className="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary transition-colors group-hover:bg-primary group-hover:text-on-primary">
                    <Download
                      size={18}
                      className={estaBaixando ? 'animate-pulse' : ''}
                    />
                  </span>
                </a>
              );
            },
          )}

          <p className="pt-2 text-center text-label-sm text-on-surface-variant/60">
            macOS chega em breve. Dúvidas? Fale com a gente no chat.
          </p>
        </div>
      </div>
    </div>
  );
}

export default function Landing() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const [legalOpen, setLegalOpen] = useState<LegalDoc | null>(null);
  const [modalAberto, setModalAberto] = useState(false);

  const [plataforma, setPlataforma] = useState<Plataforma | null>(null);

  useEffect(() => {
    setPlataforma(detectarPlataforma());
  }, []);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40);
    window.addEventListener('scroll', onScroll);
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  /** Texto auxiliar abaixo dos CTAs — detecta a plataforma do usuário */
  const textoPlataforma =
    plataforma === 'windows'
      ? `Detectamos Windows · versão ${VERSAO}`
      : plataforma === 'linux'
        ? `Detectamos Linux · versão ${VERSAO}`
        : `Disponível para Windows e Linux · versão ${VERSAO}`;

  const pilares = [
    { icon: Megaphone, title: 'Plano Transversal de Comunicação', desc: 'Estratégia de mídia, linguagem e identidade visual para as 4 macrorregiões do Piauí.' },
    { icon: Radio, title: 'Kit de Comunicação Regional', desc: 'Cards, podcasts e vídeos curtos sobre riscos climáticos, prontos para redes sociais e WhatsApp.' },
    { icon: GraduationCap, title: 'Educação Permanente', desc: 'Oficinas de letramento em saúde e combate às fake news para ACS e profissionais da rede.' },
    { icon: BookOpenCheck, title: 'AVA — Cursos Autoinstrucionais', desc: 'Módulos sobre clima e saúde pública, disponíveis a qualquer momento para a equipe da rede.' },
    { icon: MapPinned, title: 'Territorialização', desc: 'Acompanhamento por macrorregião (Meio-Norte, Semiárido, Cerrado e Entre Rios).' },
    { icon: BarChart3, title: 'Indicadores de Alcance', desc: 'Relatórios de tráfego, engajamento e downloads para prestação de contas à SESAPI.' },
  ];

  const etapaAtual = cronograma.find((s) => s.status === 'andamento') ?? cronograma[0];

  return (
    <div style={{ background: 'var(--bg)', color: 'var(--text-primary)', overflowX: 'hidden' }}>
      {/* Navbar */}
      <nav style={{ position: 'fixed', top: 16, left: 0, width: '100%', zIndex: 50, padding: '0 20px' }}>
        <div
          className={scrolled ? '' : 'glass-card'}
          style={{
            maxWidth: 1180, margin: '0 auto', borderRadius: 9999,
            padding: scrolled ? '10px 20px' : '8px 20px',
            display: 'flex', alignItems: 'center', justifyContent: 'space-between',
            background: scrolled ? 'rgba(255,255,255,0.92)' : undefined,
            boxShadow: scrolled ? 'var(--shadow-md)' : undefined,
            transition: 'all 0.3s ease',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <img src="petsaudeclima_icon.png" alt="PET-Saúde Clima Logo" style={{ height: 32 }} />
            <span style={{ fontSize: 13, fontWeight: 700, color: 'var(--primary)' }}>PET-Saúde Clima</span>
          </div>

          <div className="hidden md:flex items-center gap-7">
            <a href="#pilares" style={{ fontSize: 11, textTransform: 'uppercase', letterSpacing: 1, color: 'var(--text-secondary)' }}>Pilares</a>
            <a href="#cronograma" style={{ fontSize: 11, textTransform: 'uppercase', letterSpacing: 1, color: 'var(--text-secondary)' }}>Cronograma</a>
            <a href="#territorios" style={{ fontSize: 11, textTransform: 'uppercase', letterSpacing: 1, color: 'var(--text-secondary)' }}>Territórios</a>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <button onClick={() => navigate('/login')} className="btn-outline hidden sm:inline-flex" style={{ fontSize: 12, padding: '8px 16px' }}>
              Entrar
            </button>
            <button onClick={() => navigate('/cadastro')} className="btn-primary inline-flex items-center gap-1.5" style={{ fontSize: 12, padding: '8px 16px' }}>
              Cadastrar <ArrowRight size={14} />
            </button>
            <button onClick={() => setMenuOpen(v => !v)} style={{ background: 'transparent', border: 'none', cursor: 'pointer' }} className="md:hidden">
              {menuOpen ? <X size={22} /> : <Menu size={22} />}
            </button>
          </div>
        </div>

        {menuOpen && (
          <div className="glass-card md:hidden" style={{ maxWidth: 1180, margin: '8px auto 0', borderRadius: 16, padding: 16, display: 'flex', flexDirection: 'column', gap: 10, background: 'rgba(255,255,255,0.95)' }}>
            <a href="#pilares" onClick={() => setMenuOpen(false)} style={{ fontSize: 13, color: 'var(--text-secondary)' }}>Pilares</a>
            <a href="#cronograma" onClick={() => setMenuOpen(false)} style={{ fontSize: 13, color: 'var(--text-secondary)' }}>Cronograma</a>
            <a href="#territorios" onClick={() => setMenuOpen(false)} style={{ fontSize: 13, color: 'var(--text-secondary)' }}>Territórios</a>
            <button onClick={() => navigate('/login')} style={{ fontSize: 13, textAlign: 'left', background: 'none', border: 'none', color: 'var(--primary)', fontWeight: 600, cursor: 'pointer' }}>Entrar</button>
          </div>
        )}
      </nav>

      {/* Hero */}
      <header className="hero-gradient" style={{ paddingTop: 140, paddingBottom: 100, position: 'relative' }}>
        {/* Background Image */}
        <div className="absolute inset-0 z-0">
          <img
            alt="Piauí Landscape"
            className="w-full h-full object-cover transform transition-transform duration-[20s] ease-out hover:scale-110"
            src="https://ipiranganews.inf.br/wp-content/uploads/VIAGEM-4-03-08-21.jpg"
          />
          <div className="absolute inset-0 bg-gradient-to-r from-black/85 via-black/50 to-transparent mix-blend-multiply"></div>
          {/* <div className="absolute inset-0 bg-gradient-to-t from-[#f7faf8] via-transparent to-transparent"></div> */}
        </div>
        <div style={{ maxWidth: 900, margin: '0 auto', padding: '0 20px', textAlign: 'center', position: 'relative', zIndex: 10 }}>
          <span style={{ display: 'inline-block', fontSize: 11, fontWeight: 700, letterSpacing: 1.5, textTransform: 'uppercase', color: 'var(--primary-fixed)', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', padding: '6px 14px', borderRadius: 9999, marginBottom: 20 }}>
            Territórios do Piauí
          </span>
          <h1 style={{ fontSize: 'clamp(32px, 5vw, 52px)', fontWeight: 800, color: '#fff', lineHeight: 1.1, marginBottom: 20 }}>
            PET-Saúde Clima
          </h1>
          <p style={{ fontSize: 16, color: 'var(--text-white-85)', lineHeight: 1.6, maxWidth: 640, margin: '0 auto 32px' }}>
            O portal do GAT 5 (Comunicação, Tecnologias Digitais e Educação Permanente) reúne os materiais,
            oficinas e indicadores produzidos para levar comunicação de risco climático às quatro macrorregiões do Piauí.
          </p>
          <div style={{ display: 'flex', gap: 14, justifyContent: 'center', flexWrap: 'wrap' }}>
            <button onClick={() => navigate('/cadastro')} className="btn-primary" style={{ fontSize: 14, padding: '14px 28px', background: 'var(--secondary-fixed)', color: '#002113' }}>
              Solicitar acesso <ArrowRight size={16} style={{ display: 'inline', marginLeft: 6 }} />
            </button>
            <button onClick={() => navigate('/login')} style={{ fontSize: 14, padding: '14px 28px', borderRadius: 9999, background: 'transparent', border: '2px solid rgba(255,255,255,0.3)', color: '#fff', cursor: 'pointer' }}>
              Já sou da equipe
            </button>
          </div>
        </div>
      </header>

      {/* Etapa atual */}
      <section style={{ maxWidth: 1180, margin: '-48px auto 0', padding: '0 20px', position: 'relative', zIndex: 20 }}>
        <div className="glass-card premium-hover" style={{ borderRadius: 20, padding: 28, background: 'var(--bg-card)', display: 'flex', flexWrap: 'wrap', gap: 20, alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <span style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', color: 'var(--primary)', letterSpacing: 1 }}>Etapa em andamento</span>
            <h3 style={{ fontSize: 18, fontWeight: 700, marginTop: 4 }}>{etapaAtual.titulo}</h3>
            <p style={{ fontSize: 13, color: 'var(--text-muted)', marginTop: 2 }}>{etapaAtual.periodo}</p>
          </div>
          <div style={{ display: 'flex', gap: 28 }}>
            {etapaAtual.entregas.slice(0, 2).map((entrega) => (
              <div key={entrega.nome} style={{ maxWidth: 220 }}>
                <div style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)' }}>{entrega.nome}</div>
                <div style={{ fontSize: 11, color: 'var(--text-muted)', marginTop: 2 }}>{entrega.detalhe}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Pilares */}
      <section id="pilares" style={{ maxWidth: 1180, margin: '0 auto', padding: '96px 20px 40px' }}>
        <div style={{ textAlign: 'center', marginBottom: 48 }}>
          <h2 style={{ fontSize: 30, fontWeight: 800 }}>O que o GAT 5 entrega</h2>
          <p style={{ fontSize: 14, color: 'var(--text-muted)', marginTop: 8 }}>Seis frentes de trabalho, alinhadas ao cronograma oficial do projeto.</p>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 20 }}>
          {pilares.map((p) => (
            <div key={p.title} className="glass-card premium-hover" style={{ borderRadius: 18, padding: 24, background: 'var(--bg-card)' }}>
              <div style={{ width: 42, height: 42, borderRadius: 12, background: 'var(--primary-light)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 14, color: 'var(--primary)' }}>
                <p.icon size={22} />
              </div>
              <h3 style={{ fontSize: 15, fontWeight: 700, marginBottom: 6 }}>{p.title}</h3>
              <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.5 }}>{p.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Cronograma */}
      <section id="cronograma" style={{ maxWidth: 1180, margin: '0 auto', padding: '40px 20px' }}>
        <div style={{ textAlign: 'center', marginBottom: 40 }}>
          <h2 style={{ fontSize: 30, fontWeight: 800 }}>Linha do tempo</h2>
          <p style={{ fontSize: 14, color: 'var(--text-muted)', marginTop: 8 }}>Quatro semestres, de infraestrutura digital à transferência do repositório à SESAPI.</p>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 16 }}>
          {cronograma.map((s) => (
            <div key={s.id} className="glass-card" style={{ borderRadius: 16, padding: 20, background: 'var(--bg-card)', borderTop: `3px solid ${s.status === 'andamento' ? 'var(--primary)' : 'var(--border-strong)'}` }}>
              <span style={{ fontSize: 10, fontWeight: 700, textTransform: 'uppercase', color: s.status === 'andamento' ? 'var(--primary)' : 'var(--text-muted)' }}>Semestre {s.id} · {s.periodo}</span>
              <h4 style={{ fontSize: 14, fontWeight: 700, marginTop: 6, marginBottom: 10 }}>{s.titulo}</h4>
              <ul style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                {s.entregas.map((e) => (
                  <li key={e.nome} style={{ fontSize: 12, color: 'var(--text-secondary)' }}>• {e.nome}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>

      {/* Territórios */}
      <section id="territorios" style={{ maxWidth: 1180, margin: '0 auto', padding: '40px 20px 96px' }}>
        <div style={{ textAlign: 'center', marginBottom: 40 }}>
          <h2 style={{ fontSize: 30, fontWeight: 800 }}>Macrorregiões atendidas</h2>
          <p style={{ fontSize: 14, color: 'var(--text-muted)', marginTop: 8 }}>Meta de 12 municípios cobertos até o fim do projeto.</p>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16 }}>
          {macrorregioes.map((m) => (
            <div key={m.nome} className="glass-card" style={{ borderRadius: 16, padding: 20, background: 'var(--bg-card)', textAlign: 'center' }}>
              <div style={{ fontSize: 26, fontWeight: 800, color: 'var(--primary)' }}>{m.municipiosAtendidos}<span style={{ fontSize: 14, color: 'var(--text-muted)' }}>/{m.meta}</span></div>
              <div style={{ fontSize: 12, color: 'var(--text-secondary)', marginTop: 4 }}>{m.nome}</div>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="hero-gradient" style={{ padding: '90px 20px', textAlign: 'center' }}>
        <div style={{ maxWidth: 640, margin: '0 auto' }}>
          <ShieldCheck size={32} color="var(--primary-fixed)" style={{ marginBottom: 16 }} />
          <h2 style={{ fontSize: 28, fontWeight: 800, color: '#fff', marginBottom: 14 }}>Faça parte da equipe do GAT 5</h2>
          <p style={{ fontSize: 15, color: 'var(--text-white-85)', marginBottom: 28, lineHeight: 1.6 }}>
            Cadastre-se para acompanhar materiais, oficinas e indicadores de comunicação do PET-Saúde Clima Piauí.
          </p>
          <button onClick={() => navigate('/cadastro')} className="btn-primary" style={{ fontSize: 14, padding: '14px 30px', background: 'var(--secondary-fixed)', color: '#002113' }}>
            Quero me cadastrar
          </button>
        </div>
      </section>

      {/* Footer */}
      <footer style={{ borderTop: '1px solid var(--border)', padding: '40px 20px', background: 'var(--bg-card)' }}>
        <div style={{ maxWidth: 1180, margin: '0 auto', display: 'flex', flexWrap: 'wrap', justifyContent: 'space-between', gap: 20, alignItems: 'center' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <img src="petsaudeclima_icon.png" alt="PET-Saúde Clima Logo" style={{ height: 24 }} />
            <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>© 2026 PET-Saúde Clima — Territórios do Piauí</span>
          </div>
          <div style={{ display: 'flex', gap: 20 }}>
            <button onClick={() => setLegalOpen('termos')} style={{ background: 'none', border: 'none', fontSize: 12, color: 'var(--text-muted)', cursor: 'pointer' }}>Termos de Uso</button>
            <button onClick={() => setLegalOpen('privacidade')} style={{ background: 'none', border: 'none', fontSize: 12, color: 'var(--text-muted)', cursor: 'pointer' }}>Política de Privacidade</button>
          </div>
        </div>
      </footer>

      <LegalModal open={legalOpen} onClose={() => setLegalOpen(null)} />

      {/* Modal de download — renderizado no fim pra ficar por cima de tudo */}
      <ModalDownload
        aberto={modalAberto}
        onFechar={() => setModalAberto(false)}
        plataforma={plataforma}
      />
    </div>
  );
}