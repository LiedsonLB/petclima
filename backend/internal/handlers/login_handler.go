package handlers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/petsaude/petclima-backend/internal/apperr"
	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/mailer"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/queue"
	"github.com/petsaude/petclima-backend/internal/repository"
)

type LoginHandler struct {
	usuarioRepo   *repository.UsuarioRepository
	aplicacaoRepo *repository.AplicacaoRepository
	tokenRepo     *repository.TokenRepository
	resetRepo     *repository.PasswordResetRepository
	mailer        *mailer.Mailer
	rabbit        *queue.RabbitMQ
	frontendURL   string
	appName       string
	tokenKey      string
	appKey        string
	// Chave secreta para gerar tokens de verificação (sem banco)
	verificationSecret string
}

func NewLoginHandler(
	u *repository.UsuarioRepository,
	a *repository.AplicacaoRepository,
	t *repository.TokenRepository,
	resetRepo *repository.PasswordResetRepository,
	m *mailer.Mailer,
	rabbit *queue.RabbitMQ,
	frontendURL, appName, tokenKey, appKey, verificationSecret string,
) *LoginHandler {
	return &LoginHandler{
		usuarioRepo:        u,
		aplicacaoRepo:      a,
		tokenRepo:          t,
		resetRepo:          resetRepo,
		mailer:             m,
		rabbit:             rabbit,
		frontendURL:        frontendURL,
		appName:            appName,
		tokenKey:           tokenKey,
		appKey:             appKey,
		verificationSecret: verificationSecret,
	}
}

type loginRequest struct {
	Email     string `json:"email"`
	Senha     string `json:"senha"`
	LongToken string `json:"long_token"`
}

func md5hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// makeToken mirrors Usuario::getToken($data): md5(APP_NAME . usuario_id . data)
func (h *LoginHandler) makeToken(usuarioID int64, data string) string {
	return md5hex(h.appName + strconv.FormatInt(usuarioID, 10) + data)
}

// makeLongToken mirrors Usuario::getLongToken(): 64 hex chars (two md5s concatenated)
func (h *LoginHandler) makeLongToken(usuarioID int64) string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	tk := h.makeToken(usuarioID, data)
	tk += md5hex("long" + h.tokenKey + strconv.FormatInt(usuarioID, 10) + data)
	return tk
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i != -1 {
		host = host[:i]
	}
	return host
}

// generateVerificationToken cria um token baseado no email + secret + timestamp
// O token é válido por 48 horas e NÃO precisa de banco de dados
func (h *LoginHandler) generateVerificationToken(email string) (string, error) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d", email, timestamp)

	hsh := hmac.New(sha256.New, []byte(h.verificationSecret))
	hsh.Write([]byte(data))
	token := hsh.Sum(nil)

	// Usa base64 URL-safe (sem +, / ou =)
	tokenStr := base64.URLEncoding.EncodeToString(token)

	// Formato: timestamp:token
	return fmt.Sprintf("%d:%s", timestamp, tokenStr), nil
}

// verifyVerificationToken valida o token sem precisar de banco de dados
func (h *LoginHandler) verifyVerificationToken(email, token string) (bool, error) {

	// Normaliza o email (lowercase)
	email = strings.ToLower(strings.TrimSpace(email))

	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		log.Printf("❌ Token inválido: não tem 2 partes, tem %d", len(parts))
		return false, nil
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false, nil
	}

	// Verifica se o token expirou (48 horas)
	if time.Now().Unix()-timestamp > 48*3600 {
		return false, nil
	}

	// Decodifica o token base64
	tokenBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, nil
	}

	// Recalcula o HMAC para verificar
	data := fmt.Sprintf("%s:%d", email, timestamp)
	hsh := hmac.New(sha256.New, []byte(h.verificationSecret))
	hsh.Write([]byte(data))
	expectedToken := hsh.Sum(nil)

	return hmac.Equal(tokenBytes, expectedToken), nil
}

// logoURL aponta para o ícone do PET-Saúde Clima servido pelo frontend.
func (h *LoginHandler) logoURL() string {
	return h.frontendURL + "/petsaudeclima_icon.png"
}

