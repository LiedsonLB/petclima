package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

type UsuarioHandler struct {
	repo      *repository.UsuarioRepository
	uploadDir string // diretório onde as fotos são armazenadas
}

func NewUsuarioHandler(repo *repository.UsuarioRepository, uploadDir string) *UsuarioHandler {
	return &UsuarioHandler{repo: repo, uploadDir: uploadDir}
}

func (h *UsuarioHandler) paramsFromRequest(r *http.Request) repository.ListParams {
	q := r.URL.Query()
	return repository.ListParams{
		Order:       q.Get("order"),
		Sort:        q.Get("sort"),
		Limit:       atoiDefault(q.Get("limit"), 20),
		Page:        atoiDefault(q.Get("page"), 1),
		OnlyTrashed: q.Get("onlyTrashed") == "1" || q.Get("onlyTrashed") == "sim",
		Nome:        q.Get("nome"),
		Email:       q.Get("email"),
		Papel:       q.Get("papel"),
		MunicipioID: int64(atoiDefault(q.Get("municipio_id"), 0)),
		GatID:       int64(atoiDefault(q.Get("gat_id"), 0)),
	}
}

// All handles GET /usuarios (routeList's "all")
func (h *UsuarioHandler) All(w http.ResponseWriter, r *http.Request) {
	p := h.paramsFromRequest(r)
	list, total, err := h.repo.All(p)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	lastPage := int((total + int64(p.Limit) - 1) / int64(p.Limit))
	if lastPage < 1 {
		lastPage = 1
	}
	httpx.JSON(w, 200, httpx.Paginated{
		Data: list, CurrentPage: p.Page, PerPage: p.Limit, Total: total, LastPage: lastPage,
	})
}

// Count handles GET /usuarios/count
func (h *UsuarioHandler) Count(w http.ResponseWriter, r *http.Request) {
	p := h.paramsFromRequest(r)
	total, err := h.repo.Count(p)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, map[string]any{"total": total})
}

func idFromPath(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// Heartbeat handles POST /usuarios/heartbeat — chamado pelo frontend a
// cada ~30s enquanto o usuário está com a plataforma aberta, só para
// atualizar usuario.ultimo_acesso (ver TouchAcesso).
func (h *UsuarioHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}
	if err := h.repo.TouchAcesso(session.ID); err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.Success(w, "ok", 200)
}

// ResumoPapeis handles GET /usuarios/resumo — quantos usuários existem em
// cada perfil (painel de gestão de usuários).
func (h *UsuarioHandler) ResumoPapeis(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.GroupByPapel()
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}

// Find handles GET /usuarios/{id}
func (h *UsuarioHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	u, err := h.repo.FindByID(id, false)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 200, u)
}

type usuarioPayload struct {
	Nome          string  `json:"nome"`
	Email         string  `json:"email"`
	Senha         string  `json:"senha"`
	Foto          *string `json:"foto"`
	CPF           *string `json:"cpf"`
	Telefone      *string `json:"telefone"`
	Instituicao   *string `json:"instituicao"`
	CursoFormacao *string `json:"curso_formacao"`
	Descricao     *string `json:"descricao"`
	MunicipioID   *int64  `json:"municipio_id"`
	GatID         *int64  `json:"gat_id"`
	Papel         string  `json:"papel"`
}

// nilIfEmpty evita gravar string vazia como se fosse um valor real — o
// repositório trata ponteiro nil como "não mudou esse campo" (ver
// UsuarioRepository.Update), então "" (campo ausente/limpo no formulário)
// vira nil em vez de sobrescrever com uma string vazia.
func nilIfEmpty(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

// Save handles POST /usuarios — criação de conta por um administrador (o
// "convite institucional" do planejamento). Aceita qualquer papel. A conta
// já nasce com o e-mail verificado, pois foi cadastrada por quem administra
// a plataforma; a pessoa só precisa usar "esqueci minha senha" ou receber a
// senha inicial combinada.
func (h *UsuarioHandler) Save(w http.ResponseWriter, r *http.Request) {
	var payload usuarioPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}

	papel := models.Papel(strings.ToLower(strings.TrimSpace(payload.Papel)))
	if papel == "" {
		papel = models.PapelComunidade
	}
	if !models.PapelValido(papel) {
		httpx.Error(w, "Perfil (papel) inválido.", 422)
		return
	}

	u := &models.Usuario{
		Nome:          strings.TrimSpace(payload.Nome),
		Email:         strings.ToLower(strings.TrimSpace(payload.Email)),
		Foto:          nilIfEmpty(payload.Foto),
		CPF:           nilIfEmpty(payload.CPF),
		Telefone:      nilIfEmpty(payload.Telefone),
		Instituicao:   nilIfEmpty(payload.Instituicao),
		CursoFormacao: nilIfEmpty(payload.CursoFormacao),
		Papel:         papel,
		MunicipioID:   payload.MunicipioID,
		GatID:         payload.GatID,
	}
	created, err := h.repo.Create(u, payload.Senha)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if err := h.repo.SetEmailVerifiedByEmail(created.Email); err != nil {
		writeAppErr(w, err)
		return
	}
	created, err = h.repo.FindByID(created.ID, false)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, created)
}

