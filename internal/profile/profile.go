// Package profile persists a user-selected preset to the home directory so it
// can be re-applied or shared later.
package profile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mirakuta-dev/mirakuta/internal/preset"
	"gopkg.in/yaml.v3"
)

const (
	dirName  = ".mirakuta"
	fileName = "profile.yaml"
)

// Path returns the absolute path where the user profile is stored.
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, dirName, fileName), nil
}

// Save writes the preset to ~/.mirakuta/profile.yaml, creating the directory
// if needed. On Windows the directory is marked hidden so it doesn't clutter
// Explorer.
func Save(p preset.Preset) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	hideOnWindows(dir)

	data, err := yaml.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal profile: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

// Load reads and parses ~/.mirakuta/profile.yaml if it exists.
func Load() (preset.Preset, error) {
	path, err := Path()
	if err != nil {
		return preset.Preset{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return preset.Preset{}, err
	}
	var p preset.Preset
	if err := yaml.Unmarshal(data, &p); err != nil {
		return preset.Preset{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return p, nil
}

func hideOnWindows(dir string) {
	if runtime.GOOS != "windows" {
		return
	}
	// Best-effort: ignore errors. The dot-prefix convention alone is acceptable
	// if attrib is missing or fails.
	_ = exec.Command("attrib", "+h", dir).Run()
}
