package sdptool

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

func SelectProject(ctx context.Context, p Project, id string, o PreviewOptions) (Result, error) {
	_, source, e := p.model(id, false)
	if e != nil {
		return Result{}, e
	}
	if o.URI == "" || o.Revision == "" {
		return Result{}, failure("arguments", fmt.Errorf("selection requires URI and expected source revision"))
	}
	selected, e := viewpoint.ParseURI(o.URI)
	if e != nil {
		return Result{}, failure("selection", e)
	}
	if selected.Project != p.Registration.ProjectID {
		return Result{}, failure("selection", fmt.Errorf("URI project does not match selected registration"))
	}
	o.Source = source
	o.operation = "select"
	return Preview(ctx, o)
}
