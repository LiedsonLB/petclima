import { useEffect, useRef, useState } from 'react';
import { Radio, Users, Mic, MicOff, Video as VideoIcon, VideoOff, PhoneOff, Loader2 } from 'lucide-react';
import Card from './Card';
import { listarSalas, entrarNaSala, type Sala } from '../lib/api';

// URL pública do servidor LiveKit (ver .env.example / VITE_LIVEKIT_URL).
const LIVEKIT_URL = import.meta.env.VITE_LIVEKIT_URL as string | undefined;

interface ParticipanteTile {
  id: string;
  nome: string;
  isLocal: boolean;
  videoTrack?: MediaStreamTrack;
  audioTrack?: MediaStreamTrack;
}

/**
 * Oficinas ao vivo (Educação Permanente / Comunicação) usando o LiveKit já
 * configurado no backend (ver internal/handlers/sala_handler.go — POST
 * /salas/{id}/entrar devolve o token de acesso à room).
 *
 * A biblioteca `livekit-client` é importada dinamicamente para não pesar o
 * bundle de quem nunca abre uma sala.
 */
export default function SalasAoVivo() {
  const [salas, setSalas] = useState<Sala[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [entrando, setEntrando] = useState<number | null>(null);
  const [salaAtiva, setSalaAtiva] = useState<Sala | null>(null);
  const [participantes, setParticipantes] = useState<ParticipanteTile[]>([]);
  const [micLigado, setMicLigado] = useState(true);
  const [camLigada, setCamLigada] = useState(true);
  const roomRef = useRef<import('livekit-client').Room | null>(null);

  useEffect(() => {
    let ativo = true;
    listarSalas()
      .then((data) => ativo && setSalas(data.filter((s) => s.ativa)))
      .catch(() => ativo && setErro('Não foi possível carregar as oficinas ao vivo.'))
      .finally(() => ativo && setCarregando(false));
    return () => {
      ativo = false;
    };
  }, []);

  async function handleEntrar(sala: Sala) {
    if (!LIVEKIT_URL) {
      setErro('VITE_LIVEKIT_URL não configurada — peça para o time técnico definir a URL do servidor LiveKit.');
      return;
    }
    setErro('');
    setEntrando(sala.id);
    try {
      const { token } = await entrarNaSala(sala.id);
      const { Room, RoomEvent, Track } = await import('livekit-client');
      const room = new Room();
      roomRef.current = room;

      const upsert = (id: string, patch: Partial<ParticipanteTile>) => {
        setParticipantes((prev) => {
          const existe = prev.find((p) => p.id === id);
          if (existe) return prev.map((p) => (p.id === id ? { ...p, ...patch } : p));
          return [...prev, { id, nome: id, isLocal: false, ...patch }];
        });
      };

      room.on(RoomEvent.TrackSubscribed, (track, _pub, participant) => {
        if (track.kind === Track.Kind.Video) {
          upsert(participant.identity, { nome: participant.name || participant.identity, videoTrack: track.mediaStreamTrack });
        } else if (track.kind === Track.Kind.Audio) {
          upsert(participant.identity, { audioTrack: track.mediaStreamTrack });
        }
      });
      room.on(RoomEvent.ParticipantDisconnected, (participant) => {
        setParticipantes((prev) => prev.filter((p) => p.id !== participant.identity));
      });

      await room.connect(LIVEKIT_URL, token);
      await room.localParticipant.setCameraEnabled(true);
      await room.localParticipant.setMicrophoneEnabled(true);
      setParticipantes([{ id: room.localParticipant.identity, nome: 'Você', isLocal: true }]);
      setSalaAtiva(sala);
    } catch {
      setErro('Não foi possível entrar na sala. Verifique sua câmera/microfone e tente novamente.');
    } finally {
      setEntrando(null);
    }
  }

  async function handleSair() {
    await roomRef.current?.disconnect();
    roomRef.current = null;
    setSalaAtiva(null);
    setParticipantes([]);
  }

  async function toggleMic() {
    const room = roomRef.current;
    if (!room) return;
    const proximo = !micLigado;
    await room.localParticipant.setMicrophoneEnabled(proximo);
    setMicLigado(proximo);
  }

  async function toggleCam() {
    const room = roomRef.current;
    if (!room) return;
    const proximo = !camLigada;
    await room.localParticipant.setCameraEnabled(proximo);
    setCamLigada(proximo);
  }

  if (salaAtiva) {
    return (
      <Card title={salaAtiva.nome} icon={<Radio size={16} />}>
        <div style={{
          display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 10, marginBottom: 14,
        }}>
          {participantes.map((p) => (
            <VideoTile key={p.id} participante={p} />
          ))}
        </div>
        <div style={{ display: 'flex', gap: 10, justifyContent: 'center' }}>
          <button onClick={toggleMic} className="btn-outline" style={{ padding: 10, borderRadius: 999 }} aria-label="Microfone">
            {micLigado ? <Mic size={16} /> : <MicOff size={16} />}
          </button>
          <button onClick={toggleCam} className="btn-outline" style={{ padding: 10, borderRadius: 999 }} aria-label="Câmera">
            {camLigada ? <VideoIcon size={16} /> : <VideoOff size={16} />}
          </button>
          <button
            onClick={handleSair}
            style={{ padding: '10px 18px', borderRadius: 999, background: 'var(--danger-bg)', color: 'var(--danger-text)', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 6, fontSize: 13, fontWeight: 600 }}
          >
            <PhoneOff size={16} /> Sair da oficina
          </button>
        </div>
      </Card>
    );
  }

  return (
    <Card title="Oficinas ao vivo" icon={<Radio size={16} />}>
      {erro && (
        <div style={{ fontSize: 12, color: 'var(--danger-text)', background: 'var(--danger-bg)', padding: '8px 12px', borderRadius: 10, marginBottom: 12 }}>
          {erro}
        </div>
      )}
      {carregando ? (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 13, color: 'var(--text-muted)' }}>
          <Loader2 size={14} className="animate-spin" /> Carregando oficinas...
        </div>
      ) : salas.length === 0 ? (
        <p style={{ fontSize: 13, color: 'var(--text-muted)' }}>Nenhuma oficina ao vivo no momento.</p>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
          {salas.map((sala) => (
            <div key={sala.id} style={{
              display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 12,
              padding: '12px 14px', borderRadius: 12, border: '1px solid var(--border)',
            }}>
              <div style={{ minWidth: 0 }}>
                <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>{sala.nome}</div>
                {sala.descricao && (
                  <div style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 2 }}>{sala.descricao}</div>
                )}
                <div style={{ display: 'flex', alignItems: 'center', gap: 4, marginTop: 4, fontSize: 11, color: 'var(--text-muted)' }}>
                  <Users size={12} /> {sala.participantes_online ?? 0} conectado(s)
                  {sala.categoria && <span> · {sala.categoria}</span>}
                </div>
              </div>
              <button
                onClick={() => handleEntrar(sala)}
                disabled={entrando === sala.id}
                className="btn-primary"
                style={{ fontSize: 12, padding: '8px 16px', flexShrink: 0 }}
              >
                {entrando === sala.id ? 'Entrando...' : 'Entrar'}
              </button>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}

function VideoTile({ participante }: { participante: ParticipanteTile }) {
  const ref = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (ref.current && participante.videoTrack) {
      ref.current.srcObject = new MediaStream([participante.videoTrack]);
    }
  }, [participante.videoTrack]);

  return (
    <div style={{ position: 'relative', borderRadius: 12, overflow: 'hidden', background: '#0f1513', aspectRatio: '4/3' }}>
      {participante.videoTrack ? (
        <video ref={ref} autoPlay playsInline muted={participante.isLocal} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
      ) : (
        <div style={{ width: '100%', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#79fabf', fontSize: 24, fontWeight: 700 }}>
          {participante.nome.charAt(0).toUpperCase()}
        </div>
      )}
      <span style={{ position: 'absolute', bottom: 6, left: 8, fontSize: 11, color: '#fff', background: 'rgba(0,0,0,0.5)', padding: '2px 8px', borderRadius: 999 }}>
        {participante.nome}
      </span>
    </div>
  );
}
