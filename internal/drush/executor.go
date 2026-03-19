package drush

type Result struct {
	Output   string
	ExitCode int
}

func Execute(command string, args []string) (*Result, error) {
	return &Result{}, nil
}
