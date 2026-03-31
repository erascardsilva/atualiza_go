// Atualiza GO - Atualizador de Sistema Linux
// Erasmo Cardoso - Dev

package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// OBRIGA O STARTUP GTK A CORTAR TODA RENDERIZAÇÃO WAYLAND/HARDWARE EGL
	// Essencial para sobreviver ao bloqueio do Webkit/Core22 ACPI e EGL.
	os.Setenv("WEBKIT_DISABLE_COMPOSITING_MODE", "1")
	os.Setenv("LIBGL_ALWAYS_SOFTWARE", "1")
	
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Atualiza GO",
		Width:     960,
		Height:    640,
		MinWidth:  800,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 25, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
		},
	})

	if err != nil {
		println("Erro:", err.Error())
	}
}
