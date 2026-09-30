import { createContext, useContext, useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

type Theme = 'light' | 'dark';
interface ThemeCtx { theme: Theme; toggle: () => void; }

const ThemeContext = createContext<ThemeCtx>({ theme: 'light', toggle: () => {} });

// Páginas públicas: sempre tema claro (sem dark mode).
const PUBLIC_PATHS = ['/', '/login', '/cadastro', '/esqueci-senha', '/redefinir-senha', '/confirmar-email'];

// Deve ser usado DENTRO do <BrowserRouter>.
export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const { pathname } = useLocation();
  const isPublic = PUBLIC_PATHS.includes(pathname.replace(/\/+$/, '') || '/');

  // Preferência do usuário: só vale para a área logada e só muda por clique
  // (não segue mais o prefers-color-scheme do sistema).
  const [theme, setTheme] = useState<Theme>(() =>
    localStorage.getItem('petsaudeclima-theme') === 'dark' ? 'dark' : 'light'
  );

  const effective: Theme = isPublic ? 'light' : theme;

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', effective);
    document.documentElement.style.colorScheme = effective;
  }, [effective]);

  useEffect(() => {
    localStorage.setItem('petsaudeclima-theme', theme);
  }, [theme]);

  const toggle = () => setTheme(t => (t === 'light' ? 'dark' : 'light'));

  return <ThemeContext.Provider value={{ theme: effective, toggle }}>{children}</ThemeContext.Provider>;
}

export const useTheme = () => useContext(ThemeContext);
