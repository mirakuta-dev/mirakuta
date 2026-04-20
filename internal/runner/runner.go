package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Executor runs a single command. Injectable so tests can stub it.
type Executor interface {
	Run(bin string, args []string, verbose bool) error
}

type execRunner struct{}

func (execRunner) Run(bin string, args []string, verbose bool) error {
	cmd := exec.Command(bin, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// RealExecutor runs commands via os/exec. Use for production; tests pass fakes.
func RealExecutor() Executor { return execRunner{} }

type Runner struct {
	DryRun   bool
	Verbose  bool
	Out      io.Writer
	Executor Executor
}

// New returns a Runner with sensible defaults. Out defaults to os.Stdout and
// Executor to the real os/exec-backed one.
func New() *Runner {
	return &Runner{
		Out:      os.Stdout,
		Executor: RealExecutor(),
	}
}

// Result summarizes an Execute run.
type Result struct {
	Total     int
	Skipped   int
	Succeeded int
	Failed    int
	Errors    []StepError
}

type StepError struct {
	Step Step
	Err  error
}

// Execute runs each step in order. A failure in one step is recorded but does
// not abort subsequent steps — tool installs are largely independent.
func (r *Runner) Execute(steps []Step) Result {
	res := Result{Total: len(steps)}
	for i, s := range steps {
		fmt.Fprintf(r.Out, "[%d/%d] %s\n", i+1, len(steps), s.Desc)

		if s.Bin == "" {
			res.Skipped++
			continue
		}

		if r.DryRun {
			fmt.Fprintf(r.Out, "       $ %s %s\n", s.Bin, strings.Join(s.Args, " "))
			res.Skipped++
			continue
		}

		if err := r.Executor.Run(s.Bin, s.Args, r.Verbose); err != nil {
			fmt.Fprintf(r.Out, "       failed: %v\n", err)
			res.Failed++
			res.Errors = append(res.Errors, StepError{Step: s, Err: err})
			continue
		}
		res.Succeeded++
	}
	return res
}
