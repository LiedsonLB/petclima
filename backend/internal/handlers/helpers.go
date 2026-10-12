package handlers

import "strconv"

// atoiDefault parses s as an int, returning fallback when it's empty or
// invalid. Shared by handlers that read pagination/filter query params.
func atoiDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

// parseInt64 converte s em int64 (helper para path params extras).
func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
