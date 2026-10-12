// Package mailer envia e-mails transacionais simples (texto puro) quando o
// RabbitMQ/worker-email não está disponível — é o "plano B" usado por
// handlers/login_handler.go para confirmação de e-mail e redefinição de
// senha. Os e-mails com identidade visual do PET-Saúde Clima são
// renderizados pelo worker-email a partir de templates/email/*.html.
//
// Se MAIL_HOST não estiver configurado, Send apenas registra a mensagem no
// log, para o fluxo continuar funcionando de ponta a ponta em desenvolvimento
// (o link de confirmação/redefinição aparece no console do backend).
package mailer

import (
	"fmt"
	"log"
	"mime"
	"net/smtp"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

type Mailer struct {
	cfg Config
}

func New(cfg Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) configured() bool {
	return m.cfg.Host != "" && m.cfg.Port != "" && m.cfg.From != ""
}

// Send entrega um e-mail em texto puro, ou o registra no log quando o SMTP
// não está configurado.
func (m *Mailer) Send(to, subject, body string) error {
	if !m.configured() {
		log.Printf("[mailer] SMTP não configurado — e-mail não enviado de verdade.\nPara: %s\nAssunto: %s\n%s\n", to, subject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)
	msg := []byte("From: " + m.cfg.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + mime.QEncoding.Encode("UTF-8", subject) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	var auth smtp.Auth
	if m.cfg.User != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Pass, m.cfg.Host)
	}

	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg); err != nil {
		log.Printf("[mailer] falha ao enviar e-mail para %s: %v", to, err)
		return err
	}

	log.Printf("[mailer] ✅ E-mail enviado para %s", to)
	return nil
}
