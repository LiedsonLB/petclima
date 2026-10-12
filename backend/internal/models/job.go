package models

import "time"

// Tipos de email conhecidos pelo worker-email. Cada tipo mapeia para um
// template em templates/email/<tipo>.html. Novos tipos podem ser
// adicionados livremente — o worker só precisa do arquivo de template
// correspondente (cai para "generico.html" quando não encontra um
// específico, ver internal/service/email_service.go).
const (
	EmailTipoBoasVindas       = "boas_vindas"
	EmailTipoVerificacaoEmail = "verificacao_email"
	EmailTipoRedefinirSenha   = "redefinir_senha"
	EmailTipoGenerico         = "generico"
)

// EmailJob é o payload publicado na fila `petclima.emails`. `Dados` carrega
// as variáveis usadas pelo template HTML (nome, link, título da atividade,
// etc) — fica a critério de quem publica o job preencher o que o template
// daquele Tipo espera.
type EmailJob struct {
	Tipo         string            `json:"tipo"`
	Destinatario string            `json:"destinatario"`
	NomeDestino  string            `json:"nome_destinatario"`
	Assunto      string            `json:"assunto"`
	Dados        map[string]string `json:"dados"`
	Anexos       []EmailAnexo      `json:"anexos,omitempty"`
	CriadoEm     time.Time         `json:"criado_em"`
}

// EmailAnexo é um anexo do e-mail (ex.: certificado de conclusão de curso
// em PDF). Conteudo traz os bytes já prontos; Caminho é usado apenas se o
// worker precisar ler algo do disco de uploads.
type EmailAnexo struct {
	NomeArquivo string `json:"nome_arquivo"`
	Conteudo    []byte `json:"conteudo,omitempty"`
	Caminho     string `json:"caminho,omitempty"`
}
