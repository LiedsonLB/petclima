import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react';
import {
  login as apiLogin,
  logout as apiLogout,
  me as apiMe,
  cadastrar as apiCadastrar,
  esqueciSenha as apiEsqueciSenha,
  redefinirSenha as apiRedefinirSenha,
  confirmarEmail as apiConfirmarEmail,
  reenviarConfirmacao as apiReenviarConfirmacao,
  limparSessaoLegada, SESSION_EXPIRED_EVENT,
  type Usuario, type CadastroPayload,
} from '../lib/api';

interface AuthContextValue {
  usuario: Usuario | null;
  autenticado: boolean;
  /** true enquanto o site confirma com o servidor se já existe sessão (cookie). */
  verificando: boolean;
  carregando: boolean;
  entrar: (email: string, senha: string, manterConectado?: boolean) => Promise<Usuario>;
  cadastrar: (payload: CadastroPayload) => Promise<Usuario>;
  esqueciSenha: (email: string) => Promise<void>;
  redefinirSenha: (email: string, token: string, novaSenha: string) => Promise<void>;
  confirmarEmail: (email: string, token: string) => Promise<void>;
  reenviarConfirmacao: (email: string) => Promise<void>;
  sair: () => void;
  /** Atualiza o usuário em memória (ex.: depois de editar o perfil). */
  atualizarUsuario: (u: Usuario) => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  // A sessão vive num cookie HttpOnly que o JavaScript não enxerga; aqui só
  // guardamos em memória quem é o usuário, perguntando ao servidor (/acesso/me).
  const [usuario, setUsuario] = useState<Usuario | null>(null);
  const [verificando, setVerificando] = useState(true);
  const [carregando, setCarregando] = useState(false);

  useEffect(() => {
    limparSessaoLegada();

    let ativo = true;
    apiMe()
      .then((u) => ativo && setUsuario(u))
      .catch(() => ativo && setUsuario(null))
      .finally(() => ativo && setVerificando(false));

    // Qualquer rota autenticada que responda 401 derruba a sessão local.
    const aoExpirar = () => setUsuario(null);
    window.addEventListener(SESSION_EXPIRED_EVENT, aoExpirar);
    return () => {
      ativo = false;
      window.removeEventListener(SESSION_EXPIRED_EVENT, aoExpirar);
    };
  }, []);

  const entrar = useCallback(async (email: string, senha: string, manterConectado = true) => {
    setCarregando(true);
    try {
      const result = await apiLogin(email, senha, manterConectado);
      setUsuario(result);
      return result;
    } finally {
      setCarregando(false);
    }
  }, []);

  const cadastrar = useCallback(async (payload: CadastroPayload) => {
    setCarregando(true);
    try {
      return await apiCadastrar(payload);
    } finally {
      setCarregando(false);
    }
  }, []);

  const esqueciSenha = useCallback(async (email: string) => {
    await apiEsqueciSenha(email);
  }, []);

  const redefinirSenha = useCallback(async (email: string, token: string, novaSenha: string) => {
    await apiRedefinirSenha(email, token, novaSenha);
  }, []);

  const confirmarEmail = useCallback(async (email: string, token: string) => {
    await apiConfirmarEmail(email, token);
  }, []);

  const reenviarConfirmacao = useCallback(async (email: string) => {
    await apiReenviarConfirmacao(email);
  }, []);

  // Sai da tela na hora e avisa o servidor, que revoga o token e apaga o cookie.
  const sair = useCallback(() => {
    setUsuario(null);
    apiLogout().catch(() => {
      /* sessão já expirada ou sem rede: o cookie some quando expirar */
    });
  }, []);

  return (
    <AuthContext.Provider
      value={{
        usuario,
        autenticado: usuario !== null,
        verificando,
        carregando,
        entrar,
        cadastrar,
        esqueciSenha,
        redefinirSenha,
        confirmarEmail,
        reenviarConfirmacao,
        sair,
        atualizarUsuario: setUsuario,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth precisa estar dentro de <AuthProvider>');
  return ctx;
}