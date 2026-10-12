package repository

import (
	"database/sql"
	"time"
)

// Metas do projeto (planejamento PET-Saúde Clima) usadas como denominador
// dos indicadores do Painel. Ficam aqui, em um único lugar, para serem
// facilmente ajustadas quando o plano de trabalho mudar.
const (
	MetaProfissionaisCapacitados = 200
	MetaMateriaisProduzidos      = 10
)

type PainelRepository struct {
	db *sql.DB
}

func NewPainelRepository(db *sql.DB) *PainelRepository {
	return &PainelRepository{db: db}
}

// Indicadores é o payload de GET /painel/indicadores: tudo que o Dashboard
// e a tela de Relatórios precisam em uma única chamada.
type Indicadores struct {
	MunicipiosAtendidos       int `json:"municipios_atendidos"`
	MunicipiosMeta            int `json:"municipios_meta"`
	ProfissionaisCapacitados  int `json:"profissionais_capacitados"`
	ProfissionaisMeta         int `json:"profissionais_meta"`
	MateriaisProduzidos       int `json:"materiais_produzidos"`
	MateriaisMeta             int `json:"materiais_meta"`
	MacrorregioesCobertas     int `json:"macrorregioes_cobertas"`
	MacrorregioesTotal        int `json:"macrorregioes_total"`
	TotalUsuarios             int `json:"total_usuarios"`
	TotalDownloads            int `json:"total_downloads"`
	TotalAtividades           int `json:"total_atividades"`
	TotalInscritosAtividades  int `json:"total_inscritos_atividades"`

	EngajamentoMensal      []EngajamentoMes    `json:"engajamento_mensal"`
	AlcancePorMacrorregiao []AlcanceMacro      `json:"alcance_por_macrorregiao"`
}

// EngajamentoMes é um ponto do gráfico de linha de Relatórios.
type EngajamentoMes struct {
	Mes       string `json:"mes"` // "Jan", "Fev", ...
	Acessos   int    `json:"acessos"`
	Downloads int    `json:"downloads"`
}

// AlcanceMacro é uma barra do gráfico "Acessos por macrorregião". Os
// acessos são os logins autenticados feitos por usuários vinculados a
// municípios da macrorregião — o que é mensurável hoje (ainda não há
// analytics de tráfego público).
type AlcanceMacro struct {
	Macrorregiao string `json:"macrorregiao"`
	Acessos      int    `json:"acessos"`
}

var mesesPT = [...]string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}

func (r *PainelRepository) queryInt(query string, args ...any) (int, error) {
	var n int
	if err := r.db.QueryRow(query, args...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// Calcular monta todos os indicadores do Painel.
func (r *PainelRepository) Calcular(mesesHistorico int) (*Indicadores, error) {
	if mesesHistorico <= 0 {
		mesesHistorico = 6
	}
	ind := &Indicadores{
		ProfissionaisMeta: MetaProfissionaisCapacitados,
		MateriaisMeta:     MetaMateriaisProduzidos,
	}
	var err error

	if ind.MunicipiosAtendidos, err = r.queryInt("SELECT COUNT(id) FROM cidade WHERE status = 'ativo'"); err != nil {
		return nil, err
	}
	if ind.MunicipiosMeta, err = r.queryInt("SELECT COALESCE(SUM(meta_municipios), 0) FROM macrorregiao"); err != nil {
		return nil, err
	}
	if ind.MacrorregioesTotal, err = r.queryInt("SELECT COUNT(id) FROM macrorregiao"); err != nil {
		return nil, err
	}
	if ind.MacrorregioesCobertas, err = r.queryInt("SELECT COUNT(DISTINCT macrorregiao_id) FROM cidade WHERE status = 'ativo' AND macrorregiao_id IS NOT NULL"); err != nil {
		return nil, err
	}
	// Capacitado = concluiu um curso do AVA OU teve presença confirmada em
	// uma oficina/capacitação (contado uma única vez por pessoa).
	if ind.ProfissionaisCapacitados, err = r.queryInt(`
		SELECT COUNT(DISTINCT uid) FROM (
			SELECT usuario_id AS uid FROM matricula WHERE status = 'concluido'
			UNION
			SELECT usuario_id AS uid FROM atividade_inscricao WHERE presente = 1
		) t`); err != nil {
		return nil, err
	}
	if ind.MateriaisProduzidos, err = r.queryInt("SELECT COUNT(id) FROM material WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	if ind.TotalUsuarios, err = r.queryInt("SELECT COUNT(id) FROM usuario WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	if ind.TotalDownloads, err = r.queryInt("SELECT COALESCE(SUM(downloads), 0) FROM material WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	if ind.TotalAtividades, err = r.queryInt("SELECT COUNT(id) FROM atividade WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	if ind.TotalInscritosAtividades, err = r.queryInt("SELECT COUNT(id) FROM atividade_inscricao"); err != nil {
		return nil, err
	}

	if ind.EngajamentoMensal, err = r.engajamentoMensal(mesesHistorico); err != nil {
		return nil, err
	}
	if ind.AlcancePorMacrorregiao, err = r.alcancePorMacrorregiao(); err != nil {
		return nil, err
	}
	return ind, nil
}

func (r *PainelRepository) engajamentoMensal(meses int) ([]EngajamentoMes, error) {
	// Monta a série com todos os meses do período (mesmo os sem dados),
	// para o gráfico não "pular" meses.
	now := time.Now()
	primeiro := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -(meses - 1), 0)

	type key struct{ y, m int }
	acessos := map[key]int{}
	downloads := map[key]int{}

	rows, err := r.db.Query(`SELECT YEAR(created_at), MONTH(created_at), COUNT(id) FROM acesso_log WHERE created_at >= ? GROUP BY YEAR(created_at), MONTH(created_at)`, primeiro)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var y, m, c int
		if err := rows.Scan(&y, &m, &c); err != nil {
			rows.Close()
			return nil, err
		}
		acessos[key{y, m}] = c
	}
	rows.Close()

	rows, err = r.db.Query(`SELECT YEAR(created_at), MONTH(created_at), COUNT(id) FROM material_download_log WHERE created_at >= ? GROUP BY YEAR(created_at), MONTH(created_at)`, primeiro)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var y, m, c int
		if err := rows.Scan(&y, &m, &c); err != nil {
			rows.Close()
			return nil, err
		}
		downloads[key{y, m}] = c
	}
	rows.Close()

	out := make([]EngajamentoMes, 0, meses)
	for i := 0; i < meses; i++ {
		d := primeiro.AddDate(0, i, 0)
		k := key{d.Year(), int(d.Month())}
		out = append(out, EngajamentoMes{
			Mes:       mesesPT[d.Month()-1],
			Acessos:   acessos[k],
			Downloads: downloads[k],
		})
	}
	return out, nil
}

func (r *PainelRepository) alcancePorMacrorregiao() ([]AlcanceMacro, error) {
	rows, err := r.db.Query(`
		SELECT mr.nome, COUNT(al.id) AS acessos
		FROM macrorregiao mr
		LEFT JOIN cidade c ON c.macrorregiao_id = mr.id
		LEFT JOIN usuario u ON u.municipio_id = c.id
		LEFT JOIN acesso_log al ON al.usuario_id = u.id
		GROUP BY mr.id, mr.nome
		ORDER BY mr.nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AlcanceMacro
	for rows.Next() {
		var a AlcanceMacro
		if err := rows.Scan(&a.Macrorregiao, &a.Acessos); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
