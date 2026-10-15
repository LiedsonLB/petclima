package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/repository"
)

type ctxKey string

const userCtxKey ctxKey = "usuario_session"

// SessionCookieName é o cookie HttpOnly que guarda a sessão no navegador.
// O frontend não enxerga o valor (só o envia com credentials: 'include').
const SessionCookieName = "pet_session"

// SessionParts extrai (usuarioID, token, appKey) da sessão da requisição.
// Prioridade: cookie HttpOnly `pet_session` (site) e, se não houver, o header
// `TokenUser` (apps móveis/integrações). Ambos usam o formato
// "userId:token:appKey".
func SessionParts(r *http.Request) (int64, string, string, bool) {
	raw := r.Header.Get("TokenUser")
	if c, err := r.Cookie(SessionCookieName); err == nil && c.Value != "" {
		raw = c.Value
	}
	parts := strings.Split(raw, ":")
	if len(parts) < 3 {
		return 0, "", "", false
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	token, appKey := parts[1], parts[2]
	if err != nil || userID <= 0 || token == "" || appKey == "" {
		return 0, "", "", false
	}
	return userID, token, appKey, true
}

// authenticate valida a sessão (cookie ou header TokenUser) contra
// aplicacao_acesso_token e devolve o usuário logado. Retorna (nil, erro)
// quando a sessão está ausente/inválida/expirada. Falhas de credencial viram
// sempre 401 — é esse status que o frontend usa para derrubar a sessão local.
func authenticate(r *http.Request, usuarioRepo *repository.UsuarioRepository, tokenRepo *repository.TokenRepository) (*models.Usuario, *apperr.AppError) {
	userID, token, appKey, ok := SessionParts(r)
	if !ok {
		return nil, apperr.New("Token inválido!", 401)
	}

	if _, err := tokenRepo.Check(userID, token, appKey); err != nil {
		if ae, ok := apperr.As(err); ok {
			if ae.Status == 401 || ae.Status == 403 {
				return nil, apperr.New(ae.Message, 401)
			}
			return nil, ae
		}
		return nil, apperr.New("Erro ao validar token.", 500)
	}

	usuario, err := usuarioRepo.FindByID(userID, false)
	if err != nil {
		return nil, apperr.New("Usuário não encontrado(a)!", 401)
	}
	return usuario, nil
}

// RequireAuth exige um usuário autenticado (área logada da plataforma).
func RequireAuth(usuarioRepo *repository.UsuarioRepository, tokenRepo *repository.TokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			usuario, aerr := authenticate(r, usuarioRepo, tokenRepo)
			if aerr != nil {
				httpx.Error(w, aerr.Message, aerr.Status)
				return
			}
			ctx := context.WithValue(r.Context(), userCtxKey, usuario)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth serve endpoints públicos (portal, biblioteca aberta) que
// ganham informação extra quando há alguém logado — por exemplo, marcar em
// quais atividades o usuário já está inscrito. Token ausente ou inválido
// simplesmente segue como visitante (ROLE_PUBLIC do planejamento).
func OptionalAuth(usuarioRepo *repository.UsuarioRepository, tokenRepo *repository.TokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, _, _, ok := SessionParts(r); ok {
				if usuario, aerr := authenticate(r, usuarioRepo, tokenRepo); aerr == nil {
					ctx := context.WithValue(r.Context(), userCtxKey, usuario)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func forbid(w http.ResponseWriter) {
	httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
}

// RequireAdmin: somente papel "admin". Deve vir depois de RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r)
		if u == nil || !u.IsAdmin() {
			forbid(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireCoordenador: coordenador ou admin (indicadores, relatórios,
// gestão de territórios/GATs). Deve vir depois de RequireAuth.
func RequireCoordenador(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r)
		if u == nil || !u.IsCoordenadorOuMais() {
			forbid(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireStaff: preceptor, coordenador ou admin — quem pode publicar
// materiais/notícias, criar cursos e registrar atividades. Deve vir depois
// de RequireAuth.
func RequireStaff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r)
		if u == nil || !u.IsStaff() {
			forbid(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserFromContext devolve o usuário autenticado (ou nil se anônimo).
func UserFromContext(r *http.Request) *models.Usuario {
	u, _ := r.Context().Value(userCtxKey).(*models.Usuario)
	return u
}
