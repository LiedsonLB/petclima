package models

import "time"

// Tipos de material aceitos pela Biblioteca/Repositório e pelo módulo de
// Comunicação (cards, vídeos e podcasts para WhatsApp/rádio).
const (
	MaterialCartilha    = "cartilha"
	MaterialProtocolo   = "protocolo"
	MaterialGuia        = "guia"
	MaterialVideo       = "video"
	MaterialPodcast     = "podcast"
	MaterialCard        = "card"
	MaterialInfografico = "infografico"
	MaterialRadio       = "radio"
	MaterialPDF         = "pdf"
)

// Material é qualquer item da Biblioteca/Repositório ou do kit de
// Comunicação: PDFs, cartilhas, protocolos, guias, vídeos, podcasts, cards
// etc., buscáveis por tema, região, público, tipo e GAT.
type Material struct {
	ID             int64   `json:"id" db:"id"`
	Titulo         string  `json:"titulo" db:"titulo"`
	Descricao      *string `json:"descricao" db:"descricao"`
	Tipo           string  `json:"tipo" db:"tipo"`
	Tema           *string `json:"tema" db:"tema"`
	PublicoAlvo    *string `json:"publico_alvo" db:"publico_alvo"`
	ArquivoURL     string  `json:"arquivo_url" db:"arquivo_url"`
	CapaURL        *string `json:"capa_url" db:"capa_url"`
	MacrorregiaoID *int64  `json:"macrorregiao_id" db:"macrorregiao_id"`
	MunicipioID    *int64  `json:"municipio_id" db:"municipio_id"`
	GatID          *int64  `json:"gat_id" db:"gat_id"`
	Downloads      int     `json:"downloads" db:"downloads"`
	PublicadoPor   *int64  `json:"publicado_por" db:"publicado_por"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`

	MacrorregiaoNome *string `json:"macrorregiao_nome,omitempty" db:"-"`
	MunicipioNome    *string `json:"municipio_nome,omitempty" db:"-"`
	GatNome          *string `json:"gat_nome,omitempty" db:"-"`
	AutorNome        *string `json:"autor_nome,omitempty" db:"-"`
}

func (Material) TableName() string { return "material" }

// Noticia alimenta o Portal Público (🏠) com informações sobre o projeto,
// clima e saúde, materiais e ações.
type Noticia struct {
	ID          int64      `json:"id" db:"id"`
	Titulo      string     `json:"titulo" db:"titulo"`
	Resumo      *string    `json:"resumo" db:"resumo"`
	Conteudo    string     `json:"conteudo" db:"conteudo"`
	CapaURL     *string    `json:"capa_url" db:"capa_url"`
	AutorID     *int64     `json:"autor_id" db:"autor_id"`
	PublicadoEm *time.Time `json:"publicado_em" db:"publicado_em"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`

	AutorNome *string `json:"autor_nome,omitempty" db:"-"`
}

func (Noticia) TableName() string { return "noticia" }
