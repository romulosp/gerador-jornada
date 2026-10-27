package mcp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"gerador-jornada.local/editor/internal/adapter/buildinfo"
	"gerador-jornada.local/editor/internal/adapter/workspace"
	"gerador-jornada.local/editor/internal/application"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config contém somente autorizações fornecidas na composição do processo pelo operador.
// Nenhuma ferramenta MCP recebe ou modifica esta configuração.
type Config struct {
	// Workspace é pasta existente que delimita jornadas regulares da raiz.
	Workspace string
	// LegacyRoots/ImportRoots separam finalidades de leitura autorizadas no startup.
	LegacyRoots, ImportRoots []string
}

// runtime reúne adaptadores e registro compartilhado, sem criar domínio MCP paralelo.
type runtime struct {
	// authoring/source coordenam rascunhos sobre a mesma fonte usada pelo editor manual.
	authoring *application.Authoring
	source    *application.AuthoringWorkspace
	// materials fixa autorizações por finalidade; pages autentica cursores deste processo.
	materials *workspace.MaterialScopes
	pages     *pagination
	// audit escreve apenas metadados em stderr; tools reflete ferramentas realmente registradas.
	audit *audit
	tools []string
}

// newRuntime valida dependências sem criar preferências, logs, Git ou pacotes no workspace.
func newRuntime(config Config, diagnostic io.Writer) (*runtime, error) {
	if config.Workspace == "" {
		return nil, fmt.Errorf("--workspace é obrigatório no modo MCP")
	}
	// source captura a raiz física e reutiliza codec/construtor do armazenamento manual.
	source, err := workspace.OpenWorkspace(config.Workspace)
	if err != nil {
		return nil, fmt.Errorf("workspace MCP indisponível")
	}
	materials, err := workspace.NewMaterialScopes(config.LegacyRoots, config.ImportRoots)
	if err != nil {
		return nil, fmt.Errorf("raiz de materiais MCP indisponível")
	}
	authoring := application.NewAuthoring(application.DefaultAuthoringLimits(), nil)
	return &runtime{authoring: authoring, source: application.NewAuthoringWorkspace(authoring, source), materials: materials, pages: newPagination(), audit: newAudit(diagnostic), tools: []string{}}, nil
}

// installedLimits descreve os limites operacionais sem serializar structs de domínio diretamente.
type installedLimits struct {
	// FrameBytes/JSONDepth limitam dados antes de decodificar ou emitir resposta.
	FrameBytes int `json:"frameBytes"`
	JSONDepth  int `json:"jsonDepth"`
	// Sessions/ConcurrentCalls/EditRequests/Proposals/Commits delimitam recursos de autoria.
	Sessions        int `json:"sessions"`
	ConcurrentCalls int `json:"concurrentCalls"`
	EditRequests    int `json:"editRequests"`
	Proposals       int `json:"proposals"`
	Commits         int `json:"commits"`
	// RetainedBytes/StagingBytes/TargetStagingBytes são reservas lógicas de finalização.
	RetainedBytes      int64 `json:"retainedBytes"`
	StagingBytes       int64 `json:"stagingBytes"`
	TargetStagingBytes int64 `json:"targetStagingBytes"`
	// FinalizationWorkBytes/FinalizationCalls expõem a reserva de um encode/verify e seu limite de concorrência.
	FinalizationWorkBytes              int64 `json:"finalizationWorkBytes"`
	FinalizationCalls                  int   `json:"finalizationCalls"`
	FinalizationReserveBytesPerSession int64 `json:"finalizationReserveBytesPerSession"`
	// PageItems/TextCodePoints/NameCodePoints têm unidades explícitas e não negociáveis.
	PageItems      int `json:"pageItems"`
	TextCodePoints int `json:"textCodePoints"`
	NameCodePoints int `json:"nameCodePoints"`
	// IdleSeconds/ReviewSeconds/CallSeconds/PackageSeconds são prazos separados da espera humana.
	IdleSeconds    int64 `json:"idleSeconds"`
	ReviewSeconds  int64 `json:"reviewSeconds"`
	CallSeconds    int   `json:"callSeconds"`
	PackageSeconds int   `json:"packageSeconds"`
}

// approvalCapability comunica indisponibilidade verificável, sem aceitar declaração do agente.
type approvalCapability struct {
	// Name identifica o contrato projetado; Qualified depende de homologação real.
	Name      string `json:"name"`
	Qualified bool   `json:"qualified"`
	// ReasonCode informa bloqueio sem supor produto/versão do cliente conectado.
	ReasonCode string `json:"reasonCode"`
}

// installedCapabilities descreve somente o software oferecido nesta composição.
type installedCapabilities struct {
	// ApplicationVersion/ProfileVersion/ProtocolVersion são versões independentes.
	ApplicationVersion string `json:"applicationVersion"`
	ProfileVersion     int    `json:"profileVersion"`
	ProtocolVersion    string `json:"protocolVersion"`
	// Formats expõe pares aceitos sem sugerir combinações inexistentes de versões.
	Formats        []formatVersion `json:"formats"`
	CatalogVersion int             `json:"catalogVersion"`
	Families       []catalogItem   `json:"families"`
	// Tools corresponde ao registro real e não ao backlog da feature.
	Tools []string `json:"tools"`
	// PersistentWrites nunca se torna verdadeiro por uma mensagem de aprovação do agente.
	PersistentWrites bool               `json:"persistentWrites"`
	ApprovalProfile  approvalCapability `json:"approvalProfile"`
	// StartupProbe falso diferencia a composição com núcleo da prova mínima G0 histórica.
	StartupProbe bool `json:"startupProbe"`
	// Stage/Limits tornam explícito o incremento instalado e seus orçamentos.
	Stage  string          `json:"stage"`
	Limits installedLimits `json:"limits"`
}

