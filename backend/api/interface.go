package api

import "helprepet/internal"

type Server interface {
	Init() error
	MapHandlers(app *internal.App) error
	Run() error
}
