-- Módulos "📝 Atividades/ações" e "🧑‍⚕️ Educação Permanente": oficinas,
-- capacitações, intervenções, visitas e reuniões, com inscrição/presença.
-- Também liga `sala` (LiveKit) a uma atividade remota/híbrida.

CREATE TABLE IF NOT EXISTS `atividade` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(255) NOT NULL,
  `tipo` varchar(20) NOT NULL DEFAULT 'oficina',
  `descricao` text DEFAULT NULL,
  `modalidade` varchar(20) NOT NULL DEFAULT 'presencial',
  `municipio_id` bigint unsigned DEFAULT NULL,
  `macrorregiao_id` bigint unsigned DEFAULT NULL,
  `gat_id` bigint unsigned DEFAULT NULL,
  `responsavel_id` bigint unsigned NOT NULL,
  `data_realizacao` datetime DEFAULT NULL,
  `carga_horaria` int unsigned NOT NULL DEFAULT '0',
  `vagas` int unsigned DEFAULT NULL,
  `sala_id` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `atividade_municipio_fk` (`municipio_id`),
  KEY `atividade_macrorregiao_fk` (`macrorregiao_id`),
  KEY `atividade_gat_fk` (`gat_id`),
  KEY `atividade_sala_fk` (`sala_id`),
  CONSTRAINT `atividade_municipio_fk` FOREIGN KEY (`municipio_id`) REFERENCES `cidade` (`id`) ON DELETE SET NULL,
  CONSTRAINT `atividade_macrorregiao_fk` FOREIGN KEY (`macrorregiao_id`) REFERENCES `macrorregiao` (`id`) ON DELETE SET NULL,
  CONSTRAINT `atividade_gat_fk` FOREIGN KEY (`gat_id`) REFERENCES `gat` (`id`) ON DELETE SET NULL,
  CONSTRAINT `atividade_responsavel_fk` FOREIGN KEY (`responsavel_id`) REFERENCES `usuario` (`id`) ON DELETE CASCADE,
  CONSTRAINT `atividade_sala_fk` FOREIGN KEY (`sala_id`) REFERENCES `sala` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `atividade_inscricao` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `atividade_id` bigint unsigned NOT NULL,
  `usuario_id` bigint unsigned NOT NULL,
  `presente` tinyint unsigned NOT NULL DEFAULT '0',
  `inscrito_em` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `atividade_inscricao_unique` (`atividade_id`, `usuario_id`),
  KEY `atividade_inscricao_usuario_fk` (`usuario_id`),
  CONSTRAINT `atividade_inscricao_atividade_fk` FOREIGN KEY (`atividade_id`) REFERENCES `atividade` (`id`) ON DELETE CASCADE,
  CONSTRAINT `atividade_inscricao_usuario_fk` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `sala`
  ADD CONSTRAINT `sala_atividade_fk` FOREIGN KEY (`atividade_id`) REFERENCES `atividade` (`id`) ON DELETE SET NULL;
