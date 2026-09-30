// Package templates embute os arquivos HTML de e-mail no binário, para que
// tanto a API quanto o worker-email consigam renderizá-los sem depender de
// arquivos no disco do container.
package templates

import "embed"

//go:embed email/*.html
var FS embed.FS
