package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/models"
)

// ListParams controla filtros/paginação de GET /usuarios (usado pela
// "Gestão de usuários" do painel administrativo).
type ListParams struct {
	Order       string
	Sort        string
	Limit       int
	Page        int
	OnlyTrashed bool

	Nome        string
	Email       string
	Papel       string
	MunicipioID int64
	GatID       int64
}

var allowedOrderColumns = map[string]bool{
	"id": true, "nome": true, "email": true, "papel": true,
	"created_at": true, "updated_at": true,
}

type UsuarioRepository struct {
	db *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) buildWhere(p ListParams) (string, []any) {
	clauses := []string{}
	args := []any{}

	if p.OnlyTrashed {
		clauses = append(clauses, "u.deleted_at IS NOT NULL")
	} else {
		clauses = append(clauses, "u.deleted_at IS NULL")
	}

	if p.Nome != "" {
		clauses = append(clauses, "u.nome LIKE ?")
		args = append(args, "%"+p.Nome+"%")
	}
	if p.Email != "" {
		clauses = append(clauses, "u.email = ?")
		args = append(args, p.Email)
	}
	if p.Papel != "" {
		clauses = append(clauses, "u.papel = ?")
		args = append(args, p.Papel)
	}
	if p.MunicipioID > 0 {
		clauses = append(clauses, "u.municipio_id = ?")
		args = append(args, p.MunicipioID)
	}
	if p.GatID > 0 {
		clauses = append(clauses, "u.gat_id = ?")
		args = append(args, p.GatID)
	}

	return strings.Join(clauses, " AND "), args
}

// usuarioSelect junta cidade/macrorregiao/gat para já devolver os nomes
// legíveis usados pelo frontend (perfil, cadastro, gestão de usuários), sem
// exigir uma segunda chamada.
const usuarioSelect = `
	SELECT u.id, u.nome, u.email, u.senha, u.foto, u.cpf, u.telefone, u.instituicao,
	       u.curso_formacao, u.descricao, u.papel, u.municipio_id, u.gat_id,
	       u.email_verified_at, u.ultimo_acesso, u.aceite_termos_em,
	       u.created_at, u.updated_at, u.deleted_at,
	       c.nome, mr.id, mr.nome, g.nome
	FROM usuario u
	LEFT JOIN cidade c ON c.id = u.municipio_id
	LEFT JOIN macrorregiao mr ON mr.id = c.macrorregiao_id
	LEFT JOIN gat g ON g.id = u.gat_id
`

func scanUsuario(row interface{ Scan(...any) error }) (*models.Usuario, error) {
	u := &models.Usuario{}
	var municipioNome, macrorregiaoNome, gatNome sql.NullString
	var macrorregiaoID sql.NullInt64
	err := row.Scan(&u.ID, &u.Nome, &u.Email, &u.Senha, &u.Foto, &u.CPF, &u.Telefone, &u.Instituicao,
		&u.CursoFormacao, &u.Descricao, &u.Papel, &u.MunicipioID, &u.GatID,
		&u.EmailVerifiedAt, &u.UltimoAcesso, &u.AceiteTermosEm,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		&municipioNome, &macrorregiaoID, &macrorregiaoNome, &gatNome)
	if err != nil {
		return nil, err
	}
	if municipioNome.Valid {
		u.MunicipioNome = &municipioNome.String
	}
	if macrorregiaoID.Valid {
		u.MacrorregiaoID = &macrorregiaoID.Int64
	}
	if macrorregiaoNome.Valid {
		u.MacrorregiaoNome = &macrorregiaoNome.String
	}
	if gatNome.Valid {
		u.GatNome = &gatNome.String
	}
	return u, nil
}

