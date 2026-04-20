package preset

// languageTools maps a high-level language key to the concrete tools that back it.
// Preset authors write `languages: [go]` and the loader expands it using this table.
var languageTools = map[string][]Tool{
	"go":     {{ID: "go", Via: "winget", Pkg: "GoLang.Go"}},
	"node":   {{ID: "node", Via: "winget", Pkg: "OpenJS.NodeJS.LTS"}},
	"python": {{ID: "python", Via: "winget", Pkg: "Python.Python.3.12"}},
	"rust":   {{ID: "rustup", Via: "winget", Pkg: "Rustlang.Rustup"}},
}

var terminalTools = map[string][]Tool{
	"windows-terminal": {{ID: "windows-terminal", Via: "winget", Pkg: "Microsoft.WindowsTerminal"}},
	"pwsh":             {{ID: "pwsh", Via: "winget", Pkg: "Microsoft.PowerShell"}},
}

var editorTools = map[string][]Tool{
	"vscode":   {{ID: "vscode", Via: "winget", Pkg: "Microsoft.VisualStudioCode"}},
	"neovim":   {{ID: "neovim", Via: "winget", Pkg: "Neovim.Neovim"}},
	"intellij": {{ID: "intellij", Via: "winget", Pkg: "JetBrains.IntelliJIDEA.Community"}},
}

// expandHighLevel turns a preset's high-level keys (languages, terminals, editors)
// into concrete Tool entries, merged with the explicit tools list.
// Later entries with the same Tool.ID override earlier ones to keep the list unique.
func expandHighLevel(p Preset) []Tool {
	seen := map[string]int{}
	var out []Tool

	add := func(t Tool) {
		if idx, ok := seen[t.ID]; ok {
			out[idx] = t
			return
		}
		seen[t.ID] = len(out)
		out = append(out, t)
	}

	for _, lang := range p.Languages {
		for _, t := range languageTools[lang] {
			add(t)
		}
	}
	for _, term := range p.Terminals {
		for _, t := range terminalTools[term] {
			add(t)
		}
	}
	for _, ed := range p.Editors {
		for _, t := range editorTools[ed.ID] {
			add(t)
		}
	}
	for _, t := range p.Tools {
		add(t)
	}
	return out
}
