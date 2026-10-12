-- Dados de exemplo para desenvolvimento/demonstração local. Nenhuma destas
-- linhas é necessária em produção — sinta-se livre para apagá-las (ou usar
-- 0006_seed_demo.down.sql) antes de publicar o ambiente real.
--
-- Usuários de teste (senha para todos: "petclima123" — troque/remova antes
-- de ir para produção):
--   admin@petsaudeclima.pi.gov.br        (admin)
--   coordenador@petsaudeclima.pi.gov.br  (coordenador)
--   preceptor@petsaudeclima.pi.gov.br    (preceptor)
--   acs.teresina@petsaudeclima.pi.gov.br (profissional / ACS)
--   estudante@petsaudeclima.pi.gov.br    (estudante)

INSERT INTO `aplicacao` (`nome`, `codigo`, `tipo`, `versao`, `created_at`, `updated_at`) VALUES
  ('Portal Web PET-Saúde Clima', 'PETCLIMA_WEB', 1, 1, NOW(), NOW());

INSERT INTO `macrorregiao` (`nome`, `descricao`, `meta_municipios`) VALUES
  ('Meio-Norte', 'Litoral e região metropolitana de Teresina.', 3),
  ('Semiárido', 'Território do semiárido piauiense.', 3),
  ('Cerrado', 'Região de cerrado do sudeste do estado.', 3),
  ('Entre Rios', 'Território entre os rios Parnaíba e Gurguéia.', 3);

INSERT INTO `cidade` (`nome`, `uf`, `macrorregiao_id`, `status`, `latitude`, `longitude`) VALUES
  ('Teresina', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'ativo', -5.089200, -42.801600),
  ('Parnaíba', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'ativo', -2.905600, -41.775400),
  ('Campo Maior', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'planejado', -4.827700, -42.168800),
  ('Picos', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Cerrado'), 'ativo', -7.076900, -41.466900),
  ('Floriano', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Cerrado'), 'planejado', -6.766400, -43.022500),
  ('Oeiras', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), 'ativo', -7.025100, -42.130500),
  ('São Raimundo Nonato', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), 'planejado', -9.012400, -42.698700),
  ('Bom Jesus', 'PI', (SELECT id FROM macrorregiao WHERE nome = 'Entre Rios'), 'planejado', -9.071400, -44.359000);

INSERT INTO `gat` (`nome`, `macrorregiao_id`, `descricao`) VALUES
  ('GAT Meio-Norte', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'Grupo de Aprendizagem Tutorial da macrorregião Meio-Norte.'),
  ('GAT Semiárido', (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), 'Grupo de Aprendizagem Tutorial da macrorregião Semiárido.'),
  ('GAT Cerrado', (SELECT id FROM macrorregiao WHERE nome = 'Cerrado'), 'Grupo de Aprendizagem Tutorial da macrorregião Cerrado.'),
  ('GAT Entre Rios', (SELECT id FROM macrorregiao WHERE nome = 'Entre Rios'), 'Grupo de Aprendizagem Tutorial da macrorregião Entre Rios.'),
  ('GAT Comunicação e Educação Permanente', NULL, 'Grupo transversal responsável pelo Plano de Comunicação, pela plataforma digital e pela Educação Permanente do PET-Saúde Clima.');

