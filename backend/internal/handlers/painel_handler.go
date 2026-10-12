package handlers

import (
	"net/http"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// PainelHandler alimenta o módulo "📊 Painel": indicadores de acesso,
// downloads, municípios alcançados, capacitações etc., calculados a partir
// dos dados reais da plataforma.
type PainelHandler struct {
	repo *repository.PainelRepository
}

func NewPainelHandler(repo *repository.PainelRepository) *PainelHandler {
	return &PainelHandler{repo: repo}
}

// Resumo handles GET /painel/resumo (qualquer usuário autenticado): KPIs
// agregados do projeto, sem as séries de acesso — que ficam restritas à
// coordenação (ver Indicadores).
func (h *PainelHandler) Resumo(w http.ResponseWriter, r *http.Request) {
	ind, err := h.repo.Calcular(6)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	ind.EngajamentoMensal = []repository.EngajamentoMes{}
	ind.AlcancePorMacrorregiao = []repository.AlcanceMacro{}
	ind.TotalUsuarios = 0
	httpx.JSON(w, 200, ind)
}

// Indicadores handles GET /painel/indicadores (coordenador ou admin):
// tudo do resumo + engajamento mensal + acessos por macrorregião.
// ?meses=N controla o tamanho da série mensal (padrão 6, máximo 24).
func (h *PainelHandler) Indicadores(w http.ResponseWriter, r *http.Request) {
	meses := atoiDefault(r.URL.Query().Get("meses"), 6)
	if meses < 1 {
		meses = 6
	}
	if meses > 24 {
		meses = 24
	}
	ind, err := h.repo.Calcular(meses)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.JSON(w, 200, ind)
}
