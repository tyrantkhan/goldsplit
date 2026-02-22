package main

import (
	"embed"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed VERSION
var rawVersion string

var versionSuffix = "-dev" // overridden to "" via ldflags in release builds

var version = strings.TrimSpace(rawVersion) + versionSuffix

func main() {
	app := NewApp(version)

	err := wails.Run(&options.App{
		Title:     "Goldsplit",
		Width:     400,
		Height:    650,
		MinWidth:  400,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 24, A: 1},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Goldsplit",
				Message: "Version " + version,
			},
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
