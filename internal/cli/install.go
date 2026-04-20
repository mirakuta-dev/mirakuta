package cli

import (
	"fmt"

	"github.com/mirakuta-dev/mirakuta/internal/preset"
	"github.com/mirakuta-dev/mirakuta/internal/profile"
	"github.com/mirakuta-dev/mirakuta/internal/runner"
	"github.com/spf13/cobra"
)

var (
	installFile    string
	installDryRun  bool
	installVerbose bool
)

var installCmd = &cobra.Command{
	Use:   "install [preset]",
	Short: "Install a dev environment preset",
	Long: `Install a dev environment preset.

Three modes:
  mirakuta install                  interactive wizard, saves your choices to ~/.mirakuta/profile.yaml
  mirakuta install <preset>         install one of the built-in presets (minimal, backend, frontend, fullstack, data, all-rounder)
  mirakuta install --file <path>    install from an external YAML file`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installFile, "file", "", "Load preset from an external YAML file")
	installCmd.Flags().BoolVar(&installDryRun, "dry-run", false, "Print commands without executing")
	installCmd.Flags().BoolVar(&installVerbose, "verbose", false, "Print command output")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	p, err := resolvePreset(args)
	if err != nil {
		return err
	}

	steps := runner.Plan(p)
	if len(steps) == 0 {
		fmt.Println("Nothing to install.")
		return nil
	}

	r := runner.New()
	r.DryRun = installDryRun
	r.Verbose = installVerbose

	fmt.Printf("Preset: %s (%d steps)\n", p.Name, len(steps))
	if installDryRun {
		fmt.Println("Dry-run mode: commands will be printed, not executed.")
	}
	fmt.Println()

	res := r.Execute(steps)

	fmt.Printf("\nDone: %d succeeded, %d failed, %d skipped\n",
		res.Succeeded, res.Failed, res.Skipped)
	if res.Failed > 0 {
		return fmt.Errorf("%d steps failed", res.Failed)
	}
	return nil
}

func resolvePreset(args []string) (preset.Preset, error) {
	loader := preset.NewEmbeddedLoader()

	switch {
	case installFile != "":
		return loader.LoadFile(installFile)

	case len(args) == 1:
		return loader.Load(args[0])

	default:
		p, confirm, err := runWizard()
		if err != nil {
			return preset.Preset{}, err
		}
		if !confirm {
			return preset.Preset{}, fmt.Errorf("cancelled")
		}
		// Persist the wizard output before resolving tools so the user can
		// re-run the same profile later.
		if path, err := profile.Save(p); err != nil {
			fmt.Printf("Warning: could not save profile: %v\n", err)
		} else {
			fmt.Printf("Profile saved to %s\n", path)
		}
		return preset.ExpandForWizard(p), nil
	}
}
