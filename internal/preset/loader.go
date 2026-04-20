package preset

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mirakuta-dev/mirakuta/presets"
	"gopkg.in/yaml.v3"
)

type Loader struct {
	fsys fs.FS
}

func NewEmbeddedLoader() *Loader {
	return &Loader{fsys: presets.FS}
}

// NewFSLoader loads presets from an arbitrary fs.FS. Used in tests.
func NewFSLoader(fsys fs.FS) *Loader {
	return &Loader{fsys: fsys}
}

// List returns the names of every preset known to the loader (without the .yaml suffix).
func (l *Loader) List() ([]string, error) {
	entries, err := fs.ReadDir(l.fsys, ".")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".yaml") {
			continue
		}
		names = append(names, strings.TrimSuffix(n, ".yaml"))
	}
	return names, nil
}

// Load resolves a preset by name, recursively applying `extends` and expanding
// high-level language/terminal/editor keys into concrete tools.
func (l *Loader) Load(name string) (Preset, error) {
	return l.loadResolved(name, map[string]bool{})
}

// LoadFile resolves a preset from an arbitrary filesystem path.
// Only the root preset is loaded from disk; its `extends` entries must still
// resolve against the embedded set.
func (l *Loader) LoadFile(path string) (Preset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Preset{}, err
	}
	var raw Preset
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Preset{}, fmt.Errorf("parse %s: %w", path, err)
	}
	resolved, err := l.resolveExtends(raw, map[string]bool{filepath.Base(path): true})
	if err != nil {
		return Preset{}, err
	}
	resolved.Tools = expandHighLevel(resolved)
	return resolved, nil
}

func (l *Loader) loadResolved(name string, seen map[string]bool) (Preset, error) {
	if seen[name] {
		return Preset{}, fmt.Errorf("preset %q forms an extends cycle", name)
	}
	seen[name] = true
	defer delete(seen, name)

	raw, err := l.loadRaw(name)
	if err != nil {
		return Preset{}, err
	}
	resolved, err := l.resolveExtends(raw, seen)
	if err != nil {
		return Preset{}, err
	}
	resolved.Tools = expandHighLevel(resolved)
	return resolved, nil
}

func (l *Loader) loadRaw(name string) (Preset, error) {
	data, err := fs.ReadFile(l.fsys, name+".yaml")
	if err != nil {
		return Preset{}, fmt.Errorf("preset %q not found: %w", name, err)
	}
	var p Preset
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Preset{}, fmt.Errorf("parse %s.yaml: %w", name, err)
	}
	return p, nil
}

// resolveExtends merges each parent preset's fields into the child. Child values
// win on scalar conflicts; list fields are concatenated with deduplication.
func (l *Loader) resolveExtends(child Preset, seen map[string]bool) (Preset, error) {
	if len(child.Extends) == 0 {
		return child, nil
	}
	merged := Preset{}
	for _, parentName := range child.Extends {
		parent, err := l.loadResolved(parentName, seen)
		if err != nil {
			return Preset{}, err
		}
		merged = mergePresets(merged, parent)
	}
	// Child has no Extends anymore at the merge step — strip to avoid re-resolving.
	childNoExt := child
	childNoExt.Extends = nil
	return mergePresets(merged, childNoExt), nil
}

func mergePresets(base, over Preset) Preset {
	out := base
	if over.Name != "" {
		out.Name = over.Name
	}
	if over.Description != "" {
		out.Description = over.Description
	}
	out.Languages = dedupStrings(append(out.Languages, over.Languages...))
	out.Terminals = dedupStrings(append(out.Terminals, over.Terminals...))
	out.Editors = mergeEditors(out.Editors, over.Editors)
	out.Tools = mergeTools(out.Tools, over.Tools)
	if over.Git.Name != "" {
		out.Git.Name = over.Git.Name
	}
	if over.Git.Email != "" {
		out.Git.Email = over.Git.Email
	}
	return out
}

func mergeEditors(base, over []Editor) []Editor {
	idx := map[string]int{}
	out := make([]Editor, 0, len(base)+len(over))
	for _, list := range [][]Editor{base, over} {
		for _, e := range list {
			if i, ok := idx[e.ID]; ok {
				out[i].Extensions = dedupStrings(append(out[i].Extensions, e.Extensions...))
				continue
			}
			idx[e.ID] = len(out)
			out = append(out, Editor{ID: e.ID, Extensions: dedupStrings(e.Extensions)})
		}
	}
	return out
}

func mergeTools(base, over []Tool) []Tool {
	idx := map[string]int{}
	out := make([]Tool, 0, len(base)+len(over))
	for _, list := range [][]Tool{base, over} {
		for _, t := range list {
			if i, ok := idx[t.ID]; ok {
				out[i] = t
				continue
			}
			idx[t.ID] = len(out)
			out = append(out, t)
		}
	}
	return out
}

func dedupStrings(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
