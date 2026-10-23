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

// MaterialHandler cobre a Biblioteca/Repositório e o kit de Comunicação
// (mesma tabela, filtrada por tipo). Leitura e download são públicos — o
// planejamento prevê um repositório aberto — e criação/remoção exigem staff.
type MaterialHandler struct {
	repo *repository.MaterialRepository
}

func NewMaterialHandler(repo *repository.MaterialRepository) *MaterialHandler {
	return &MaterialHandler{repo: repo}
}

var tiposMaterialValidos = map[string]bool{
	models.MaterialCartilha:    true,
	models.MaterialProtocolo:   true,
	models.MaterialGuia:        true,
	models.MaterialVideo:       true,
	models.MaterialPodcast:     true,
	models.MaterialCard:        true,
	models.MaterialInfografico: true,
	models.MaterialRadio:       true,
	models.MaterialPDF:         true,
}

// splitCSV transforma "card,video, podcast" em ["card","video","podcast"].
func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func materialParams(r *http.Request) repository.MaterialListParams {
	q := r.URL.Query()
	return repository.MaterialListParams{
		Q:              strings.TrimSpace(q.Get("q")),
		Tipos:          splitCSV(q.Get("tipo")),
		Tema:           q.Get("tema"),
		PublicoAlvo:    q.Get("publico_alvo"),
		MacrorregiaoID: int64(atoiDefault(q.Get("macrorregiao_id"), 0)),
		MunicipioID:    int64(atoiDefault(q.Get("municipio_id"), 0)),
		GatID:          int64(atoiDefault(q.Get("gat_id"), 0)),
		Limit:          atoiDefault(q.Get("limit"), 50),
		Page:           atoiDefault(q.Get("page"), 1),
	}
}

// All handles GET /materiais
func (h *MaterialHandler) All(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.All(materialParams(r))
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}

// Find handles GET /materiais/{id}
func (h *MaterialHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	m, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 200, m)
}

type materialPayload struct {
	Titulo         string  `json:"titulo"`
	Descricao      *string `json:"descricao"`
	Tipo           string  `json:"tipo"`
	Tema           *string `json:"tema"`
	PublicoAlvo    *string `json:"publico_alvo"`
	ArquivoURL     string  `json:"arquivo_url"`
	CapaURL        *string `json:"capa_url"`
	MacrorregiaoID *int64  `json:"macrorregiao_id"`
	MunicipioID    *int64  `json:"municipio_id"`
	GatID          *int64  `json:"gat_id"`
}

// Save handles POST /materiais (staff).
func (h *MaterialHandler) Save(w http.ResponseWriter, r *http.Request) {
	var p materialPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	tipo := strings.ToLower(strings.TrimSpace(p.Tipo))
	if tipo == "" {
		tipo = models.MaterialPDF
	}
	if !tiposMaterialValidos[tipo] {
		httpx.Error(w, "Tipo de material inválido.", 422)
		return
	}

	session := middleware.UserFromContext(r)
	m := &models.Material{
		Titulo:         strings.TrimSpace(p.Titulo),
		Descricao:      nilIfEmpty(p.Descricao),
		Tipo:           tipo,
		Tema:           nilIfEmpty(p.Tema),
		PublicoAlvo:    nilIfEmpty(p.PublicoAlvo),
		ArquivoURL:     strings.TrimSpace(p.ArquivoURL),
		CapaURL:        nilIfEmpty(p.CapaURL),
		MacrorregiaoID: p.MacrorregiaoID,
		MunicipioID:    p.MunicipioID,
		GatID:          p.GatID,
	}
	if session != nil {
		m.PublicadoPor = &session.ID
	}
	created, err := h.repo.Create(m)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, created)
}

// Delete handles DELETE /materiais/{id} (staff).
func (h *MaterialHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// Download handles POST /materiais/{id}/download — registra o download
// (contador + log usado no Painel) e devolve a URL do arquivo. Público: se
// houver login, o download fica associado ao usuário.
func (h *MaterialHandler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	m, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	var usuarioID int64
	if u := middleware.UserFromContext(r); u != nil {
		usuarioID = u.ID
	}
	if err := h.repo.RegistrarDownload(id, usuarioID); err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, map[string]any{
		"id":          m.ID,
		"arquivo_url": m.ArquivoURL,
		"downloads":   m.Downloads + 1,
	})
}
