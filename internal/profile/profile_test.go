package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mirakuta-dev/mirakuta/internal/preset"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp) // Windows
	t.Setenv("HOME", tmp)        // Unix

	in := preset.Preset{
		Name:      "custom",
		Languages: []string{"go", "node"},
		Git:       preset.Git{Name: "Mark", Email: "mark@example.com"},
	}

	path, err := Save(in)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("path not absolute: %s", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("saved file missing: %v", err)
	}

	out, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if out.Name != in.Name {
		t.Errorf("name = %q, want %q", out.Name, in.Name)
	}
	if out.Git.Email != in.Git.Email {
		t.Errorf("email = %q, want %q", out.Git.Email, in.Git.Email)
	}
	if len(out.Languages) != 2 {
		t.Errorf("languages = %v, want 2", out.Languages)
	}
}