-- Senha de todas as contas abaixo: petclima123
INSERT INTO `usuario` (`nome`, `email`, `senha`, `papel`, `municipio_id`, `gat_id`, `instituicao`, `curso_formacao`, `email_verified_at`, `created_at`, `updated_at`) VALUES
  ('Administração PET-Saúde Clima', 'admin@petsaudeclima.pi.gov.br', '$2b$10$EwpQZNmfA/6HktIjlI2Nvuxp4sMZnCl1O3v0p7HicuZRIjRX6zM2S', 'admin', NULL, NULL, 'SESAPI', NULL, NOW(), NOW(), NOW()),
  ('Coordenação GAT Comunicação', 'coordenador@petsaudeclima.pi.gov.br', '$2b$10$EwpQZNmfA/6HktIjlI2Nvuxp4sMZnCl1O3v0p7HicuZRIjRX6zM2S', 'coordenador', (SELECT id FROM cidade WHERE nome = 'Teresina'), (SELECT id FROM gat WHERE nome = 'GAT Comunicação e Educação Permanente'), 'UFPI', 'Saúde Coletiva', NOW(), NOW(), NOW()),
  ('Preceptoria de Campo', 'preceptor@petsaudeclima.pi.gov.br', '$2b$10$EwpQZNmfA/6HktIjlI2Nvuxp4sMZnCl1O3v0p7HicuZRIjRX6zM2S', 'preceptor', (SELECT id FROM cidade WHERE nome = 'Teresina'), (SELECT id FROM gat WHERE nome = 'GAT Meio-Norte'), 'SESAPI', NULL, NOW(), NOW(), NOW()),
  ('ACS Teresina Norte', 'acs.teresina@petsaudeclima.pi.gov.br', '$2b$10$EwpQZNmfA/6HktIjlI2Nvuxp4sMZnCl1O3v0p7HicuZRIjRX6zM2S', 'profissional', (SELECT id FROM cidade WHERE nome = 'Teresina'), (SELECT id FROM gat WHERE nome = 'GAT Meio-Norte'), 'Prefeitura Municipal de Teresina', NULL, NOW(), NOW(), NOW()),
  ('Estudante PET-Saúde', 'estudante@petsaudeclima.pi.gov.br', '$2b$10$EwpQZNmfA/6HktIjlI2Nvuxp4sMZnCl1O3v0p7HicuZRIjRX6zM2S', 'estudante', (SELECT id FROM cidade WHERE nome = 'Parnaíba'), (SELECT id FROM gat WHERE nome = 'GAT Meio-Norte'), 'UFPI', 'Enfermagem', NOW(), NOW(), NOW());

INSERT INTO `curso` (`titulo`, `descricao`, `categoria`, `publico_alvo`, `carga_horaria`, `modulos_total`, `status`, `criado_por`, `created_at`, `updated_at`) VALUES
  ('Fundamentos de Clima e Saúde Pública', 'Curso introdutório sobre a relação entre mudanças climáticas e saúde pública.', 'Clima e Saúde', 'Profissionais de saúde e estudantes', 20, 5, 'publicado', (SELECT id FROM usuario WHERE email = 'admin@petsaudeclima.pi.gov.br'), NOW(), NOW()),
  ('Identificando e combatendo fake news em saúde', 'Estratégias de letramento midiático aplicadas à saúde e ao clima.', 'Comunicação', 'ACS e profissionais de saúde', 16, 4, 'em breve', (SELECT id FROM usuario WHERE email = 'admin@petsaudeclima.pi.gov.br'), NOW(), NOW()),
  ('Telessaúde no contexto de eventos climáticos', 'Uso de telessaúde e regulação a distância em emergências climáticas.', 'Telessaúde', 'Profissionais de saúde', 24, 6, 'planejado', (SELECT id FROM usuario WHERE email = 'admin@petsaudeclima.pi.gov.br'), NOW(), NOW());

INSERT INTO `matricula` (`usuario_id`, `curso_id`, `progresso`, `status`, `iniciado_em`, `created_at`, `updated_at`) VALUES
  ((SELECT id FROM usuario WHERE email = 'estudante@petsaudeclima.pi.gov.br'), (SELECT id FROM curso WHERE titulo = 'Fundamentos de Clima e Saúde Pública'), 40, 'em_andamento', NOW(), NOW(), NOW());

