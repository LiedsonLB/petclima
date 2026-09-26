// src/data/Data.ts
// Dados mockados alinhados ao escopo do GAT 5 — Comunicação, Tecnologias
// Digitais e Educação Permanente (PET-Saúde Clima Piauí).

export type StatusEtapa = 'concluido' | 'andamento' | 'planejado';

export interface Semestre {
  id: number;
  titulo: string;
  periodo: string;
  status: StatusEtapa;
  entregas: { nome: string; status: StatusEtapa; detalhe: string }[];
}

export const cronograma: Semestre[] = [
  {
    id: 1,
    titulo: 'Infraestrutura Digital e Produção Inicial',
    periodo: 'Set/2026 – Fev/2027',
    status: 'andamento',
    entregas: [
      { nome: 'Plano Transversal de Comunicação em Saúde e Clima', status: 'andamento', detalhe: 'Estratégia de mídia, linguagem e identidade visual para as 4 macrorregiões' },
      { nome: 'Portal/Plataforma Digital — Versão 1.0', status: 'andamento', detalhe: 'Repositório funcional para cadastro de usuários e hospedagem de materiais' },
      { nome: 'Kit de Comunicação Regional', status: 'planejado', detalhe: '10 materiais iniciais (cards, podcasts e vídeos curtos) sobre riscos climáticos' },
    ],
  },
  {
    id: 2,
    titulo: 'Educação Permanente e Mídias Educativas',
    periodo: 'Mar/2027 – Ago/2027',
    status: 'planejado',
    entregas: [
      { nome: 'Formação em Letramento em Saúde e Combate às Fake News', status: 'planejado', detalhe: 'Oficinas para capacitar no mínimo 200 profissionais e ACS' },
      { nome: 'Pacote de Mídias Educativas', status: 'planejado', detalhe: '30 novos materiais: cartilhas, infográficos para WhatsApp e programas de rádio' },
    ],
  },
  {
    id: 3,
    titulo: 'Expansão da Plataforma e Telessaúde',
    periodo: 'Set/2027 – Fev/2028',
    status: 'planejado',
    entregas: [
      { nome: 'Plataforma Digital Consolidada — Versão 2.0', status: 'planejado', detalhe: 'AVA com módulos de cursos autoinstrucionais sobre Clima e Saúde' },
      { nome: 'Relatório de Alcance de Mídia e Engajamento', status: 'planejado', detalhe: 'Indicadores de tráfego, acessos regionais e downloads por município' },
    ],
  },
  {
    id: 4,
    titulo: 'Repositório Estadual e Transferência',
    periodo: 'Mar/2028 – Ago/2028',
    status: 'planejado',
    entregas: [
      { nome: 'Repositório Aberto de Tecnologias Sociais e Digitais', status: 'planejado', detalhe: 'Entrega da plataforma finalizada à SESAPI' },
      { nome: 'Relatório Final de Educação Permanente', status: 'planejado', detalhe: 'Meta: 12 municípios e as 4 macrorregiões do Piauí atendidos' },
    ],
  },
];

export const macrorregioes = [
  { nome: 'Meio-Norte', municipiosAtendidos: 2, meta: 3 },
  { nome: 'Semiárido', municipiosAtendidos: 1, meta: 3 },
  { nome: 'Cerrado', municipiosAtendidos: 1, meta: 3 },
  { nome: 'Entre Rios', municipiosAtendidos: 0, meta: 3 },
];

export const municipiosMapa = [
  { name: 'Teresina', lat: -5.0892, lng: -42.8016, macrorregiao: 'Meio-Norte', status: 'ativo' as const },
  { name: 'Parnaíba', lat: -2.9056, lng: -41.7754, macrorregiao: 'Meio-Norte', status: 'ativo' as const },
  { name: 'Picos', lat: -7.0769, lng: -41.4669, macrorregiao: 'Cerrado', status: 'ativo' as const },
  { name: 'Floriano', lat: -6.7664, lng: -43.0225, macrorregiao: 'Cerrado', status: 'planejado' as const },
  { name: 'Campo Maior', lat: -4.8277, lng: -42.1688, macrorregiao: 'Meio-Norte', status: 'planejado' as const },
  { name: 'Oeiras', lat: -7.0251, lng: -42.1305, macrorregiao: 'Semiárido', status: 'ativo' as const },
  { name: 'São Raimundo Nonato', lat: -9.0124, lng: -42.6987, macrorregiao: 'Semiárido', status: 'planejado' as const },
  { name: 'Bom Jesus', lat: -9.0714, lng: -44.3590, macrorregiao: 'Entre Rios', status: 'planejado' as const },
];

