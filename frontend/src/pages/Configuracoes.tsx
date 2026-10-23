import { useState } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { useAuth } from '../contexts/AuthContext';
import Card from '../components/Card';
import {
  PAPEL_LABEL, assetUrl, atualizarPerfil, enviarFoto, listarMunicipios, mensagemDeErro,
} from '../lib/api';
import { useRecurso } from '../lib/useRecurso';
import { alertaErro, alertaSucesso, botaoPrimario, campo, rotulo } from '../lib/ui';

export default function Configuracoes() {
  const { theme, toggle } = useTheme();
  const { usuario, atualizarUsuario } = useAuth();
  const municipios = useRecurso(listarMunicipios);
  const [f, setF] = useState({
    nome: usuario?.nome ?? '',
    telefone: usuario?.telefone ?? '',
    instituicao: usuario?.instituicao ?? '',
    curso: usuario?.curso_formacao ?? '',
    descricao: usuario?.descricao ?? '',
    municipio: usuario?.municipio_id ? String(usuario.municipio_id) : '',
  });
  const [foto, setFoto] = useState<File | null>(null);
  const [erro, setErro] = useState('');
  const [ok, setOk] = useState(false);
  const [salvando, setSalvando] = useState(false);

  async function salvar(e: React.FormEvent) {
    e.preventDefault();
    if (!usuario) return;
    setErro(''); setOk(false); setSalvando(true);
    try {
      const novaFoto = foto ? (await enviarFoto(foto, usuario.foto)).url : undefined;
      const atualizado = await atualizarPerfil(usuario.id, {
        nome: f.nome, email: usuario.email, telefone: f.telefone, instituicao: f.instituicao,
        curso_formacao: f.curso, descricao: f.descricao,
        municipio_id: f.municipio ? Number(f.municipio) : null,
        foto: novaFoto,
      });
      atualizarUsuario(atualizado);
      setFoto(null);
      setOk(true);
    } catch (err) {
      setErro(mensagemDeErro(err));
    } finally {
      setSalvando(false);
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20, maxWidth: 560 }}>
      {usuario && (
        <Card title="Meu perfil">
          <form onSubmit={salvar} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              {usuario.foto && <img src={assetUrl(usuario.foto)} alt="" style={{ width: 48, height: 48, borderRadius: '50%', objectFit: 'cover' }} />}
              <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>
                {usuario.email} · <strong style={{ color: 'var(--text-secondary)' }}>{PAPEL_LABEL[usuario.papel]}</strong>
              </div>
            </div>
            <div><label style={rotulo}>Nome</label><input required style={campo} value={f.nome} onChange={(e) => setF({ ...f, nome: e.target.value })} /></div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              <div><label style={rotulo}>Telefone</label><input style={campo} value={f.telefone} onChange={(e) => setF({ ...f, telefone: e.target.value })} /></div>
              <div>
                <label style={rotulo}>Município</label>
                <select style={campo} value={f.municipio} onChange={(e) => setF({ ...f, municipio: e.target.value })}>
                  <option value="">—</option>
                  {(municipios.data ?? []).map((m) => <option key={m.id} value={m.id}>{m.nome}/{m.uf}</option>)}
                </select>
              </div>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              <div><label style={rotulo}>Instituição</label><input style={campo} value={f.instituicao} onChange={(e) => setF({ ...f, instituicao: e.target.value })} /></div>
              <div><label style={rotulo}>Curso / vínculo</label><input style={campo} value={f.curso} onChange={(e) => setF({ ...f, curso: e.target.value })} /></div>
            </div>
            <div><label style={rotulo}>Sobre mim</label><textarea rows={3} style={campo} value={f.descricao} onChange={(e) => setF({ ...f, descricao: e.target.value })} /></div>
            <div><label style={rotulo}>Foto (jpg, png, webp — até 5MB)</label><input type="file" accept="image/*" onChange={(e) => setFoto(e.target.files?.[0] ?? null)} /></div>
            {erro && <div style={alertaErro}>{erro}</div>}
            {ok && <div style={alertaSucesso}>Perfil atualizado.</div>}
            <button type="submit" disabled={salvando} style={{ ...botaoPrimario, justifyContent: 'center' }}>{salvando ? 'Salvando…' : 'Salvar alterações'}</button>
          </form>
        </Card>
      )}

      <Card title="Aparência">
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <div style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-primary)' }}>Tema</div>
            <div style={{ fontSize: 12, color: 'var(--text-muted)' }}>Alternar entre modo claro e escuro</div>
          </div>
          <button onClick={toggle} style={botaoPrimario}>{theme === 'light' ? 'Ativar escuro' : 'Ativar claro'}</button>
        </div>
      </Card>

      <Card title="Sobre o PET-Saúde Clima">
        <p style={{ fontSize: 13, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          Grupo de Atuação Tutorial responsável por Comunicação, Tecnologias Digitais e Educação
          Permanente dentro do PET-Saúde Clima Piauí, com foco em produção de tecnologias sociais e
          digitais, enfrentamento da desinformação e formação contínua de profissionais de saúde.
        </p>
      </Card>
    </div>
  );
}