INSERT INTO `material` (`titulo`, `descricao`, `tipo`, `tema`, `macrorregiao_id`, `arquivo_url`, `downloads`, `publicado_por`, `created_at`, `updated_at`) VALUES
  ('Como se hidratar em ondas de calor', 'Card de prevenção para períodos de calor extremo.', 'card', 'Calor extremo', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'https://petclima.pi.gov.br/materiais/hidratacao-calor.png', 18, (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-09-10 09:00:00', NOW()),
  ('Podcast: Clima e saúde mental', 'Episódio sobre os efeitos das mudanças climáticas na saúde mental.', 'podcast', 'Saúde mental', (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), 'https://petclima.pi.gov.br/materiais/podcast-saude-mental.mp3', 9, (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-09-15 09:00:00', NOW()),
  ('Vídeo curto: reconhecendo insolação', 'Vídeo educativo sobre sinais de insolação e primeiros socorros.', 'video', 'Emergências climáticas', (SELECT id FROM macrorregiao WHERE nome = 'Cerrado'), 'https://petclima.pi.gov.br/materiais/video-insolacao.mp4', 14, (SELECT id FROM usuario WHERE email = 'preceptor@petsaudeclima.pi.gov.br'), '2026-09-18 09:00:00', NOW()),
  ('Cartilha: água segura em período de seca', 'Orientações sobre armazenamento e tratamento de água.', 'cartilha', 'Água e saneamento', (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), 'https://petclima.pi.gov.br/materiais/cartilha-agua-segura.pdf', 6, (SELECT id FROM usuario WHERE email = 'preceptor@petsaudeclima.pi.gov.br'), '2026-09-20 09:00:00', NOW()),
  ('Infográfico WhatsApp: sinais de desidratação', 'Infográfico pronto para compartilhamento em grupos de WhatsApp.', 'infografico', 'Calor extremo', (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), 'https://petclima.pi.gov.br/materiais/infografico-desidratacao.png', 11, (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-09-21 09:00:00', NOW()),
  ('Card: fake news sobre vacinas e clima', 'Card de checagem sobre desinformação envolvendo clima e vacinação.', 'card', 'Desinformação', (SELECT id FROM macrorregiao WHERE nome = 'Entre Rios'), 'https://petclima.pi.gov.br/materiais/card-fake-news.png', 4, (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-09-22 09:00:00', NOW());

INSERT INTO `atividade` (`titulo`, `tipo`, `descricao`, `modalidade`, `municipio_id`, `macrorregiao_id`, `gat_id`, `responsavel_id`, `data_realizacao`, `carga_horaria`, `vagas`, `created_at`, `updated_at`) VALUES
  ('Letramento em saúde para ACS', 'oficina', 'Oficina de letramento em saúde e clima voltada a Agentes Comunitários de Saúde.', 'presencial', (SELECT id FROM cidade WHERE nome = 'Teresina'), (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), (SELECT id FROM gat WHERE nome = 'GAT Meio-Norte'), (SELECT id FROM usuario WHERE email = 'preceptor@petsaudeclima.pi.gov.br'), '2026-10-14 08:00:00', 8, 40, NOW(), NOW()),
  ('Combate às fake news em saúde', 'oficina', 'Oficina prática de checagem de informação em saúde e clima.', 'presencial', (SELECT id FROM cidade WHERE nome = 'Parnaíba'), (SELECT id FROM macrorregiao WHERE nome = 'Meio-Norte'), (SELECT id FROM gat WHERE nome = 'GAT Meio-Norte'), (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-10-28 08:00:00', 6, 35, NOW(), NOW()),
  ('Comunicação de risco climático', 'capacitacao', 'Capacitação sobre comunicação de risco em eventos climáticos extremos.', 'presencial', (SELECT id FROM cidade WHERE nome = 'Picos'), (SELECT id FROM macrorregiao WHERE nome = 'Cerrado'), (SELECT id FROM gat WHERE nome = 'GAT Cerrado'), (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-11-11 08:00:00', 6, 30, NOW(), NOW()),
  ('Uso da plataforma digital do PET', 'capacitacao', 'Oficina remota de apresentação da plataforma digital do PET-Saúde Clima.', 'remota', (SELECT id FROM cidade WHERE nome = 'Oeiras'), (SELECT id FROM macrorregiao WHERE nome = 'Semiárido'), (SELECT id FROM gat WHERE nome = 'GAT Semiárido'), (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), '2026-11-25 14:00:00', 3, 25, NOW(), NOW());

-- A oficina remota ganha uma sala de videoconferência (LiveKit) vinculada,
-- para exercitar o fluxo "Entrar na sala virtual" do frontend. Usamos
-- variáveis de sessão em vez de subqueries aninhadas na mesma tabela, que o
-- MySQL rejeita ("can't specify target table for update in FROM clause").
SET @atividade_pet_id = (SELECT id FROM atividade WHERE titulo = 'Uso da plataforma digital do PET' LIMIT 1);

INSERT INTO `sala` (`nome`, `codigo`, `tipo`, `descricao`, `atividade_id`, `criado_por`, `ativa`, `created_at`, `updated_at`) VALUES
  ('Uso da plataforma digital do PET', CONCAT('oficina-', @atividade_pet_id, '-', UNIX_TIMESTAMP()), 'oficina', 'Sala virtual da oficina remota sobre a plataforma digital.', @atividade_pet_id, (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), 1, NOW(), NOW());

SET @sala_pet_id = LAST_INSERT_ID();

UPDATE `atividade` SET `sala_id` = @sala_pet_id WHERE `id` = @atividade_pet_id;

INSERT INTO `atividade_inscricao` (`atividade_id`, `usuario_id`, `inscrito_em`) VALUES
  ((SELECT id FROM atividade WHERE titulo = 'Letramento em saúde para ACS'), (SELECT id FROM usuario WHERE email = 'acs.teresina@petsaudeclima.pi.gov.br'), NOW()),
  ((SELECT id FROM atividade WHERE titulo = 'Letramento em saúde para ACS'), (SELECT id FROM usuario WHERE email = 'estudante@petsaudeclima.pi.gov.br'), NOW()),
  ((SELECT id FROM atividade WHERE titulo = 'Uso da plataforma digital do PET'), (SELECT id FROM usuario WHERE email = 'acs.teresina@petsaudeclima.pi.gov.br'), NOW());

INSERT INTO `noticia` (`titulo`, `resumo`, `conteudo`, `autor_id`, `publicado_em`, `created_at`) VALUES
  ('PET-Saúde Clima Piauí lança portal digital', 'Plataforma reúne materiais, cursos e territórios das 4 macrorregiões do estado.', 'O PET-Saúde Clima Piauí lança sua plataforma digital, reunindo em um só lugar o repositório de materiais, os cursos autoinstrucionais do AVA e os dados territoriais das quatro macrorregiões de saúde do estado.', (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), NOW(), NOW()),
  ('Kit de Comunicação Regional chega aos primeiros municípios', 'Cards, podcasts e vídeos curtos sobre riscos climáticos já estão disponíveis.', 'Os primeiros dez materiais do Kit de Comunicação Regional já estão disponíveis na Biblioteca, com conteúdos sobre calor extremo, saúde mental, emergências climáticas e desinformação.', (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), NOW(), NOW());

-- Indicadores: alguns acessos e downloads recentes para o Painel/Relatórios
-- não começarem completamente zerados em um ambiente de demonstração.
INSERT INTO `acesso_log` (`usuario_id`, `created_at`)
SELECT (SELECT id FROM usuario WHERE email = 'coordenador@petsaudeclima.pi.gov.br'), DATE_SUB(NOW(), INTERVAL n DAY)
FROM (SELECT 1 AS n UNION SELECT 3 UNION SELECT 5 UNION SELECT 8 UNION SELECT 12 UNION SELECT 20) AS dias;

INSERT INTO `material_download_log` (`material_id`, `usuario_id`, `created_at`)
SELECT (SELECT id FROM material WHERE titulo = 'Como se hidratar em ondas de calor'), (SELECT id FROM usuario WHERE email = 'acs.teresina@petsaudeclima.pi.gov.br'), DATE_SUB(NOW(), INTERVAL n DAY)
FROM (SELECT 2 AS n UNION SELECT 6 UNION SELECT 9 UNION SELECT 15 UNION SELECT 18 UNION SELECT 25) AS dias;
