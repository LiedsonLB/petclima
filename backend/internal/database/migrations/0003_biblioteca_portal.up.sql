-- Módulos "📚 Biblioteca/Repositório", "📢 Comunicação" (mesma tabela,
-- filtrada por `tipo`) e "🏠 Portal público" (notícias).

CREATE TABLE IF NOT EXISTS `material` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(255) NOT NULL,
  `descricao` text DEFAULT NULL,
  `tipo` varchar(20) NOT NULL,
  `tema` varchar(120) DEFAULT NULL,
  `publico_alvo` varchar(255) DEFAULT NULL,
  `arquivo_url` varchar(500) NOT NULL,
  `capa_url` varchar(500) DEFAULT NULL,
  `macrorregiao_id` bigint unsigned DEFAULT NULL,
  `municipio_id` bigint unsigned DEFAULT NULL,
  `gat_id` bigint unsigned DEFAULT NULL,
  `downloads` int unsigned NOT NULL DEFAULT '0',
  `publicado_por` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `material_tipo_idx` (`tipo`),
  KEY `material_macrorregiao_fk` (`macrorregiao_id`),
  KEY `material_municipio_fk` (`municipio_id`),
  KEY `material_gat_fk` (`gat_id`),
  CONSTRAINT `material_macrorregiao_fk` FOREIGN KEY (`macrorregiao_id`) REFERENCES `macrorregiao` (`id`) ON DELETE SET NULL,
  CONSTRAINT `material_municipio_fk` FOREIGN KEY (`municipio_id`) REFERENCES `cidade` (`id`) ON DELETE SET NULL,
  CONSTRAINT `material_gat_fk` FOREIGN KEY (`gat_id`) REFERENCES `gat` (`id`) ON DELETE SET NULL,
  CONSTRAINT `material_publicado_por_fk` FOREIGN KEY (`publicado_por`) REFERENCES `usuario` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `noticia` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(255) NOT NULL,
  `resumo` varchar(500) DEFAULT NULL,
  `conteudo` text NOT NULL,
  `capa_url` varchar(500) DEFAULT NULL,
  `autor_id` bigint unsigned DEFAULT NULL,
  `publicado_em` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `noticia_autor_fk` FOREIGN KEY (`autor_id`) REFERENCES `usuario` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
