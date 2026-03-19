package drush

type NamespaceGroup struct {
	Namespace string
	Commands  []Command
}

type Command struct {
	Name        string
	Description string
	Aliases     []string
}

func ListCommands() ([]NamespaceGroup, error) {
	return nil, nil
}
