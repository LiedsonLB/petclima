-- Campos institucionais do PET-Saúde Clima: instituição, município e
-- profissão/vínculo informados no cadastro (ver Cadastro.tsx), usados para
-- segmentar conteúdo por território e por perfil (ACS, profissional,
-- estudante, gestor, preceptor, coordenador...).
ALTER TABLE `usuario`
  ADD COLUMN `instituicao` varchar(255) DEFAULT NULL AFTER `descricao`,
  ADD COLUMN `municipio` varchar(255) DEFAULT NULL AFTER `instituicao`,
  ADD COLUMN `profissao` varchar(255) DEFAULT NULL AFTER `municipio`;

-- Perfis (ver models.Usuario / models.Perfil*):
--   1 = Administrador     4 = Gestor Público
--   2 = Estudante/Bolsista 5 = Preceptor
--   3 = Profissional       6 = Coordenador de GAT
ALTER TABLE `usuario`
  MODIFY COLUMN `perfil` tinyint unsigned NOT NULL DEFAULT '2';
