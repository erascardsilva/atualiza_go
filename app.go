// Atualiza GO - Atualizador de Sistema Linux
// Erasmo Cardoso - Dev

package main

import (
	"context"
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
