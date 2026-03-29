// Atualiza GO - Motor de Atualização
// Erasmo Cardoso - Dev

package main

import (
	"bufio"
	"fmt"
	"os"
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
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (APT)",
			Command:  "apt clean",
			NeedRoot: true,
		})
		if !distro.HasFlatpak {
			steps = append(steps, UpdateStep{
				ID:       "install_flatpak",
				Label:    "Instalar Suporte Flatpak",
				Command:  "apt install -y flatpak",
				NeedRoot: true,
			})
		}
		if !distro.HasSnap {
			steps = append(steps, UpdateStep{
				ID:       "install_snap",
				Label:    "Instalar Suporte Snap",
				Command:  "apt install -y snapd",
				NeedRoot: true,
			})
		}
	case "arch":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (Pacman)",
			Command:  "pacman -Syu --noconfirm",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (Pacman)",
			Command:  "pacman -Sc --noconfirm",
			NeedRoot: true,
		})
		if !distro.HasFlatpak {
			steps = append(steps, UpdateStep{
				ID:       "install_flatpak",
				Label:    "Instalar Suporte Flatpak",
				Command:  "pacman -S --noconfirm flatpak",
				NeedRoot: true,
			})
		}
		if !distro.HasSnap {
			steps = append(steps, UpdateStep{
				ID:       "install_snap",
				Label:    "Instalar Suporte Snap",
				Command:  "pacman -S --noconfirm snapd && systemctl enable --now snapd.socket",
				NeedRoot: true,
			})
		}
	case "fedora":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (DNF)",
			Command:  "dnf upgrade -y",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (DNF)",
			Command:  "dnf clean all",
			NeedRoot: true,
		})
		if !distro.HasFlatpak {
			steps = append(steps, UpdateStep{
				ID:       "install_flatpak",
				Label:    "Instalar Suporte Flatpak",
				Command:  "dnf install -y flatpak",
				NeedRoot: true,
			})
		}
		if !distro.HasSnap {
			steps = append(steps, UpdateStep{
				ID:       "install_snap",
				Label:    "Instalar Suporte Snap",
				Command:  "dnf install -y snapd && systemctl enable --now snapd.socket && ln -s /var/lib/snapd/snap /snap",
				NeedRoot: true,
			})
		}
	case "suse":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (Zypper)",
			Command:  "zypper --non-interactive update --auto-agree-with-licenses",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (Zypper)",
			Command:  "zypper clean --all",
			NeedRoot: true,
		})
		if !distro.HasFlatpak {
			steps = append(steps, UpdateStep{
				ID:       "install_flatpak",
				Label:    "Instalar Suporte Flatpak",
				Command:  "zypper --non-interactive install flatpak",
				NeedRoot: true,
			})
		}
		if !distro.HasSnap {
			steps = append(steps, UpdateStep{
				ID:       "install_snap",
				Label:    "Instalar Suporte Snap",
				Command:  "zypper --non-interactive install snapd && systemctl enable --now snapd.socket",
				NeedRoot: true,
			})
		}
	case "alpine":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (APK)",
			Command:  "apk update && apk upgrade",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (APK)",
			Command:  "apk cache clean",
			NeedRoot: true,
		})
	case "void":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (XBPS)",
			Command:  "xbps-install -Syu --yes",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (XBPS)",
			Command:  "xbps-remove -O",
			NeedRoot: true,
		})
	case "solus":
		steps = append(steps, UpdateStep{
			ID:       "system_update",
			Label:    "Atualizar Sistema (EOPKG)",
			Command:  "eopkg upgrade -y",
			NeedRoot: true,
		})
		steps = append(steps, UpdateStep{
			ID:       "system_cleanup",
			Label:    "Limpeza de Sistema (EOPKG)",
			Command:  "eopkg delete-cache",
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
			Command:  "snap wait system seed && snap refresh",
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

	inSandbox := false
	if _, err := os.Stat("/.flatpak-info"); err == nil {
		inSandbox = true
	}

	command := step.Command
	if step.NeedRoot {
		command = "pkexec bash -c '" + step.Command + "'"
	}

	if inSandbox {
		// No sandbox, usamos flatpak-spawn para rodar o comando no host
		// Nota: flatpak-spawn --host já lida com a saída e entrada se configurado
		cmd = exec.Command("flatpak-spawn", "--host", "bash", "-c", command)
	} else {
		if step.NeedRoot {
			cmd = exec.Command("pkexec", "bash", "-c", step.Command)
		} else {
			cmd = exec.Command("bash", "-c", step.Command)
		}
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

	// Início / Sincronização
	if strings.Contains(lower, "reading package") ||
		strings.Contains(lower, ":: synchronizing") ||
		strings.Contains(lower, "metadata expiration check") ||
		strings.Contains(lower, "fetching") ||
		strings.Contains(lower, "updating repository") {
		return 10
	}

	// Dependências / Resolução
	if strings.Contains(lower, "building dependency") ||
		strings.Contains(lower, "resolving dependencies") ||
		strings.Contains(lower, "checking for conflicts") ||
		strings.Contains(lower, "transaction check") {
		return 20
	}

	// Confirmação / Download
	if strings.Contains(lower, "the following packages will be upgraded") ||
		strings.Contains(lower, "packages to install") ||
		strings.Contains(lower, "total download size") ||
		strings.Contains(lower, "downloading") {
		return 35
	}

	// Instalação / Descompactação
	if strings.Contains(lower, "unpacking") ||
		strings.Contains(lower, "installing") ||
		strings.Contains(lower, "upgrading") ||
		strings.Contains(lower, "extracting") {
		return 50 + float64(lineCount%30)
	}

	// Configuração / Triggers
	if strings.Contains(lower, "setting up") ||
		strings.Contains(lower, "configuring") ||
		strings.Contains(lower, "running transaction") ||
		strings.Contains(lower, "verifying") {
		return 80 + float64(lineCount%15)
	}

	if strings.Contains(lower, "processing triggers") ||
		strings.Contains(lower, "cleanup") ||
		strings.Contains(lower, "complete!") {
		return 95
	}

	base := float64(lineCount) * 0.5
	if base > 95 {
		base = 95
	}
	return base
}

func (a *App) RunSystemAction(action string) error {
	var cmd *exec.Cmd

	inSandbox := false
	if _, err := os.Stat("/.flatpak-info"); err == nil {
		inSandbox = true
	}

	command := ""
	switch action {
	case "reboot":
		command = "systemctl reboot"
	case "shutdown":
		command = "systemctl poweroff"
	default:
		return fmt.Errorf("ação desconhecida: %s", action)
	}

	if inSandbox {
		cmd = exec.Command("flatpak-spawn", "--host", "pkexec", "bash", "-c", command)
	} else {
		cmd = exec.Command("pkexec", "bash", "-c", command)
	}

	return cmd.Run()
}
