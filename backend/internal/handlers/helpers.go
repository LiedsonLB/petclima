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

// cpfValido confere o CPF pelos dois dígitos verificadores (aceita com ou sem
// máscara). Rejeita sequências repetidas como 111.111.111-11.
func cpfValido(cpf string) bool {
	d := onlyDigits(cpf)
	if len(d) != 11 {
		return false
	}
	repetido := true
	for i := 1; i < 11; i++ {
		if d[i] != d[0] {
			repetido = false
			break
		}
	}
	if repetido {
		return false
	}
	for n := 9; n <= 10; n++ {
		soma := 0
		for i := 0; i < n; i++ {
			soma += int(d[i]-'0') * (n + 1 - i)
		}
		dv := (soma * 10) % 11
		if dv == 10 {
			dv = 0
		}
		if dv != int(d[n]-'0') {
			return false
		}
	}
	return true
}
