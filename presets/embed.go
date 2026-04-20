// Package presets embeds the built-in YAML preset files.
package presets

import "embed"

//go:embed *.yaml
var FS embed.FS
