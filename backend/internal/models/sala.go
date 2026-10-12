package models

import "time"

// Tipos de sala virtual (LiveKit) suportados pela plataforma. Todas usam a
// mesma infraestrutura de videoconferência; o tipo só muda o rótulo/ícone
// exibido no frontend.
const (
	SalaTipoReuniao    = "reuniao"    // reunião de equipe (GAT/macrorregião)
	SalaTipoOficina    = "oficina"    // oficina/capacitação remota (Educação Permanente)
	SalaTipoTelessaude = "telessaude" // teleconsultoria/regulação a distância
)

// Sala é uma sala de videoconferência (LiveKit) que pode existir sozinha
// (reunião avulsa) ou vinculada a uma Atividade (oficina/capacitação
// remota ou híbrida cadastrada em "Educação Permanente"), permitindo que os
// inscritos entrem direto da tela da atividade.
type Sala struct {
	ID          int64   `json:"id" db:"id"`
	Nome        string  `json:"nome" db:"nome"`
	Codigo      string  `json:"codigo" db:"codigo"` // nome da room no LiveKit
	Tipo        string  `json:"tipo" db:"tipo"`
	Descricao   *string `json:"descricao" db:"descricao"`
	AtividadeID *int64  `json:"atividade_id" db:"atividade_id"`

	CriadoPor int64      `json:"criado_por" db:"criado_por"`
	Ativa     bool       `json:"ativa" db:"ativa"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`

	// Preenchido em memória pelo handler (não é coluna) — consultado ao
	// vivo no LiveKit via ListParticipants, não persistido.
	ParticipantesOnline int `json:"participantes_online,omitempty" db:"-"`
}

func (Sala) TableName() string { return "sala" }

// ParticipanteInfo descreve quem está conectado agora em uma sala (lido do
// LiveKit em tempo real via livekit.ListParticipants).
type ParticipanteInfo struct {
	Identity    string `json:"identity"`
	Nome        string `json:"nome"`
	Foto        string `json:"foto,omitempty"`
	MicEnabled  bool   `json:"micEnabled"`
	Deafened    bool   `json:"deafened,omitempty"`
	CameraOn    bool   `json:"cameraOn,omitempty"`
	ScreenShare bool   `json:"screenShare,omitempty"`
}
