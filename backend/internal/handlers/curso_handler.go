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

// CursoHandler cobre o módulo "🎓 Cursos / AVA": catálogo (público),
// matrícula, progresso e certificado (autenticados).
type CursoHandler struct {
	repo     *repository.CursoRepository
	notifier *Notifier
}

func NewCursoHandler(repo *repository.CursoRepository, notifier *Notifier) *CursoHandler {
	return &CursoHandler{repo: repo, notifier: notifier}
}

// All handles GET /cursos (?status=&categoria=&macrorregiao_id=). Se houver
// usuário logado, cada curso vem com `minha_matricula`.
func (h *CursoHandler) All(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.repo.All(repository.CursoListParams{
		Status:         q.Get("status"),
		Categoria:      q.Get("categoria"),
		MacrorregiaoID: int64(atoiDefault(q.Get("macrorregiao_id"), 0)),
		Limit:          atoiDefault(q.Get("limit"), 50),
		Page:           atoiDefault(q.Get("page"), 1),
	})
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if u := middleware.UserFromContext(r); u != nil {
		minhas, err := h.repo.ListMinhas(u.ID)
		if err == nil {
			byCurso := map[int64]*models.Matricula{}
			for _, m := range minhas {
				byCurso[m.CursoID] = m
			}
			for _, c := range list {
				c.MinhaMatricula = byCurso[c.ID]
			}
		}
	}
	if list == nil {
		list = []*models.Curso{}
	}
	httpx.JSON(w, 200, list)
}

// Find handles GET /cursos/{id}
func (h *CursoHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	c, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if u := middleware.UserFromContext(r); u != nil {
		if m, err := h.repo.FindMatricula(u.ID, id); err == nil {
			c.MinhaMatricula = m
		}
	}
	httpx.JSON(w, 200, c)
}

type cursoPayload struct {
	Titulo         string  `json:"titulo"`
	Descricao      *string `json:"descricao"`
	Categoria      *string `json:"categoria"`
	PublicoAlvo    *string `json:"publico_alvo"`
	CargaHoraria   int     `json:"carga_horaria"`
	ModulosTotal   int     `json:"modulos_total"`
	CapaURL        *string `json:"capa_url"`
	Status         string  `json:"status"`
	MacrorregiaoID *int64  `json:"macrorregiao_id"`
}

// Save handles POST /cursos (staff).
func (h *CursoHandler) Save(w http.ResponseWriter, r *http.Request) {
	var p cursoPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	status := strings.ToLower(strings.TrimSpace(p.Status))
	switch status {
	case models.CursoStatusPublicado, models.CursoStatusEmBreve, models.CursoStatusPlanejado:
	case "":
		status = models.CursoStatusPlanejado
	default:
		httpx.Error(w, "Status de curso inválido.", 422)
		return
	}
	c := &models.Curso{
		Titulo:         strings.TrimSpace(p.Titulo),
		Descricao:      nilIfEmpty(p.Descricao),
		Categoria:      nilIfEmpty(p.Categoria),
		PublicoAlvo:    nilIfEmpty(p.PublicoAlvo),
		CargaHoraria:   p.CargaHoraria,
		ModulosTotal:   p.ModulosTotal,
		CapaURL:        nilIfEmpty(p.CapaURL),
		Status:         status,
		MacrorregiaoID: p.MacrorregiaoID,
	}
	if u := middleware.UserFromContext(r); u != nil {
		c.CriadoPor = &u.ID
	}
	created, err := h.repo.Create(c)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, created)
}

// Matricular handles POST /cursos/{id}/matricula (autenticado).
func (h *CursoHandler) Matricular(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	id, err := idFromPath(r)
	if err != nil || u == nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	curso, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if curso.Status != models.CursoStatusPublicado {
		httpx.Error(w, "Este curso ainda não está aberto para matrículas.", 422)
		return
	}
	m, err := h.repo.Matricular(u.ID, id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, m)
}

type progressoPayload struct {
	Progresso int `json:"progresso"`
}

// AtualizarProgresso handles PUT /cursos/{id}/progresso (autenticado).
// Ao chegar a 100% a matrícula é concluída, o certificado é gerado e a
// pessoa recebe um e-mail de parabéns com o código de validação.
func (h *CursoHandler) AtualizarProgresso(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	id, err := idFromPath(r)
	if err != nil || u == nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	var p progressoPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}

	antes, _ := h.repo.FindMatricula(u.ID, id)
	jaConcluida := antes != nil && antes.Status == models.MatriculaConcluido

	m, err := h.repo.AtualizarProgresso(u.ID, id, p.Progresso)
	if err != nil {
		writeAppErr(w, err)
		return
	}

	if !jaConcluida && m.Status == models.MatriculaConcluido && m.CertificadoCodigo != nil {
		titulo := "o curso"
		if m.CursoTitulo != nil {
			titulo = "“" + *m.CursoTitulo + "”"
		}
		h.notifier.Generico(u,
			"Certificado de conclusão",
			"Parabéns pela conclusão!",
			"Você concluiu "+titulo+". Seu código de validação do certificado é "+*m.CertificadoCodigo+".",
			"Ver meus cursos",
			h.notifier.FrontendURL()+"/ava",
		)
	}
	httpx.JSON(w, 200, m)
}

// Minhas handles GET /cursos/minhas-matriculas (autenticado).
func (h *CursoHandler) Minhas(w http.ResponseWriter, r *http.Request) {
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
		list = []*models.Matricula{}
	}
	httpx.JSON(w, 200, list)
}

// ValidarCertificado handles GET /certificados/{codigo} (público): confirma
// que um certificado existe, sem expor dados sensíveis da pessoa.
func (h *CursoHandler) ValidarCertificado(w http.ResponseWriter, r *http.Request) {
	codigo := strings.TrimSpace(r.PathValue("codigo"))
	if codigo == "" {
		httpx.Error(w, "Código inválido.", 422)
		return
	}
	m, err := h.repo.FindMatriculaPorCertificado(codigo)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 200, map[string]any{
		"valido":       true,
		"codigo":       codigo,
		"curso_titulo": m.CursoTitulo,
		"concluido_em": m.ConcluidoEm,
	})
}
