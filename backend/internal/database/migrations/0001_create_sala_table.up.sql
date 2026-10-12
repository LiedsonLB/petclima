-- Tabela `sala`: cada linha é uma sala de videoconferência (LiveKit) — de
-- reunião de equipe, oficina remota ou telessaúde. Pode existir avulsa ou
-- vinculada a uma atividade (ver 0004_atividades.up.sql, que adiciona a FK
-- `atividade_id` depois que a tabela `atividade` já existe).

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

ALTER TABLE `sala`
  ADD KEY `sala_atividade_fk` (`atividade_id`),
  ADD CONSTRAINT `sala_atividade_fk` FOREIGN KEY (`atividade_id`) REFERENCES `atividade` (`id`) ON DELETE SET NULL;