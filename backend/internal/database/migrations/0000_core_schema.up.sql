-- Schema "core" do PET-Saúde Clima: territórios (macrorregião/município/GAT),
-- contas de acesso (usuario) e o básico de autenticação (aplicacao/tokens).

CREATE TABLE IF NOT EXISTS `macrorregiao` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(120) NOT NULL,
  `descricao` varchar(500) DEFAULT NULL,
  `meta_municipios` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `macrorregiao_nome_unique` (`nome`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `cidade` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) NOT NULL,
  `uf` char(2) NOT NULL,
  `macrorregiao_id` bigint unsigned DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'planejado',
  `latitude` decimal(9,6) DEFAULT NULL,
  `longitude` decimal(9,6) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cidade_macrorregiao_fk` (`macrorregiao_id`),
  CONSTRAINT `cidade_macrorregiao_fk` FOREIGN KEY (`macrorregiao_id`) REFERENCES `macrorregiao` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `gat` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(120) NOT NULL,
  `macrorregiao_id` bigint unsigned DEFAULT NULL,
  `descricao` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `gat_macrorregiao_fk` (`macrorregiao_id`),
  CONSTRAINT `gat_macrorregiao_fk` FOREIGN KEY (`macrorregiao_id`) REFERENCES `macrorregiao` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `usuario` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `senha` varchar(255) NOT NULL,
  `foto` varchar(255) DEFAULT NULL,
  `cpf` varchar(14) DEFAULT NULL,
  `telefone` varchar(20) DEFAULT NULL,
  `instituicao` varchar(255) DEFAULT NULL,
  `curso_formacao` varchar(255) DEFAULT NULL,
  `descricao` varchar(500) DEFAULT NULL,
  `papel` varchar(20) NOT NULL DEFAULT 'comunidade',
  `municipio_id` bigint unsigned DEFAULT NULL,
  `gat_id` bigint unsigned DEFAULT NULL,
  `email_verified_at` timestamp NULL DEFAULT NULL,
  `ultimo_acesso` timestamp NULL DEFAULT NULL,
  `aceite_termos_em` timestamp NULL DEFAULT NULL,
  `remember_token` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `usuario_email_unique` (`email`),
  KEY `usuario_municipio_fk` (`municipio_id`),
  KEY `usuario_gat_fk` (`gat_id`),
  CONSTRAINT `usuario_municipio_fk` FOREIGN KEY (`municipio_id`) REFERENCES `cidade` (`id`) ON DELETE SET NULL,
  CONSTRAINT `usuario_gat_fk` FOREIGN KEY (`gat_id`) REFERENCES `gat` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `aplicacao` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) NOT NULL,
  `codigo` varchar(100) NOT NULL,
  `tipo` tinyint unsigned NOT NULL DEFAULT '1',
  `versao` int unsigned NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `aplicacao_codigo_unique` (`codigo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `aplicacao_acesso_token` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `token` varchar(64) NOT NULL,
  `ip` varchar(64) NOT NULL DEFAULT '',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `aplicacao_id` bigint unsigned NOT NULL,
  `usuario_id` bigint unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `token_idx` (`token`),
  KEY `token_usuario_fk` (`usuario_id`),
  KEY `token_aplicacao_fk` (`aplicacao_id`),
  CONSTRAINT `token_usuario_fk` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`) ON DELETE CASCADE,
  CONSTRAINT `token_aplicacao_fk` FOREIGN KEY (`aplicacao_id`) REFERENCES `aplicacao` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `password_reset_tokens` (
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