// formatVersion descreve um par histórico aceito pelo codec .jornada existente.
type formatVersion struct {
	// SchemaVersion/CatalogVersion são versões do formato, independentes do software/protocolo.
	SchemaVersion  int `json:"schemaVersion"`
	CatalogVersion int `json:"catalogVersion"`
}

// capabilities constrói uma fotografia das operações realmente registradas.
func (runtime *runtime) capabilities() installedCapabilities {
	// limits é política local fixa, não configuração escolhida pelo cliente.
	limits := application.DefaultAuthoringLimits()
	return installedCapabilities{ApplicationVersion: buildinfo.Version(), ProfileVersion: 1, ProtocolVersion: "2025-11-25", Tools: append([]string{}, runtime.tools...), Formats: []formatVersion{{1, 1}, {2, 1}, {3, 2}, {3, 3}, {3, 4}}, CatalogVersion: 4, Families: installedCatalog(),
		PersistentWrites: false, ApprovalProfile: approvalCapability{Name: "HostApproval-v1", Qualified: false, ReasonCode: "authorizationUnavailable"}, StartupProbe: false, Stage: "authoringDraftInProgress",
		Limits: installedLimits{FrameBytes: MaxFrameBytes, JSONDepth: maxJSONDepth, Sessions: limits.Sessions, ConcurrentCalls: limits.ConcurrentCalls, EditRequests: limits.EditRequests, Proposals: limits.Proposals, Commits: limits.Commits, RetainedBytes: limits.RetainedBytes, StagingBytes: limits.StagingBytes, TargetStagingBytes: limits.TargetStagingBytes, FinalizationWorkBytes: limits.FinalizationWorkBytes, FinalizationCalls: 1, FinalizationReserveBytesPerSession: limits.FinalizationReserveBytes, PageItems: 100, TextCodePoints: 65536, NameCodePoints: 256, IdleSeconds: int64(limits.Idle.Seconds()), ReviewSeconds: int64(limits.ReviewTTL.Seconds()), CallSeconds: 30, PackageSeconds: 120}}
}

// registerTool conserva descoberta e auditoria vinculadas ao registro instalado.
func registerTool[In any](server *sdk.Server, runtime *runtime, tool *sdk.Tool, handler func(context.Context, In) (*sdk.CallToolResult, error)) {
	runtime.tools = append(runtime.tools, tool.Name)
	addTypedTool(server, runtime.audit, tool, handler)
}

// newServer compõe SDK com transporte/erros limitados, sem logging remoto ou sampling.
func newServer(runtime *runtime) *sdk.Server {
	// server restringe negociação ao perfil testado e descarta logs internos da SDK.
	server := sdk.NewServer(&sdk.Implementation{Name: "JornadaStudio", Version: buildinfo.Version()}, &sdk.ServerOptions{SupportedProtocolVersions: []string{"2025-11-25"}, Capabilities: &sdk.ServerCapabilities{}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
	server.AddReceivingMiddleware(callBudget)
	registerTool(server, runtime, &sdk.Tool{Name: "aplicacao_capacidades", Description: "Consulta operações instaladas, limites e disponibilidade de consentimento.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, _ struct{}) (*sdk.CallToolResult, error) {
			return successResult(runtime.capabilities(), ""), nil
		})
	registerDiscovery(server, runtime)
	registerProjects(server, runtime)
	registerScreens(server, runtime)
	registerComponents(server, runtime)
	registerPresets(server, runtime)
	registerNavigation(server, runtime)
	registerValidation(server, runtime)
	registerSessions(server, runtime)
	return server
}

// Run serve uma conexão stdio no mesmo EXE, descartando apenas rascunhos ao terminar.
// input/output são pipes de protocolo; diagnostic é stderr separado. Não abre GUI nem rede.
func Run(ctx context.Context, input io.ReadCloser, output io.WriteCloser, diagnostic io.Writer, config Config) error {
	// runtime é criado antes de aceitar ferramentas, fixando todas as raízes operacionais.
	runtime, err := newRuntime(config, diagnostic)
	if err != nil {
		return err
	}
	defer runtime.authoring.CloseAll()
	// sweepContext/sweepDone encerram a varredura antes de liberar rascunhos no EOF.
	sweepContext, stopSweep := context.WithCancel(ctx)
	sweepDone := make(chan struct{})
	go sweepIdle(sweepContext, runtime.authoring, 30*time.Second, sweepDone)
	defer func() { stopSweep(); <-sweepDone }()
	if err := newServer(runtime).Run(ctx, &sdk.IOTransport{Reader: newFrameReader(input), Writer: &frameWriter{destination: output}, MaxLineLength: MaxFrameBytes}); err != nil {
		return fmt.Errorf("transporte MCP encerrado com falha")
	}
	return nil
}
