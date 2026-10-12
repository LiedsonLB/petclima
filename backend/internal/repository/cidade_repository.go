package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

// CidadeListParams filtra a listagem de municípios (Mapa de territórios,
// filtros de busca por região).
type CidadeListParams struct {
	Order string
	Sort  string
	Limit int
	Page  int

	Nome           string
	UF             string
	MacrorregiaoID int64
	Status         string
}

type CidadeRepository struct {
	db *sql.DB
}

func NewCidadeRepository(db *sql.DB) *CidadeRepository {
	return &CidadeRepository{db: db}
}

const cidadeSelect = `
	SELECT c.id, c.nome, c.uf, c.macrorregiao_id, c.status, c.latitude, c.longitude, mr.nome
	FROM cidade c
	LEFT JOIN macrorregiao mr ON mr.id = c.macrorregiao_id
`

func scanCidade(row interface{ Scan(...any) error }) (*models.Cidade, error) {
	c := &models.Cidade{}
	var macrorregiaoNome sql.NullString
	err := row.Scan(&c.ID, &c.Nome, &c.UF, &c.MacrorregiaoID, &c.Status, &c.Latitude, &c.Longitude, &macrorregiaoNome)
	if err != nil {
		return nil, err
	}
	if macrorregiaoNome.Valid {
		c.MacrorregiaoNome = &macrorregiaoNome.String
	}
	return c, nil
}

func (r *CidadeRepository) buildWhere(p CidadeListParams) (string, []any) {
	clauses := []string{"1=1"}
	args := []any{}
	if p.Nome != "" {
		clauses = append(clauses, "c.nome LIKE ?")
		args = append(args, "%"+p.Nome+"%")
	}
	if p.UF != "" {
		clauses = append(clauses, "c.uf = ?")
		args = append(args, p.UF)
	}
	if p.MacrorregiaoID > 0 {
		clauses = append(clauses, "c.macrorregiao_id = ?")
		args = append(args, p.MacrorregiaoID)
	}
	if p.Status != "" {
		clauses = append(clauses, "c.status = ?")
		args = append(args, p.Status)
	}
	return strings.Join(clauses, " AND "), args
}

// All lista municípios — usado pelo Mapa de territórios e pelos selects de
// cadastro/atividades/materiais.
func (r *CidadeRepository) All(p CidadeListParams) ([]*models.Cidade, int64, error) {
	order := "nome"
	if p.Order != "" {
		order = p.Order
	}
	sortDir := "asc"
	if strings.EqualFold(p.Sort, "desc") {
		sortDir = "desc"
	}

	where, args := r.buildWhere(p)

	var total int64
	if err := r.db.QueryRow("SELECT COUNT(id) FROM cidade c WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 100
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf("%s WHERE %s ORDER BY c.%s %s LIMIT ? OFFSET ?", cidadeSelect, where, order, sortDir)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.Cidade
	for rows.Next() {
		c, err := scanCidade(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}
	return list, total, rows.Err()
}

func (r *CidadeRepository) FindByID(id int64) (*models.Cidade, error) {
	row := r.db.QueryRow(cidadeSelect+" WHERE c.id = ? LIMIT 1", id)
	c, err := scanCidade(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Município não encontrado.", 404)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
