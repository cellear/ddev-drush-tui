package ddev

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNotInProject     = errors.New("not in a DDEV project directory")
	ErrDrushUnavailable = errors.New("drush is not available in this DDEV project")
)

type Context struct {
	ProjectName string
	AppRoot     string
}

type describeEnvelope struct {
	Name    string         `json:"name"`
	AppRoot string         `json:"approot"`
	Raw     describeResult `json:"raw"`
}

type describeResult struct {
	Name    string `json:"name"`
	AppRoot string `json:"approot"`
}

func Detect() (*Context, error) {
	describe, err := describeProject()
	if err != nil {
		return nil, err
	}

	if err := verifyDrush(); err != nil {
		return nil, err
	}

	return &Context{
		ProjectName: describe.Name,
		AppRoot:     describe.AppRoot,
	}, nil
}

func describeProject() (*describeResult, error) {
	output, err := runDescribe()
	if err != nil {
		return nil, err
	}

	var envelope describeEnvelope
	if err := json.Unmarshal(output, &envelope); err != nil {
		return nil, fmt.Errorf("parse ddev describe output: %w", err)
	}

	result := envelope.Raw
	if result.Name == "" {
		result.Name = envelope.Name
	}
	if result.AppRoot == "" {
		result.AppRoot = envelope.AppRoot
	}
	if result.Name == "" {
		return nil, fmt.Errorf("parse ddev describe output: missing project name")
	}

	return &result, nil
}

func runDescribe() ([]byte, error) {
	output, err := runDDEV("describe", "--json-output")
	if err == nil {
		return output, nil
	}

	if strings.Contains(err.Error(), "unknown flag: --json-output") {
		return runDDEV("describe", "--json")
	}

	if isNotInProjectError(err.Error()) {
		return nil, ErrNotInProject
	}

	return nil, fmt.Errorf("ddev describe failed: %w", err)
}

func verifyDrush() error {
	if _, err := runDDEV("drush", "version", "--format=string"); err != nil {
		return ErrDrushUnavailable
	}

	return nil
}

func runDDEV(args ...string) ([]byte, error) {
	cmd := exec.Command("ddev", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = strings.TrimSpace(stdout.String())
		}
		if message == "" {
			message = err.Error()
		}
		return nil, errors.New(message)
	}

	return stdout.Bytes(), nil
}

func isNotInProjectError(message string) bool {
	return strings.Contains(message, "No DDEV project was found") ||
		strings.Contains(message, "No project was found") ||
		strings.Contains(message, "Please run 'ddev config' to configure a project")
}
