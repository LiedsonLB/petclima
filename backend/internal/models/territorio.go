package models

// Macrorregiao é uma das macrorregiões de saúde do Piauí em que o PET-Saúde
// Clima atua (Meio-Norte, Semiárido, Cerrado, Entre Rios), cada uma reunindo
// vários municípios e um ou mais GATs (Grupos de Aprendizagem Tutorial).
type Macrorregiao struct {
	ID             int64  `json:"id" db:"id"`
	Nome           string `json:"nome" db:"nome"`
	Descricao      *string `json:"descricao" db:"descricao"`
	MetaMunicipios int    `json:"meta_municipios" db:"meta_municipios"`

	// Computados via agregação, preenchidos pelo repositório.
	MunicipiosAtivos int `json:"municipios_ativos,omitempty" db:"-"`
	MunicipiosTotal  int `json:"municipios_total,omitempty" db:"-"`
}

func (Macrorregiao) TableName() string { return "macrorregiao" }

// Gat é um Grupo de Aprendizagem Tutorial, a unidade de equipe/atuação do
// PET-Saúde dentro de uma macrorregião.
type Gat struct {
	ID             int64   `json:"id" db:"id"`
	Nome           string  `json:"nome" db:"nome"`
	MacrorregiaoID *int64  `json:"macrorregiao_id" db:"macrorregiao_id"`
	Descricao      *string `json:"descricao" db:"descricao"`

	MacrorregiaoNome *string `json:"macrorregiao_nome,omitempty" db:"-"`
}

func (Gat) TableName() string { return "gat" }
