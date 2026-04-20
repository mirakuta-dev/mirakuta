package runner

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/mirakuta-dev/mirakuta/internal/preset"
)

type fakeExec struct {
	calls   []call
	errFor  map[string]error
	verbose bool
}

type call struct {
	Bin  string
	Args []string
}

func (f *fakeExec) Run(bin string, args []string, verbose bool) error {
	f.calls = append(f.calls, call{Bin: bin, Args: args})
	f.verbose = verbose
	if err, ok := f.errFor[bin]; ok {
		return err
	}
	return nil
}

func TestPlanWingetStep(t *testing.T) {
	p := preset.Preset{
		Tools: []preset.Tool{{ID: "git", Via: "winget", Pkg: "Git.Git"}},
	}
	steps := Plan(p)
	if len(steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(steps))
	}
	s := steps[0]
	if s.Bin != "winget" {
		t.Errorf("bin = %q, want winget", s.Bin)
	}
	if !containsArg(s.Args, "Git.Git") || !containsArg(s.Args, "--exact") {
		t.Errorf("args missing expected flags: %v", s.Args)
	}
}

func TestPlanSkipsUnsupportedVia(t *testing.T) {
	p := preset.Preset{
		Tools: []preset.Tool{{ID: "weird", Via: "scoop", Pkg: "whatever"}},
	}
	steps := Plan(p)
	if len(steps) != 1 || steps[0].Bin != "" {
		t.Errorf("expected one skip-marker step, got %+v", steps)
	}
}

func TestPlanGitConfig(t *testing.T) {
	p := preset.Preset{
		Git: preset.Git{Name: "Mark", Email: "mark@example.com"},
	}
	steps := Plan(p)
	if len(steps) != 2 {
		t.Fatalf("expected 2 git config steps, got %d", len(steps))
	}
	if steps[0].Bin != "git" || !containsArg(steps[0].Args, "user.name") {
		t.Errorf("first step not user.name: %+v", steps[0])
	}
	if !containsArg(steps[1].Args, "user.email") {
		t.Errorf("second step not user.email: %+v", steps[1])
	}
}

func TestPlanVSCodeExtensions(t *testing.T) {
	p := preset.Preset{
		Editors: []preset.Editor{
			{ID: "vscode", Extensions: []string{"golang.Go", "esbenp.prettier-vscode"}},
			{ID: "neovim", Extensions: []string{"ignored"}},
		},
	}
	steps := Plan(p)
	if len(steps) != 2 {
		t.Fatalf("expected 2 vscode ext steps, got %d", len(steps))
	}
	for _, s := range steps {
		if s.Bin != "code" {
			t.Errorf("bin = %q, want code", s.Bin)
		}
	}
}

func TestExecuteDryRunDoesNotCallExecutor(t *testing.T) {
	f := &fakeExec{}
	r := &Runner{DryRun: true, Out: &bytes.Buffer{}, Executor: f}
	steps := []Step{{Desc: "x", Bin: "winget", Args: []string{"install"}}}

	res := r.Execute(steps)

	if len(f.calls) != 0 {
		t.Errorf("dry-run called executor: %+v", f.calls)
	}
	if res.Skipped != 1 {
		t.Errorf("skipped = %d, want 1", res.Skipped)
	}
}

func TestExecuteContinuesOnError(t *testing.T) {
	f := &fakeExec{errFor: map[string]error{"fails": errors.New("boom")}}
	out := &bytes.Buffer{}
	r := &Runner{Out: out, Executor: f}
	steps := []Step{
		{Desc: "ok", Bin: "ok"},
		{Desc: "bad", Bin: "fails"},
		{Desc: "ok2", Bin: "ok2"},
	}

	res := r.Execute(steps)

	if res.Succeeded != 2 {
		t.Errorf("succeeded = %d, want 2", res.Succeeded)
	}
	if res.Failed != 1 {
		t.Errorf("failed = %d, want 1", res.Failed)
	}
	if len(f.calls) != 3 {
		t.Errorf("executor called %d times, want 3 (no abort on fail)", len(f.calls))
	}
	if !strings.Contains(out.String(), "failed: boom") {
		t.Errorf("stderr message missing from output: %q", out.String())
	}
}

func TestExecuteSkipsEmptyBin(t *testing.T) {
	f := &fakeExec{}
	r := &Runner{Out: &bytes.Buffer{}, Executor: f}
	steps := []Step{
		{Desc: "skip me", Bin: ""},
		{Desc: "run me", Bin: "ok"},
	}

	res := r.Execute(steps)

	if len(f.calls) != 1 {
		t.Errorf("executor calls = %d, want 1", len(f.calls))
	}
	if res.Skipped != 1 || res.Succeeded != 1 {
		t.Errorf("counts wrong: %+v", res)
	}
}

func containsArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}
