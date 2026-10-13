package handlers

import (
	"net/http"
	"strings"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// BuscaHandler implementa o módulo "🔎 Busca": procura materiais, cursos e
// atividades de uma vez, filtrando por tema, região, público, tipo e GAT.
type BuscaHandler struct {
	materialRepo  *repository.MaterialRepository
	cursoRepo     *repository.CursoRepository
	atividadeRepo *repository.AtividadeRepository
}

func NewBuscaHandler(m *repository.MaterialRepository, c *repository.CursoRepository, a *repository.AtividadeRepository) *BuscaHandler {
	return &BuscaHandler{materialRepo: m, cursoRepo: c, atividadeRepo: a}
}

type buscaResultado struct {
	Materiais  []*models.Material  `json:"materiais"`
	Cursos     []*models.Curso     `json:"cursos"`
	Atividades []*models.Atividade `json:"atividades"`
}

// containsFold: busca simples, sem acentos-insensível, suficiente para
// filtrar a lista curta de cursos em memória.
func containsFold(s *string, term string) bool {
	if s == nil {
		return false
	}
	return strings.Contains(strings.ToLower(*s), term)
}

// All handles GET /busca?q=&tipo=&tema=&publico_alvo=&macrorregiao_id=&municipio_id=&gat_id=
func (h *BuscaHandler) All(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	term := strings.ToLower(strings.TrimSpace(q.Get("q")))
	macroID := int64(atoiDefault(q.Get("macrorregiao_id"), 0))
	munID := int64(atoiDefault(q.Get("municipio_id"), 0))
	gatID := int64(atoiDefault(q.Get("gat_id"), 0))

	materiais, err := h.materialRepo.All(repository.MaterialListParams{
		Q:              strings.TrimSpace(q.Get("q")),
		Tipos:          splitCSV(q.Get("tipo")),
		Tema:           q.Get("tema"),
		PublicoAlvo:    q.Get("publico_alvo"),
		MacrorregiaoID: macroID,
		MunicipioID:    munID,
		GatID:          gatID,
		Limit:          atoiDefault(q.Get("limit"), 30),
	})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}

	atividades, err := h.atividadeRepo.All(repository.AtividadeListParams{
		Q:              strings.TrimSpace(q.Get("q")),
		MunicipioID:    munID,
		MacrorregiaoID: macroID,
		GatID:          gatID,
		Limit:          atoiDefault(q.Get("limit"), 30),
	})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}

	cursosAll, err := h.cursoRepo.All(repository.CursoListParams{MacrorregiaoID: macroID, Limit: 100})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	cursos := make([]*models.Curso, 0, len(cursosAll))
	for _, c := range cursosAll {
		if term != "" && !strings.Contains(strings.ToLower(c.Titulo), term) && !containsFold(c.Descricao, term) && !containsFold(c.Categoria, term) {
			continue
		}
		if pa := strings.ToLower(q.Get("publico_alvo")); pa != "" && !containsFold(c.PublicoAlvo, pa) {
			continue
		}
		cursos = append(cursos, c)
	}

	// Filtros de "tipo" e "tema" só fazem sentido para materiais — quando
	// usados, escondemos cursos/atividades para não devolver ruído.
	if q.Get("tipo") != "" || q.Get("tema") != "" {
		cursos = []*models.Curso{}
		atividades = []*models.Atividade{}
	}

	httpx.JSON(w, 200, buscaResultado{Materiais: materiais, Cursos: cursos, Atividades: atividades})
}
