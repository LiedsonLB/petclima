package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// AtividadeHandler cobre os módulos "📝 Atividades/ações" e "🧑‍⚕️ Educação
// Permanente": cadastro de oficinas, capacitações, intervenções, visitas e
// reuniões, inscrição de participantes e registro de presença. Atividades
// remotas/híbridas ganham automaticamente uma sala de videoconferência.
type AtividadeHandler struct {
	repo     *repository.AtividadeRepository
	salaRepo *repository.SalaRepository
	notifier *Notifier
}

func NewAtividadeHandler(repo *repository.AtividadeRepository, salaRepo *repository.SalaRepository, notifier *Notifier) *AtividadeHandler {
	return &AtividadeHandler{repo: repo, salaRepo: salaRepo, notifier: notifier}
}

var tiposAtividadeValidos = map[string]bool{
	models.AtividadeOficina:     true,
	models.AtividadeCapacitacao: true,
	models.AtividadeIntervencao: true,
	models.AtividadeVisita:      true,
	models.AtividadeReuniao:     true,
}

var modalidadesValidas = map[string]bool{
	models.ModalidadePresencial: true,
	models.ModalidadeRemota:     true,
	models.ModalidadeHibrida:    true,
}

// All handles GET /atividades
// Filtros: ?q=&tipo=oficina,capacitacao&municipio_id=&macrorregiao_id=&gat_id=&futuras=1
func (h *AtividadeHandler) All(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.repo.All(repository.AtividadeListParams{
		Q:              strings.TrimSpace(q.Get("q")),
		Tipos:          splitCSV(q.Get("tipo")),
		MunicipioID:    int64(atoiDefault(q.Get("municipio_id"), 0)),
		MacrorregiaoID: int64(atoiDefault(q.Get("macrorregiao_id"), 0)),
		GatID:          int64(atoiDefault(q.Get("gat_id"), 0)),
		ApenasFuturas:  q.Get("futuras") == "1",
		Limit:          atoiDefault(q.Get("limit"), 50),
		Page:           atoiDefault(q.Get("page"), 1),
	})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if list == nil {
		list = []*models.Atividade{}
	}
	if u := middleware.UserFromContext(r); u != nil {
		if minhas, err := h.repo.InscricoesDoUsuario(u.ID); err == nil {
			for _, a := range list {
				a.MinhaInscricao = minhas[a.ID]
			}
		}
	}
	httpx.JSON(w, 200, list)
}

// Find handles GET /atividades/{id}
func (h *AtividadeHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	a, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if u := middleware.UserFromContext(r); u != nil {
		if insc, err := h.repo.FindInscricao(id, u.ID); err == nil {
			a.MinhaInscricao = insc
		}
	}
	httpx.JSON(w, 200, a)
}

// Minhas handles GET /atividades/minhas (autenticado).
func (h *AtividadeHandler) Minhas(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	if u == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}
	list, err := h.repo.ListMinhas(u.ID)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if list == nil {
		list = []*models.Atividade{}
	}
	httpx.JSON(w, 200, list)
}

type atividadePayload struct {
	Titulo         string  `json:"titulo"`
	Tipo           string  `json:"tipo"`
	Descricao      *string `json:"descricao"`
	Modalidade     string  `json:"modalidade"`
	MunicipioID    *int64  `json:"municipio_id"`
	MacrorregiaoID *int64  `json:"macrorregiao_id"`
	GatID          *int64  `json:"gat_id"`
	DataRealizacao string  `json:"data_realizacao"` // RFC3339 ou "2006-01-02" / "2006-01-02T15:04"
	CargaHoraria   int     `json:"carga_horaria"`
	Vagas          *int    `json:"vagas"`
}

func parseData(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02 15:04", "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("data inválida")
}

