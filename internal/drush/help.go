package drush

type CommandHelp struct{}

func Help(command string) (*CommandHelp, error) {
	return &CommandHelp{}, nil
}
