// Atualiza GO - Telemetria de Sistema
// Erasmo Cardoso - Dev

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SystemStats struct {
	MemTotal     uint64  `json:"memTotal"`
	MemUsed      uint64  `json:"memUsed"`
	MemPercent   float64 `json:"memPercent"`
	DiskTotal    string  `json:"diskTotal"`
	DiskUsed     string  `json:"diskUsed"`
	DiskPercent  float64 `json:"diskPercent"`
	DiskMessage  string  `json:"diskMessage"`
}

func (a *App) GetSystemStats() SystemStats {
	stats := SystemStats{}

	// RAM Info
	stats.MemTotal, stats.MemUsed, stats.MemPercent = getMemInfo()

	// Disk Info
	stats.DiskTotal, stats.DiskUsed, stats.DiskPercent, stats.DiskMessage = getDiskInfo()

	return stats
}

func getMemInfo() (uint64, uint64, float64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0
	}
	defer file.Close()

	var total, free, available uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			total = val
		case "MemFree:":
			free = val
		case "MemAvailable:":
			available = val
		}
	}

	if total == 0 {
		return 0, 0, 0
	}

	// MemAvailable é mais preciso que MemFree
	used := total - available
	if available == 0 {
		used = total - free
	}

	percent := (float64(used) / float64(total)) * 100
	return total / 1024, used / 1024, percent // Convert to MB
}

func getDiskInfo() (string, string, float64, string) {
	// Usamos df para facilitar o parse em vez de syscall.Statfs (que requer caminhos e mountpoints)
	// df / --output=size,used,pcent
	out, err := exec.Command("df", "/", "--output=size,used,pcent").Output()
	if err != nil {
		return "—", "—", 0, "Erro ao ler disco"
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return "—", "—", 0, "Formato desconhecido"
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return "—", "—", 0, "Dados incompletos"
	}

	sizeKB, _ := strconv.ParseFloat(fields[0], 64)
	usedKB, _ := strconv.ParseFloat(fields[1], 64)
	pctStr := strings.TrimSuffix(fields[2], "%")
	percent, _ := strconv.ParseFloat(pctStr, 64)

	sizeGB := sizeKB / 1024 / 1024
	usedGB := usedKB / 1024 / 1024

	return fmt.Sprintf("%.1f GB", sizeGB), fmt.Sprintf("%.1f GB", usedGB), percent, ""
}
