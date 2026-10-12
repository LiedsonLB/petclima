package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/petsaude/petclima-backend/internal/httpx"
	"github.com/petsaude/petclima-backend/internal/middleware"
)

// UploadHandler salva arquivos enviados pelos usuários (fotos de perfil,
// capas e materiais da Biblioteca) em disco, numa pasta fora de backend/ e
// frontend/ — na raiz do projeto — para que o nginx sirva direto
// (location /uploads/ -> essa pasta) sem passar pela API Go.
type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir}
}

const (
	maxImagemSize  = 5 << 20  // 5MB
	maxMaterialSize = 50 << 20 // 50MB — PDFs, podcasts e vídeos curtos
)

var allowedImageExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".jfif": true,
}

// Materiais da Biblioteca/Comunicação: documentos, áudio, vídeo e imagens.
var allowedMaterialExt = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".ppt": true, ".pptx": true,
	".mp3": true, ".wav": true, ".ogg": true, ".m4a": true,
	".mp4": true, ".webm": true, ".mov": true,
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
}

// pastasImagemPermitidas restringe onde o cliente pode gravar imagens,
// evitando path traversal via o campo "pasta" do form (ex.: "../../etc").
var pastasImagemPermitidas = map[string]bool{
	"fotos":    true, // foto de perfil do usuário
	"capas":    true, // capa de curso/material
	"noticias": true, // imagem de notícia (staff)
}

func randomFileName() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// salvar grava o arquivo do form em <UploadDir>/<pasta>/<token><ext> e
// devolve o nome gerado.
func (h *UploadHandler) salvar(file io.Reader, pasta, ext string) (string, error) {
	folder := filepath.Join(h.uploadDir, pasta)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}
	name, err := randomFileName()
	if err != nil {
		return "", err
	}
	name += ext

	dst, err := os.Create(filepath.Join(folder, name))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return name, nil
}

// Foto handles POST /upload/foto (multipart/form-data):
//   - campo "arquivo": obrigatório, a imagem.
//   - campo "pasta": opcional, uma de pastasImagemPermitidas (default "fotos").
//   - campo "foto_antiga": opcional, URL antiga (na mesma pasta) a apagar.
//
// Devolve `{url}` (`/uploads/<pasta>/<token>.<ext>`), para ser gravado via
// PUT /usuarios/{id} (foto) ou nos formulários de curso/material/notícia
// (capa_url). Pastas diferentes de "fotos" exigem staff.
func (h *UploadHandler) Foto(w http.ResponseWriter, r *http.Request) {
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImagemSize+1<<20)
	if err := r.ParseMultipartForm(maxImagemSize); err != nil {
		httpx.Error(w, "Arquivo inválido ou muito grande (máximo 5MB).", 422)
		return
	}

	file, header, err := r.FormFile("arquivo")
	if err != nil {
		httpx.Error(w, "Envie o arquivo no campo \"arquivo\".", 422)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExt[ext] {
		httpx.Error(w, "Formato não suportado. Use jpg, png, webp ou gif.", 422)
		return
	}

	pasta := r.FormValue("pasta")
	if pasta == "" {
		pasta = "fotos"
	}
	if !pastasImagemPermitidas[pasta] {
		httpx.Error(w, "Pasta de upload inválida.", 422)
		return
	}
	if pasta != "fotos" && !session.IsStaff() {
		httpx.Error(w, "Você não tem permissão para acessar esse recurso.", 403)
		return
	}

	name, err := h.salvar(file, pasta, ext)
	if err != nil {
		httpx.Error(w, "Erro interno ao salvar o arquivo.", 500)
		return
	}

	// Apaga o arquivo antigo (se existir, na mesma pasta e diferente do novo).
	if oldURL := r.FormValue("foto_antiga"); oldURL != "" {
		prefix := "/uploads/" + pasta + "/"
		oldName := strings.TrimPrefix(oldURL, prefix)
		if oldName != "" && oldName != oldURL && oldName != name && !strings.ContainsAny(oldName, `/\`) {
			_ = os.Remove(filepath.Join(h.uploadDir, pasta, oldName))
		}
	}

	httpx.JSON(w, 201, map[string]any{"url": "/uploads/" + pasta + "/" + name})
}

// Material handles POST /upload/material (multipart/form-data, campo
// "arquivo") — staff. Salva PDFs, cartilhas, áudios, vídeos e cards em
// <UploadDir>/materiais/ e devolve `{url, nome, tamanho}`; a URL vai no
// campo arquivo_url de POST /materiais.
func (h *UploadHandler) Material(w http.ResponseWriter, r *http.Request) {
	session := middleware.UserFromContext(r)
	if session == nil {
		httpx.Error(w, "Não autenticado.", 401)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxMaterialSize+1<<20)
	if err := r.ParseMultipartForm(maxMaterialSize); err != nil {
		httpx.Error(w, "Arquivo inválido ou muito grande (máximo 50MB).", 422)
		return
	}

	file, header, err := r.FormFile("arquivo")
	if err != nil {
		httpx.Error(w, "Envie o arquivo no campo \"arquivo\".", 422)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedMaterialExt[ext] {
		httpx.Error(w, "Formato não suportado. Use PDF, Office, mp3/wav/ogg/m4a, mp4/webm/mov ou imagem.", 422)
		return
	}

	name, err := h.salvar(file, "materiais", ext)
	if err != nil {
		httpx.Error(w, "Erro interno ao salvar o arquivo.", 500)
		return
	}

	httpx.JSON(w, 201, map[string]any{
		"url":     "/uploads/materiais/" + name,
		"nome":    header.Filename,
		"tamanho": header.Size,
	})
}
