package models

import "time"

// Perfil values, mirroring App\Config\Constantes — estendidos para os
// papéis do PET-Saúde Clima (ver documento de escopo da plataforma).
const (
	PerfilAdmin        = 1 // Administrador (gestão total da plataforma)
	PerfilAluno        = 2 // mantido por compatibilidade — equivale a Estudante/Bolsista
	PerfilEstudante    = PerfilAluno
	PerfilProfissional = 3 // Profissional de saúde / ACS
	PerfilGestor       = 4 // Gestor público municipal/estadual
	PerfilPreceptor    = 5 // Preceptor / orientador de estudantes em campo
	PerfilCoordenador  = 6 // Coordenador de GAT
)

// PerfisAutoCadastro são os perfis que uma pessoa pode escolher sozinha na
// tela de Cadastro (Cadastro.tsx). Preceptor, Coordenador e Admin exigem
// convite institucional (ver CadastroAluno em login_handler.go).
var PerfisAutoCadastro = map[int]bool{
	PerfilEstudante:    true,
	PerfilProfissional: true,
	PerfilGestor:       true,
}

// PerfilLabel devolve o nome amigável do perfil, usado em e-mails e
// respostas de API que exibem o papel do usuário.
func PerfilLabel(perfil int) string {
	switch perfil {
	case PerfilAdmin:
		return "Administrador"
	case PerfilEstudante:
		return "Estudante/Bolsista"
	case PerfilProfissional:
		return "Profissional de Saúde"
	case PerfilGestor:
		return "Gestor Público"
	case PerfilPreceptor:
		return "Preceptor"
	case PerfilCoordenador:
		return "Coordenador de GAT"
	default:
		return "Usuário"
	}
}

// Usuario maps 1:1 to the `usuario` table exactly as it exists in the
// production database (see backend/database/hmgmobieduca_webleia_usuario.sql
// in the original Laravel project). Column names/types were not changed so
// the same MySQL database can keep being used without any migration.
//
//	CREATE TABLE `usuario` (
//	  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	  `nome` varchar(255) NOT NULL,
//	  `email` varchar(255) NOT NULL,
//	  `senha` varchar(255) NOT NULL,
//	  `foto` varchar(255) DEFAULT NULL,
//	  `perfil` tinyint unsigned NOT NULL,
//	  `email_verified_at` timestamp NULL DEFAULT NULL,
//	  `remember_token` varchar(255) DEFAULT NULL,
//	  `created_at` timestamp NULL DEFAULT NULL,
//	  `updated_at` timestamp NULL DEFAULT NULL,
//	  `deleted_at` timestamp NULL DEFAULT NULL,
//	  `aluno_id` bigint unsigned DEFAULT NULL,
//	  PRIMARY KEY (`id`),
//	  UNIQUE KEY `usuario_email_unique` (`email`),
//	  FOREIGN KEY (`aluno_id`) REFERENCES `aluno` (`id`)
//	)
type Usuario struct {
	ID    int64   `json:"id" db:"id"`
	Nome  string  `json:"nome" db:"nome"`
	Email string  `json:"email" db:"email"`
	Senha string  `json:"-" db:"senha"`
	Foto  *string `json:"foto" db:"foto"`
	// Banner mostrado no topo do perfil (Perfil.tsx). Moldura é o
	// identificador (preset) do anel decorativo em volta do avatar, ver
	// frontend/src/components/Avatar.tsx — nenhum dos dois faz parte do
	// schema original do WebLEIA, adicionados em 0003_comunidades.up.sql.
	Banner  *string `json:"banner" db:"banner"`
	Moldura *string `json:"moldura" db:"moldura"`
	// Descricao é a "bio" que a pessoa pode escrever no perfil dela — ver
	// EditarPerfilModal.tsx e PerfilUsuarioModal.tsx (ao ver o perfil de
	// outra pessoa). Adicionado em 0005_visibilidade_som.up.sql.
	Descricao *string `json:"descricao" db:"descricao"`

	// ---- Dados institucionais (0009_perfil_institucional) ---------------
	// Preenchidos no Cadastro.tsx do portal e usados para segmentar
	// conteúdo por território/perfil (ex.: materiais do município do ACS).
	Instituicao *string `json:"instituicao" db:"instituicao"`
	Municipio   *string `json:"municipio" db:"municipio"`
	Profissao   *string `json:"profissao" db:"profissao"`

	// ---- Perfil rico (0007_perfil_rico) ---------------------------------
	// Links e Jogos são guardados como JSON puro na coluna (TEXT) — o
	// front manda/recebe já como array, o Go só passa a string adiante
	// (ver LinksList/JogosList abaixo pra quem quiser manipular em Go).
	Links *string `json:"links" db:"links"` // JSON: [{"label":"GitHub","url":"..."}]
	Jogos *string `json:"jogos" db:"jogos"` // JSON: ["Minecraft","Valorant"]
	// StatusCustomizado é a "bio curta" mostrada junto do nome (estilo
	// status do Discord/WhatsApp), ex.: "fazendo código e clima".
	StatusCustomizado *string `json:"status_customizado" db:"status_customizado"`
	// Atividade/AtividadeTipo alimentam a presença rica (ver
	// PresencaBadge.tsx): AtividadeTipo é "jogo" | "voz" | "" (vazio =
	// só "disponível", sem atividade específica agora).
	Atividade     *string `json:"atividade" db:"atividade"`
	AtividadeTipo *string `json:"atividade_tipo" db:"atividade_tipo"`

	Perfil          int        `json:"perfil" db:"perfil"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"-" db:"deleted_at"`
	AlunoID         *int64     `json:"aluno_id" db:"aluno_id"`

	// Computed/derived fields, added at response time (not DB columns),
	// mirroring Usuario::toArray() adding foto_url/thumb.
	FotoURL *string `json:"foto_url,omitempty" db:"-"`
	Thumb   *string `json:"thumb,omitempty" db:"-"`
}

func (u *Usuario) IsAdmin() bool        { return u.Perfil == PerfilAdmin }
func (u *Usuario) IsAluno() bool        { return u.Perfil == PerfilAluno }
func (u *Usuario) IsProfissional() bool { return u.Perfil == PerfilProfissional }
func (u *Usuario) IsGestor() bool       { return u.Perfil == PerfilGestor }
func (u *Usuario) IsPreceptor() bool    { return u.Perfil == PerfilPreceptor }
func (u *Usuario) IsCoordenador() bool  { return u.Perfil == PerfilCoordenador }

// TableName is kept explicit (rather than pluralizing/guessing) so it is
// always obvious which physical table this struct reads/writes.
func (Usuario) TableName() string { return "usuario" }
