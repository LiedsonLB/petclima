package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/liedsonlb/petsaude-clima/internal/apperr"
	"github.com/liedsonlb/petsaude-clima/internal/httpx"
	"github.com/liedsonlb/petsaude-clima/internal/models"
	"github.com/liedsonlb/petsaude-clima/internal/repository"
)

type ctxKey string

const userCtxKey ctxKey = "usuario_session"

// SessionCookieName é o cookie HttpOnly que carrega a sessão do portal.
// Valor: "{usuario_id}:{token}:{codigo_da_aplicacao}" — o mesmo formato do
// header legado `TokenUser`, só que guardado pelo navegador e invisível
// para JavaScript (protege o token contra roubo via XSS).
const SessionCookieName = "pet_session"

// Credentials são as credenciais extraídas de uma requisição.
type Credentials struct {
	UserID     int64
	Token      string
	AppKey     string
	FromCookie bool
}

// ExtractCredentials lê a sessão do cookie HttpOnly (preferencial) ou, na
// falta dele, do header legado `TokenUser` (clientes que não são o site).
func ExtractCredentials(r *http.Request) (*Credentials, bool) {
	raw := ""
	fromCookie := false
	if c, err := r.Cookie(SessionCookieName); err == nil && c.Value != "" {
		raw, fromCookie = c.Value, true
	} else {
		raw = r.Header.Get("TokenUser")
	}

	parts := strings.Split(raw, ":")
	if len(parts) < 3 {
		return nil, false
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || userID <= 0 || parts[1] == "" || parts[2] == "" {
		return nil, false
	}
	return &Credentials{UserID: userID, Token: parts[1], AppKey: parts[2], FromCookie: fromCookie}, true
}

// IsSecureRequest diz se o cookie deve ir com o atributo Secure: sim quando
// o site é servido por HTTPS (FRONTEND_URL https:// ou proxy repassando
// X-Forwarded-Proto=https). Em desenvolvimento (http://localhost) fica falso.
func IsSecureRequest(r *http.Request, frontendURL string) bool {
	return strings.HasPrefix(strings.ToLower(frontendURL), "https://") ||
		r.TLS != nil ||
		strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

// SetSessionCookie grava a sessão em um cookie HttpOnly. Tokens "longos"
// (opção "manter conectado") duram 30 dias; os curtos viram cookie de
// sessão e somem ao fechar o navegador.
func SetSessionCookie(w http.ResponseWriter, userID int64, token, appKey string, long, secure bool) {
	c := &http.Cookie{
		Name:     SessionCookieName,
		Value:    fmt.Sprintf("%d:%s:%s", userID, token, appKey),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	if long {
		c.MaxAge = models.TempoLimiteLongoAPIMinutes * 60
	}
	http.SetCookie(w, c)
}

// ClearSessionCookie apaga o cookie de sessão no navegador.
func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func hostnameOf(hostport string) string {
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		return strings.ToLower(h)
	}
	return strings.ToLower(hostport)
}

// OriginAllowed é a defesa contra CSRF para requisições autenticadas por
// cookie e que alteram dados (POST/PUT/PATCH/DELETE): o header Origin, que
// o navegador sempre envia nesses métodos, precisa ser do próprio site (ou
// do FRONTEND_URL). Requisições sem Origin não vêm de um navegador, que é
// o único cliente capaz de anexar o cookie automaticamente.
func OriginAllowed(r *http.Request, frontendURL string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	host := hostnameOf(u.Host)
	if host == hostnameOf(r.Host) {
		return true
	}
	if f, err := url.Parse(frontendURL); err == nil && f.Host != "" && host == hostnameOf(f.Host) {
		return true
	}
	return false
}

func isUnsafeMethod(m string) bool {
	switch m {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	return true
}

// RequireAuth valida a sessão (cookie HttpOnly ou header legado TokenUser)
// contra aplicacao_acesso_token. Qualquer falha de autenticação responde
// 401 — o frontend usa isso para mandar o usuário de volta ao login — e,
// se a sessão veio de cookie, o cookie inválido é apagado.
func RequireAuth(usuarioRepo *repository.UsuarioRepository, tokenRepo *repository.TokenRepository, frontendURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cred, ok := ExtractCredentials(r)
			if !ok {
				httpx.Error(w, "Não autenticado!", 401)
				return
			}

			if cred.FromCookie && isUnsafeMethod(r.Method) && !OriginAllowed(r, frontendURL) {
				httpx.Error(w, "Origem não permitida.", 403)
				return
			}

			if _, err := tokenRepo.Check(cred.UserID, cred.Token, cred.AppKey); err != nil {
				ae, isApp := apperr.As(err)
				if !isApp {
					httpx.Error(w, "Erro ao validar token.", 500)
					return
				}
				if cred.FromCookie {
					ClearSessionCookie(w, IsSecureRequest(r, frontendURL))
				}
				httpx.Error(w, ae.Message, 401)
				return
			}

			usuario, err := usuarioRepo.FindByID(cred.UserID, false)
			if err != nil {
				if cred.FromCookie {
					ClearSessionCookie(w, IsSecureRequest(r, frontendURL))
				}
				httpx.Error(w, "Usuário não encontrado(a)!", 401)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, usuario)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin mirrors the `if(!$this->userSession->is_admin) throw ...`
// guards sprinkled through the Laravel controllers (GrauInstrucaoController,
// InstituicaoEnsinoController, GeneroTextualController's publish endpoints,
// etc). Must be chained after RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r)
		if u == nil || !u.IsAdmin() {
			httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserFromContext mirrors ApiController::getUserSession()
func UserFromContext(r *http.Request) *models.Usuario {
	u, _ := r.Context().Value(userCtxKey).(*models.Usuario)
	return u
}
