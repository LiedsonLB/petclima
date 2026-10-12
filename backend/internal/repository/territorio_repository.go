package repository

import (
	"database/sql"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

// TerritorioRepository agrupa o acesso a Macrorregiao e Gat — tabelas
// pequenas, majoritariamente de leitura, que sustentam o módulo
// "🗺️ Territórios" do planejamento.
type TerritorioRepository struct {
	db *sql.DB
}

func NewTerritorioRepository(db *sql.DB) *TerritorioRepository {
	return &TerritorioRepository{db: db}
}

// ListMacrorregioes devolve todas as macrorregiões já com a contagem de
// municípios ativos/total, usada pelos cards do Dashboard e do Mapa.
func (r *TerritorioRepository) ListMacrorregioes() ([]*models.Macrorregiao, error) {
	rows, err := r.db.Query(`
		SELECT mr.id, mr.nome, mr.descricao, mr.meta_municipios,
		       COUNT(c.id) AS total,
		       SUM(CASE WHEN c.status = 'ativo' THEN 1 ELSE 0 END) AS ativos
		FROM macrorregiao mr
		LEFT JOIN cidade c ON c.macrorregiao_id = mr.id
		GROUP BY mr.id, mr.nome, mr.descricao, mr.meta_municipios
		ORDER BY mr.nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Macrorregiao
	for rows.Next() {
		m := &models.Macrorregiao{}
		if err := rows.Scan(&m.ID, &m.Nome, &m.Descricao, &m.MetaMunicipios, &m.MunicipiosTotal, &m.MunicipiosAtivos); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *TerritorioRepository) FindMacrorregiaoByID(id int64) (*models.Macrorregiao, error) {
	m := &models.Macrorregiao{}
	err := r.db.QueryRow("SELECT id, nome, descricao, meta_municipios FROM macrorregiao WHERE id = ?", id).
		Scan(&m.ID, &m.Nome, &m.Descricao, &m.MetaMunicipios)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Macrorregião não encontrada.", 404)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListGats devolve os GATs, opcionalmente filtrados por macrorregião.
func (r *TerritorioRepository) ListGats(macrorregiaoID int64) ([]*models.Gat, error) {
	query := `
		SELECT g.id, g.nome, g.macrorregiao_id, g.descricao, mr.nome
		FROM gat g
		LEFT JOIN macrorregiao mr ON mr.id = g.macrorregiao_id
	`
	args := []any{}
	if macrorregiaoID > 0 {
		query += " WHERE g.macrorregiao_id = ?"
		args = append(args, macrorregiaoID)
	}
	query += " ORDER BY g.nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Gat
	for rows.Next() {
		g := &models.Gat{}
		var macrorregiaoNome sql.NullString
		if err := rows.Scan(&g.ID, &g.Nome, &g.MacrorregiaoID, &g.Descricao, &macrorregiaoNome); err != nil {
			return nil, err
		}
		if macrorregiaoNome.Valid {
			g.MacrorregiaoNome = &macrorregiaoNome.String
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
