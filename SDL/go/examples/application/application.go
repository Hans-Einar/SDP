// Package application composes the simulated EditAptCell prototype. The domain
// implementation is handwritten Go; both languages only declare their contracts.
package application

import (
	"context"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/simulation"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type App struct {
	UI       *ui.Session
	Document *parser.Document
	SDL      *sdl.Engine
	Domain   *simulation.Apt
	Bindings *bridge.Bridge
}

func Load(actions, source string) (*App, error) {
	p, err := sdlparser.CompileActions(actions)
	if err != nil {
		return nil, err
	}
	doc, roots, err := parser.Compile(source)
	if err != nil {
		return nil, err
	}
	return New(p, doc, roots["page"])
}
func New(program *sdlparser.Program, doc *parser.Document, root *parser.Instance) (*App, error) {
	domain := simulation.NewApt()
	engine, err := sdl.New(program, domain.Registry())
	if err != nil {
		return nil, err
	}
	session, err := ui.New("apt-prototype", root)
	if err != nil {
		return nil, err
	}
	bindings, err := bridge.Bind(context.Background(), session, doc, map[string]*sdl.Engine{"apt": engine}, simulation.Plans(), simulation.InitialContext())
	if err != nil {
		return nil, err
	}
	return &App{session, doc, engine, domain, bindings}, nil
}
func (a *App) Close() { a.UI.Close(); a.SDL.Close() }
