package startup

// Preference governa a exibição por usuário; não pertence ao documento de autoria.
// O valor zero exibe a primeira apresentação e mantém a repetição desabilitada.
type Preference struct {
	// Presented registra uma apresentação concluída, sem consumir tentativas canceladas.
	Presented bool `json:"presented"`
	// ShowOnStartup é a escolha explícita de repetir a apresentação nas próximas aberturas.
	ShowOnStartup bool `json:"show_on_startup"`
}

// ShouldShow apresenta uma vez por padrão e depois respeita a escolha do usuário.
func (p Preference) ShouldShow() bool { return !p.Presented || p.ShowOnStartup }
