-- Conserto de esquema para bancos que já tinham tabelas de versões
-- anteriores do projeto: os CREATE TABLE IF NOT EXISTS das migrations
-- anteriores não alteram tabela existente, então colunas novas nunca
-- chegavam lá (ex.: "Unknown column 'macrorregiao_id'" no seed).
--
-- Para cada coluna prevista no esquema, adiciona-a SOMENTE se estiver
-- faltando. Em banco novo (ou já correto) não faz nada.

-- macrorregiao
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `macrorregiao` ADD COLUMN `nome` varchar(120) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'macrorregiao' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `macrorregiao` ADD COLUMN `descricao` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'macrorregiao' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `macrorregiao` ADD COLUMN `meta_municipios` int unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'macrorregiao' AND COLUMN_NAME = 'meta_municipios');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- cidade
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `nome` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `uf` char(2) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'uf');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `macrorregiao_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'macrorregiao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `status` varchar(20) NOT NULL DEFAULT ''planejado''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'status');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `latitude` decimal(9,6) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'latitude');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `cidade` ADD COLUMN `longitude` decimal(9,6) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cidade' AND COLUMN_NAME = 'longitude');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- gat
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `gat` ADD COLUMN `nome` varchar(120) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'gat' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `gat` ADD COLUMN `macrorregiao_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'gat' AND COLUMN_NAME = 'macrorregiao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `gat` ADD COLUMN `descricao` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'gat' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- usuario
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `nome` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `email` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'email');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `senha` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'senha');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `foto` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'foto');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `cpf` varchar(14) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'cpf');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `telefone` varchar(20) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'telefone');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `instituicao` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'instituicao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `curso_formacao` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'curso_formacao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `descricao` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `papel` varchar(20) NOT NULL DEFAULT ''comunidade''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'papel');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `municipio_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'municipio_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `gat_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'gat_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `email_verified_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'email_verified_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `ultimo_acesso` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'ultimo_acesso');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `aceite_termos_em` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'aceite_termos_em');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `remember_token` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'remember_token');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `usuario` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'usuario' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- aplicacao
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `nome` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `codigo` varchar(100) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'codigo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `tipo` tinyint unsigned NOT NULL DEFAULT ''1''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'tipo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `versao` int unsigned NOT NULL DEFAULT ''1''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'versao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- aplicacao_acesso_token
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `token` varchar(64) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'token');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `ip` varchar(64) NOT NULL DEFAULT ''''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'ip');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `status` tinyint unsigned NOT NULL DEFAULT ''1''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'status');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `aplicacao_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'aplicacao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `usuario_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'usuario_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `aplicacao_acesso_token` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'aplicacao_acesso_token' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- password_reset_tokens
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `password_reset_tokens` ADD COLUMN `email` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'password_reset_tokens' AND COLUMN_NAME = 'email');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `password_reset_tokens` ADD COLUMN `token` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'password_reset_tokens' AND COLUMN_NAME = 'token');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `password_reset_tokens` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'password_reset_tokens' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- sala
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `nome` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'nome');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `codigo` varchar(40) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'codigo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `tipo` varchar(20) NOT NULL DEFAULT ''reuniao''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'tipo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `descricao` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `atividade_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'atividade_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `criado_por` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'criado_por');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `ativa` tinyint unsigned NOT NULL DEFAULT ''1''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'ativa');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- curso
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `titulo` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'titulo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `descricao` text DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `categoria` varchar(120) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'categoria');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `publico_alvo` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'publico_alvo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `carga_horaria` int unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'carga_horaria');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `modulos_total` int unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'modulos_total');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `capa_url` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'capa_url');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `status` varchar(20) NOT NULL DEFAULT ''planejado''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'status');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `macrorregiao_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'macrorregiao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `criado_por` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'criado_por');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `curso` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'curso' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- matricula
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `usuario_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'usuario_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `curso_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'curso_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `progresso` tinyint unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'progresso');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `status` varchar(20) NOT NULL DEFAULT ''inscrito''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'status');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `certificado_codigo` varchar(40) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'certificado_codigo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `iniciado_em` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'iniciado_em');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `concluido_em` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'concluido_em');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `matricula` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'matricula' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- material
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `titulo` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'titulo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `descricao` text DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `tipo` varchar(20) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'tipo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `tema` varchar(120) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'tema');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `publico_alvo` varchar(255) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'publico_alvo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `arquivo_url` varchar(500) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'arquivo_url');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `capa_url` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'capa_url');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `macrorregiao_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'macrorregiao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `municipio_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'municipio_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `gat_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'gat_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `downloads` int unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'downloads');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `publicado_por` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'publicado_por');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- noticia
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `titulo` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'titulo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `resumo` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'resumo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `conteudo` text NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'conteudo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `capa_url` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'capa_url');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `autor_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'autor_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `publicado_em` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'publicado_em');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `noticia` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'noticia' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- atividade
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `titulo` varchar(255) NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'titulo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `tipo` varchar(20) NOT NULL DEFAULT ''oficina''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'tipo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `descricao` text DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `modalidade` varchar(20) NOT NULL DEFAULT ''presencial''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'modalidade');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `municipio_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'municipio_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `macrorregiao_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'macrorregiao_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `gat_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'gat_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `responsavel_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'responsavel_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `data_realizacao` datetime DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'data_realizacao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `carga_horaria` int unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'carga_horaria');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `vagas` int unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'vagas');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `sala_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'sala_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- atividade_inscricao
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade_inscricao` ADD COLUMN `atividade_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade_inscricao' AND COLUMN_NAME = 'atividade_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade_inscricao` ADD COLUMN `usuario_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade_inscricao' AND COLUMN_NAME = 'usuario_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade_inscricao` ADD COLUMN `presente` tinyint unsigned NOT NULL DEFAULT ''0''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade_inscricao' AND COLUMN_NAME = 'presente');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `atividade_inscricao` ADD COLUMN `inscrito_em` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atividade_inscricao' AND COLUMN_NAME = 'inscrito_em');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- acesso_log
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `acesso_log` ADD COLUMN `usuario_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'acesso_log' AND COLUMN_NAME = 'usuario_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `acesso_log` ADD COLUMN `ip` varchar(64) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'acesso_log' AND COLUMN_NAME = 'ip');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `acesso_log` ADD COLUMN `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'acesso_log' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- material_download_log
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material_download_log` ADD COLUMN `material_id` bigint unsigned NOT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material_download_log' AND COLUMN_NAME = 'material_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material_download_log` ADD COLUMN `usuario_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material_download_log' AND COLUMN_NAME = 'usuario_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `material_download_log` ADD COLUMN `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'material_download_log' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;