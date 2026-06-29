package main

import (
	"embed"
	"log"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/wailsapp/wails/v3/pkg/application"

	"ydsterm/internal/boot"
	"ydsterm/internal/database"
	"ydsterm/internal/service"
	_ "ydsterm/packed"
)

//go:embed all:frontend/dist
var assets embed.FS

// Register terminal event payload types for Wails event bridge.
// Events MUST be registered before they can be emitted.
func init() {
	application.RegisterEvent[service.TermOutputPayload]("term-output")
	application.RegisterEvent[service.TermDisconnectedPayload]("term-disconnected")
	application.RegisterEvent[service.SftpTransferPayload]("sftp-transfer-progress")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("YDSterm starting...")

	dbPath := boot.DefaultDBPath()
	log.Printf("DB path: %s", dbPath)

	if err := boot.EnsureDataDir(dbPath); err != nil {
		log.Fatalf("create data dir: %v", err)
	}

	if err := database.Init(dbPath); err != nil {
		log.Fatalf("init database: %v", err)
	}
	log.Println("Database initialized")

	if err := boot.RunMigrations(); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	appInstance := boot.NewApp()
	if err := appInstance.UserService.InitActiveUser(); err != nil {
		log.Printf("init active user: %v", err)
	}
	log.Println("Services registered (Host, Snippet, PortForward, Terminal, Settings, User, Sync)")

	app := application.New(application.Options{
		Name:        "YDSterm",
		Description: "A cross-platform SSH terminal client",
		Services: []application.Service{
			application.NewService(appInstance.HostService),
			application.NewService(appInstance.SnippetService),
			application.NewService(appInstance.PortForwardService),
			application.NewService(appInstance.TerminalService),
			application.NewService(appInstance.SFTPService),
			application.NewService(appInstance.SettingsService),
			application.NewService(appInstance.UserService),
			application.NewService(appInstance.SyncService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "YDSterm",
		Width:     1280,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		URL:       "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
