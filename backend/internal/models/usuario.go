package models

import "time"

// Papel represents the profile/role a Usuario has inside the PET-Saúde
// Clima platform. This mirrors the hierarchy described in the project's
// planning document (perfil ≠ cargo): each papel unlocks a growing set of
// permissions, from a public visitor up to full platform administration.
type Papel string

const (
	PapelComunidade  Papel = "comunidade"  // público geral cadastrado, sem vínculo institucional
	PapelEstudante   Papel = "estudante"   // estudante/bolsista do PET-Saúde
	PapelPesquisador Papel = "pesquisador" // pesquisador vinculado ao projeto
	PapelProfissional Papel = "profissional" // profissional de saúde / ACS
	PapelGestor      Papel = "gestor"      // gestor municipal/estadual
	PapelPreceptor   Papel = "preceptor"   // preceptor/docente/orientador
	PapelCoordenador Papel = "coordenador" // coordenador de GAT/macrorregião
	PapelAdmin       Papel = "admin"       // administração completa da plataforma
)

// papelNivel define a hierarquia usada pelos helpers abaixo e pelas rotas
// que exigem "papel mínimo X ou superior".
var papelNivel = map[Papel]int{
	PapelComunidade:   1,
	PapelEstudante:    1,
	PapelPesquisador:  1,
	PapelProfissional: 2,
	PapelGestor:       2,
	PapelPreceptor:    3,
	PapelCoordenador:  4,
	PapelAdmin:        5,
}

// PapeisAutoCadastro são os únicos papéis que uma pessoa pode escolher no
// cadastro público (/cadastro). Papéis institucionais (preceptor,
// coordenador, admin) só são atribuídos por um administrador — o "convite
// institucional" descrito no planejamento.
var PapeisAutoCadastro = map[Papel]bool{
	PapelComunidade:   true,
	PapelEstudante:    true,
	PapelPesquisador:  true,
	PapelProfissional: true,
	PapelGestor:       true,
}

// PapelValido confirma que a string recebida é um dos papéis conhecidos.
func PapelValido(p Papel) bool {
	_, ok := papelNivel[p]
	return ok
}

// Usuario é a conta de acesso à plataforma. A tabela original (herdada de
// um projeto de rede social) trazia colunas de perfil "gamer" (banner,
// moldura, jogos, status customizado, atividade de voz) que não fazem
// sentido aqui e foram removidas; no lugar entraram os dados de cadastro
// institucional descritos no documento de planejamento do PET-Saúde Clima
// (CPF, instituição, formação, papel, município/GAT de vínculo).
type Usuario struct {
	ID    int64  `json:"id" db:"id"`
	Nome  string `json:"nome" db:"nome"`
	Email string `json:"email" db:"email"`
	Senha string `json:"-" db:"senha"`
	Foto  *string `json:"foto" db:"foto"`

	CPF            *string `json:"cpf" db:"cpf"`
	Telefone       *string `json:"telefone" db:"telefone"`
	Instituicao    *string `json:"instituicao" db:"instituicao"`
	CursoFormacao  *string `json:"curso_formacao" db:"curso_formacao"`
	Descricao      *string `json:"descricao" db:"descricao"` // bio curta do perfil

	Papel        Papel  `json:"papel" db:"papel"`
	MunicipioID  *int64 `json:"municipio_id" db:"municipio_id"`
	GatID        *int64 `json:"gat_id" db:"gat_id"`

	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	UltimoAcesso    *time.Time `json:"ultimo_acesso" db:"ultimo_acesso"`
	AceiteTermosEm  *time.Time `json:"aceite_termos_em" db:"aceite_termos_em"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"-" db:"deleted_at"`

	// Campos computados na resposta (não são colunas), preenchidos pelo
	// handler a partir de joins com cidade/gat/macrorregiao.
	FotoURL       *string `json:"foto_url,omitempty" db:"-"`
	MunicipioNome *string `json:"municipio_nome,omitempty" db:"-"`
	MacrorregiaoID   *int64  `json:"macrorregiao_id,omitempty" db:"-"`
	MacrorregiaoNome *string `json:"macrorregiao_nome,omitempty" db:"-"`
	GatNome       *string `json:"gat_nome,omitempty" db:"-"`
}

func (u *Usuario) IsAdmin() bool       { return u.Papel == PapelAdmin }
func (u *Usuario) IsCoordenadorOuMais() bool { return papelNivel[u.Papel] >= papelNivel[PapelCoordenador] }
func (u *Usuario) IsStaff() bool       { return papelNivel[u.Papel] >= papelNivel[PapelPreceptor] } // preceptor/coordenador/admin: podem publicar conteúdo e registrar atividades
func (u *Usuario) IsProfissionalOuMais() bool { return papelNivel[u.Papel] >= papelNivel[PapelProfissional] }

func (Usuario) TableName() string { return "usuario" }
