//go:build !linux

package governancehost

import (
	"context"
	"fmt"
)

func OpenSocket(context.Context, string) (*Client, error) {
	return nil, fmt.Errorf("private host socket profile requires Linux peer identity")
}