// linkComQuery monta um link do frontend com e-mail/token devidamente
// escapados (e-mails com "+" e tokens com ":" quebrariam a URL sem isso).
func (h *LoginHandler) linkComQuery(path, email, token string) string {
	q := url.Values{}
	q.Set("email", email)
	q.Set("token", token)
	return fmt.Sprintf("%s%s?%s", h.frontendURL, path, q.Encode())
}

// enviarEmailJob publica um job de e-mail na fila; retorna false se o
// RabbitMQ não estiver disponível (o chamador decide o fallback).
func (h *LoginHandler) enviarEmailJob(job models.EmailJob) bool {
	if h.rabbit == nil {
		return false
	}
	job.CriadoEm = time.Now()
	if err := h.rabbit.Publish(queue.QueueEmail, job); err != nil {
		log.Println("erro ao publicar e-mail na fila:", err)
		return false
	}
	return true
}

// sendVerificationEmail envia e-mail de confirmação usando tokens auto-validáveis
// SEM tabela email_verification_tokens
func (h *LoginHandler) sendVerificationEmail(email, nome string) error {
	token, err := h.generateVerificationToken(email)
	if err != nil {
		return err
	}

	link := h.linkComQuery("/confirmar-email", email, token)

	if h.enviarEmailJob(models.EmailJob{
		Tipo:         models.EmailTipoVerificacaoEmail,
		Destinatario: email,
		NomeDestino:  nome,
		Assunto:      "Confirme seu e-mail — " + h.appName,
		Dados: map[string]string{
			"nome_destinatario": nome,
			"link_verificacao":  link,
			"nome_app":          h.appName,
			"logo_url":          h.logoURL(),
		},
	}) {
		return nil
	}
	log.Println("RabbitMQ indisponível, tentando envio direto do e-mail de verificação")

	body := fmt.Sprintf(
		"Olá, %s!\n\nPara ativar sua conta no %s, confirme seu e-mail clicando no link abaixo (válido por 48 horas):\n\n%s\n\nSe você não fez este cadastro, ignore este e-mail.",
		nome, h.appName, link,
	)
	return h.mailer.Send(email, "Confirme seu e-mail — "+h.appName, body)
}

