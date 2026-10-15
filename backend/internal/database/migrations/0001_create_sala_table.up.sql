-- Tabela `sala`: cada linha é uma sala de videoconferência (LiveKit) — de
-- reunião de equipe, oficina remota ou telessaúde. Pode existir avulsa ou
-- vinculada a uma atividade. A coluna `atividade_id` nasce aqui sem FK: a
-- constraint `sala_atividade_fk` é adicionada em 0004_atividades.up.sql,
-- depois que a tabela `atividade` já existe.

CREATE TABLE IF NOT EXISTS `sala` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) NOT NULL,
  `codigo` varchar(40) NOT NULL,
  `tipo` varchar(20) NOT NULL DEFAULT 'reuniao',
  `descricao` varchar(500) DEFAULT NULL,
  `atividade_id` bigint unsigned DEFAULT NULL,
  `criado_por` bigint unsigned NOT NULL,
  `ativa` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `sala_codigo_unique` (`codigo`),
  KEY `sala_criado_por_fk` (`criado_por`),
  CONSTRAINT `sala_criado_por_fk` FOREIGN KEY (`criado_por`) REFERENCES `usuario` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Bancos que já tinham uma tabela `sala` de uma versão anterior (o CREATE
-- acima não altera tabela existente): garante as colunas que a API usa.
-- Idempotente — só adiciona o que estiver faltando.
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `tipo` varchar(20) NOT NULL DEFAULT ''reuniao''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'tipo');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `descricao` varchar(500) DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'descricao');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `atividade_id` bigint unsigned DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'atividade_id');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `ativa` tinyint unsigned NOT NULL DEFAULT ''1''', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'ativa');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `created_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'created_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `updated_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'updated_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl := (SELECT IF(COUNT(*) = 0, 'ALTER TABLE `sala` ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL', 'SELECT 1') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sala' AND COLUMN_NAME = 'deleted_at');
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;