package api

import "backend/internal"

type Server interface {
	Init() error
	MapHandlers(app *internal.App) error
	Run() error
}
