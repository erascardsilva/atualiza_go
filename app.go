// Atualiza GO - Atualizador de Sistema Linux
// Erasmo Cardoso - Software Engineer |  Electronics Technician
package main

import (
	"context"
	"os"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDistroInfo() DistroInfo {
	return detectDistro()
}

func (a *App) IsRestrictedSandbox() bool {
	snapEnv := os.Getenv("SNAP")
	return snapEnv != ""
}
