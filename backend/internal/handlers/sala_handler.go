package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/livekit"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

// SalaHandler gerencia as salas de videoconferência (LiveKit) do PET-Saúde
// Clima: reuniões de equipe, oficinas remotas de Educação Permanente e
// teleconsultoria. Quando a sala está vinculada a uma atividade, só entram
// os inscritos, o responsável e a coordenação — e quem é inscrito tem a
// presença registrada automaticamente ao entrar.
type SalaHandler struct {
	repo             *repository.SalaRepository
	atividadeRepo    *repository.AtividadeRepository
	liveKitAPIKey    string
	liveKitAPISecret string
	liveKitURL       string
	liveKitPublicURL string
}

func NewSalaHandler(
	repo *repository.SalaRepository,
	atividadeRepo *repository.AtividadeRepository,
	apiKey, apiSecret, url, publicURL string,
) *SalaHandler {
	if publicURL == "" {
		publicURL = url
	}
	return &SalaHandler{
		repo: repo, atividadeRepo: atividadeRepo,
		liveKitAPIKey: apiKey, liveKitAPISecret: apiSecret,
		liveKitURL: url, liveKitPublicURL: publicURL,
	}
}

func (h *SalaHandler) liveKitConfigurado() bool {
	return h.liveKitAPIKey != "" && h.liveKitAPISecret != ""
}

// contarOnline consulta o LiveKit para saber quantas pessoas estão na sala
// agora. Falhas viram 0 (o painel não deve quebrar porque o LiveKit está
// fora do ar).
func (h *SalaHandler) contarOnline(s *models.Sala) int {
	if !h.liveKitConfigurado() {
		return 0
	}
	ps, err := livekit.ListParticipants(h.liveKitURL, h.liveKitAPIKey, h.liveKitAPISecret, s.Codigo)
	if err != nil {
		log.Printf("livekit: erro ao listar participantes da sala %d: %v", s.ID, err)
		return 0
	}
	return len(ps)
}

// All handles GET /salas — lista as salas ativas (?tipo=reuniao|oficina|telessaude).
func (h *SalaHandler) All(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.All(r.URL.Query().Get("tipo"), true)
	if err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if list == nil {
		list = []*models.Sala{}
	}
	var wg sync.WaitGroup
	for _, s := range list {
		wg.Add(1)
		go func(s *models.Sala) {
			defer wg.Done()
			s.ParticipantesOnline = h.contarOnline(s)
		}(s)
	}
	wg.Wait()
	httpx.JSON(w, 200, list)
}

// Find handles GET /salas/{id}
func (h *SalaHandler) Find(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	s, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	s.ParticipantesOnline = h.contarOnline(s)
	httpx.JSON(w, 200, s)
}

type salaPayload struct {
	Nome      string  `json:"nome"`
	Tipo      string  `json:"tipo"` // reuniao | oficina | telessaude
	Descricao *string `json:"descricao"`
}

// Save handles POST /salas (staff) — cria uma sala avulsa (reunião de
// equipe, teleconsultoria). Salas de oficina são criadas junto com a
// atividade remota/híbrida (ver AtividadeHandler.Save).
func (h *SalaHandler) Save(w http.ResponseWriter, r *http.Request) {
	var payload salaPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}
	if payload.Nome == "" {
		payload.Nome = fmt.Sprintf("Sala de %s - %s", session.Nome, time.Now().Format("02/01 15:04"))
	}
	if payload.Tipo == "" {
		payload.Tipo = models.SalaTipoReuniao
	}
	created, err := h.repo.Create(payload.Nome, payload.Tipo, nilIfEmpty(payload.Descricao), nil, session.ID)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	httpx.JSON(w, 201, created)
}

