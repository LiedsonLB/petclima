-- Suporte ao módulo "📊 Painel": logs simples usados para calcular
-- indicadores reais de acesso e download ao longo do tempo, em vez de
-- números fixos. `acesso_log` registra cada login bem-sucedido;
-- `material_download_log` registra cada download de material da
-- Biblioteca/Comunicação (usuario_id fica NULL para downloads públicos,
-- feitos sem login, já que a Biblioteca também é aberta ao público).

CREATE TABLE IF NOT EXISTS `acesso_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `usuario_id` bigint unsigned NOT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `acesso_log_usuario_fk` (`usuario_id`),
  KEY `acesso_log_created_idx` (`created_at`),
  CONSTRAINT `acesso_log_usuario_fk` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `material_download_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `material_id` bigint unsigned NOT NULL,
  `usuario_id` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `material_download_log_material_fk` (`material_id`),
  KEY `material_download_log_created_idx` (`created_at`),
  CONSTRAINT `material_download_log_material_fk` FOREIGN KEY (`material_id`) REFERENCES `material` (`id`) ON DELETE CASCADE,
  CONSTRAINT `material_download_log_usuario_fk` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