// Save handles POST /atividades (staff). Se a modalidade for remota ou
// híbrida, cria e vincula uma sala de videoconferência automaticamente.
func (h *AtividadeHandler) Save(w http.ResponseWriter, r *http.Request) {
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}
	var p atividadePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}

	tipo := strings.ToLower(strings.TrimSpace(p.Tipo))
	if tipo == "" {
		tipo = models.AtividadeOficina
	}
	if !tiposAtividadeValidos[tipo] {
		httpx.Error(w, "Tipo de atividade inválido.", 422)
		return
	}
	modalidade := strings.ToLower(strings.TrimSpace(p.Modalidade))
	if modalidade == "" {
		modalidade = models.ModalidadePresencial
	}
	if !modalidadesValidas[modalidade] {
		httpx.Error(w, "Modalidade inválida.", 422)
		return
	}
	data, err := parseData(p.DataRealizacao)
	if err != nil {
		httpx.ErrorWithFields(w, "Dados inválidos.", map[string]string{"data_realizacao": "formato inválido"}, 422)
		return
	}

	a := &models.Atividade{
		Titulo:         strings.TrimSpace(p.Titulo),
		Tipo:           tipo,
		Descricao:      nilIfEmpty(p.Descricao),
		Modalidade:     modalidade,
		MunicipioID:    p.MunicipioID,
		MacrorregiaoID: p.MacrorregiaoID,
		GatID:          p.GatID,
		ResponsavelID:  session.ID,
		DataRealizacao: data,
		CargaHoraria:   p.CargaHoraria,
		Vagas:          p.Vagas,
	}
	created, err := h.repo.Create(a)
	if err != nil {
		writeAppErr(w, err)
		return
	}

	if modalidade != models.ModalidadePresencial {
		tipoSala := models.SalaTipoOficina
		if tipo == models.AtividadeReuniao {
			tipoSala = models.SalaTipoReuniao
		}
		sala, err := h.salaRepo.Create(created.Titulo, tipoSala, created.Descricao, &created.ID, session.ID)
		if err == nil {
			if err := h.repo.VincularSala(created.ID, sala.ID); err == nil {
				created, _ = h.repo.FindByID(created.ID)
			}
		}
	}
	httpx.JSON(w, 201, created)
}

// podeGerenciar: o responsável pela atividade ou a coordenação.
func podeGerenciar(u *models.Usuario, a *models.Atividade) bool {
	return u != nil && (u.ID == a.ResponsavelID || u.IsCoordenadorOuMais())
}

// Delete handles DELETE /atividades/{id}
func (h *AtividadeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	a, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if !podeGerenciar(middleware.UserFromContext(r), a) {
		httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
		return
	}
	if a.SalaID != nil {
		_ = h.salaRepo.Encerrar(*a.SalaID)
	}
	if err := h.repo.SoftDelete(id); err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.Success(w, "Sucesso!", 200)
}

// Inscrever handles POST /atividades/{id}/inscricao (autenticado).
func (h *AtividadeHandler) Inscrever(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	id, err := idFromPath(r)
	if err != nil || u == nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	a, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	jaInscrito, _ := h.repo.FindInscricao(id, u.ID)
	insc, err := h.repo.Inscrever(id, u.ID)
	if err != nil {
		writeAppErr(w, err)
		return
	}

	if jaInscrito == nil {
		quando := "em data a definir"
		if a.DataRealizacao != nil {
			quando = "em " + a.DataRealizacao.Format("02/01/2006 às 15:04")
		}
		msg := fmt.Sprintf("Sua inscrição em “%s” foi confirmada. A atividade acontece %s.", a.Titulo, quando)
		acao, link := "Ver atividades", h.notifier.FrontendURL()+"/formacao"
		if a.SalaID != nil {
			msg += " Como ela é remota/híbrida, você entrará pela sala virtual, direto na plataforma."
			acao = "Abrir a atividade"
		}
		h.notifier.Generico(u, "Inscrição confirmada", "Inscrição confirmada!", msg, acao, link)
	}
	httpx.JSON(w, 201, insc)
}

// CancelarInscricao handles DELETE /atividades/{id}/inscricao (autenticado).
func (h *AtividadeHandler) CancelarInscricao(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	id, err := idFromPath(r)
	if err != nil || u == nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	if err := h.repo.CancelarInscricao(id, u.ID); err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.Success(w, "Inscrição cancelada.", 200)
}

// Inscritos handles GET /atividades/{id}/inscritos (responsável/coordenação).
func (h *AtividadeHandler) Inscritos(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	a, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if !podeGerenciar(middleware.UserFromContext(r), a) {
		httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
		return
	}
	list, err := h.repo.ListInscritos(id)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if list == nil {
		list = []*models.AtividadeInscricao{}
	}
	httpx.JSON(w, 200, list)
}

type presencaPayload struct {
	Presente bool `json:"presente"`
}

// MarcarPresenca handles PUT /atividades/{id}/presenca/{usuarioId}
// (responsável/coordenação) — o "registrar atividades e resultados".
func (h *AtividadeHandler) MarcarPresenca(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	usuarioID, err := parseInt64(r.PathValue("usuarioId"))
	if err != nil {
		httpx.Error(w, "Id de usuário inválido.", 422)
		return
	}
	a, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if !podeGerenciar(middleware.UserFromContext(r), a) {
		httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
		return
	}
	var p presencaPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	if err := h.repo.MarcarPresenca(id, usuarioID, p.Presente); err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.Success(w, "Presença atualizada.", 200)
}
