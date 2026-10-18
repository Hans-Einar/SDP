package fynehost

import (
	"errors"
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestPreparedPreviewNativeOwnerLoss(t *testing.T) {
	for _, modal := range []bool{true, false} {
		for _, reason := range []string{"parent-hidden", "parent-closed"} {
			t.Run(fmt.Sprintf("modal=%v/%s", modal, reason), func(t *testing.T) {
				h, window := lifecycleHost(t)
				source := fmt.Sprintf(`sdui 0.3;ref: art "application-art";page=[open=button("Open",effect="open",target="detail");main=svg(art.Chart.@resource,description="Main artwork",fallback="reject");detail=dialog("Preview detail",modal=%v)[image=svg(art.Chart.@resource,label="Detail",description="Detail artwork",fallback="reject") {x=fill,y=fill};prose=markdown("# Detail prose",description="Detail summary",fallback="label")] {scale-x=0.7,scale-y=0.7}] {x=fill,y=fill};`, modal)
				r := previewRequest(t, source, 1)
				bindPreview(&r, "page/main", []byte(previewSVG))
				bindPreview(&r, "page/detail/image", []byte(previewSVG))
				if err := h.Adopt(r); err != nil {
					t.Fatal(err)
				}
				window.SetContent(h.Container)
				window.Resize(fyne.NewSize(800, 500))
				window.Show()
				b := h.Current()
				gateCalls, resourceCalls := 0, 0
				reject := false
				if err := h.Mutate(func(s *ui.Session) error {
					return s.CheckPresentationWith(func(snapshot ui.Snapshot) (ui.PresentationState, error) {
						gateCalls++
						if reject {
							return ui.PresentationState{}, errors.New("gate unavailable")
						}
						return b.stage(snapshot)
					})
				}); err != nil {
					t.Fatal(err)
				}
				b.request.PrepareResources = func(ui.Snapshot) error {
					resourceCalls++
					if reject {
						return errors.New("resources unavailable")
					}
					return nil
				}
				open := b.Controls()["page/open"].(*commandButton)
				open.Tapped(&fyne.PointEvent{})
				old := b.surfaces["page/detail"]
				if old == nil {
					t.Fatal("detail did not open")
				}
				image := b.Controls()["page/detail/image"].(*previewControl)
				prose := b.Controls()["page/detail/prose"].(*previewControl)
				main := b.Controls()["page/main"].(*previewControl)
				mainResource := main.image.Resource
				if image.image.Resource == nil || image.text.Resource == nil {
					t.Fatal("test requires mounted native images")
				}
				var receipts []ui.DialogResult
				h.OnDialogResult = func(r ui.DialogResult) { receipts = append(receipts, r) }
				reject = true
				gateCalls, resourceCalls = 0, 0
				var err error
				if reason == "parent-hidden" {
					err = h.NativeParentHidden()
				} else {
					err = h.NativeParentClosed()
				}
				if err != nil || gateCalls != 0 || resourceCalls != 0 {
					t.Fatal("owner loss was gated", err, gateCalls, resourceCalls)
				}
				if b.Session.Snapshot().Surfaces["page/detail"].Open || len(receipts) != 1 || receipts[0].Reason != reason || receipts[0].Sequence != 0 {
					t.Fatal("owner revocation/receipt", receipts)
				}
				for _, path := range []string{"page/detail/image", "page/detail/prose"} {
					info := b.inspectPreviews()[path]
					if info["mounted"] != false || info["visible"] != false {
						t.Errorf("retired preview still reported mounted: %s %+v", path, info)
					}
					c := b.view.controls[path]
					preview := c.widget.(*previewControl)
					if c.clip.Visible() || preview.frame.box != (layout.Rect{}) || preview.image.Resource != nil || preview.image.Image != nil || preview.text.Resource != nil || preview.text.Image != nil {
						t.Errorf("retired preview retained native content: %s", path)
					}
					if _, exists := b.presentation.previews[path]; exists {
						t.Errorf("retired preview geometry retained: %s", path)
					}
				}
				if main.image.Resource != mainResource || main.frame.box == (layout.Rect{}) {
					t.Fatal("owner cleanup affected main preview")
				}
				if t.Failed() {
					return
				}
				if _, exists := b.presentation.canvases["page/detail"]; exists {
					t.Fatal("retired native frame retained")
				}
				if old.image.Resource != nil || old.image.Image != nil || len(old.content.Objects) != 0 {
					t.Fatal("retired surface retained background pixels or mounted objects")
				}
				if b.previewResources["page/detail/image"] == nil {
					t.Fatal("owner loss discarded bundle declaration")
				}
				reject = false
				open.Tapped(&fyne.PointEvent{})
				next := b.surfaces["page/detail"]
				if next == nil || next.target == old.target || b.Controls()["page/detail/image"] != image || b.Controls()["page/detail/prose"] != prose {
					t.Fatal("reopen lost fresh target/control identity")
				}
				b.retireSurfacePreviews(old) // independent stale cleanup after replacement apply
				old.nativeClosed()
				old.nativeCloseRequested()
				if b.surfaces["page/detail"] != next || image.image.Resource == nil || len(receipts) != 1 {
					t.Fatal("obsolete native callback disposed replacement")
				}
				for _, path := range []string{"page/detail/image", "page/detail/prose"} {
					info := b.inspectPreviews()[path]
					if info["mounted"] != true || info["visible"] != true {
						t.Fatal("reopened preview absent", path, info)
					}
				}
			})
		}
	}
}
