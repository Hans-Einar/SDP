//go:build bundled

package sdptool

import _ "embed"

//go:embed payload/sdui-fyne.gz
var bundledNativeHost []byte
