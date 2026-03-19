package tui

import (
	"github.com/cellear/ddev-drush-tui/internal/ddev"
	"github.com/cellear/ddev-drush-tui/internal/drush"
)

type App struct {
	Context  *ddev.Context
	Commands []drush.NamespaceGroup
}

func NewApp(ctx *ddev.Context, commands []drush.NamespaceGroup) *App {
	return &App{
		Context:  ctx,
		Commands: commands,
	}
}

func (a *App) Run() error {
	return nil
}