type papelPayload struct {
	Papel string `json:"papel"`
}

// AlterarPapel handles PUT /usuarios/{id}/papel — só administrador. É o
// caminho para promover alguém a preceptor/coordenador/admin.
func (h *UsuarioHandler) AlterarPapel(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	var payload papelPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	session := middleware.UserFromContext(r)
	if session != nil && session.ID == id {
		httpx.Error(w, "Você não pode alterar o seu próprio perfil de acesso.", 422)
		return
	}
	papel := models.Papel(strings.ToLower(strings.TrimSpace(payload.Papel)))
	if err := h.repo.AlterarPapel(id, papel); err != nil {
		writeAppErr(w, err)
		return
	}
	updated, err := h.repo.FindByID(id, false)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 200, updated)
}

// deleteOldPhoto apaga a foto antiga do disco se existir e for diferente da nova
func (h *UsuarioHandler) deleteOldPhoto(oldFotoURL string, newFotoURL *string) {
	// Se não tem foto antiga, não faz nada
	if oldFotoURL == "" {
		return
	}

	// Se a nova foto é igual à antiga, não apaga
	if newFotoURL != nil && *newFotoURL == oldFotoURL {
		return
	}

	// Converte a URL para caminho físico
	// Ex: "/uploads/fotos/abc123.jpg" -> "fotos/abc123.jpg"
	relativePath := strings.TrimPrefix(oldFotoURL, "/uploads/")
	if relativePath == oldFotoURL {
		// Se não começa com /uploads/, tenta remover apenas o prefixo /uploads
		relativePath = strings.TrimPrefix(oldFotoURL, "uploads/")
	}

	if relativePath == "" || relativePath == oldFotoURL {
		// Se não conseguiu extrair o caminho, não apaga por segurança
		return
	}

	// Monta o caminho completo no sistema de arquivos
	fullPath := filepath.Join(h.uploadDir, relativePath)

	// Verifica se o arquivo existe antes de tentar apagar
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return // arquivo já foi apagado ou não existe
	}

	// Tenta apagar o arquivo (ignora erros)
	_ = os.Remove(fullPath)
}

// Update handles PUT /usuarios/{id}. Cada pessoa edita o próprio perfil; só
// administradores editam o de outras pessoas. A troca de papel NÃO passa
// por aqui (ver AlterarPapel).
func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}

	session := middleware.UserFromContext(r)
	if session == nil || (session.ID != id && !session.IsAdmin()) {
		httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
		return
	}

	// Busca o usuário atual para obter a foto antiga
	oldUser, err := h.repo.FindByID(id, false)
	if err != nil {
		writeAppErr(w, err)
		return
	}

	var payload usuarioPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}

	u := &models.Usuario{
		Nome:          strings.TrimSpace(payload.Nome),
		Email:         strings.ToLower(strings.TrimSpace(payload.Email)),
		Foto:          nilIfEmpty(payload.Foto),
		CPF:           nilIfEmpty(payload.CPF),
		Telefone:      nilIfEmpty(payload.Telefone),
		Instituicao:   nilIfEmpty(payload.Instituicao),
		CursoFormacao: nilIfEmpty(payload.CursoFormacao),
		Descricao:     nilIfEmpty(payload.Descricao),
		MunicipioID:   payload.MunicipioID,
		GatID:         payload.GatID,
	}

	updated, err := h.repo.Update(id, u, payload.Senha)
	if err != nil {
		writeAppErr(w, err)
		return
	}

	// Apaga a foto antiga do disco se foi substituída
	if oldUser != nil && oldUser.Foto != nil {
		h.deleteOldPhoto(*oldUser.Foto, u.Foto)
	}

	httpx.JSON(w, 200, updated)
}

// Delete handles DELETE /usuarios/{id}
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}

	// Busca o usuário para obter a foto antes de deletar
	user, err := h.repo.FindByID(id, false)
	if err == nil && user != nil && user.Foto != nil {
		// Apaga a foto do disco
		h.deleteOldPhoto(*user.Foto, nil)
	}

	if err := h.repo.SoftDelete(id); err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.Success(w, "Sucesso!", 200)
}

// Restore handles POST /usuarios/restore/{id}
func (h *UsuarioHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	if err := h.repo.Restore(id); err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.Success(w, "Sucesso!", 200)
}

func writeAppErr(w http.ResponseWriter, err error) {
	if ae, ok := apperr.As(err); ok {
		httpx.Error(w, ae.Message, ae.Status)
		return
	}
	httpx.Error(w, "Erro interno.", 500)
}
