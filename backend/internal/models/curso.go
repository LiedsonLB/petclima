package models

import "time"

// Status de um curso autoinstrucional do AVA (Ambiente Virtual de
// Aprendizagem) do PET-Saúde Clima.
const (
	CursoStatusPublicado = "publicado"
	CursoStatusEmBreve   = "em breve"
	CursoStatusPlanejado = "planejado"
)

// Curso é um curso autoinstrucional sobre clima e saúde, correspondendo ao
// módulo "🎓 Cursos / AVA" do planejamento.
type Curso struct {
	ID             int64   `json:"id" db:"id"`
	Titulo         string  `json:"titulo" db:"titulo"`
	Descricao      *string `json:"descricao" db:"descricao"`
	Categoria      *string `json:"categoria" db:"categoria"`
	PublicoAlvo    *string `json:"publico_alvo" db:"publico_alvo"`
	CargaHoraria   int     `json:"carga_horaria" db:"carga_horaria"`
	ModulosTotal   int     `json:"modulos_total" db:"modulos_total"`
	CapaURL        *string `json:"capa_url" db:"capa_url"`
	Status         string  `json:"status" db:"status"`
	MacrorregiaoID *int64  `json:"macrorregiao_id" db:"macrorregiao_id"`
	CriadoPor      *int64  `json:"criado_por" db:"criado_por"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`

	// Computados via agregação de matriculas.
	TotalMatriculados int     `json:"total_matriculados,omitempty" db:"-"`
	ProgressoMedio    float64 `json:"progresso_medio,omitempty" db:"-"`
	// Preenchido apenas quando a listagem é feita no contexto de um
	// usuário autenticado (ex.: "meus cursos").
	MinhaMatricula *Matricula `json:"minha_matricula,omitempty" db:"-"`
}

func (Curso) TableName() string { return "curso" }

// Status de matrícula de um usuário em um curso.
const (
	MatriculaInscrito    = "inscrito"
	MatriculaEmAndamento = "em_andamento"
	MatriculaConcluido   = "concluido"
)

// Matricula é a inscrição de um usuário em um curso, com progresso (%) e
// data de conclusão — atende ao requisito "com progresso e conclusão" do
// planejamento.
type Matricula struct {
	ID                 int64      `json:"id" db:"id"`
	UsuarioID          int64      `json:"usuario_id" db:"usuario_id"`
	CursoID            int64      `json:"curso_id" db:"curso_id"`
	Progresso          int        `json:"progresso" db:"progresso"` // 0-100
	Status             string     `json:"status" db:"status"`
	CertificadoCodigo  *string    `json:"certificado_codigo" db:"certificado_codigo"`
	IniciadoEm         *time.Time `json:"iniciado_em" db:"iniciado_em"`
	ConcluidoEm        *time.Time `json:"concluido_em" db:"concluido_em"`
	CreatedAt          *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at" db:"updated_at"`

	CursoTitulo *string `json:"curso_titulo,omitempty" db:"-"`
}

func (Matricula) TableName() string { return "matricula" }
