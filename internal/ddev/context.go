package ddev

type Context struct {
	ProjectName string
	AppRoot     string
}

func Detect() (*Context, error) {
	return &Context{}, nil
}
