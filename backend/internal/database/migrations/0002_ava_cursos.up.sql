-- Módulo "🎓 Cursos / AVA": cursos autoinstrucionais com matrícula,
-- progresso (%) e conclusão (gera código de certificado).

CREATE TABLE IF NOT EXISTS `curso` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(255) NOT NULL,
  `descricao` text DEFAULT NULL,
  `categoria` varchar(120) DEFAULT NULL,
  `publico_alvo` varchar(255) DEFAULT NULL,
  `carga_horaria` int unsigned NOT NULL DEFAULT '0',
  `modulos_total` int unsigned NOT NULL DEFAULT '0',
  `capa_url` varchar(500) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'planejado',
  `macrorregiao_id` bigint unsigned DEFAULT NULL,
  `criado_por` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `curso_macrorregiao_fk` (`macrorregiao_id`),
  CONSTRAINT `curso_macrorregiao_fk` FOREIGN KEY (`macrorregiao_id`) REFERENCES `macrorregiao` (`id`) ON DELETE SET NULL,
  CONSTRAINT `curso_criado_por_fk` FOREIGN KEY (`criado_por`) REFERENCES `usuario` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `matricula` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `usuario_id` bigint unsigned NOT NULL,
  `curso_id` bigint unsigned NOT NULL,
  `progresso` tinyint unsigned NOT NULL DEFAULT '0',
  `status` varchar(20) NOT NULL DEFAULT 'inscrito',
  `certificado_codigo` varchar(40) DEFAULT NULL,
  `iniciado_em` timestamp NULL DEFAULT NULL,
  `concluido_em` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `matricula_usuario_curso_unique` (`usuario_id`, `curso_id`),
  KEY `matricula_curso_fk` (`curso_id`),
  CONSTRAINT `matricula_usuario_fk` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`) ON DELETE CASCADE,
  CONSTRAINT `matricula_curso_fk` FOREIGN KEY (`curso_id`) REFERENCES `curso` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
