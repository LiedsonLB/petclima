package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/petsaude/petclima-backend/internal/models"
	"github.com/petsaude/petclima-backend/internal/queue"
	"github.com/petsaude/petclima-backend/internal/service"
)

// maxTentativas é quantas vezes um e-mail é reenviado antes de ser
// descartado (evita laço infinito quando o SMTP está errado/fora do ar).
const maxTentativas = 5

// StartEmailWorker consome a fila petclima.emails e envia cada EmailJob por
// SMTP (com o template HTML correspondente a job.Tipo — ver
// templates/email/). Falhas são reenfileiradas com espera crescente, até
// maxTentativas.
func StartEmailWorker(rabbit *queue.RabbitMQ, emailService *service.EmailService) {
	log.Println("📧 Worker de Email: conectando à fila", queue.QueueEmail)

	if !emailService.Configurado() {
		log.Println("⚠️  SMTP não configurado (MAIL_HOST / MAIL_PORT / MAIL_FROM_ADDRESS): os e-mails da fila serão apenas registrados no log, sem envio.")
	}

	msgs, err := rabbit.Consume(queue.QueueEmail)
	if err != nil {
		log.Fatalf("❌ Erro ao consumir fila %s: %v", queue.QueueEmail, err)
	}

	log.Println("✅ Worker de Email aguardando mensagens...")

	for msg := range msgs {
		var job models.EmailJob
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("❌ Erro ao parsear EmailJob: %v", err)
			msg.Nack(false, false) // descarta mensagem malformada (não requeue)
			continue
		}

		log.Printf("📨 Processando email tipo=%s destinatario=%s", job.Tipo, job.Destinatario)

		if !emailService.Configurado() {
			log.Printf("📝 [sem SMTP] e-mail tipo=%s para %s assunto=%q dados=%v", job.Tipo, job.Destinatario, job.Assunto, job.Dados)
			msg.Ack(false)
			continue
		}

		if err := emailService.Enviar(job); err != nil {
			job.Tentativas++
			log.Printf("❌ Erro ao enviar email para %s (tentativa %d/%d): %v", job.Destinatario, job.Tentativas, maxTentativas, err)
			if job.Tentativas >= maxTentativas {
				log.Printf("🗑️  Descartando e-mail para %s após %d tentativas", job.Destinatario, maxTentativas)
				msg.Ack(false)
				continue
			}
			// Republica com o contador atualizado, depois de uma espera
			// crescente (falha temporária de SMTP), e confirma a original.
			time.Sleep(time.Duration(job.Tentativas) * 5 * time.Second)
			if perr := rabbit.Publish(queue.QueueEmail, job); perr != nil {
				log.Printf("❌ Não foi possível reenfileirar: %v", perr)
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
			continue
		}

		log.Printf("✅ Email enviado para %s (tipo=%s)", job.Destinatario, job.Tipo)
		msg.Ack(false)
	}
}
