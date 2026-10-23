package handlers

import (
	"net/http"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// TerritorioHandler expõe os dados do módulo "🗺️ Territórios". Todos os
// endpoints são públicos: são usados pelo mapa, pelos filtros de busca e
// pelo formulário de cadastro (que precisa listar municípios antes de a
// pessoa ter uma conta).
type TerritorioHandler struct {
	territorioRepo *repository.TerritorioRepository
	cidadeRepo     *repository.CidadeRepository
}

func NewTerritorioHandler(t *repository.TerritorioRepository, c *repository.CidadeRepository) *TerritorioHandler {
	return &TerritorioHandler{territorioRepo: t, cidadeRepo: c}
}

// Macrorregioes handles GET /territorios/macrorregioes
func (h *TerritorioHandler) Macrorregioes(w http.ResponseWriter, r *http.Request) {
	list, err := h.territorioRepo.ListMacrorregioes()
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}

// Municipios handles GET /territorios/municipios
// Filtros: ?macrorregiao_id=&status=ativo|planejado&nome=&uf=
func (h *TerritorioHandler) Municipios(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, _, err := h.cidadeRepo.All(repository.CidadeListParams{
		Nome:           q.Get("nome"),
		UF:             q.Get("uf"),
		Status:         q.Get("status"),
		MacrorregiaoID: int64(atoiDefault(q.Get("macrorregiao_id"), 0)),
		Limit:          atoiDefault(q.Get("limit"), 500),
		Page:           atoiDefault(q.Get("page"), 1),
	})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}

// Gats handles GET /territorios/gats  (?macrorregiao_id=)
func (h *TerritorioHandler) Gats(w http.ResponseWriter, r *http.Request) {
	list, err := h.territorioRepo.ListGats(int64(atoiDefault(r.URL.Query().Get("macrorregiao_id"), 0)))
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, list)
}
