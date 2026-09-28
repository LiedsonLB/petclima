import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import {
  login as apiLogin,
  cadastrar as apiCadastrar,
  esqueciSenha as apiEsqueciSenha,
  redefinirSenha as apiRedefinirSenha,
  confirmarEmail as apiConfirmarEmail,
  reenviarConfirmacao as apiReenviarConfirmacao,
  saveSession, clearSession, getStoredUser,
  isAuthenticated, type Usuario, type CadastroPayload,
} from '../lib/api';

interface AuthContextValue {
  usuario: Usuario | null;
  autenticado: boolean;
  carregando: boolean;
  entrar: (email: string, senha: string, manterConectado?: boolean) => Promise<Usuario>;
  cadastrar: (payload: CadastroPayload) => Promise<Usuario>;
  esqueciSenha: (email: string) => Promise<void>;
  redefinirSenha: (email: string, token: string, novaSenha: string) => Promise<void>;
  confirmarEmail: (email: string, token: string) => Promise<void>;
  reenviarConfirmacao: (email: string) => Promise<void>;
  sair: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [usuario, setUsuario] = useState<Usuario | null>(() => getStoredUser());
  const [carregando, setCarregando] = useState(false);

  const entrar = useCallback(async (email: string, senha: string, manterConectado = true) => {
    setCarregando(true);
    try {
      const result = await apiLogin(email, senha, manterConectado);
      saveSession(result);
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

  // ---- NOVOS ----
  const confirmarEmail = useCallback(async (email: string, token: string) => {
    await apiConfirmarEmail(email, token);
  }, []);

  const reenviarConfirmacao = useCallback(async (email: string) => {
    await apiReenviarConfirmacao(email);
  }, []);

  const sair = useCallback(() => {
    clearSession();
    setUsuario(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        usuario,
        autenticado: isAuthenticated(),
        carregando,
        entrar,
        cadastrar,
        esqueciSenha,
        redefinirSenha,
        confirmarEmail,
        reenviarConfirmacao,
        sair,
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