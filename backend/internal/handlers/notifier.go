package handlers

import (
	"log"
	"time"

	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/queue"
)

// Notifier publica e-mails de notificação (template "generico.html") para
// eventos da plataforma — inscrição em atividade, conclusão de curso etc.
// É "melhor esforço": se o RabbitMQ não estiver disponível, a notificação
// é simplesmente ignorada e a ação principal do usuário segue normalmente.
type Notifier struct {
	rabbit      *queue.RabbitMQ
	frontendURL string
	appName     string
}

func NewNotifier(rabbit *queue.RabbitMQ, frontendURL, appName string) *Notifier {
	return &Notifier{rabbit: rabbit, frontendURL: frontendURL, appName: appName}
}

// Generico envia um e-mail com título, mensagem e (opcionalmente) um botão
// de ação. textoAcao/linkAcao podem ser vazios.
func (n *Notifier) Generico(u *models.Usuario, assunto, titulo, mensagem, textoAcao, linkAcao string) {
	if n == nil || n.rabbit == nil || u == nil {
		return
	}
	job := models.EmailJob{
		Tipo:         models.EmailTipoGenerico,
		Destinatario: u.Email,
		NomeDestino:  u.Nome,
		Assunto:      assunto + " — " + n.appName,
		Dados: map[string]string{
			"nome_destinatario": u.Nome,
			"titulo":            titulo,
			"mensagem":          mensagem,
			"texto_acao":        textoAcao,
			"link_acao":         linkAcao,
			"nome_app":          n.appName,
			"logo_url":          n.frontendURL + "/petsaudeclima_icon.png",
		},
		CriadoEm: time.Now(),
	}
	if err := n.rabbit.Publish(queue.QueueEmail, job); err != nil {
		log.Println("notifier: erro ao publicar e-mail:", err)
	}
}

// FrontendURL expõe a URL base do frontend para montar links de ação.
func (n *Notifier) FrontendURL() string {
	if n == nil {
		return ""
	}
	return n.frontendURL
}
