package fynehost

import (
	"errors"
	"reflect"
	goruntime "runtime"
	"strconv"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

const admissionCommandBody = `run=command("Run",key="F6");invoke=button(command="run");edit=input("Draft",value="saved");`
const admissionIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16"><rect width="16" height="16" fill="red"/></svg>`

func admissionCommandRequest(t *testing.T, sequence uint64, extra string, calls *int) DocumentRequest {
	t.Helper()
	doc, err := parser.Parse(`sdui 0.3;page=[` + admissionCommandBody + extra + `] {x=fill,y=fill};`)
	lifecycleOK(t, err)
	return DocumentRequest{Document: doc, Entry: "page", SessionID: "admission-m2", SourceRevision: "candidate-" + strconv.FormatUint(sequence, 10), Sequence: sequence, Mode: preparation.Connected,
		Bind: func(s *ui.Session, _ *parser.Document) error {
			handle, ok := s.Command("page/run")
			if !ok {
				return errors.New("missing test command")
			}
			return s.BindInteraction(handle, func(ui.Event) (ui.InteractionReply, error) {
				*calls++
				return ui.InteractionReply{Domain: ui.DomainSucceeded}, nil
			})
		},
	}
}

// The public software canvas exposes TypedShortcut; test.Window's wrapper hides
// that method. Use the real canvas registration table with the actual host.
func admissionCommandHost(t *testing.T) (*DocumentHost, software.WindowlessCanvas) {
	t.Helper()
	test.NewTempApp(t)
	c := software.NewCanvas()
	c.SetPadded(false)
	c.Resize(fyne.NewSize(800, 500))
	h := NewDocumentHost(c, layout.Size{W: 800, H: 500})
	t.Cleanup(h.Close)
	return h, c
}

func TestCommandAdmissionFailedReloadPreservesNativeShortcut(t *testing.T) {
	valid := func() fyne.Resource { return fyne.NewStaticResource("valid.svg", []byte(admissionIconSVG)) }
	hidden := `closed=dialog("Closed")[hidden=command("Hidden",toggle=true,icon="missing") {visible=false}];`
	cases := []struct {
		name, extra, diagnostic string
		icons                   map[string]fyne.Resource
		guard                   func() error
		nativeAlias             bool
	}{
		{name: "hidden-icon-provider-capability", extra: hidden, diagnostic: "provider icon/1"},
		{name: "hidden-icon-unresolved", extra: hidden, diagnostic: `icon-resource: missing "missing"`, icons: map[string]fyne.Resource{"other": valid()}},
		{name: "hidden-icon-invalid", extra: hidden, diagnostic: "icon-resource:", icons: map[string]fyne.Resource{"missing": fyne.NewStaticResource("broken.svg", []byte("not SVG"))}},
		{name: "hidden-ordinary-command-unbound", extra: `absent=command("Missing") {visible=false};`, diagnostic: "unbound-command: page/absent"},
		{name: "hidden-native-key-alias", extra: `a=command("A",toggle=true,key="Primary+K");b=command("B",toggle=true,key="Ctrl+K") {visible=false};`, diagnostic: "command-key: duplicate", nativeAlias: true},
		{name: "late-guard", diagnostic: "test candidate guard", guard: func() error { return errors.New("test candidate guard") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.nativeAlias && goruntime.GOOS == "darwin" {
				t.Skip("Primary and Ctrl are distinct native modifiers on Darwin")
			}
			h, c := admissionCommandHost(t)
			domainCalls := 0
			lifecycleOK(t, h.Adopt(admissionCommandRequest(t, 1, "", &domainCalls)))
			c.SetContent(h.Container)
			c.Resize(fyne.NewSize(800, 500))
			old := h.Current()
			edit := old.Controls()["page/edit"].(*Input)
			edit.SetText("preserved draft")
			// Exercise the installed canvas registration, not Bundle.invokeCommand or
			// the adapter's direct shortcut forwarding path.
			canvas, ok := c.(fyne.Shortcutable)
			if !ok {
				t.Fatal("test canvas cannot deliver a registered native shortcut")
			}
			shortcut, err := nativeShortcut("F6")
			lifecycleOK(t, err)
			if domainCalls != 0 {
				t.Fatal("admission executed domain action", domainCalls)
			}
			canvas.TypedShortcut(shortcut)
			if domainCalls != 1 {
				t.Fatal("baseline native shortcut unavailable", domainCalls)
			}
			before, presentation := old.Session.Snapshot(), old.presentation
			request := admissionCommandRequest(t, 2, tc.extra, &domainCalls)
			request.Icons, request.Guard = tc.icons, tc.guard
			candidate, err := h.Prepare(request)
			if candidate != nil {
				candidate.Close()
				t.Fatal("rejected reload returned candidate")
			}
			if err == nil || !strings.Contains(err.Error(), tc.diagnostic) {
				t.Fatalf("wrong rejection: %v; want %q", err, tc.diagnostic)
			}
			if domainCalls != 1 {
				t.Fatal("failed preparation invoked domain", domainCalls)
			}
			if h.Current() != old || old.closed || old.Session.Closed() || old.presentation != presentation || old.Controls()["page/edit"] != edit || edit.Text != "preserved draft" || !reflect.DeepEqual(before, old.Session.Snapshot()) {
				t.Fatal("failed reload changed published bundle/draft/snapshot")
			}
			canvas.TypedShortcut(shortcut)
			if domainCalls != 2 {
				t.Fatal("failed reload removed or duplicated the live canvas shortcut", domainCalls)
			}
		})
	}
}

func TestCommandAdmissionConnectedHiddenIconSuccessIsInert(t *testing.T) {
	h, c := admissionCommandHost(t)
	calls := 0
	extra := `closed=dialog("Closed")[hidden=command("Hidden",toggle=true,icon="provided") {visible=false}];`
	request := admissionCommandRequest(t, 1, extra, &calls)
	data := []byte(admissionIconSVG)
	request.Icons = map[string]fyne.Resource{"provided": fyne.NewStaticResource("provided.svg", data)}
	lifecycleOK(t, h.Adopt(request))
	c.SetContent(h.Container)
	c.Resize(fyne.NewSize(800, 500))
	b := h.Current()
	if calls != 0 || len(b.surfaces) != 0 || b.icons["provided"] == nil {
		t.Fatal("valid hidden resource triggered work or was not prepared", calls)
	}
	data[0] = 'x'
	if string(b.icons["provided"].Content()) != admissionIconSVG {
		t.Fatal("prepared hidden icon aliases caller bytes")
	}
	shortcut, err := nativeShortcut("F6")
	lifecycleOK(t, err)
	c.(fyne.Shortcutable).TypedShortcut(shortcut)
	if calls != 1 {
		t.Fatal("connected command not bound after successful admission", calls)
	}
}
