DELETE FROM `material_download_log`;
DELETE FROM `acesso_log`;
DELETE FROM `noticia`;
DELETE FROM `atividade_inscricao`;
UPDATE `atividade` SET `sala_id` = NULL;
DELETE FROM `sala`;
DELETE FROM `atividade`;
DELETE FROM `material`;
DELETE FROM `matricula`;
DELETE FROM `curso`;
DELETE FROM `usuario` WHERE `email` IN (
  'admin@petsaudeclima.pi.gov.br',
  'coordenador@petsaudeclima.pi.gov.br',
  'preceptor@petsaudeclima.pi.gov.br',
  'acs.teresina@petsaudeclima.pi.gov.br',
  'estudante@petsaudeclima.pi.gov.br'
);
DELETE FROM `gat`;
DELETE FROM `cidade`;
DELETE FROM `macrorregiao`;
DELETE FROM `aplicacao` WHERE `codigo` = 'PETCLIMA_WEB';
