package drush

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Result holds the output and exit code of a Drush command execution.
type Result struct {
	Command  string // The full command string that was run.
	Output   string // Combined stdout and stderr.
	ExitCode int
}

// DefaultTimeout is the maximum time a command can run before being killed.
const DefaultTimeout = 60 * time.Second

// Execute runs a Drush command via ddev with the given arguments and options.
// Arguments are positional values; options are "--key=value" or "--flag" strings.
func Execute(command string, args []string, opts []string) (*Result, error) {
	// Build the command: ddev drush <command> [args...] [--opts...]
	cmdArgs := []string{"drush", command}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, opts...)

	fullCommand := "ddev " + strings.Join(cmdArgs, " ")

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ddev", cmdArgs...)
	out, err := cmd.CombinedOutput()

	result := &Result{
		Command: fullCommand,
		Output:  string(out),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.ExitCode = 124 // Standard timeout exit code.
		result.Output += "\n\nCommand timed out after " + DefaultTimeout.String()
		return result, nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		return nil, fmt.Errorf("execute %s: %w", command, err)
	}

	return result, nil
}
