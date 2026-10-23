package models

import "time"

// Tipos de atividade registrável, cobrindo tanto o módulo "📝
// Atividades/ações" quanto "🧑‍⚕️ Educação Permanente" do planejamento.
const (
	AtividadeOficina      = "oficina"
	AtividadeCapacitacao  = "capacitacao"
	AtividadeIntervencao  = "intervencao"
	AtividadeVisita       = "visita"
	AtividadeReuniao      = "reuniao"
)

// Modalidade de realização da atividade — quando remota ou híbrida, a
// atividade pode ter uma sala virtual (LiveKit) vinculada para
// videoconferência (telessaúde/educação permanente a distância).
const (
	ModalidadePresencial = "presencial"
	ModalidadeRemota     = "remota"
	ModalidadeHibrida    = "hibrida"
)

// Atividade é uma oficina, capacitação, intervenção, visita ou reunião
// registrada por um preceptor/coordenador/profissional, associada a um
// território (município/macrorregião) e opcionalmente a um GAT.
type Atividade struct {
	ID             int64      `json:"id" db:"id"`
	Titulo         string     `json:"titulo" db:"titulo"`
	Tipo           string     `json:"tipo" db:"tipo"`
	Descricao      *string    `json:"descricao" db:"descricao"`
	Modalidade     string     `json:"modalidade" db:"modalidade"`
	MunicipioID    *int64     `json:"municipio_id" db:"municipio_id"`
	MacrorregiaoID *int64     `json:"macrorregiao_id" db:"macrorregiao_id"`
	GatID          *int64     `json:"gat_id" db:"gat_id"`
	ResponsavelID  int64      `json:"responsavel_id" db:"responsavel_id"`
	DataRealizacao *time.Time `json:"data_realizacao" db:"data_realizacao"`
	CargaHoraria   int        `json:"carga_horaria" db:"carga_horaria"`
	Vagas          *int       `json:"vagas" db:"vagas"`
	SalaID         *int64     `json:"sala_id" db:"sala_id"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`

	MunicipioNome    *string `json:"municipio_nome,omitempty" db:"-"`
	MacrorregiaoNome *string `json:"macrorregiao_nome,omitempty" db:"-"`
	GatNome          *string `json:"gat_nome,omitempty" db:"-"`
	ResponsavelNome  *string `json:"responsavel_nome,omitempty" db:"-"`
	SalaCodigo       *string `json:"sala_codigo,omitempty" db:"-"`
	Inscritos        int     `json:"inscritos,omitempty" db:"-"`
	// Preenchido apenas no contexto de um usuário autenticado.
	MinhaInscricao *AtividadeInscricao `json:"minha_inscricao,omitempty" db:"-"`
}

func (Atividade) TableName() string { return "atividade" }

// AtividadeInscricao registra a inscrição (e presença) de um usuário em uma
// atividade — o "registrar atividades e resultados" do planejamento.
type AtividadeInscricao struct {
	ID          int64      `json:"id" db:"id"`
	AtividadeID int64      `json:"atividade_id" db:"atividade_id"`
	UsuarioID   int64      `json:"usuario_id" db:"usuario_id"`
	Presente    bool       `json:"presente" db:"presente"`
	InscritoEm  *time.Time `json:"inscrito_em" db:"inscrito_em"`

	UsuarioNome *string `json:"usuario_nome,omitempty" db:"-"`
}

func (AtividadeInscricao) TableName() string { return "atividade_inscricao" }
