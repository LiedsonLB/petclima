// Cliente HTTP fino para a API do PET-Saúde Clima (backend Go).
//
// Autenticação: o backend grava a sessão em um cookie HttpOnly (`pet_session`)
// no login. O navegador envia esse cookie sozinho em toda requisição
// (`credentials: 'include'`); o JavaScript NÃO consegue ler o token, então um
// eventual XSS não leva a sessão embora. Por isso este arquivo não guarda
// token nenhum em localStorage/sessionStorage.
//
// A `AppKey` identifica a aplicação no login (não é segredo: vai no bundle) e
// precisa ser igual ao `codigo` da tabela `aplicacao` (seed: PETCLIMA_WEB).
export const API_URL = (import.meta.env.VITE_API_URL || '/api').replace(/\/$/, '');
export const APP_KEY = import.meta.env.VITE_APP_KEY || 'PETCLIMA_WEB';

// Disparado quando uma rota autenticada responde 401 (sessão expirada,
// revogada ou inexistente). O AuthContext escuta e leva o usuário ao login.
export const SESSION_EXPIRED_EVENT = 'petclima:session-expired';

// ---- papéis ------------------------------------------------------------
export type Papel =
  | 'comunidade' | 'estudante' | 'pesquisador' | 'profissional'
  | 'gestor' | 'preceptor' | 'coordenador' | 'admin';

export const PAPEL_LABEL: Record<Papel, string> = {
  comunidade: 'Comunidade',
  estudante: 'Estudante / Bolsista',
  pesquisador: 'Pesquisador',
  profissional: 'Profissional de saúde',
  gestor: 'Gestor público',
  preceptor: 'Preceptor',
  coordenador: 'Coordenador',
  admin: 'Administrador',
};

// Espelham IsStaff / IsCoordenadorOuMais do backend (models/usuario.go).
export const ehStaff = (p?: Papel | null) => p === 'preceptor' || p === 'coordenador' || p === 'admin';
export const ehCoordenador = (p?: Papel | null) => p === 'coordenador' || p === 'admin';

export interface Usuario {
  id: number;
  nome: string;
  email: string;
  foto?: string | null;
  cpf?: string | null;
  telefone?: string | null;
  instituicao?: string | null;
  curso_formacao?: string | null;
  descricao?: string | null;
  papel: Papel;
  municipio_id?: number | null;
  municipio_nome?: string | null;
  macrorregiao_nome?: string | null;
  gat_id?: number | null;
  gat_nome?: string | null;
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

export function mensagemDeErro(e: unknown, padrao = 'Não foi possível completar a solicitação.') {
  if (e instanceof ApiError) return e.message;
  if (e instanceof TypeError) return 'Não foi possível conectar ao servidor. Tente novamente.';
  return padrao;
}

async function request<T>(path: string, options: RequestInit = {}, auth = false): Promise<T> {
  const headers: Record<string, string> = {
    AppKey: APP_KEY,
    ...(options.headers as Record<string, string> | undefined),
  };
  // Em upload (FormData) o navegador define o Content-Type com o boundary.
  if (!(options.body instanceof FormData)) headers['Content-Type'] = 'application/json';

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
    // O backend devolve os erros de validação em `errors` (campo -> mensagem).
    throw new ApiError(message, res.status, body?.errors ?? body?.fields);
  }
  return body as T;
}

const json = (data: unknown) => JSON.stringify(data);

function qs(params: Record<string, string | number | boolean | undefined | null>) {
  const q = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== null && v !== '' && v !== false) q.set(k, String(v));
  }
  const s = q.toString();
  return s ? `?${s}` : '';
}

// Arquivos enviados ficam em /uploads/... no backend; o site os acessa pelo
// mesmo prefixo da API (/api/uploads/...). URLs absolutas passam direto.
export function assetUrl(path?: string | null) {
  if (!path) return '';
  if (/^(https?:|data:|blob:)/i.test(path)) return path;
  if (path.startsWith('/uploads/')) return `${API_URL}${path}`;
  return path;
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
    body: json({ email, senha, long_token: longToken ? 'sim' : 'nao' }),
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
  papel: Papel;
  instituicao?: string;
  municipio_id?: number | null;
  curso_formacao?: string;
  aceite_termos: boolean;
}

