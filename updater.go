// Atualiza GO - Motor de Atualização
// Erasmo Cardoso - Dev

package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UpdateStep struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Command  string `json:"command"`
	NeedRoot bool   `json:"needRoot"`
}

type UpdateProgress struct {
	StepID   string  `json:"stepId"`
	StepName string  `json:"stepName"`
	Line     string  `json:"line"`
	Percent  float64 `json:"percent"`
	Done     bool    `json:"done"`
	Error    string  `json:"error"`
}

func (a *App) GetUpdateSteps() []UpdateStep {
	distro := detectDistro()
	steps := []UpdateStep{}

	switch distro.Family {
	case "debian":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (APT)",
			Command:  "apt update -y && apt upgrade -y",
			NeedRoot: true,
		})
	case "arch":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (Pacman)",
			Command:  "pacman -Syu --noconfirm",
			NeedRoot: true,
		})
	}

	if distro.HasFlatpak {
		steps = append(steps, UpdateStep{
			ID:       "flatpak_update",
			Label:    "Atualizar Flatpak",
			Command:  "flatpak update --assumeyes",
			NeedRoot: false,
		})
	}

	if distro.HasSnap {
		steps = append(steps, UpdateStep{
			ID:       "snap_update",
			Label:    "Atualizar Snap",
			Command:  "snap refresh",
			NeedRoot: true,
		})
	}

	return steps
}

func (a *App) RunUpdate(stepIDs []string) {
	go func() {
		steps := a.GetUpdateSteps()
		stepMap := map[string]UpdateStep{}
		for _, s := range steps {
			stepMap[s.ID] = s
		}

		for _, id := range stepIDs {
			step, ok := stepMap[id]
			if !ok {
				continue
			}
			a.executeStep(step)
		}

		runtime.EventsEmit(a.ctx, "update:complete", nil)
	}()
}

func (a *App) executeStep(step UpdateStep) {
	var cmd *exec.Cmd

	if step.NeedRoot {
		cmd = exec.Command("pkexec", "bash", "-c", step.Command)
	} else {
		cmd = exec.Command("bash", "-c", step.Command)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		runtime.EventsEmit(a.ctx, "update:progress", UpdateProgress{
			StepID:   step.ID,
			StepName: step.Label,
			Error:    fmt.Sprintf("Erro ao criar pipe: %v", err),
			Done:     true,
		})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		runtime.EventsEmit(a.ctx, "update:progress", UpdateProgress{
			StepID:   step.ID,
			StepName: step.Label,
			Error:    fmt.Sprintf("Erro ao criar pipe stderr: %v", err),
			Done:     true,
		})
		return
	}

	runtime.EventsEmit(a.ctx, "update:progress", UpdateProgress{
		StepID:   step.ID,
		StepName: step.Label,
		Line:     "Iniciando " + step.Label + "...",
		Percent:  0,
	})

	if err := cmd.Start(); err != nil {
		runtime.EventsEmit(a.ctx, "update:progress", UpdateProgress{
			StepID:   step.ID,
			StepName: step.Label,
			Error:    fmt.Sprintf("Erro ao iniciar: %v", err),
			Done:     true,
		})
		return
	}

	var wg sync.WaitGroup
	lineCount := 0

	readPipe := func(scanner *bufio.Scanner) {
		defer wg.Done()
		for scanner.Scan() {
			line := scanner.Text()
			lineCount++
			pct := estimateProgress(step.ID, lineCount, line)

			runtime.EventsEmit(a.ctx, "update:progress", UpdateProgress{
				StepID:   step.ID,
				StepName: step.Label,
				Line:     line,
				Percent:  pct,
			})
		}
	}

	wg.Add(2)
	go readPipe(bufio.NewScanner(stdout))
	go readPipe(bufio.NewScanner(stderr))
	wg.Wait()

	err = cmd.Wait()

	progress := UpdateProgress{
		StepID:   step.ID,
		StepName: step.Label,
		Done:     true,
		Percent:  100,
	}

	if err != nil {
		progress.Error = fmt.Sprintf("Falha: %v", err)
		progress.Percent = 0
	} else {
		progress.Line = step.Label + " concluído"
	}

	runtime.EventsEmit(a.ctx, "update:progress", progress)
}

func estimateProgress(stepID string, lineCount int, line string) float64 {
	lower := strings.ToLower(line)

	if strings.Contains(lower, "reading package") || strings.Contains(lower, ":: synchronizing") {
		return 10
	}
	if strings.Contains(lower, "building dependency") || strings.Contains(lower, "resolving dependencies") {
		return 20
	}
	if strings.Contains(lower, "the following packages will be upgraded") ||
		strings.Contains(lower, "packages to install") {
		return 30
	}
	if strings.Contains(lower, "unpacking") || strings.Contains(lower, "installing") {
		return 50 + float64(lineCount%30)
	}
	if strings.Contains(lower, "setting up") || strings.Contains(lower, "configuring") {
		return 70 + float64(lineCount%20)
	}
	if strings.Contains(lower, "processing triggers") {
		return 90
	}

	base := float64(lineCount) * 0.5
	if base > 95 {
		base = 95
	}
	return base
}

func (a *App) RunSystemAction(action string) error {
	var cmd *exec.Cmd

	switch action {
	case "reboot":
		cmd = exec.Command("pkexec", "systemctl", "reboot")
	case "shutdown":
		cmd = exec.Command("pkexec", "systemctl", "poweroff")
	default:
		return fmt.Errorf("ação desconhecida: %s", action)
	}

	return cmd.Run()
}
