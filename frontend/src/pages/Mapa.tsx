import { MapContainer, TileLayer, CircleMarker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import Card from '../components/Card';
import { listarMacrorregioes, listarMunicipios } from '../lib/api';
import { useRecurso } from '../lib/useRecurso';

const statusColor = { ativo: '#1baf7a', planejado: '#9aaba7' } as const;
const corDe = (status: string) => (status === 'ativo' ? statusColor.ativo : statusColor.planejado);

export default function Mapa() {
  const macrorregioes = useRecurso(listarMacrorregioes).data ?? [];
  const { data, erro } = useRecurso(listarMunicipios);
  const municipios = (data ?? []).filter((m) => m.latitude != null && m.longitude != null);
  const ativos = (data ?? []).filter((m) => m.status === 'ativo').length;
  const meta = macrorregioes.reduce((acc, m) => acc + m.meta_municipios, 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 14 }}>
        {macrorregioes.map(m => (
          <Card key={m.id} style={{ padding: 0 }}>
            <div style={{ padding: 14 }}>
              <div style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)' }}>{m.nome}</div>
              <div style={{ fontSize: 20, fontWeight: 700, color: 'var(--primary)', marginTop: 4 }}>{m.municipios_ativos ?? 0}/{m.meta_municipios}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>municípios atendidos</div>
            </div>
          </Card>
        ))}
      </div>

      <Card title={`Municípios no Piauí — ${ativos} ativos de ${meta || 12} (meta final)`} noPad>
        <div style={{ height: 600, position: 'relative' }}>
          <MapContainer center={[-6.5, -42.8]} zoom={7} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false}>
            <TileLayer
              attribution='&copy; OpenStreetMap contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {municipios.map(m => (
              <CircleMarker
                key={m.id}
                center={[m.latitude as number, m.longitude as number]}
                radius={10}
                pathOptions={{ color: corDe(m.status), fillColor: corDe(m.status), fillOpacity: 0.75, weight: 2 }}
              >
                <Popup>
                  <strong>{m.nome}</strong><br />
                  {m.macrorregiao_nome ?? '—'}<br />
                  Status: {m.status === 'ativo' ? 'Ações em andamento' : 'Planejado'}
                </Popup>
              </CircleMarker>
            ))}
          </MapContainer>
        </div>
      </Card>

      {erro && <div style={{ fontSize: 12, color: 'var(--danger-text)' }}>{erro}</div>}

      <div style={{ display: 'flex', gap: 20, fontSize: 12, color: 'var(--text-secondary)' }}>
        <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <span style={{ width: 10, height: 10, borderRadius: '50%', background: statusColor.ativo, display: 'inline-block' }} /> Ações em andamento
        </span>
        <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <span style={{ width: 10, height: 10, borderRadius: '50%', background: statusColor.planejado, display: 'inline-block' }} /> Planejado
        </span>
      </div>
    </div>
  );
}
