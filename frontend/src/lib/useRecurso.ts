import { useEffect, useState } from 'react';
import { mensagemDeErro } from './api';

interface Estado<T> {
  data: T | null;
  erro: string;
  carregando: boolean;
}

/**
 * Busca um recurso da API ao montar (e sempre que os argumentos mudarem).
 * `recarregar()` refaz a busca — use depois de uma ação que altera os dados.
 *
 * Passe uma função estável (as de lib/api.ts) e os argumentos dela:
 *   const { data, carregando, erro, recarregar } = useRecurso(listarMateriais, { tipo: 'card' });
 */
export function useRecurso<A extends unknown[], T>(fn: (...args: A) => Promise<T>, ...args: A) {
  const chave = JSON.stringify(args);
  const [versao, setVersao] = useState(0);
  const [estado, setEstado] = useState<Estado<T>>({ data: null, erro: '', carregando: true });

  useEffect(() => {
    let ativo = true;
    fn(...(JSON.parse(chave) as A))
      .then((data) => ativo && setEstado({ data, erro: '', carregando: false }))
      .catch((e) => ativo && setEstado((s) => ({ ...s, erro: mensagemDeErro(e), carregando: false })));
    return () => {
      ativo = false;
    };
  }, [fn, chave, versao]);

  return { ...estado, recarregar: () => setVersao((v) => v + 1) };
}
