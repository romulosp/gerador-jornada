// Package github implementa OAuth desktop e consultas GitHub sem coletar a senha do usuário.
package github

// ClientID identifica o cadastro público distribuído pelo build.
var ClientID string

// ClientSecret é o identificador recuperável do cliente nativo; não é credencial do usuário.
var ClientSecret string

// Configured informa se o build recebeu o par completo do cadastro do produto.
func Configured() bool { return ClientID != "" && ClientSecret != "" }
