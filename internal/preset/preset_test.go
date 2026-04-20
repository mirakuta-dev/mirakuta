package preset

import (
	"testing"
	"testing/fstest"
)

func TestLoadMinimal(t *testing.T) {
	l := NewEmbeddedLoader()
	p, err := l.Load("minimal")
	if err != nil {
		t.Fatalf("load minimal: %v", err)
	}
	if p.Name != "minimal" {
		t.Errorf("name = %q, want minimal", p.Name)
	}
	if !hasToolID(p.Tools, "git") {
		t.Errorf("minimal must include git; got %v", toolIDs(p.Tools))
	}
	if !hasToolID(p.Tools, "windows-terminal") || !hasToolID(p.Tools, "pwsh") {
		t.Errorf("minimal terminals not expanded; got %v", toolIDs(p.Tools))
	}
	if !hasToolID(p.Tools, "vscode") {
		t.Errorf("minimal editor (vscode) not expanded; got %v", toolIDs(p.Tools))
	}
}

func TestLoadBackendExtendsMinimal(t *testing.T) {
	l := NewEmbeddedLoader()
	p, err := l.Load("backend")
	if err != nil {
		t.Fatalf("load backend: %v", err)
	}
	// Inherited from minimal.
	for _, id := range []string{"git", "windows-terminal", "pwsh", "vscode"} {
		if !hasToolID(p.Tools, id) {
			t.Errorf("backend missing inherited tool %q", id)
		}
	}
	// From backend itself.
	for _, id := range []string{"go", "node", "docker-desktop", "postman"} {
		if !hasToolID(p.Tools, id) {
			t.Errorf("backend missing own tool %q", id)
		}
	}
	// vscode should carry extensions from backend.
	var vs Editor
	for _, e := range p.Editors {
		if e.ID == "vscode" {
			vs = e
		}
	}
	if !containsString(vs.Extensions, "golang.Go") {
		t.Errorf("vscode extensions missing golang.Go; got %v", vs.Extensions)
	}
}

func TestLoadFullstackMultiExtends(t *testing.T) {
	l := NewEmbeddedLoader()
	p, err := l.Load("fullstack")
	if err != nil {
		t.Fatalf("load fullstack: %v", err)
	}
	for _, id := range []string{"go", "node", "docker-desktop", "postman", "pnpm", "chrome", "figma"} {
		if !hasToolID(p.Tools, id) {
			t.Errorf("fullstack missing tool %q", id)
		}
	}
	// node must appear exactly once despite being in both parents.
	count := 0
	for _, t := range p.Tools {
		if t.ID == "node" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("node appears %d times in fullstack tools, want 1", count)
	}
	// vscode extensions from both parents should be merged without duplicates.
	var vs Editor
	for _, e := range p.Editors {
		if e.ID == "vscode" {
			vs = e
		}
	}
	if !containsString(vs.Extensions, "golang.Go") || !containsString(vs.Extensions, "esbenp.prettier-vscode") {
		t.Errorf("fullstack vscode extensions not merged; got %v", vs.Extensions)
	}
	eslintCount := 0
	for _, e := range vs.Extensions {
		if e == "dbaeumer.vscode-eslint" {
			eslintCount++
		}
	}
	if eslintCount != 1 {
		t.Errorf("eslint extension appears %d times, want 1", eslintCount)
	}
}

func TestLoadAllRounder(t *testing.T) {
	l := NewEmbeddedLoader()
	p, err := l.Load("all-rounder")
	if err != nil {
		t.Fatalf("load all-rounder: %v", err)
	}
	for _, id := range []string{"go", "node", "python", "uv", "docker-desktop", "chrome"} {
		if !hasToolID(p.Tools, id) {
			t.Errorf("all-rounder missing tool %q", id)
		}
	}
}

func TestList(t *testing.T) {
	l := NewEmbeddedLoader()
	names, err := l.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	want := []string{"all-rounder", "backend", "data", "frontend", "fullstack", "minimal"}
	for _, w := range want {
		if !containsString(names, w) {
			t.Errorf("List() missing %q; got %v", w, names)
		}
	}
}

func TestExtendsCycleDetected(t *testing.T) {
	fsys := fstest.MapFS{
		"a.yaml": &fstest.MapFile{Data: []byte("name: a\nextends: [b]\n")},
		"b.yaml": &fstest.MapFile{Data: []byte("name: b\nextends: [a]\n")},
	}
	l := NewFSLoader(fsys)
	if _, err := l.Load("a"); err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestUnknownPreset(t *testing.T) {
	l := NewEmbeddedLoader()
	if _, err := l.Load("does-not-exist"); err == nil {
		t.Fatal("expected error for unknown preset")
	}
}

func hasToolID(tools []Tool, id string) bool {
	for _, t := range tools {
		if t.ID == id {
			return true
		}
	}
	return false
}

func toolIDs(tools []Tool) []string {
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		out = append(out, t.ID)
	}
	return out
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