export function cadastrar(payload: CadastroPayload) {
  return request<Usuario>('/cadastro', {
    method: 'POST',
    body: json({ ...payload, run_login: 'nao' }),
  });
}

export function esqueciSenha(email: string) {
  return request<{ message: string }>('/acesso/esqueci-senha', {
    method: 'POST',
    body: json({ email }),
  });
}

export function redefinirSenha(email: string, token: string, novaSenha: string) {
  return request<{ message: string }>('/acesso/redefinir-senha', {
    method: 'POST',
    body: json({ email, token, nova_senha: novaSenha }),
  });
}

export function confirmarEmail(email: string, token: string) {
  return request<{ message: string }>('/acesso/confirmar-email', {
    method: 'POST',
    body: json({ email, token }),
  });
}

export function reenviarConfirmacao(email: string) {
  return request<{ message: string }>('/acesso/reenviar-confirmacao', {
    method: 'POST',
    body: json({ email }),
  });
}

// ---- perfil ------------------------------------------------------------
export interface PerfilPayload {
  nome: string;
  email: string;
  telefone?: string;
  instituicao?: string;
  curso_formacao?: string;
  descricao?: string;
  municipio_id?: number | null;
  foto?: string;
}

export function atualizarPerfil(id: number, payload: PerfilPayload) {
  return request<Usuario>(`/usuarios/${id}`, { method: 'PUT', body: json(payload) }, true);
}

// Envia a foto de perfil (jpg/png/webp/gif, até 5MB) e devolve a URL para
// gravar em `foto` via atualizarPerfil. `fotoAntiga` é apagada do disco.
export function enviarFoto(arquivo: File, fotoAntiga?: string | null) {
  const form = new FormData();
  form.append('arquivo', arquivo);
  if (fotoAntiga) form.append('foto_antiga', fotoAntiga);
  return request<{ url: string }>('/upload/foto', { method: 'POST', body: form }, true);
}

// ---- territórios (público) ---------------------------------------------
export interface Macrorregiao {
  id: number;
  nome: string;
  descricao?: string | null;
  meta_municipios: number;
  municipios_ativos?: number;
  municipios_total?: number;
}

export interface Municipio {
  id: number;
  nome: string;
  uf: string;
  macrorregiao_id?: number | null;
  macrorregiao_nome?: string | null;
  status: 'ativo' | 'planejado' | string;
  latitude?: number | null;
  longitude?: number | null;
}

export function listarMacrorregioes() {
  return request<Macrorregiao[]>('/territorios/macrorregioes');
}

export function listarMunicipios(params: { macrorregiao_id?: number; status?: string } = {}) {
  return request<Municipio[]>(`/territorios/municipios${qs(params)}`);
}

// ---- biblioteca / comunicação ------------------------------------------
export type TipoMaterial =
  | 'cartilha' | 'protocolo' | 'guia' | 'video' | 'podcast'
  | 'card' | 'infografico' | 'radio' | 'pdf';

export interface Material {
  id: number;
  titulo: string;
  descricao?: string | null;
  tipo: TipoMaterial;
  tema?: string | null;
  publico_alvo?: string | null;
  arquivo_url: string;
  capa_url?: string | null;
  macrorregiao_id?: number | null;
  macrorregiao_nome?: string | null;
  municipio_nome?: string | null;
  downloads: number;
  autor_nome?: string | null;
  created_at?: string | null;
}

export interface MateriaisFiltro {
  q?: string;
  tipo?: string; // lista separada por vírgula: "card,podcast,video"
  macrorregiao_id?: number;
  limit?: number;
}

export function listarMateriais(filtro: MateriaisFiltro = {}) {
  return request<Material[]>(`/materiais${qs({ ...filtro })}`);
}

// Registra o download (contador + Painel) e devolve a URL do arquivo.
export function baixarMaterial(id: number) {
  return request<{ id: number; arquivo_url: string; downloads: number }>(
    `/materiais/${id}/download`,
    { method: 'POST' },
  );
}

export interface MaterialPayload {
  titulo: string;
  descricao?: string;
  tipo: TipoMaterial;
  tema?: string;
  arquivo_url: string;
  macrorregiao_id?: number | null;
}

export function criarMaterial(payload: MaterialPayload) {
  return request<Material>('/materiais', { method: 'POST', body: json(payload) }, true);
}

