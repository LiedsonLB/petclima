// Cliente HTTP fino para a API do PET-Saúde Clima (backend Go).
//
// Autenticação: o backend grava a sessão em um cookie HttpOnly (`pet_session`)
// no login. O navegador envia esse cookie sozinho em toda requisição
// (`credentials: 'include'`); o JavaScript NÃO consegue ler o token, então um
// eventual XSS não leva a sessão embora. Por isso este arquivo não guarda
// token nenhum em localStorage/sessionStorage.
//
// A `AppKey` identifica a aplicação no login (não é segredo: vai no bundle).
export const API_URL = (import.meta.env.VITE_API_URL || '/api').replace(/\/$/, '');
export const APP_KEY = import.meta.env.VITE_APP_KEY ?? 'PETSAUDEWEB';

// Disparado quando uma rota autenticada responde 401 (sessão expirada,
// revogada ou inexistente). O AuthContext escuta e leva o usuário ao login.
export const SESSION_EXPIRED_EVENT = 'petclima:session-expired';

export interface Usuario {
  id: number;
  nome: string;
  email: string;
  foto?: string | null;
  perfil: number;
  perfil_label?: string;
  instituicao?: string | null;
  municipio?: string | null;
  profissao?: string | null;
  email_verified_at?: string | null;
}

export class ApiError extends Error {
  status: number;
  fields?: Record<string, string>;
  constructor(message: string, status: number, fields?: Record<string, string>) {
    super(message);
    this.status = status;
    this.fields = fields;
  }
}

async function request<T>(path: string, options: RequestInit = {}, auth = false): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    AppKey: APP_KEY,
    ...(options.headers as Record<string, string> | undefined),
  };

  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
    credentials: 'include', // envia/recebe o cookie de sessão
  });
  const isJson = res.headers.get('content-type')?.includes('application/json');
  const body = isJson ? await res.json().catch(() => ({})) : {};

  if (!res.ok) {
    if (auth && res.status === 401) {
      window.dispatchEvent(new Event(SESSION_EXPIRED_EVENT));
    }
    const message = body?.message || body?.error || 'Não foi possível completar a solicitação.';
    throw new ApiError(message, res.status, body?.fields);
  }
  return body as T;
}

// Remove o token/usuário que versões antigas do site guardavam em
// localStorage, para que não fiquem esquecidos no navegador.
export function limparSessaoLegada() {
  try {
    localStorage.removeItem('petsaudeclima:token');
    localStorage.removeItem('petsaudeclima:user');
  } catch {
    /* localStorage indisponível: nada a limpar */
  }
}

// ---- acesso / cadastro -----------------------------------------------
export function login(email: string, senha: string, longToken = true) {
  return request<Usuario>('/acesso/login', {
    method: 'POST',
    body: JSON.stringify({ email, senha, long_token: longToken ? 'sim' : 'nao' }),
  });
}

// Usuário da sessão atual (lida do cookie). 401 => ninguém logado.
export function me() {
  return request<Usuario>('/acesso/me', { method: 'GET' }, true);
}

// Revoga o token no servidor e apaga o cookie.
export function logout() {
  return request<{ message: string }>('/acesso/logout', { method: 'POST' }, true);
}

export interface CadastroPayload {
  nome: string;
  email: string;
  senha: string;
  instituicao?: string;
  municipio?: string;
  profissao?: string;
  perfil: number;
}

export function cadastrar(payload: CadastroPayload) {
  return request<Usuario>('/cadastro', {
    method: 'POST',
    body: JSON.stringify({ ...payload, run_login: 'nao' }),
  });
}

export function esqueciSenha(email: string) {
  return request<{ message: string }>('/acesso/esqueci-senha', {
    method: 'POST',
    body: JSON.stringify({ email }),
  });
}

export function redefinirSenha(email: string, token: string, novaSenha: string) {
  return request<{ message: string }>('/acesso/redefinir-senha', {
    method: 'POST',
    body: JSON.stringify({ email, token, nova_senha: novaSenha }),
  });
}

export function confirmarEmail(email: string, token: string) {
  return request<{ message: string }>('/acesso/confirmar-email', {
    method: 'POST',
    body: JSON.stringify({ email, token }),
  });
}

export function reenviarConfirmacao(email: string) {
  return request<{ message: string }>('/acesso/reenviar-confirmacao', {
    method: 'POST',
    body: JSON.stringify({ email }),
  });
}

// ---- salas (LiveKit) — Oficinas ao vivo / Educação Permanente ---------
export interface Sala {
  id: number;
  nome: string;
  codigo: string;
  tipo: string;
  descricao?: string | null;
  categoria?: string | null;
  ativa: boolean;
  criado_por: number;
  participantes_online?: number;
}

export interface EntrarNaSalaResponse {
  token: string;
  url: string;
  room: string;
  sala: Sala;
}

export function listarSalas() {
  return request<Sala[]>('/salas', { method: 'GET' }, true);
}

export function entrarNaSala(id: number) {
  return request<EntrarNaSalaResponse>(`/salas/${id}/entrar`, { method: 'POST' }, true);
}

export { request as apiRequest };