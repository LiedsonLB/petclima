package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// NoticiaHandler alimenta o Portal Público (🏠). Leitura é aberta; criar e
// remover exige staff.
type NoticiaHandler struct {
	repo *repository.NoticiaRepository
}

func NewNoticiaHandler(repo *repository.NoticiaRepository) *NoticiaHandler {
	return &NoticiaHandler{repo: repo}
}

// All handles GET /noticias (?limit=)
func (h *NoticiaHandler) All(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.All(atoiDefault(r.URL.Query().Get("limit"), 20))
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}

// Find handles GET /noticias/{id}
func (h *NoticiaHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	n, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 200, n)
}

type noticiaPayload struct {
	Titulo   string  `json:"titulo"`
	Resumo   *string `json:"resumo"`
	Conteudo string  `json:"conteudo"`
	CapaURL  *string `json:"capa_url"`
}

// Save handles POST /noticias (staff).
func (h *NoticiaHandler) Save(w http.ResponseWriter, r *http.Request) {
	var p noticiaPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	n := &models.Noticia{
		Titulo:   strings.TrimSpace(p.Titulo),
		Resumo:   nilIfEmpty(p.Resumo),
		Conteudo: strings.TrimSpace(p.Conteudo),
		CapaURL:  nilIfEmpty(p.CapaURL),
	}
	if u := middleware.UserFromContext(r); u != nil {
		n.AutorID = &u.ID
	}
	created, err := h.repo.Create(n)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, created)
}

// Delete handles DELETE /noticias/{id} (staff).
func (h *NoticiaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	if err := h.repo.SoftDelete(id); err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.Success(w, "Sucesso!", 200)
}
