// Package emailtpl renderiza os templates HTML de e-mail (embutidos no
// binário). Usado pelo worker-email e também pelo envio direto da API quando
// a fila do RabbitMQ não está disponível.
package emailtpl

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/liedsonlb/petsaude-clima/templates"
)

const tipoGenerico = "generico"

// Render renderiza templates/email/<tipo>.html com os dados informados. Se o
// tipo não existir, cai para "generico".
func Render(tipo string, data map[string]string) (string, error) {
	path := "email/" + tipo + ".html"
	if _, err := templates.FS.ReadFile(path); err != nil {
		path = "email/" + tipoGenerico + ".html"
	}

	tmpl, err := template.ParseFS(templates.FS, path)
	if err != nil {
		return "", fmt.Errorf("erro ao carregar template de email (%s): %w", path, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("erro ao renderizar template de email: %w", err)
	}
	return buf.String(), nil
}
