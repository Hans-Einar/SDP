package sdptool

import "fmt"

const (
	maxPreviewBundleFiles = 512
	maxPreviewBundleBytes = 128 << 20
)

// Count the sealed output, including generated metadata and retained sources.
// Preview and source reading use the same publication budget.
func checkPreviewBundle(files map[string][]byte, label string) error {
	var size int64
	for _, data := range files {
		size += int64(len(data))
	}
	if len(files) > maxPreviewBundleFiles || size > maxPreviewBundleBytes {
		return failure("limit", fmt.Errorf("%s exceeds %d files or %d MiB (actual: %d files, %d bytes); select a diagram or narrower focus", label, maxPreviewBundleFiles, maxPreviewBundleBytes>>20, len(files), size))
	}
	return nil
}
