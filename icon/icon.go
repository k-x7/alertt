//go:build !windows

package icon

import _ "embed"

var (
	//go:embed gopher.png
	IconTray   []byte
	IconPoster []byte = IconTray
)
