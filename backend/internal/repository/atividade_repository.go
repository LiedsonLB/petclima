package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

type AtividadeListParams struct {
	Q              string
	Tipos          []string
	MunicipioID    int64
	MacrorregiaoID int64
	GatID          int64
	ApenasFuturas  bool
	Limit          int
	Page           int
}

type AtividadeRepository struct {
	db *sql.DB
}

func NewAtividadeRepository(db *sql.DB) *AtividadeRepository {
	return &AtividadeRepository{db: db}
}

const atividadeSelect = `
	SELECT a.id, a.titulo, a.tipo, a.descricao, a.modalidade, a.municipio_id, a.macrorregiao_id, a.gat_id,
	       a.responsavel_id, a.data_realizacao, a.carga_horaria, a.vagas, a.sala_id, a.created_at, a.updated_at,
	       c.nome, mr.nome, g.nome, u.nome, s.codigo,
	       (SELECT COUNT(i.id) FROM atividade_inscricao i WHERE i.atividade_id = a.id) AS inscritos
	FROM atividade a
	LEFT JOIN cidade c ON c.id = a.municipio_id
	LEFT JOIN macrorregiao mr ON mr.id = a.macrorregiao_id
	LEFT JOIN gat g ON g.id = a.gat_id
	LEFT JOIN usuario u ON u.id = a.responsavel_id
	LEFT JOIN sala s ON s.id = a.sala_id
`

func scanAtividade(row interface{ Scan(...any) error }) (*models.Atividade, error) {
	a := &models.Atividade{}
	var cNome, mrNome, gNome, uNome, sCodigo sql.NullString
	err := row.Scan(&a.ID, &a.Titulo, &a.Tipo, &a.Descricao, &a.Modalidade, &a.MunicipioID, &a.MacrorregiaoID, &a.GatID,
		&a.ResponsavelID, &a.DataRealizacao, &a.CargaHoraria, &a.Vagas, &a.SalaID, &a.CreatedAt, &a.UpdatedAt,
		&cNome, &mrNome, &gNome, &uNome, &sCodigo, &a.Inscritos)
	if err != nil {
		return nil, err
	}
	if cNome.Valid {
		a.MunicipioNome = &cNome.String
	}
	if mrNome.Valid {
		a.MacrorregiaoNome = &mrNome.String
	}
	if gNome.Valid {
		a.GatNome = &gNome.String
	}
	if uNome.Valid {
		a.ResponsavelNome = &uNome.String
	}
	if sCodigo.Valid {
		a.SalaCodigo = &sCodigo.String
	}
	return a, nil
}