// Run handles POST /acesso/login, mirroring LoginController::run.
func (h *LoginHandler) Run(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	if req.Email == "" || req.Senha == "" {
		httpx.ErrorWithFields(w, "Dados inválidos.", map[string]string{
			"email": "obrigatório", "senha": "obrigatório",
		}, 422)
		return
	}

	appKey := r.Header.Get("AppKey")
	if appKey == "" {
		httpx.Error(w, "Aplicação não encontrada!", 401)
		return
	}
	app, err := h.aplicacaoRepo.FindByCodigo(appKey)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if app == nil {
		httpx.Error(w, "Aplicação não encontrada!", 401)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	usuario, err := h.usuarioRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if usuario == nil || !repository.CheckPassword(usuario.Senha, req.Senha) {
		httpx.Error(w, "Usuário ou Senha Inválidos!", 401)
		return
	}
	// Verifica se o e-mail foi confirmado (campo email_verified_at)
	if usuario.EmailVerifiedAt == nil {
		httpx.Error(w, "E-mail não verificado, verifique o link de ativação no seu e-mail!", 401)
		return
	}

	var token string
	if req.LongToken == "sim" || req.LongToken == "1" {
		token = h.makeLongToken(usuario.ID)
	} else {
		token = h.makeToken(usuario.ID, fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	at, err := h.tokenRepo.CreateOrReuse(usuario.ID, app.ID, clientIP(r), token)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}

	// Métricas do Painel: cada login vira uma linha em acesso_log e atualiza
	// "visto por último". Falhas aqui não devem impedir o login.
	if err := h.usuarioRepo.RegistrarAcesso(usuario.ID, clientIP(r)); err != nil {
		log.Println("erro ao registrar acesso:", err)
	}
	_ = h.usuarioRepo.TouchAcesso(usuario.ID)

	httpx.JSON(w, 200, loginResponse{Usuario: usuario, Token: at.Token})
}

// loginResponse devolve o perfil completo (papel, município, GAT...) junto
// com o token, para o frontend montar o dashboard conforme o perfil sem uma
// segunda chamada.
type loginResponse struct {
	*models.Usuario
	Token string `json:"token"`
}

// Me handles GET /acesso/me — devolve o usuário do token atual.
func (h *LoginHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r)
	if u == nil {
		httpx.Error(w, "Token inválido!", 401)
		return
	}
	httpx.JSON(w, 200, u)
}

type cadastroRequest struct {
	Nome          string `json:"nome"`
	Email         string `json:"email"`
	Senha         string `json:"senha"`
	RunLogin      string `json:"run_login"`
	CPF           string `json:"cpf"`
	Telefone      string `json:"telefone"`
	Instituicao   string `json:"instituicao"`
	CursoFormacao string `json:"curso_formacao"`
	Papel         string `json:"papel"`
	MunicipioID   *int64 `json:"municipio_id"`
	AceiteTermos  bool   `json:"aceite_termos"`
}

// onlyDigits mantém só os dígitos de s (usado para normalizar CPF/telefone).
func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func strPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// Cadastro handles POST /cadastro — autocadastro público de profissionais,
// estudantes, pesquisadores, gestores e comunidade. Papéis institucionais
// (preceptor, coordenador, admin) NÃO podem ser escolhidos aqui: só um
// administrador os atribui depois (convite institucional).
func (h *LoginHandler) Cadastro(w http.ResponseWriter, r *http.Request) {
	var req cadastroRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	req.Nome = strings.TrimSpace(req.Nome)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Nome == "" || req.Email == "" {
		httpx.ErrorWithFields(w, "Dados inválidos.", map[string]string{
			"nome": "obrigatório", "email": "obrigatório",
		}, 422)
		return
	}
	if len(req.Senha) < 6 {
		httpx.Error(w, "A senha deve ter no mínimo 6 caracteres.", 422)
		return
	}

	papel := models.Papel(strings.ToLower(strings.TrimSpace(req.Papel)))
	if papel == "" {
		papel = models.PapelComunidade
	}
	if !models.PapeisAutoCadastro[papel] {
		httpx.Error(w, "Perfil não permitido no cadastro público. Perfis institucionais são atribuídos pela coordenação.", 422)
		return
	}
	cpf := onlyDigits(req.CPF)
	if cpf != "" && len(cpf) != 11 {
		httpx.ErrorWithFields(w, "Dados inválidos.", map[string]string{"cpf": "deve ter 11 dígitos"}, 422)
		return
	}
	if !req.AceiteTermos {
		httpx.Error(w, "É necessário aceitar os Termos de Uso e a Política de Privacidade.", 422)
		return
	}

	appKey := r.Header.Get("AppKey")
	if appKey == "" {
		httpx.Error(w, "Aplicação não encontrada!", 401)
		return
	}
	app, err := h.aplicacaoRepo.FindByCodigo(appKey)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if app == nil {
		httpx.Error(w, "Aplicação não encontrada!", 401)
		return
	}

	created, err := h.usuarioRepo.Create(&models.Usuario{
		Nome:          req.Nome,
		Email:         req.Email,
		Papel:         papel,
		CPF:           strPtr(cpf),
		Telefone:      strPtr(onlyDigits(req.Telefone)),
		Instituicao:   strPtr(req.Instituicao),
		CursoFormacao: strPtr(req.CursoFormacao),
		MunicipioID:   req.MunicipioID,
	}, req.Senha)
	if err != nil {
		if ae, ok := apperr.As(err); ok {
			httpx.Error(w, ae.Message, ae.Status)
		} else {
			log.Println(err)
			httpx.Error(w, "Erro interno.", 500)
		}
		return
	}

	// Envia e-mail de confirmação (SEM tabela extra)
	if err := h.sendVerificationEmail(created.Email, created.Nome); err != nil {
		log.Println("erro ao enviar e-mail de confirmação:", err)
	}

	if req.RunLogin != "sim" && req.RunLogin != "1" {
		httpx.JSON(w, 201, created)
		return
	}

	// A conta recém-criada ainda não foi confirmada por e-mail — login
	// automático só é possível depois da confirmação.
	if created.EmailVerifiedAt == nil {
		httpx.JSON(w, 201, created)
		return
	}

	token := h.makeToken(created.ID, fmt.Sprintf("%d", time.Now().UnixNano()))
	at, err := h.tokenRepo.CreateOrReuse(created.ID, app.ID, clientIP(r), token)
	if err != nil {
		httpx.Error(w, "Erro ao gerar token", 500)
		return
	}

	httpx.JSON(w, 201, loginResponse{Usuario: created, Token: at.Token})
}

// CheckToken handles POST /acesso/check_token
func (h *LoginHandler) CheckToken(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.Header.Get("TokenUser"), ":")
	if len(parts) < 3 {
		httpx.Error(w, "Token inválido!", 401)
		return
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || userID <= 0 || parts[1] == "" || parts[2] == "" {
		httpx.Error(w, "Token inválido!", 401)
		return
	}
	if _, err := h.tokenRepo.Check(userID, parts[1], parts[2]); err != nil {
		if ae, ok := apperr.As(err); ok {
			httpx.Error(w, ae.Message, ae.Status)
		} else {
			log.Println(err)
			httpx.Error(w, "Erro interno.", 500)
		}
		return
	}
	httpx.Success(w, "Sucesso!", 200)
}

// RefreshToken handles POST /acesso/refresh-token. It piggybacks on the
// same validation as CheckToken/RequireAuth (TokenRepository.Check), which
// already slides the token's expiration window forward on every valid
// call. Having a dedicated endpoint lets the frontend call it on an
// interval (or right when the app regains focus) purely to keep the
// session alive, without that traffic being tied to any particular
// feature request.
func (h *LoginHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.Header.Get("TokenUser"), ":")
	if len(parts) < 3 {
		httpx.Error(w, "Token inválido!", 401)
		return
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || userID <= 0 || parts[1] == "" || parts[2] == "" {
		httpx.Error(w, "Token inválido!", 401)
		return
	}

	at, err := h.tokenRepo.Check(userID, parts[1], parts[2])
	if err != nil {
		if ae, ok := apperr.As(err); ok {
			httpx.Error(w, ae.Message, ae.Status)
		} else {
			log.Println(err)
			httpx.Error(w, "Erro interno.", 500)
		}
		return
	}

	limiteMin := models.TempoLimiteAPIMinutes
	if models.IsLongToken(at.Token) {
		limiteMin = models.TempoLimiteLongoAPIMinutes
	}

	httpx.JSON(w, 200, map[string]any{
		"token":             at.Token,
		"expira_em_minutos": limiteMin,
		"renovado_em":       time.Now(),
	})
}

type esqueciSenhaRequest struct {
	Email string `json:"email"`
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// EsqueciSenha handles POST /acesso/esqueci-senha
func (h *LoginHandler) EsqueciSenha(w http.ResponseWriter, r *http.Request) {
	var req esqueciSenhaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		httpx.Error(w, "Informe um e-mail válido.", 422)
		return
	}

	const generic = "Se este e-mail estiver cadastrado, enviaremos um link para redefinir a senha."

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	usuario, err := h.usuarioRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if usuario == nil {
		httpx.Success(w, generic, 200)
		return
	}

	token, err := randomToken()
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if err := h.resetRepo.Create(req.Email, token); err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}

	link := h.linkComQuery("/redefinir-senha", req.Email, token)

	sent := h.enviarEmailJob(models.EmailJob{
		Tipo:         models.EmailTipoRedefinirSenha,
		Destinatario: req.Email,
		NomeDestino:  usuario.Nome,
		Assunto:      "Redefinição de senha — " + h.appName,
		Dados: map[string]string{
			"nome_destinatario": usuario.Nome,
			"link_redefinicao":  link,
			"nome_app":          h.appName,
			"logo_url":          h.logoURL(),
		},
	})
	if !sent {
		body := fmt.Sprintf(
			"Olá, %s!\n\nRecebemos um pedido para redefinir sua senha no %s.\nClique no link abaixo para escolher uma nova senha (válido por 1 hora):\n\n%s\n\nSe você não pediu isso, pode ignorar este e-mail.",
			usuario.Nome, h.appName, link,
		)
		if err := h.mailer.Send(req.Email, "Redefinição de senha — "+h.appName, body); err != nil {
			log.Println(err)
		}
	}

	httpx.Success(w, generic, 200)
}

