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

// authenticate lê o header `TokenUser` (formato "userId:token:appKey"), valida
// contra aplicacao_acesso_token e devolve o usuário logado. Retorna
// (nil, erro) quando o token está ausente/inválido.
func authenticate(r *http.Request, usuarioRepo *repository.UsuarioRepository, tokenRepo *repository.TokenRepository) (*models.Usuario, *apperr.AppError) {
	parts := strings.Split(r.Header.Get("TokenUser"), ":")
	if len(parts) < 3 {
		return nil, apperr.New("Token inválido!", 401)
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	token := parts[1]
	appKey := parts[2]
	if err != nil || userID <= 0 || token == "" || appKey == "" {
		return nil, apperr.New("Token inválido!", 401)
	}

	if _, err := tokenRepo.Check(userID, token, appKey); err != nil {
		if ae, ok := apperr.As(err); ok {
			return nil, ae
		}
		return nil, apperr.New("Erro ao validar token.", 500)
	}

	usuario, err := usuarioRepo.FindByID(userID, false)
	if err != nil {
		return nil, apperr.New("Usuário não encontrado(a)!", 404)
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
			if r.Header.Get("TokenUser") != "" {
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
