package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/petsaude/petclima-backend/internal/config"
	"github.com/petsaude/petclima-backend/internal/emailtpl"
	"github.com/petsaude/petclima-backend/internal/mailer"
	"github.com/petsaude/petclima-backend/internal/models"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// Configurado informa se há SMTP configurado (MAIL_HOST e MAIL_FROM_ADDRESS).
func (e *EmailService) Configurado() bool {
	return e.cfg.SMTPHost != "" && e.cfg.SMTPPort != 0 && e.cfg.SMTPFrom != ""
}

// Enviar renderiza o template embutido no binário correspondente ao Tipo do
// job (templates/email/<tipo>.html, com fallback para "generico") e o entrega
// por SMTP, com anexos opcionais.
func (e *EmailService) Enviar(job models.EmailJob) error {
	data := map[string]string{}
	for k, v := range job.Dados {
		data[k] = v
	}
	if _, ok := data["nome_destinatario"]; !ok {
		data["nome_destinatario"] = job.NomeDestino
	}
	data["assunto"] = job.Assunto

	html, err := emailtpl.Render(job.Tipo, data)
	if err != nil {
		return err
	}

	msg, err := buildMIMEMessage(e.cfg, job.Destinatario, job.Assunto, html, job.Anexos)
	if err != nil {
		return err
	}

	return mailer.Deliver(
		e.cfg.SMTPHost, strconv.Itoa(e.cfg.SMTPPort),
		e.cfg.SMTPUser, e.cfg.SMTPPass,
		e.cfg.SMTPFrom, job.Destinatario, msg,
	)
}

// buildMIMEMessage monta o e-mail: corpo HTML em base64 e, se houver,
// anexos (multipart/mixed).
func buildMIMEMessage(cfg *config.Config, destino, assunto, htmlBody string, anexos []models.EmailAnexo) ([]byte, error) {
	var buf bytes.Buffer
	mailer.WriteHeaders(&buf, cfg.AppName, cfg.SMTPFrom, destino, assunto)

	if len(anexos) == 0 {
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		buf.WriteString(mailer.ChunkedBase64([]byte(htmlBody)))
		return buf.Bytes(), nil
	}

	boundary := "petclima-boundary-42"
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%q\r\n\r\n", boundary)

	fmt.Fprintf(&buf, "--%s\r\n", boundary)
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(mailer.ChunkedBase64([]byte(htmlBody)))
	buf.WriteString("\r\n")

	for _, anexo := range anexos {
		content := anexo.Conteudo
		if len(content) == 0 && anexo.Caminho != "" {
			var err error
			content, err = os.ReadFile(anexo.Caminho)
			if err != nil {
				return nil, fmt.Errorf("erro ao ler anexo %s: %w", anexo.Caminho, err)
			}
		}
		if len(content) == 0 {
			log.Printf("anexo %q vazio, ignorado", anexo.NomeArquivo)
			continue
		}
		fmt.Fprintf(&buf, "--%s\r\n", boundary)
		fmt.Fprintf(&buf, "Content-Type: application/octet-stream; name=%q\r\n", anexo.NomeArquivo)
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		fmt.Fprintf(&buf, "Content-Disposition: attachment; filename=%q\r\n\r\n", anexo.NomeArquivo)
		buf.WriteString(mailer.ChunkedBase64(content))
		buf.WriteString("\r\n")
	}
	fmt.Fprintf(&buf, "--%s--\r\n", boundary)

	return buf.Bytes(), nil
}
