package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

type CursoListParams struct {
	Status         string
	Categoria      string
	MacrorregiaoID int64
	Limit          int
	Page           int
}

type CursoRepository struct {
	db *sql.DB
}

func NewCursoRepository(db *sql.DB) *CursoRepository {
	return &CursoRepository{db: db}
}

const cursoSelect = `
	SELECT cu.id, cu.titulo, cu.descricao, cu.categoria, cu.publico_alvo, cu.carga_horaria,
	       cu.modulos_total, cu.capa_url, cu.status, cu.macrorregiao_id, cu.criado_por,
	       cu.created_at, cu.updated_at,
	       COUNT(DISTINCT m.id) AS total_matriculados,
	       COALESCE(AVG(m.progresso), 0) AS progresso_medio
	FROM curso cu
	LEFT JOIN matricula m ON m.curso_id = cu.id
`

func scanCurso(row interface{ Scan(...any) error }) (*models.Curso, error) {
	c := &models.Curso{}
	err := row.Scan(&c.ID, &c.Titulo, &c.Descricao, &c.Categoria, &c.PublicoAlvo, &c.CargaHoraria,
		&c.ModulosTotal, &c.CapaURL, &c.Status, &c.MacrorregiaoID, &c.CriadoPor,
		&c.CreatedAt, &c.UpdatedAt, &c.TotalMatriculados, &c.ProgressoMedio)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CursoRepository) All(p CursoListParams) ([]*models.Curso, error) {
	clauses := []string{"cu.deleted_at IS NULL"}
	args := []any{}
	if p.Status != "" {
		clauses = append(clauses, "cu.status = ?")
		args = append(args, p.Status)
	}
	if p.Categoria != "" {
		clauses = append(clauses, "cu.categoria = ?")
		args = append(args, p.Categoria)
	}
	if p.MacrorregiaoID > 0 {
		clauses = append(clauses, "cu.macrorregiao_id = ?")
		args = append(args, p.MacrorregiaoID)
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

	query := cursoSelect + " WHERE " + strings.Join(clauses, " AND ") +
		" GROUP BY cu.id ORDER BY cu.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Curso
	for rows.Next() {
		c, err := scanCurso(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CursoRepository) FindByID(id int64) (*models.Curso, error) {
	row := r.db.QueryRow(cursoSelect+" WHERE cu.id = ? AND cu.deleted_at IS NULL GROUP BY cu.id", id)
	c, err := scanCurso(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Curso não encontrado.", 404)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CursoRepository) Create(c *models.Curso) (*models.Curso, error) {
	if c.Titulo == "" {
		return nil, apperr.New("O campo título é obrigatório.", 422)
	}
	if c.Status == "" {
		c.Status = models.CursoStatusPlanejado
	}
	res, err := r.db.Exec(
		`INSERT INTO curso (titulo, descricao, categoria, publico_alvo, carga_horaria, modulos_total, capa_url, status, macrorregiao_id, criado_por, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		c.Titulo, c.Descricao, c.Categoria, c.PublicoAlvo, c.CargaHoraria, c.ModulosTotal, c.CapaURL, c.Status, c.MacrorregiaoID, c.CriadoPor,
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

// ---------------------------------------------------------------------
// MATRÍCULA
// ---------------------------------------------------------------------

const matriculaSelect = `
	SELECT ma.id, ma.usuario_id, ma.curso_id, ma.progresso, ma.status, ma.certificado_codigo,
	       ma.iniciado_em, ma.concluido_em, ma.created_at, ma.updated_at, cu.titulo
	FROM matricula ma
	JOIN curso cu ON cu.id = ma.curso_id
`

func scanMatricula(row interface{ Scan(...any) error }) (*models.Matricula, error) {
	m := &models.Matricula{}
	err := row.Scan(&m.ID, &m.UsuarioID, &m.CursoID, &m.Progresso, &m.Status, &m.CertificadoCodigo,
		&m.IniciadoEm, &m.ConcluidoEm, &m.CreatedAt, &m.UpdatedAt, &m.CursoTitulo)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListMinhas devolve as matrículas de um usuário — alimenta "meus cursos" e
// o cálculo de "profissionais capacitados" do Painel.
func (r *CursoRepository) ListMinhas(usuarioID int64) ([]*models.Matricula, error) {
	rows, err := r.db.Query(matriculaSelect+" WHERE ma.usuario_id = ? ORDER BY ma.created_at DESC", usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Matricula
	for rows.Next() {
		m, err := scanMatricula(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *CursoRepository) FindMatricula(usuarioID, cursoID int64) (*models.Matricula, error) {
	row := r.db.QueryRow(matriculaSelect+" WHERE ma.usuario_id = ? AND ma.curso_id = ? LIMIT 1", usuarioID, cursoID)
	m, err := scanMatricula(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Matricular inscreve o usuário no curso (idempotente: se já existir,
// devolve a matrícula existente em vez de duplicar).
func (r *CursoRepository) Matricular(usuarioID, cursoID int64) (*models.Matricula, error) {
	existing, err := r.FindMatricula(usuarioID, cursoID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	_, err = r.db.Exec(
		`INSERT INTO matricula (usuario_id, curso_id, progresso, status, iniciado_em, created_at, updated_at)
		 VALUES (?, ?, 0, ?, NOW(), NOW(), NOW())`,
		usuarioID, cursoID, models.MatriculaInscrito,
	)
	if err != nil {
		return nil, err
	}
	return r.FindMatricula(usuarioID, cursoID)
}

// AtualizarProgresso move o progresso (0-100) de uma matrícula. Ao alcançar
// 100, marca como concluída e gera um código de certificado — atende o
// requisito "com progresso e conclusão" do módulo AVA.
func (r *CursoRepository) AtualizarProgresso(usuarioID, cursoID int64, progresso int) (*models.Matricula, error) {
	if progresso < 0 {
		progresso = 0
	}
	if progresso > 100 {
		progresso = 100
	}
	existing, err := r.FindMatricula(usuarioID, cursoID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		if _, err := r.Matricular(usuarioID, cursoID); err != nil {
			return nil, err
		}
	}

	status := models.MatriculaEmAndamento
	var concluidoEm any = nil
	var certificado any = nil
	if progresso >= 100 {
		status = models.MatriculaConcluido
		concluidoEm = time.Now()
		codigo, err := gerarCodigoCertificado()
		if err != nil {
			return nil, err
		}
		certificado = codigo
	}

	_, err = r.db.Exec(
		`UPDATE matricula SET progresso = ?, status = ?, concluido_em = COALESCE(?, concluido_em),
		 certificado_codigo = COALESCE(certificado_codigo, ?), updated_at = NOW()
		 WHERE usuario_id = ? AND curso_id = ?`,
		progresso, status, concluidoEm, certificado, usuarioID, cursoID,
	)
	if err != nil {
		return nil, err
	}
	return r.FindMatricula(usuarioID, cursoID)
}

// FindMatriculaPorCertificado valida um certificado pelo código — permite
// que qualquer pessoa confirme a autenticidade de um certificado emitido.
func (r *CursoRepository) FindMatriculaPorCertificado(codigo string) (*models.Matricula, error) {
	row := r.db.QueryRow(matriculaSelect+" WHERE ma.certificado_codigo = ? LIMIT 1", codigo)
	m, err := scanMatricula(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Certificado não encontrado.", 404)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// TotalConcluidos conta matrículas concluídas — insumo do indicador
// "profissionais capacitados" do Painel.
func (r *CursoRepository) TotalConcluidos() (int64, error) {
	var total int64
	err := r.db.QueryRow("SELECT COUNT(id) FROM matricula WHERE status = ?", models.MatriculaConcluido).Scan(&total)
	return total, err
}

func gerarCodigoCertificado() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("PETCLIMA-%s", strings.ToUpper(hex.EncodeToString(b))), nil
}
