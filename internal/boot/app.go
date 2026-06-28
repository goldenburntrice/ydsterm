package boot

import (
	"ydsterm/internal/service"
)

type App struct {
	HostService        *service.HostServiceImpl
	SnippetService     *service.SnippetServiceImpl
	PortForwardService *service.PortForwardServiceImpl
	TerminalService    *service.TerminalServiceImpl
	SFTPService        *service.SFTPServiceImpl
}

func NewApp() *App {
	return &App{
		HostService:        service.NewHostService(),
		SnippetService:     service.NewSnippetService(),
		PortForwardService: service.NewPortForwardService(),
		TerminalService:    service.NewTerminalService(),
		SFTPService:        service.NewSFTPService(),
	}
}