type redefinirSenhaRequest struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	NovaSenha string `json:"nova_senha"`
}

// RedefinirSenha handles POST /acesso/redefinir-senha
func (h *LoginHandler) RedefinirSenha(w http.ResponseWriter, r *http.Request) {
	var req redefinirSenhaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "Requisição inválida.", 422)
		return
	}
	if req.Email == "" || req.Token == "" || len(req.NovaSenha) < 6 {
		httpx.Error(w, "Dados inválidos. A senha deve ter ao menos 6 caracteres.", 422)
		return
	}

	ok, err := h.resetRepo.Verify(req.Email, req.Token)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if !ok {
		httpx.Error(w, "Link inválido ou expirado. Solicite a redefinição novamente.", 401)
		return
	}

	if err := h.usuarioRepo.SetSenhaByEmail(req.Email, req.NovaSenha); err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	_ = h.resetRepo.Delete(req.Email)

	httpx.Success(w, "Senha redefinida com sucesso! Você já pode entrar.", 200)
}

type confirmarEmailRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// ConfirmarEmail handles POST /acesso/confirmar-email
// Agora valida o token SEM banco de dados e atualiza apenas email_verified_at
// ConfirmarEmail handles POST /acesso/confirmar-email
func (h *LoginHandler) ConfirmarEmail(w http.ResponseWriter, r *http.Request) {
	var req confirmarEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "Dados inválidos.", 422)
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Valida o token sem banco de dados
	ok, err := h.verifyVerificationToken(req.Email, req.Token)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if !ok {
		httpx.Error(w, "Link inválido ou expirado. Solicite um novo e-mail de confirmação.", 401)
		return
	}

	// Descobre se esta é a primeira confirmação antes de gravar a data.
	usuario, _ := h.usuarioRepo.FindByEmail(req.Email)
	primeiraConfirmacao := usuario != nil && usuario.EmailVerifiedAt == nil

	// Atualiza apenas o campo email_verified_at com a data/hora atual
	if err := h.usuarioRepo.SetEmailVerifiedByEmail(req.Email); err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}

	// Boas-vindas: só a primeira confirmação dispara o e-mail (evita spam
	// se a pessoa clicar no link mais de uma vez).
	if primeiraConfirmacao {
		h.enviarEmailJob(models.EmailJob{
			Tipo:         models.EmailTipoBoasVindas,
			Destinatario: usuario.Email,
			NomeDestino:  usuario.Nome,
			Assunto:      "Bem-vindo(a) ao " + h.appName,
			Dados: map[string]string{
				"nome_destinatario": usuario.Nome,
				"link_login":        h.frontendURL + "/login",
				"nome_app":          h.appName,
				"logo_url":          h.logoURL(),
			},
		})
	}

	httpx.Success(w, "E-mail confirmado com sucesso! Você já pode entrar.", 200)
}

type reenviarConfirmacaoRequest struct {
	Email string `json:"email"`
}

// ReenviarConfirmacao handles POST /acesso/reenviar-confirmacao
// Agora usa sendVerificationEmail sem banco
func (h *LoginHandler) ReenviarConfirmacao(w http.ResponseWriter, r *http.Request) {
	var req reenviarConfirmacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		httpx.Error(w, "Informe um e-mail válido.", 422)
		return
	}

	const generic = "Se este e-mail estiver cadastrado e pendente de confirmação, reenviaremos o link."

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	usuario, err := h.usuarioRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println(err)
		httpx.Error(w, "Erro interno.", 500)
		return
	}
	if usuario == nil || usuario.EmailVerifiedAt != nil {
		httpx.Success(w, generic, 200)
		return
	}

	// Reenvia o e-mail (SEM tabela extra)
	if err := h.sendVerificationEmail(usuario.Email, usuario.Nome); err != nil {
		log.Println("erro ao reenviar e-mail de confirmação:", err)
	}

	httpx.Success(w, generic, 200)
}
