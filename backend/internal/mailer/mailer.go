// Package mailer entrega e-mails por SMTP. É usado em dois lugares:
//
//   - pela API, como "plano B" quando o RabbitMQ/worker-email não está
//     disponível (confirmação de e-mail e redefinição de senha), e
//   - pelo worker-email (internal/service), via Deliver.
//
// Se MAIL_HOST não estiver configurado, Send/SendHTML apenas registram a
// mensagem no log, para o fluxo funcionar de ponta a ponta em desenvolvimento
// (o link de confirmação/redefinição aparece no console do backend).
package mailer

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Pass     string
	From     string // endereço do remetente (aceita também "Nome <a@b.c>")
	FromName string // nome exibido (opcional)
}

type Mailer struct {
	cfg Config
}

func New(cfg Config) *Mailer {
	return &Mailer{cfg: cfg}
}

// Configured informa se há SMTP configurado (host, porta e remetente).
func (m *Mailer) Configured() bool {
	return m.cfg.Host != "" && m.cfg.Port != "" && m.cfg.From != ""
}

// Send entrega um e-mail em texto puro.
func (m *Mailer) Send(to, subject, body string) error {
	return m.send(to, subject, "text/plain", body)
}

// SendHTML entrega um e-mail HTML (ex.: templates com a identidade do projeto).
func (m *Mailer) SendHTML(to, subject, html string) error {
	return m.send(to, subject, "text/html", html)
}

func (m *Mailer) send(to, subject, contentType, body string) error {
	if !m.Configured() {
		log.Printf("[mailer] SMTP não configurado — e-mail não enviado de verdade.\nPara: %s\nAssunto: %s\n%s\n", to, subject, body)
		return nil
	}

	msg := BuildMessage(m.cfg.FromName, m.cfg.From, to, subject, contentType, body)
	if err := Deliver(m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Pass, m.cfg.From, to, msg); err != nil {
		log.Printf("[mailer] falha ao enviar e-mail para %s: %v", to, err)
		return err
	}
	log.Printf("[mailer] ✅ E-mail enviado para %s", to)
	return nil
}

// BuildMessage monta uma mensagem MIME simples (corpo em base64, UTF-8) com
// os cabeçalhos que provedores como Gmail/Outlook exigem (Date, Message-ID).
func BuildMessage(fromName, from, to, subject, contentType, body string) []byte {
	var buf bytes.Buffer
	writeHeaders(&buf, fromName, from, to, subject)
	fmt.Fprintf(&buf, "Content-Type: %s; charset=\"UTF-8\"\r\n", contentType)
	buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(ChunkedBase64([]byte(body)))
	return buf.Bytes()
}

// WriteHeaders escreve From/To/Subject/Date/Message-ID/MIME-Version.
// Exportada para o worker montar mensagens com anexos.
func WriteHeaders(buf *bytes.Buffer, fromName, from, to, subject string) {
	writeHeaders(buf, fromName, from, to, subject)
}

func writeHeaders(buf *bytes.Buffer, fromName, from, to, subject string) {
	addr := envelopeAddress(from)
	if fromName != "" {
		fmt.Fprintf(buf, "From: %s\r\n", (&mail.Address{Name: fromName, Address: addr}).String())
	} else {
		fmt.Fprintf(buf, "From: %s\r\n", from)
	}
	fmt.Fprintf(buf, "To: %s\r\n", to)
	fmt.Fprintf(buf, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject))
	fmt.Fprintf(buf, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(buf, "Message-ID: <%s@%s>\r\n", randomID(), domainOf(addr))
	buf.WriteString("MIME-Version: 1.0\r\n")
}

// ChunkedBase64 codifica em base64 quebrando em linhas de 76 caracteres
// (RFC 2045).
func ChunkedBase64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	var out strings.Builder
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		out.WriteString(encoded[i:end])
		out.WriteString("\r\n")
	}
	return out.String()
}

func envelopeAddress(from string) string {
	if a, err := mail.ParseAddress(from); err == nil {
		return a.Address
	}
	return from
}

func domainOf(addr string) string {
	if i := strings.LastIndex(addr, "@"); i >= 0 && i < len(addr)-1 {
		return addr[i+1:]
	}
	return "localhost"
}

func randomID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Deliver entrega a mensagem já montada ao servidor SMTP.
//   - porta 465: TLS implícito;
//   - demais portas (587, 25, 1025...): STARTTLS quando o servidor oferece;
//   - autentica (PLAIN) só se user estiver preenchido — servidores de teste
//     locais (MailHog, Mailpit) não pedem login.
func Deliver(host, port, user, pass, from, to string, msg []byte) error {
	addr := net.JoinHostPort(host, port)
	dialer := net.Dialer{Timeout: 15 * time.Second}

	var conn net.Conn
	var err error
	if port == "465" {
		conn, err = tls.DialWithDialer(&dialer, "tcp", addr, &tls.Config{ServerName: host})
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("conexão SMTP com %s: %w", addr, err)
	}
	_ = conn.SetDeadline(time.Now().Add(60 * time.Second))

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("cumprimento SMTP: %w", err)
	}
	defer c.Close()

	if port != "465" {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
				return fmt.Errorf("STARTTLS: %w", err)
			}
		}
	}
	if user != "" {
		if ok, _ := c.Extension("AUTH"); !ok {
			return fmt.Errorf("servidor SMTP não oferece autenticação, mas MAIL_USERNAME está preenchido")
		}
		if err := c.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
			return fmt.Errorf("autenticação SMTP: %w", err)
		}
	}
	if err := c.Mail(envelopeAddress(from)); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	if err := c.Rcpt(envelopeAddress(to)); err != nil {
		return fmt.Errorf("RCPT TO: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("escrevendo mensagem: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("finalizando mensagem: %w", err)
	}
	return c.Quit()
}
