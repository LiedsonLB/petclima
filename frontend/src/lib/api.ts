// Cliente HTTP fino para a API do PET-Saúde Clima (backend Go).
//
// Todas as rotas de acesso exigem o header `AppKey` (ver
// backend/internal/router/router.go); rotas autenticadas exigem também o
// header `TokenUser`, no formato "{id}:{token_curto}:{token_longo_ou_igual}"
// tal como devolvido pelo backend em `usuario.id` + `token`.
export const API_URL = (import.meta.env.VITE_API_URL || '/api').replace(/\/$/, '');
export const APP_KEY = import.meta.env.VITE_APP_KEY ?? 'PETSAUDEWEB';

const TOKEN_KEY = 'petsaudeclima:token';
const USER_KEY = 'petsaudeclima:user';

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
  token?: string;
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

function authHeader(): string | null {
  const raw = localStorage.getItem(USER_KEY);
  const token = localStorage.getItem(TOKEN_KEY);
  if (!raw || !token) return null;
  const user = JSON.parse(raw) as Usuario;
  return `${user.id}:${token}:${APP_KEY}`;
}

async function request<T>(path: string, options: RequestInit = {}, auth = false): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    AppKey: APP_KEY,
    ...(options.headers as Record<string, string> | undefined),
  };
  if (auth) {
    const tokenUser = authHeader();
    if (tokenUser) headers.TokenUser = tokenUser;
  }

  const res = await fetch(`${API_URL}${path}`, { ...options, headers });
  const isJson = res.headers.get('content-type')?.includes('application/json');
  const body = isJson ? await res.json().catch(() => ({})) : {};

  if (!res.ok) {
    const message = body?.message || body?.error || 'Não foi possível completar a solicitação.';
    throw new ApiError(message, res.status, body?.fields);
  }
  return body as T;
}

export function saveSession(usuario: Usuario) {
  const { token, ...rest } = usuario;
  if (token) localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(rest));
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

export function getStoredUser(): Usuario | null {
  const raw = localStorage.getItem(USER_KEY);
  return raw ? (JSON.parse(raw) as Usuario) : null;
}

export function isAuthenticated(): boolean {
  return Boolean(localStorage.getItem(TOKEN_KEY) && localStorage.getItem(USER_KEY));
}

// ---- acesso / cadastro -----------------------------------------------
export function login(email: string, senha: string, longToken = true) {
  return request<Usuario>('/acesso/login', {
    method: 'POST',
    body: JSON.stringify({ email, senha, long_token: longToken ? 'sim' : 'nao' }),
  });
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