// All lista usuários paginados, usado pela tela de Gestão de usuários.
func (r *UsuarioRepository) All(p ListParams) ([]*models.Usuario, int64, error) {
	order := p.Order
	if order == "" || !allowedOrderColumns[order] {
		order = "id"
	}
	sortDir := "asc"
	if strings.EqualFold(p.Sort, "desc") {
		sortDir = "desc"
	}

	where, args := r.buildWhere(p)

	total, err := r.countWithWhere(where, args)
	if err != nil {
		return nil, 0, err
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 20
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf("%s WHERE %s ORDER BY u.%s %s LIMIT ? OFFSET ?", usuarioSelect, where, order, sortDir)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.Usuario
	for rows.Next() {
		u, err := scanUsuario(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, u)
	}
	return list, total, rows.Err()
}

func (r *UsuarioRepository) countWithWhere(where string, args []any) (int64, error) {
	var total int64
	query := "SELECT COUNT(u.id) FROM usuario u WHERE " + where
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

// Count mirrors BaseListRepository::count.
func (r *UsuarioRepository) Count(p ListParams) (int64, error) {
	where, args := r.buildWhere(p)
	return r.countWithWhere(where, args)
}

// PapelCount é o formato de resposta de GroupByPapel, usado no Painel para
// mostrar quantos usuários existem em cada perfil (profissional, estudante
// etc).
type PapelCount struct {
	Papel string `json:"papel" db:"papel"`
	Total int64  `json:"total" db:"total"`
}

func (r *UsuarioRepository) GroupByPapel() ([]PapelCount, error) {
	rows, err := r.db.Query("SELECT papel, COUNT(id) AS total FROM usuario WHERE deleted_at IS NULL GROUP BY papel")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PapelCount
	for rows.Next() {
		var pc PapelCount
		if err := rows.Scan(&pc.Papel, &pc.Total); err != nil {
			return nil, err
		}
		out = append(out, pc)
	}
	return out, rows.Err()
}

// FindByID mirrors BaseListRepository::find.
func (r *UsuarioRepository) FindByID(id int64, trashed bool) (*models.Usuario, error) {
	where := "u.id = ? AND u.deleted_at IS NULL"
	if trashed {
		where = "u.id = ?"
	}
	query := fmt.Sprintf("%s WHERE %s LIMIT 1", usuarioSelect, where)
	row := r.db.QueryRow(query, id)
	u, err := scanUsuario(row)
	if err == sql.ErrNoRows {
		return nil, apperr.New("Usuário não encontrado(a)!", 404)
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail is used by the login flow. O e-mail é único na tabela, então
// não precisa de mais nenhum filtro para desambiguar.
func (r *UsuarioRepository) FindByEmail(email string) (*models.Usuario, error) {
	query := usuarioSelect + " WHERE u.email = ? AND u.deleted_at IS NULL LIMIT 1"
	row := r.db.QueryRow(query, email)
	u, err := scanUsuario(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// SetSenhaByEmail overwrites the password for the account matching email,
// used by the "esqueci minha senha" reset flow.
func (r *UsuarioRepository) SetSenhaByEmail(email, plainSenha string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainSenha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		"UPDATE usuario SET senha = ?, email_verified_at = COALESCE(email_verified_at, NOW()), updated_at = NOW() WHERE email = ? AND deleted_at IS NULL",
		string(hash), email,
	)
	return err
}

// SetEmailVerifiedByEmail marks the account as active once the person
// confirms the token sent to their e-mail (see handlers.ConfirmarEmail).
func (r *UsuarioRepository) SetEmailVerifiedByEmail(email string) error {
	_, err := r.db.Exec(
		"UPDATE usuario SET email_verified_at = NOW(), updated_at = NOW() WHERE email = ?",
		email,
	)
	return err
}

func (r *UsuarioRepository) EmailExists(email string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(id) FROM usuario WHERE email = ? AND id != ?",
		email, excludeID,
	).Scan(&count)
	return count > 0, err
}

// Create insere um novo usuário. A validação de qual `papel` é permitido em
// autocadastro público (PapeisAutoCadastro) fica a cargo do handler, que
// conhece o contexto (cadastro público vs. criado por um admin).
func (r *UsuarioRepository) Create(u *models.Usuario, plainSenha string) (*models.Usuario, error) {
	if u.Nome == "" {
		return nil, apperr.New("O campo nome é obrigatório.", 422)
	}
	if u.Email == "" {
		return nil, apperr.New("O campo email é obrigatório.", 422)
	}
	if !models.PapelValido(u.Papel) {
		return nil, apperr.New("Perfil (papel) inválido.", 422)
	}
	if len(plainSenha) < 6 {
		return nil, apperr.New("A senha deve ter no mínimo 6 caracteres.", 422)
	}
	exists, err := r.EmailExists(u.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.New("O campo email já está sendo utilizado.", 422)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainSenha), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	res, err := r.db.Exec(
		`INSERT INTO usuario (nome, email, senha, foto, cpf, telefone, instituicao, curso_formacao, papel, municipio_id, gat_id, aceite_termos_em, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())`,
		u.Nome, u.Email, string(hash), u.Foto, u.CPF, u.Telefone, u.Instituicao, u.CursoFormacao, u.Papel, u.MunicipioID, u.GatID,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(id, false)
}

// FindByIDs busca vários usuários de uma vez, preservando o mapa por id —
// útil para juntar nome/foto do responsável em listas de atividades.
func (r *UsuarioRepository) FindByIDs(ids []int64) (map[int64]*models.Usuario, error) {
	out := map[int64]*models.Usuario{}
	if len(ids) == 0 {
		return out, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := usuarioSelect + " WHERE u.id IN (" + strings.Join(placeholders, ", ") + ")"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u, err := scanUsuario(rows)
		if err != nil {
			return nil, err
		}
		out[u.ID] = u
	}
	return out, rows.Err()
}

// Update altera os dados de perfil de um usuário (nome, foto, bio, vínculo
// institucional). A troca de `papel` passa por AlterarPapel, não por aqui,
// para deixar claro no código/logs quando uma permissão está sendo elevada.
func (r *UsuarioRepository) Update(id int64, u *models.Usuario, plainSenha string) (*models.Usuario, error) {
	existing, err := r.FindByID(id, false)
	if err != nil {
		return nil, err
	}

	if u.Email != "" && u.Email != existing.Email {
		exists, err := r.EmailExists(u.Email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperr.New("O campo email já está sendo utilizado.", 422)
		}
		existing.Email = u.Email
	}
	if u.Nome != "" {
		existing.Nome = u.Nome
	}
	if u.Foto != nil {
		existing.Foto = u.Foto
	}
	if u.CPF != nil {
		existing.CPF = u.CPF
	}
	if u.Telefone != nil {
		existing.Telefone = u.Telefone
	}
	if u.Instituicao != nil {
		existing.Instituicao = u.Instituicao
	}
	if u.CursoFormacao != nil {
		existing.CursoFormacao = u.CursoFormacao
	}
	if u.Descricao != nil {
		existing.Descricao = u.Descricao
	}
	if u.MunicipioID != nil {
		existing.MunicipioID = u.MunicipioID
	}
	if u.GatID != nil {
		existing.GatID = u.GatID
	}

	senha := existing.Senha
	if plainSenha != "" {
		if len(plainSenha) < 6 {
			return nil, apperr.New("A senha deve ter no mínimo 6 caracteres.", 422)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(plainSenha), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		senha = string(hash)
	}

	_, err = r.db.Exec(
		`UPDATE usuario SET nome = ?, email = ?, senha = ?, foto = ?, cpf = ?, telefone = ?, instituicao = ?,
		 curso_formacao = ?, descricao = ?, municipio_id = ?, gat_id = ?, updated_at = NOW()
		 WHERE id = ?`,
		existing.Nome, existing.Email, senha, existing.Foto, existing.CPF, existing.Telefone, existing.Instituicao,
		existing.CursoFormacao, existing.Descricao, existing.MunicipioID, existing.GatID, id,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(id, false)
}

// AlterarPapel troca o perfil/papel de um usuário — usado pelo módulo
// "👥 Gestão de usuários" para promover alguém a preceptor/coordenador/admin
// (o "convite institucional" do planejamento).
func (r *UsuarioRepository) AlterarPapel(id int64, papel models.Papel) error {
	if !models.PapelValido(papel) {
		return apperr.New("Perfil (papel) inválido.", 422)
	}
	res, err := r.db.Exec("UPDATE usuario SET papel = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL", papel, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Usuário não encontrado(a)!", 404)
	}
	return nil
}

// SoftDelete mirrors BaseCrudRepository::delete (sets deleted_at).
func (r *UsuarioRepository) SoftDelete(id int64) error {
	res, err := r.db.Exec("UPDATE usuario SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Usuário não encontrado(a)!", 404)
	}
	return nil
}

// Restore mirrors BaseCrudRepository::restore (clears deleted_at).
func (r *UsuarioRepository) Restore(id int64) error {
	res, err := r.db.Exec("UPDATE usuario SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.New("Usuário não encontrado(a)!", 404)
	}
	return nil
}

// CheckPassword compares a plaintext password against the stored bcrypt hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// TouchAcesso registra "visto por último agora" pro usuário — chamado
// periodicamente pelo frontend enquanto o app está aberto numa aba.
func (r *UsuarioRepository) TouchAcesso(usuarioID int64) error {
	_, err := r.db.Exec(`UPDATE usuario SET ultimo_acesso = NOW() WHERE id = ?`, usuarioID)
	return err
}

// RegistrarAcesso grava uma linha em acesso_log (usada pelo Painel para
// calcular o indicador de engajamento mensal) toda vez que um login é bem
// sucedido.
func (r *UsuarioRepository) RegistrarAcesso(usuarioID int64, ip string) error {
	_, err := r.db.Exec(`INSERT INTO acesso_log (usuario_id, ip, created_at) VALUES (?, ?, NOW())`, usuarioID, ip)
	return err
}
