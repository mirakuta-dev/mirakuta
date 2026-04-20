package runner

import (
	"fmt"

	"github.com/mirakuta-dev/mirakuta/internal/preset"
)

// Step is one concrete action the runner will take. Each step is a single
// external command invocation — keep them self-contained so failures are
// isolated.
type Step struct {
	Desc string
	Bin  string
	Args []string
}

// Plan converts a resolved preset into an ordered list of executable steps.
// Order: install packages → configure git → install editor extensions.
// Tools whose `via` we don't support yet are skipped (with a marker step so
// the user sees they were noticed).
func Plan(p preset.Preset) []Step {
	var steps []Step

	for _, t := range p.Tools {
		switch t.Via {
		case "winget":
			steps = append(steps, Step{
				Desc: fmt.Sprintf("Install %s via winget (%s)", t.ID, t.Pkg),
				Bin:  "winget",
				Args: []string{
					"install", "--exact", "--id", t.Pkg,
					"--accept-package-agreements", "--accept-source-agreements",
					"--silent",
				},
			})
		default:
			steps = append(steps, Step{
				Desc: fmt.Sprintf("Skip %s: via=%q not supported yet", t.ID, t.Via),
				Bin:  "",
			})
		}
	}

	if p.Git.Name != "" {
		steps = append(steps, Step{
			Desc: fmt.Sprintf("git config --global user.name %q", p.Git.Name),
			Bin:  "git",
			Args: []string{"config", "--global", "user.name", p.Git.Name},
		})
	}
	if p.Git.Email != "" {
		steps = append(steps, Step{
			Desc: fmt.Sprintf("git config --global user.email %q", p.Git.Email),
			Bin:  "git",
			Args: []string{"config", "--global", "user.email", p.Git.Email},
		})
	}

	for _, ed := range p.Editors {
		if ed.ID != "vscode" {
			continue
		}
		for _, ext := range ed.Extensions {
			steps = append(steps, Step{
				Desc: fmt.Sprintf("Install VS Code extension %s", ext),
				Bin:  "code",
				Args: []string{"--install-extension", ext, "--force"},
			})
		}
	}

	return steps
}
