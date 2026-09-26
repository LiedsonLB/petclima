import { MapContainer, TileLayer, CircleMarker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import Card from '../components/Card';
import { municipiosMapa, macrorregioes } from '../data/Data';

const statusColor = { ativo: '#1baf7a', planejado: '#9aaba7' } as const;

export default function Mapa() {
  const ativos = municipiosMapa.filter(m => m.status === 'ativo').length;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 14 }}>
        {macrorregioes.map(m => (
          <Card key={m.nome} style={{ padding: 0 }}>
            <div style={{ padding: 14 }}>
              <div style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)' }}>{m.nome}</div>
              <div style={{ fontSize: 20, fontWeight: 700, color: 'var(--primary)', marginTop: 4 }}>{m.municipiosAtendidos}/{m.meta}</div>
              <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>municípios atendidos</div>
            </div>
          </Card>
        ))}
      </div>

      <Card title={`Municípios no Piauí — ${ativos} ativos de 12 (meta final)`} noPad>
        <div style={{ height: 600, position: 'relative' }}>
          <MapContainer center={[-6.5, -42.8]} zoom={7} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false}>
            <TileLayer
              attribution='&copy; OpenStreetMap contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {municipiosMapa.map(m => (
              <CircleMarker
                key={m.name}
                center={[m.lat, m.lng]}
                radius={10}
                pathOptions={{ color: statusColor[m.status], fillColor: statusColor[m.status], fillOpacity: 0.75, weight: 2 }}
              >
                <Popup>
                  <strong>{m.name}</strong><br />
                  {m.macrorregiao}<br />
                  Status: {m.status === 'ativo' ? 'Ações em andamento' : 'Planejado'}
                </Popup>
              </CircleMarker>
            ))}
          </MapContainer>
        </div>
      </Card>

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