// Encerrar handles POST /salas/{id}/encerrar — apenas quem criou a sala ou
// um administrador pode encerrá-la.
func (h *SalaHandler) Encerrar(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	sala, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	session := middleware.UserFromContext(r)
	if session == nil || (!session.IsAdmin() && session.ID != sala.CriadoPor) {
		httpx.Error(w, "Você não tem permissão para encerrar essa sala.", 403)
		return
	}
	if err := h.repo.Encerrar(id); err != nil {
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	httpx.Success(w, "Sala encerrada.", 200)
}

// podeEntrar aplica a regra de acesso: salas avulsas são abertas a qualquer
// usuário logado; salas de atividade exigem inscrição, ser o responsável ou
// fazer parte da coordenação. Devolve também a atividade (quando houver) e
// se a pessoa é inscrita, para o registro automático de presença.
func (h *SalaHandler) podeEntrar(sala *models.Sala, u *models.Usuario) (ok bool, atividadeID int64, inscrito bool, err error) {
	if sala.AtividadeID == nil {
		return true, 0, false, nil
	}
	atividadeID = *sala.AtividadeID
	insc, err := h.atividadeRepo.FindInscricao(atividadeID, u.ID)
	if err != nil {
		return false, atividadeID, false, err
	}
	inscrito = insc != nil
	if inscrito || u.IsCoordenadorOuMais() {
		return true, atividadeID, inscrito, nil
	}
	atividade, err := h.atividadeRepo.FindByID(atividadeID)
	if err != nil {
		return false, atividadeID, false, err
	}
	return atividade.ResponsavelID == u.ID, atividadeID, false, nil
}

// Entrar handles POST /salas/{id}/entrar — gera o token de acesso LiveKit
// (JWT) para o usuário autenticado entrar na room correspondente (ver
// internal/livekit/token.go, assinatura HS256 sem depender do SDK oficial).
func (h *SalaHandler) Entrar(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	sala, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if !sala.Ativa {
		httpx.Error(w, "Essa sala já foi encerrada.", 410)
		return
	}
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}
	if !h.liveKitConfigurado() {
		httpx.Error(w, "A videoconferência ainda não foi configurada neste ambiente.", 503)
		return
	}

	ok, atividadeID, inscrito, err := h.podeEntrar(sala, session)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	if !ok {
		httpx.Error(w, "Inscreva-se na atividade para entrar na sala virtual.", 403)
		return
	}

	meta := map[string]string{"nome": session.Nome}
	if session.Foto != nil && *session.Foto != "" {
		meta["foto"] = *session.Foto
	}
	metadataJSON := ""
	if b, err := json.Marshal(meta); err == nil {
		metadataJSON = string(b)
	}

	token, err := livekit.NewTokenBuilder(h.liveKitAPIKey, h.liveKitAPISecret).
		Identity(fmt.Sprintf("%d", session.ID)).
		Name(session.Nome).
		Metadata(metadataJSON).
		TTL(4 * time.Hour).
		Grant(livekit.VideoGrant{
			RoomJoin:             true,
			Room:                 sala.Codigo,
			CanPublish:           true,
			CanSubscribe:         true,
			CanPublishData:       true,
			CanUpdateOwnMetadata: true,
		}).
		Sign()
	if err != nil {
		httpx.Error(w, err.Error(), 500)
		return
	}

	// Presença automática: quem está inscrito e entra na sala da própria
	// atividade é marcado como presente (alimenta "profissionais
	// capacitados" no Painel).
	if inscrito && atividadeID > 0 {
		if err := h.atividadeRepo.MarcarPresenca(atividadeID, session.ID, true); err != nil {
			log.Printf("erro ao registrar presença automática: %v", err)
		}
	}

	httpx.JSON(w, 200, map[string]any{
		"token": token,
		"url":   h.liveKitPublicURL,
		"room":  sala.Codigo,
		"sala":  sala,
	})
}

// Participantes handles GET /salas/{id}/participantes — quem está na sala
// agora, direto do LiveKit.
func (h *SalaHandler) Participantes(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		httpx.Error(w, "Id inválido.", 422)
		return
	}
	sala, err := h.repo.FindByID(id)
	if err != nil {
		writeAppErr(w, err)
		return
	}
	lista := []models.ParticipanteInfo{}
	if h.liveKitConfigurado() {
		if ps, err := livekit.ListParticipants(h.liveKitURL, h.liveKitAPIKey, h.liveKitAPISecret, sala.Codigo); err != nil {
			log.Printf("livekit: erro ao listar participantes da sala %d: %v", sala.ID, err)
		} else if ps != nil {
			lista = ps
		}
	}
	httpx.JSON(w, 200, map[string]any{"sala_id": sala.ID, "total": len(lista), "participantes": lista})
}