export type TipoMaterial = 'card' | 'podcast' | 'video' | 'cartilha' | 'infografico' | 'radio';

export interface Material {
  id: string;
  titulo: string;
  tipo: TipoMaterial;
  macrorregiao: string;
  tema: string;
  data: string;
  downloads: number;
}

export const materiais: Material[] = [
  { id: 'm1', titulo: 'Como se hidratar em ondas de calor', tipo: 'card', macrorregiao: 'Meio-Norte', tema: 'Calor extremo', data: '2026-09-10', downloads: 412 },
  { id: 'm2', titulo: 'Podcast: Clima e saúde mental', tipo: 'podcast', macrorregiao: 'Semiárido', tema: 'Saúde mental', data: '2026-09-15', downloads: 198 },
  { id: 'm3', titulo: 'Vídeo curto: reconhecendo insolação', tipo: 'video', macrorregiao: 'Cerrado', tema: 'Emergências climáticas', data: '2026-09-18', downloads: 356 },
  { id: 'm4', titulo: 'Cartilha: água segura em período de seca', tipo: 'cartilha', macrorregiao: 'Semiárido', tema: 'Água e saneamento', data: '2026-09-20', downloads: 87 },
  { id: 'm5', titulo: 'Infográfico WhatsApp: sinais de desidratação', tipo: 'infografico', macrorregiao: 'Meio-Norte', tema: 'Calor extremo', data: '2026-09-21', downloads: 265 },
  { id: 'm6', titulo: 'Card: fake news sobre vacinas e clima', tipo: 'card', macrorregiao: 'Entre Rios', tema: 'Desinformação', data: '2026-09-22', downloads: 143 },
];

export const oficinas = [
  { id: 'o1', titulo: 'Letramento em saúde para ACS', municipio: 'Teresina', data: '2026-10-14', vagas: 40, inscritos: 31 },
  { id: 'o2', titulo: 'Combate às fake news em saúde', municipio: 'Parnaíba', data: '2026-10-28', vagas: 35, inscritos: 20 },
  { id: 'o3', titulo: 'Comunicação de risco climático', municipio: 'Picos', data: '2026-11-11', vagas: 30, inscritos: 12 },
  { id: 'o4', titulo: 'Uso da plataforma digital do PET', municipio: 'Oeiras', data: '2026-11-25', vagas: 25, inscritos: 5 },
];

export const cursosAVA = [
  { id: 'c1', titulo: 'Fundamentos de Clima e Saúde Pública', modulos: 5, cargaHoraria: 20, status: 'em breve' as const, progressoMedio: 0 },
  { id: 'c2', titulo: 'Identificando e combatendo fake news em saúde', modulos: 4, cargaHoraria: 16, status: 'em breve' as const, progressoMedio: 0 },
  { id: 'c3', titulo: 'Telessaúde no contexto de eventos climáticos', modulos: 6, cargaHoraria: 24, status: 'planejado' as const, progressoMedio: 0 },
];

export const engajamentoMensal = [
  { mes: 'Mai', acessos: 0, downloads: 0 },
  { mes: 'Jun', acessos: 0, downloads: 0 },
  { mes: 'Jul', acessos: 0, downloads: 0 },
  { mes: 'Ago', acessos: 0, downloads: 0 },
  { mes: 'Set', acessos: 640, downloads: 210 },
];

export const alcancePorMacrorregiao = [
  { macrorregiao: 'Meio-Norte', acessos: 380 },
  { macrorregiao: 'Semiárido', acessos: 140 },
  { macrorregiao: 'Cerrado', acessos: 90 },
  { macrorregiao: 'Entre Rios', acessos: 30 },
];
