package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type checkResult struct {
	name    string
	ok      bool
	version string
	note    string
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Diagnose your Windows development environment",
	Run:   runCheck,
}

func runCheck(cmd *cobra.Command, args []string) {
	fmt.Println("Mirakuta Environment Check")
	fmt.Println(strings.Repeat("─", 40))

	results := []checkResult{
		checkBinary("Git", "git", "--version"),
		checkBinary("Go", "go", "version"),
		checkBinary("Node.js", "node", "--version"),
		checkBinary("npm", "npm", "--version"),
		checkBinary("winget", "winget", "--version"),
		checkWSL(),
		checkEnv("GOPATH"),
		checkEnvOptional("GOROOT"),
	}

	allOk := true
	for _, r := range results {
		printResult(r)
		if !r.ok {
			allOk = false
		}
	}

	fmt.Println(strings.Repeat("─", 40))
	if allOk {
		fmt.Println("All checks passed.")
	} else {
		fmt.Println("Some checks failed. Run 'mirakuta fix' to resolve. (coming soon)")
		os.Exit(1)
	}
}

func checkBinary(name, bin, versionFlag string) checkResult {
	out, err := exec.Command(bin, versionFlag).Output()
	if err != nil {
		return checkResult{name: name, ok: false, note: "not found"}
	}
	ver := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	return checkResult{name: name, ok: true, version: ver}
}

func checkWSL() checkResult {
	out, err := exec.Command("wsl", "--status").Output()
	if err != nil {
		// wsl exists but may need elevation; try --list as fallback
		_, err2 := exec.Command("wsl", "--list").Output()
		if err2 != nil {
			return checkResult{name: "WSL2", ok: false, note: "not installed"}
		}
	}
	if strings.Contains(strings.ToLower(string(out)), "2") || err == nil {
		return checkResult{name: "WSL2", ok: true, version: "enabled"}
	}
	return checkResult{name: "WSL2", ok: false, note: "WSL1 or not configured"}
}

func checkEnv(key string) checkResult {
	val := os.Getenv(key)
	if val == "" {
		return checkResult{name: key, ok: false, note: "not set"}
	}
	return checkResult{name: key, ok: true, version: val}
}

// checkEnvOptional marks missing env vars as a warning (ok=true) rather than failure.
func checkEnvOptional(key string) checkResult {
	val := os.Getenv(key)
	if val == "" {
		return checkResult{name: key, ok: true, version: "(not set, optional)"}
	}
	return checkResult{name: key, ok: true, version: val}
}

func printResult(r checkResult) {
	mark := "✓"
	detail := r.version
	if !r.ok {
		mark = "✗"
		detail = r.note
	}
	fmt.Printf("  [%s] %-10s %s\n", mark, r.name, detail)
}