// Envia o arquivo (PDF, áudio, vídeo, imagem) e devolve a URL a ser usada em
// `arquivo_url` ao criar o material.
export function enviarArquivoMaterial(arquivo: File) {
  const form = new FormData();
  form.append('arquivo', arquivo);
  return request<{ url: string; nome: string; tamanho: number }>(
    '/upload/material',
    { method: 'POST', body: form },
    true,
  );
}

// ---- cursos / AVA ------------------------------------------------------
export interface Matricula {
  id: number;
  curso_id: number;
  progresso: number;
  status: 'inscrito' | 'em_andamento' | 'concluido' | string;
  certificado_codigo?: string | null;
}

export interface Curso {
  id: number;
  titulo: string;
  descricao?: string | null;
  categoria?: string | null;
  publico_alvo?: string | null;
  carga_horaria: number;
  modulos_total: number;
  status: 'publicado' | 'em breve' | 'planejado' | string;
  total_matriculados?: number;
  progresso_medio?: number;
  minha_matricula?: Matricula | null;
}

export function listarCursos() {
  return request<Curso[]>('/cursos');
}

export function matricular(cursoId: number) {
  return request<Matricula>(`/cursos/${cursoId}/matricula`, { method: 'POST' }, true);
}

export function atualizarProgresso(cursoId: number, progresso: number) {
  return request<Matricula>(
    `/cursos/${cursoId}/progresso`,
    { method: 'PUT', body: json({ progresso }) },
    true,
  );
}

// ---- atividades / educação permanente ----------------------------------
export interface Atividade {
  id: number;
  titulo: string;
  tipo: string;
  descricao?: string | null;
  modalidade: 'presencial' | 'remota' | 'hibrida' | string;
  municipio_id?: number | null;
  municipio_nome?: string | null;
  macrorregiao_nome?: string | null;
  data_realizacao?: string | null;
  carga_horaria: number;
  vagas?: number | null;
  sala_id?: number | null;
  inscritos?: number;
  minha_inscricao?: { id: number; presente: boolean } | null;
}

export interface AtividadesFiltro {
  tipo?: string; // "oficina,capacitacao"
  futuras?: boolean;
  limit?: number;
}

export function listarAtividades(filtro: AtividadesFiltro = {}) {
  const { futuras, ...resto } = filtro;
  return request<Atividade[]>(`/atividades${qs({ ...resto, futuras: futuras ? 1 : undefined })}`);
}

export interface AtividadePayload {
  titulo: string;
  tipo: string;
  descricao?: string;
  modalidade: string;
  municipio_id?: number | null;
  data_realizacao?: string;
  carga_horaria?: number;
  vagas?: number | null;
}

export function criarAtividade(payload: AtividadePayload) {
  return request<Atividade>('/atividades', { method: 'POST', body: json(payload) }, true);
}

export function inscreverAtividade(id: number) {
  return request<unknown>(`/atividades/${id}/inscricao`, { method: 'POST' }, true);
}

export function cancelarInscricao(id: number) {
  return request<{ message: string }>(`/atividades/${id}/inscricao`, { method: 'DELETE' }, true);
}

// ---- painel / relatórios -----------------------------------------------
export interface Indicadores {
  municipios_atendidos: number;
  municipios_meta: number;
  profissionais_capacitados: number;
  profissionais_meta: number;
  materiais_produzidos: number;
  materiais_meta: number;
  macrorregioes_cobertas: number;
  macrorregioes_total: number;
  total_usuarios: number;
  total_downloads: number;
  total_atividades: number;
  total_inscritos_atividades: number;
  engajamento_mensal: { mes: string; acessos: number; downloads: number }[];
  alcance_por_macrorregiao: { macrorregiao: string; acessos: number }[];
}

// KPIs do projeto (qualquer usuário logado).
export function painelResumo() {
  return request<Indicadores>('/painel/resumo', { method: 'GET' }, true);
}

// KPIs + séries de acesso (coordenador ou admin).
export function painelIndicadores(meses = 6) {
  return request<Indicadores>(`/painel/indicadores${qs({ meses })}`, { method: 'GET' }, true);
}

// ---- salas (LiveKit) — Oficinas ao vivo / Educação Permanente ---------
export interface Sala {
  id: number;
  nome: string;
  codigo: string;
  tipo: string;
  descricao?: string | null;
  atividade_id?: number | null;
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
