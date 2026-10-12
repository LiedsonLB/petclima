package router

import (
	"net/http"

	"github.com/petsaude/petclima-backend/internal/handlers"
	"github.com/petsaude/petclima-backend/internal/mailer"
	"github.com/petsaude/petclima-backend/internal/middleware"
	"github.com/petsaude/petclima-backend/internal/queue"
	"github.com/petsaude/petclima-backend/internal/repository"
)

type Deps struct {
	UsuarioRepo       *repository.UsuarioRepository
	AplicacaoRepo     *repository.AplicacaoRepository
	TokenRepo         *repository.TokenRepository
	CidadeRepo        *repository.CidadeRepository
	TerritorioRepo    *repository.TerritorioRepository
	PasswordResetRepo *repository.PasswordResetRepository
	CursoRepo         *repository.CursoRepository
	MaterialRepo      *repository.MaterialRepository
	NoticiaRepo       *repository.NoticiaRepository
	AtividadeRepo     *repository.AtividadeRepository
	PainelRepo        *repository.PainelRepository
	SalaRepo          *repository.SalaRepository

	AppName     string
	TokenKey    string
	AppKey      string
	FrontendURL string
	Mailer      *mailer.Mailer
	UploadDir   string

	// ---- videoconferência (LiveKit) + fila de e-mails (RabbitMQ) -------
	Rabbit           *queue.RabbitMQ
	LiveKitAPIKey    string
	LiveKitAPISecret string
	LiveKitURL       string
	LiveKitPublicURL string

	// ---- Verificação de e-mail (tokens auto-validáveis) ----------------
	VerificationSecret string
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API PET-Saúde Clima (Go)"))
	})

	// Três níveis de acesso, espelhando o planejamento:
	//   público      -> sem middleware (portal, biblioteca aberta)
	//   opcional     -> público, mas enriquecido se houver login
	//   auth + papel -> área autenticada / staff / coordenação / admin
	auth := middleware.RequireAuth(d.UsuarioRepo, d.TokenRepo)
	optional := middleware.OptionalAuth(d.UsuarioRepo, d.TokenRepo)
	staff := func(h http.HandlerFunc) http.Handler { return auth(middleware.RequireStaff(h)) }
	coord := func(h http.HandlerFunc) http.Handler { return auth(middleware.RequireCoordenador(h)) }
	admin := func(h http.HandlerFunc) http.Handler { return auth(middleware.RequireAdmin(h)) }

	notifier := handlers.NewNotifier(d.Rabbit, d.FrontendURL, d.AppName)

	// ---- acesso / cadastro (público) -------------------------------------
	login := handlers.NewLoginHandler(
		d.UsuarioRepo, d.AplicacaoRepo, d.TokenRepo, d.PasswordResetRepo,
		d.Mailer, d.Rabbit, d.FrontendURL, d.AppName, d.TokenKey, d.AppKey, d.VerificationSecret,
	)
	mux.HandleFunc("POST /acesso/login", login.Run)
	mux.HandleFunc("POST /acesso/check_token", login.CheckToken)
	// Chamado periodicamente (ou ao voltar o foco na aba) pelo frontend
	// para renovar a sessão automaticamente, sem exigir novo login.
	mux.HandleFunc("POST /acesso/refresh-token", login.RefreshToken)
	mux.HandleFunc("POST /acesso/esqueci-senha", login.EsqueciSenha)
	mux.HandleFunc("POST /acesso/redefinir-senha", login.RedefinirSenha)
	mux.HandleFunc("POST /acesso/confirmar-email", login.ConfirmarEmail)
	mux.HandleFunc("POST /acesso/reenviar-confirmacao", login.ReenviarConfirmacao)
	mux.Handle("GET /acesso/me", auth(http.HandlerFunc(login.Me)))
	mux.HandleFunc("POST /cadastro", login.Cadastro)

	// ---- 🗺️ territórios (público) -----------------------------------------
	territorio := handlers.NewTerritorioHandler(d.TerritorioRepo, d.CidadeRepo)
	mux.HandleFunc("GET /territorios/macrorregioes", territorio.Macrorregioes)
	mux.HandleFunc("GET /territorios/municipios", territorio.Municipios)
	mux.HandleFunc("GET /territorios/gats", territorio.Gats)

	// ---- 👥 gestão de usuários --------------------------------------------
	usuario := handlers.NewUsuarioHandler(d.UsuarioRepo, d.UploadDir)
	mux.Handle("GET /usuarios", coord(usuario.All))
	mux.Handle("GET /usuarios/count", coord(usuario.Count))
	mux.Handle("GET /usuarios/resumo", coord(usuario.ResumoPapeis))
	mux.Handle("POST /usuarios/heartbeat", auth(http.HandlerFunc(usuario.Heartbeat)))
	mux.Handle("GET /usuarios/{id}", auth(http.HandlerFunc(usuario.Find)))
	mux.Handle("POST /usuarios", admin(usuario.Save))
	mux.Handle("PUT /usuarios/{id}", auth(http.HandlerFunc(usuario.Update)))
	mux.Handle("PUT /usuarios/{id}/papel", admin(usuario.AlterarPapel))
	mux.Handle("DELETE /usuarios/{id}", admin(usuario.Delete))
	mux.Handle("POST /usuarios/restore/{id}", admin(usuario.Restore))

	// ---- uploads -----------------------------------------------------------
	upload := handlers.NewUploadHandler(d.UploadDir)
	mux.Handle("POST /upload/foto", auth(http.HandlerFunc(upload.Foto)))
	mux.Handle("POST /upload/material", staff(upload.Material))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(d.UploadDir))))

	// ---- 🏠 portal: notícias ------------------------------------------------
	noticia := handlers.NewNoticiaHandler(d.NoticiaRepo)
	mux.HandleFunc("GET /noticias", noticia.All)
	mux.HandleFunc("GET /noticias/{id}", noticia.Find)
	mux.Handle("POST /noticias", staff(noticia.Save))
	mux.Handle("DELETE /noticias/{id}", staff(noticia.Delete))

	// ---- 📚 biblioteca + 📢 comunicação (leitura/download abertos) ----------
	material := handlers.NewMaterialHandler(d.MaterialRepo)
	mux.HandleFunc("GET /materiais", material.All)
	mux.HandleFunc("GET /materiais/{id}", material.Find)
	mux.Handle("POST /materiais/{id}/download", optional(http.HandlerFunc(material.Download)))
	mux.Handle("POST /materiais", staff(material.Save))
	mux.Handle("DELETE /materiais/{id}", staff(material.Delete))

	// ---- 🎓 cursos / AVA ------------------------------------------------------
	curso := handlers.NewCursoHandler(d.CursoRepo, notifier)
	mux.Handle("GET /cursos", optional(http.HandlerFunc(curso.All)))
	mux.Handle("GET /cursos/minhas-matriculas", auth(http.HandlerFunc(curso.Minhas)))
	mux.Handle("GET /cursos/{id}", optional(http.HandlerFunc(curso.Find)))
	mux.Handle("POST /cursos", staff(curso.Save))
	mux.Handle("POST /cursos/{id}/matricula", auth(http.HandlerFunc(curso.Matricular)))
	mux.Handle("PUT /cursos/{id}/progresso", auth(http.HandlerFunc(curso.AtualizarProgresso)))
	mux.HandleFunc("GET /certificados/{codigo}", curso.ValidarCertificado)

	// ---- 📝 atividades/ações + 🧑‍⚕️ educação permanente ------------------------
	atividade := handlers.NewAtividadeHandler(d.AtividadeRepo, d.SalaRepo, notifier)
	mux.Handle("GET /atividades", optional(http.HandlerFunc(atividade.All)))
	mux.Handle("GET /atividades/minhas", auth(http.HandlerFunc(atividade.Minhas)))
	mux.Handle("GET /atividades/{id}", optional(http.HandlerFunc(atividade.Find)))
	mux.Handle("POST /atividades", staff(atividade.Save))
	mux.Handle("DELETE /atividades/{id}", staff(atividade.Delete))
	mux.Handle("POST /atividades/{id}/inscricao", auth(http.HandlerFunc(atividade.Inscrever)))
	mux.Handle("DELETE /atividades/{id}/inscricao", auth(http.HandlerFunc(atividade.CancelarInscricao)))
	mux.Handle("GET /atividades/{id}/inscritos", staff(atividade.Inscritos))
	mux.Handle("PUT /atividades/{id}/presenca/{usuarioId}", staff(atividade.MarcarPresenca))

	// ---- 🎥 salas virtuais (LiveKit): oficinas remotas, telessaúde, reuniões --
	sala := handlers.NewSalaHandler(d.SalaRepo, d.AtividadeRepo, d.LiveKitAPIKey, d.LiveKitAPISecret, d.LiveKitURL, d.LiveKitPublicURL)
	mux.Handle("GET /salas", auth(http.HandlerFunc(sala.All)))
	mux.Handle("GET /salas/{id}", auth(http.HandlerFunc(sala.Find)))
	mux.Handle("GET /salas/{id}/participantes", auth(http.HandlerFunc(sala.Participantes)))
	mux.Handle("POST /salas", staff(sala.Save))
	mux.Handle("POST /salas/{id}/entrar", auth(http.HandlerFunc(sala.Entrar)))
	mux.Handle("POST /salas/{id}/encerrar", auth(http.HandlerFunc(sala.Encerrar)))

	// ---- 🔎 busca -----------------------------------------------------------
	busca := handlers.NewBuscaHandler(d.MaterialRepo, d.CursoRepo, d.AtividadeRepo)
	mux.HandleFunc("GET /busca", busca.All)

	// ---- 📊 painel -----------------------------------------------------------
	painel := handlers.NewPainelHandler(d.PainelRepo)
	mux.Handle("GET /painel/resumo", auth(http.HandlerFunc(painel.Resumo)))
	mux.Handle("GET /painel/indicadores", coord(painel.Indicadores))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, AppKey, TokenUser")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
