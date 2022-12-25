//go:build windows

package icon

import _ "embed"

var (
	//go:embed gopher.ico
	IconTray []byte
	//go:embed gopher.png
	IconPoster []byte
)
