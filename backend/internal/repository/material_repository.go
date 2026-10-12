package repository

import (
	"database/sql"
	"strings"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

// MaterialListParams cobre os filtros do módulo "🔎 Busca": tema, região,
// público, tipo e GAT, além de busca textual livre (Q).
type MaterialListParams struct {
	Q              string
	Tipos          []string
	Tema           string
	PublicoAlvo    string
	MacrorregiaoID int64
	MunicipioID    int64
	GatID          int64
	Limit          int
	Page           int
}

type MaterialRepository struct {
	db *sql.DB
}

func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

const materialSelect = `
	SELECT m.id, m.titulo, m.descricao, m.tipo, m.tema, m.publico_alvo, m.arquivo_url, m.capa_url,
	       m.macrorregiao_id, m.municipio_id, m.gat_id, m.downloads, m.publicado_por,
	       m.created_at, m.updated_at,
	       mr.nome, c.nome, g.nome, u.nome
	FROM material m
	LEFT JOIN macrorregiao mr ON mr.id = m.macrorregiao_id
	LEFT JOIN cidade c ON c.id = m.municipio_id
	LEFT JOIN gat g ON g.id = m.gat_id
	LEFT JOIN usuario u ON u.id = m.publicado_por
`

func scanMaterial(row interface{ Scan(...any) error }) (*models.Material, error) {
	m := &models.Material{}
	var mrNome, cNome, gNome, uNome sql.NullString
	err := row.Scan(&m.ID, &m.Titulo, &m.Descricao, &m.Tipo, &m.Tema, &m.PublicoAlvo, &m.ArquivoURL, &m.CapaURL,
		&m.MacrorregiaoID, &m.MunicipioID, &m.GatID, &m.Downloads, &m.PublicadoPor,
		&m.CreatedAt, &m.UpdatedAt, &mrNome, &cNome, &gNome, &uNome)
	if err != nil {
		return nil, err
	}
	if mrNome.Valid {
		m.MacrorregiaoNome = &mrNome.String
	}
	if cNome.Valid {
		m.MunicipioNome = &cNome.String
	}
	if gNome.Valid {
		m.GatNome = &gNome.String
	}
	if uNome.Valid {
		m.AutorNome = &uNome.String
	}
	return m, nil
}

func (r *MaterialRepository) All(p MaterialListParams) ([]*models.Material, error) {
	clauses := []string{"m.deleted_at IS NULL"}
	args := []any{}

	if p.Q != "" {
		clauses = append(clauses, "(m.titulo LIKE ? OR m.descricao LIKE ? OR m.tema LIKE ?)")
		like := "%" + p.Q + "%"
		args = append(args, like, like, like)
	}
	if len(p.Tipos) > 0 {
		ph := make([]string, len(p.Tipos))
		for i, t := range p.Tipos {
			ph[i] = "?"
			args = append(args, t)
		}
		clauses = append(clauses, "m.tipo IN ("+strings.Join(ph, ",")+")")
	}
	if p.Tema != "" {
		clauses = append(clauses, "m.tema = ?")
		args = append(args, p.Tema)
	}
	if p.PublicoAlvo != "" {
		clauses = append(clauses, "m.publico_alvo LIKE ?")
		args = append(args, "%"+p.PublicoAlvo+"%")
	}
	if p.MacrorregiaoID > 0 {
		clauses = append(clauses, "m.macrorregiao_id = ?")
		args = append(args, p.MacrorregiaoID)
	}
	if p.MunicipioID > 0 {
		clauses = append(clauses, "m.municipio_id = ?")
		args = append(args, p.MunicipioID)
	}
	if p.GatID > 0 {
		clauses = append(clauses, "m.gat_id = ?")
		args = append(args, p.GatID)
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := materialSelect + " WHERE " + strings.Join(clauses, " AND ") + " ORDER BY m.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Material
	for rows.Next() {
		m, err := scanMaterial(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MaterialRepository) FindByID(id int64) (*models.Material, error) {
	row := r.db.QueryRow(materialSelect+" WHERE m.id = ? AND m.deleted_at IS NULL LIMIT 1", id)
	m, err := scanMaterial(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Material não encontrado.", 404)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MaterialRepository) Create(m *models.Material) (*models.Material, error) {
	if m.Titulo == "" {
		return nil, apperr.New("O campo título é obrigatório.", 422)
	}
	if m.ArquivoURL == "" {
		return nil, apperr.New("O campo arquivo_url é obrigatório.", 422)
	}
	if m.Tipo == "" {
		m.Tipo = models.MaterialPDF
	}
	res, err := r.db.Exec(
		`INSERT INTO material (titulo, descricao, tipo, tema, publico_alvo, arquivo_url, capa_url, macrorregiao_id, municipio_id, gat_id, publicado_por, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.Titulo, m.Descricao, m.Tipo, m.Tema, m.PublicoAlvo, m.ArquivoURL, m.CapaURL, m.MacrorregiaoID, m.MunicipioID, m.GatID, m.PublicadoPor,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *MaterialRepository) SoftDelete(id int64) error {
	res, err := r.db.Exec("UPDATE material SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Material não encontrado.", 404)
	}
	return nil
}

// RegistrarDownload incrementa o contador do material e grava uma linha em
// material_download_log (usuarioID zero = download público, sem login).
func (r *MaterialRepository) RegistrarDownload(materialID, usuarioID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE material SET downloads = downloads + 1 WHERE id = ?", materialID); err != nil {
		return err
	}
	var uid any
	if usuarioID > 0 {
		uid = usuarioID
	}
	if _, err := tx.Exec("INSERT INTO material_download_log (material_id, usuario_id, created_at) VALUES (?, ?, NOW())", materialID, uid); err != nil {
		return err
	}
	return tx.Commit()
}