func (r *AtividadeRepository) All(p AtividadeListParams) ([]*models.Atividade, error) {
	clauses := []string{"a.deleted_at IS NULL"}
	args := []any{}

	if p.Q != "" {
		clauses = append(clauses, "(a.titulo LIKE ? OR a.descricao LIKE ?)")
		like := "%" + p.Q + "%"
		args = append(args, like, like)
	}
	if len(p.Tipos) > 0 {
		ph := make([]string, len(p.Tipos))
		for i, t := range p.Tipos {
			ph[i] = "?"
			args = append(args, t)
		}
		clauses = append(clauses, "a.tipo IN ("+strings.Join(ph, ",")+")")
	}
	if p.MunicipioID > 0 {
		clauses = append(clauses, "a.municipio_id = ?")
		args = append(args, p.MunicipioID)
	}
	if p.MacrorregiaoID > 0 {
		clauses = append(clauses, "a.macrorregiao_id = ?")
		args = append(args, p.MacrorregiaoID)
	}
	if p.GatID > 0 {
		clauses = append(clauses, "a.gat_id = ?")
		args = append(args, p.GatID)
	}
	if p.ApenasFuturas {
		clauses = append(clauses, "(a.data_realizacao IS NULL OR a.data_realizacao >= NOW())")
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

	query := atividadeSelect + " WHERE " + strings.Join(clauses, " AND ") + " ORDER BY a.data_realizacao ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Atividade
	for rows.Next() {
		a, err := scanAtividade(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AtividadeRepository) FindByID(id int64) (*models.Atividade, error) {
	row := r.db.QueryRow(atividadeSelect+" WHERE a.id = ? AND a.deleted_at IS NULL LIMIT 1", id)
	a, err := scanAtividade(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Atividade não encontrada.", 404)
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AtividadeRepository) Create(a *models.Atividade) (*models.Atividade, error) {
	if a.Titulo == "" {
		return nil, apperr.New("O campo título é obrigatório.", 422)
	}
	if a.Tipo == "" {
		a.Tipo = models.AtividadeOficina
	}
	if a.Modalidade == "" {
		a.Modalidade = models.ModalidadePresencial
	}
	// Se o município foi informado mas a macrorregião não, herda do município.
	if a.MacrorregiaoID == nil && a.MunicipioID != nil {
		var mr sql.NullInt64
		_ = r.db.QueryRow("SELECT macrorregiao_id FROM cidade WHERE id = ?", *a.MunicipioID).Scan(&mr)
		if mr.Valid {
			a.MacrorregiaoID = &mr.Int64
		}
	}
	res, err := r.db.Exec(
		`INSERT INTO atividade (titulo, tipo, descricao, modalidade, municipio_id, macrorregiao_id, gat_id, responsavel_id, data_realizacao, carga_horaria, vagas, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		a.Titulo, a.Tipo, a.Descricao, a.Modalidade, a.MunicipioID, a.MacrorregiaoID, a.GatID, a.ResponsavelID, a.DataRealizacao, a.CargaHoraria, a.Vagas,
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

// VincularSala liga uma sala de videoconferência a uma atividade.
func (r *AtividadeRepository) VincularSala(atividadeID, salaID int64) error {
	_, err := r.db.Exec("UPDATE atividade SET sala_id = ?, updated_at = NOW() WHERE id = ?", salaID, atividadeID)
	return err
}

func (r *AtividadeRepository) SoftDelete(id int64) error {
	res, err := r.db.Exec("UPDATE atividade SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Atividade não encontrada.", 404)
	}
	return nil
}

// ---------------------------------------------------------------------
// INSCRIÇÕES
// ---------------------------------------------------------------------

func (r *AtividadeRepository) FindInscricao(atividadeID, usuarioID int64) (*models.AtividadeInscricao, error) {
	i := &models.AtividadeInscricao{}
	err := r.db.QueryRow(
		"SELECT id, atividade_id, usuario_id, presente, inscrito_em FROM atividade_inscricao WHERE atividade_id = ? AND usuario_id = ?",
		atividadeID, usuarioID,
	).Scan(&i.ID, &i.AtividadeID, &i.UsuarioID, &i.Presente, &i.InscritoEm)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return i, nil
}

// Inscrever registra a inscrição, respeitando o limite de vagas.
func (r *AtividadeRepository) Inscrever(atividadeID, usuarioID int64) (*models.AtividadeInscricao, error) {
	a, err := r.FindByID(atividadeID)
	if err != nil {
		return nil, err
	}
	if existing, err := r.FindInscricao(atividadeID, usuarioID); err != nil {
		return nil, err
	} else if existing != nil {
		return existing, nil
	}
	if a.Vagas != nil && a.Inscritos >= *a.Vagas {
		return nil, apperr.New("Não há mais vagas para esta atividade.", 422)
	}
	if a.DataRealizacao != nil && a.DataRealizacao.Before(time.Now().Add(-24*time.Hour)) {
		return nil, apperr.New("Esta atividade já foi realizada.", 422)
	}
	if _, err := r.db.Exec(
		"INSERT INTO atividade_inscricao (atividade_id, usuario_id, presente, inscrito_em) VALUES (?, ?, 0, NOW())",
		atividadeID, usuarioID,
	); err != nil {
		return nil, err
	}
	return r.FindInscricao(atividadeID, usuarioID)
}

func (r *AtividadeRepository) CancelarInscricao(atividadeID, usuarioID int64) error {
	_, err := r.db.Exec("DELETE FROM atividade_inscricao WHERE atividade_id = ? AND usuario_id = ?", atividadeID, usuarioID)
	return err
}

// ListInscritos devolve os inscritos de uma atividade com o nome de cada um
// (usado pelo preceptor/coordenador para registrar presença).
func (r *AtividadeRepository) ListInscritos(atividadeID int64) ([]*models.AtividadeInscricao, error) {
	rows, err := r.db.Query(`
		SELECT i.id, i.atividade_id, i.usuario_id, i.presente, i.inscrito_em, u.nome
		FROM atividade_inscricao i
		JOIN usuario u ON u.id = i.usuario_id
		WHERE i.atividade_id = ?
		ORDER BY u.nome
	`, atividadeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.AtividadeInscricao
	for rows.Next() {
		i := &models.AtividadeInscricao{}
		var nome sql.NullString
		if err := rows.Scan(&i.ID, &i.AtividadeID, &i.UsuarioID, &i.Presente, &i.InscritoEm, &nome); err != nil {
			return nil, err
		}
		if nome.Valid {
			i.UsuarioNome = &nome.String
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func (r *AtividadeRepository) MarcarPresenca(atividadeID, usuarioID int64, presente bool) error {
	res, err := r.db.Exec(
		"UPDATE atividade_inscricao SET presente = ? WHERE atividade_id = ? AND usuario_id = ?",
		presente, atividadeID, usuarioID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Inscrição não encontrada.", 404)
	}
	return nil
}

// TotalPresentes soma presenças confirmadas — insumo do indicador
// "profissionais capacitados" do Painel (junto com cursos concluídos).
func (r *AtividadeRepository) TotalPresentes() (int64, error) {
	var total int64
	err := r.db.QueryRow("SELECT COUNT(DISTINCT usuario_id) FROM atividade_inscricao WHERE presente = 1").Scan(&total)
	return total, err
}

// InscricoesDoUsuario devolve, indexadas por atividade, as inscrições de um
// usuário — usado para marcar "minha_inscricao" nas listagens sem uma
// consulta por atividade.
func (r *AtividadeRepository) InscricoesDoUsuario(usuarioID int64) (map[int64]*models.AtividadeInscricao, error) {
	rows, err := r.db.Query(
		"SELECT id, atividade_id, usuario_id, presente, inscrito_em FROM atividade_inscricao WHERE usuario_id = ?",
		usuarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]*models.AtividadeInscricao{}
	for rows.Next() {
		i := &models.AtividadeInscricao{}
		if err := rows.Scan(&i.ID, &i.AtividadeID, &i.UsuarioID, &i.Presente, &i.InscritoEm); err != nil {
			return nil, err
		}
		out[i.AtividadeID] = i
	}
	return out, rows.Err()
}

// ListMinhas devolve as atividades em que o usuário está inscrito.
func (r *AtividadeRepository) ListMinhas(usuarioID int64) ([]*models.Atividade, error) {
	rows, err := r.db.Query(
		atividadeSelect+" JOIN atividade_inscricao ii ON ii.atividade_id = a.id WHERE ii.usuario_id = ? AND a.deleted_at IS NULL ORDER BY a.data_realizacao ASC",
		usuarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Atividade
	for rows.Next() {
		a, err := scanAtividade(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
