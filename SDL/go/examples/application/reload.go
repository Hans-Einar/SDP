package application

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/simulation"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

// These composition methods run on the UI owner goroutine. Preview registration
// never invokes the Go domain handler. A failed preview preserves both models.
func (a *App) ValidateUI(doc *parser.Document, root *parser.Instance) error {
	temporary, err := ui.New("binding-preview", root)
	if err != nil {
		return err
	}
	defer temporary.Close()
	_, err = bridge.Bind(context.Background(), temporary, doc, map[string]*sdl.Engine{"apt": a.SDL}, simulation.Plans(), a.Bindings.Context)
	return err
}
func (a *App) ReloadUI(doc *parser.Document, root *parser.Instance) error {
	if err := a.ValidateUI(doc, root); err != nil {
		return err
	}
	if err := a.UI.Reload(root); err != nil {
		return err
	}
	return a.RebindUI(doc)
}
func (a *App) RebindUI(doc *parser.Document) error {
	if err := a.Bindings.Rebind(doc); err != nil {
		return err
	}
	a.Document = doc
	return nil
}
func (a *App) ReloadSDL(program *sdlparser.Program) error {
	if a.UI.Closed() {
		return fmt.Errorf("closed: application")
	}
	preview, err := a.SDL.Preview(program)
	if err != nil {
		return err
	}
	defer preview.Close()
	temporary, err := ui.New("binding-preview", a.UI.SnapshotRoot())
	if err != nil {
		return err
	}
	defer temporary.Close()
	if _, err = bridge.Bind(context.Background(), temporary, a.Document, map[string]*sdl.Engine{"apt": preview}, simulation.Plans(), a.Bindings.Context); err != nil {
		return err
	}
	if err = a.SDL.Reload(program); err != nil {
		return err
	}
	if err = a.UI.InvalidateEvents(); err != nil {
		return err
	}
	return a.Bindings.Rebind(a.Document)
}
