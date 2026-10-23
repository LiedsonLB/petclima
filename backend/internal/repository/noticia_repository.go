package repository

import (
	"database/sql"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

type NoticiaRepository struct {
	db *sql.DB
}

func NewNoticiaRepository(db *sql.DB) *NoticiaRepository {
	return &NoticiaRepository{db: db}
}

const noticiaSelect = `
	SELECT n.id, n.titulo, n.resumo, n.conteudo, n.capa_url, n.autor_id, n.publicado_em, n.created_at, u.nome
	FROM noticia n
	LEFT JOIN usuario u ON u.id = n.autor_id
`

func scanNoticia(row interface{ Scan(...any) error }) (*models.Noticia, error) {
	n := &models.Noticia{}
	var autor sql.NullString
	err := row.Scan(&n.ID, &n.Titulo, &n.Resumo, &n.Conteudo, &n.CapaURL, &n.AutorID, &n.PublicadoEm, &n.CreatedAt, &autor)
	if err != nil {
		return nil, err
	}
	if autor.Valid {
		n.AutorNome = &autor.String
	}
	return n, nil
}

// All devolve as notícias já publicadas (publicado_em preenchido e no
// passado), da mais recente para a mais antiga.
func (r *NoticiaRepository) All(limit int) ([]*models.Noticia, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(noticiaSelect+" WHERE n.deleted_at IS NULL AND n.publicado_em IS NOT NULL AND n.publicado_em <= NOW() ORDER BY n.publicado_em DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Noticia
	for rows.Next() {
		n, err := scanNoticia(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *NoticiaRepository) FindByID(id int64) (*models.Noticia, error) {
	row := r.db.QueryRow(noticiaSelect+" WHERE n.id = ? AND n.deleted_at IS NULL LIMIT 1", id)
	n, err := scanNoticia(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Notícia não encontrada.", 404)
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NoticiaRepository) Create(n *models.Noticia) (*models.Noticia, error) {
	if n.Titulo == "" {
		return nil, apperr.New("O campo título é obrigatório.", 422)
	}
	if n.Conteudo == "" {
		return nil, apperr.New("O campo conteúdo é obrigatório.", 422)
	}
	res, err := r.db.Exec(
		`INSERT INTO noticia (titulo, resumo, conteudo, capa_url, autor_id, publicado_em, created_at)
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		n.Titulo, n.Resumo, n.Conteudo, n.CapaURL, n.AutorID,
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

func (r *NoticiaRepository) SoftDelete(id int64) error {
	res, err := r.db.Exec("UPDATE noticia SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Notícia não encontrada.", 404)
	}
	return nil
}
