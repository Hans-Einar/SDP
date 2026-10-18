package preparation_test

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"testing"
)

func TestHiddenExtendedInputCapabilityInventory(t *testing.T) {
	_, root, err := parser.Compile(`sdui 0.3;page=[dialog=dialog("Hidden")[edit=input("",multiline=true,readOnly=true)]];`)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"input-multiline", "read-only"} {
		caps := admission.TextCapabilities()
		var filtered preparation.Capabilities
		for _, cap := range caps {
			if cap.Dimension != preparation.Host || cap.ID != id {
				filtered = append(filtered, cap)
			}
		}
		err := preparation.Check("sdui/0.3", root["page"], filtered)
		var diagnostic *preparation.Diagnostic
		if !errors.As(err, &diagnostic) || diagnostic.Path != "page/dialog/edit" || diagnostic.Capability.ID != id {
			t.Fatal("hidden input lost capability diagnostic", id, err)
		}
	}
	if err := preparation.Check("sdui/0.3", root["page"], admission.TextCapabilities()); err != nil {
		t.Fatal(err)
	}
}
