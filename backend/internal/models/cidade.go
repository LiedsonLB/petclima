package models

// StatusTerritorio reflete o estágio de atuação do projeto no município,
// usado no Mapa de territórios do frontend (Mapa.tsx).
const (
	TerritorioAtivo     = "ativo"
	TerritorioPlanejado = "planejado"
)

// Cidade representa um município do território de atuação do PET-Saúde
// Clima. Cada município pertence a uma macrorregião de saúde.
type Cidade struct {
	ID             int64    `json:"id" db:"id"`
	Nome           string   `json:"nome" db:"nome"`
	UF             string   `json:"uf" db:"uf"`
	MacrorregiaoID *int64   `json:"macrorregiao_id" db:"macrorregiao_id"`
	Status         string   `json:"status" db:"status"`
	Latitude       *float64 `json:"latitude" db:"latitude"`
	Longitude      *float64 `json:"longitude" db:"longitude"`

	// Computado a partir do join com macrorregiao.
	MacrorregiaoNome *string `json:"macrorregiao_nome,omitempty" db:"-"`
}

func (Cidade) TableName() string { return "cidade" }
