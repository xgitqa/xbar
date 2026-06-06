package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed .version
var version string

func main() {
	println("xbar", version)
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	println("xbar exited")
}

func run() error {
	app, err := newApp()
	if err != nil {
		return errors.Wrap(err, "newApp")
	}
	wailsLogLevel := logger.ERROR
	app.Verbose = true
	if app.Verbose {
		wailsLogLevel = logger.DEBUG
	}
	err = wails.Run(&options.App{
		Title:             "xbar",
		Width:             1080,
		Height:            700,
		MinWidth:          800,
		MinHeight:         600,
		StartHidden:       true,
		HideWindowOnClose: true,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			// Disabled: on macOS Sonoma/Sequoia, Wails alpha-74's
			// makeWindowBackgroundTranslucent() passes NSWindowBelow (-1)
			// to -[NSView addSubview:positioned:relativeTo:] through a 32-bit
			// `int` objc_msgSend cast. On arm64 that zero-extends to
			// 0xFFFFFFFF, tripping AppKit's new ordering-mode assertion and
			// aborting at launch. Leaving this off skips that native path.
			WindowBackgroundIsTranslucent: false,
			TitleBar:                      mac.TitleBarHiddenInset(),
			Menu:                          app.appMenu,
			ActivationPolicy:              mac.NSApplicationActivationPolicyAccessory,
			URLHandlers: map[string]func(string){
				// xbar://...
				"xbar": app.handleIncomingURL,
			},
		},
		ContextMenus: app.contextMenus,
		LogLevel:     wailsLogLevel,
		Startup:      app.Start,
		Shutdown:     app.Shutdown,
		Bind: []interface{}{
			app.PersonService,
			app.CategoriesService,
			app.PluginsService,
			app.CommandService,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
