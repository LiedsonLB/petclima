# GAT 5 Comunica — PET-Saúde Clima Piauí

Portal digital do **GAT 5 (Comunicação, Tecnologias Digitais e Educação Permanente)**,
construído reaproveitando o mesmo design system (cores, tipografia, componentes de
layout) do projeto PET-Saúde Clima já existente, adaptado ao escopo e cronograma
do GAT 5.

## Stack

- React 19 + TypeScript + Vite
- Tailwind CSS (tokens mapeados via CSS vars, igual ao projeto original)
- React Router
- Recharts (gráficos de engajamento/alcance)
- React-Leaflet (mapa de alcance por município)

## Páginas

| Rota | Conteúdo | Entrega do cronograma |
|---|---|---|
| `/login` | Acesso da equipe | — |
| `/dashboard` | Visão geral de KPIs e etapa atual | Todas |
| `/cronograma` | Linha do tempo dos 4 semestres | Todas |
| `/comunicacao` | Plano Transversal + Kit de Comunicação Regional | Semestre 1 |
| `/formacao` | Oficinas de letramento em saúde / fake news | Semestre 2 |
| `/ava` | Cursos autoinstrucionais (AVA) | Semestre 3 |
| `/mapa` | Municípios e macrorregiões atendidos | Semestre 4 |
| `/relatorios` | Indicadores de tráfego e engajamento | Semestre 3 |
| `/repositorio` | Entrega final à SESAPI | Semestre 4 |
| `/configuracoes` | Tema claro/escuro, sobre o GAT 5 | — |

## Rodando localmente

```bash
npm install
npm run dev
```

Os dados exibidos (`src/data/Data.ts`) são mockados — quando houver uma API/backend,
basta trocar as importações por chamadas reais mantendo os mesmos tipos.
