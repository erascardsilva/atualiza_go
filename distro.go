// Atualiza GO - Detecção de Distribuição Linux
// Erasmo Cardoso - Dev

package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

type DistroInfo struct {
	Name           string `json:"name"`
	ID             string `json:"id"`
	Version        string `json:"version"`
	Family         string `json:"family"`
	PackageManager string `json:"packageManager"`
	HasFlatpak     bool   `json:"hasFlatpak"`
	HasSnap        bool   `json:"hasSnap"`
}

var debianIDs = map[string]bool{
	"debian": true, "ubuntu": true, "linuxmint": true,
	"pop": true, "elementary": true, "zorin": true,
	"kali": true, "mx": true, "lmde": true,
	"neon": true, "deepin": true,
}

var archIDs = map[string]bool{
	"arch": true, "manjaro": true, "endeavouros": true,
	"garuda": true, "arcolinux": true, "artix": true,
	"cachyos": true,
}

func detectDistro() DistroInfo {
	info := DistroInfo{}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		info.Name = "Desconhecido"
		info.ID = "unknown"
		info.Family = "unknown"
		return info
	}
	defer file.Close()

	data := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := strings.Trim(parts[1], "\"")
			data[key] = val
		}
	}

	info.Name = data["PRETTY_NAME"]
	info.ID = strings.ToLower(data["ID"])
	info.Version = data["VERSION_ID"]

	idLike := strings.ToLower(data["ID_LIKE"])

	switch {
	case debianIDs[info.ID] || strings.Contains(idLike, "debian") || strings.Contains(idLike, "ubuntu"):
		info.Family = "debian"
		info.PackageManager = "apt"
	case archIDs[info.ID] || strings.Contains(idLike, "arch"):
		info.Family = "arch"
		info.PackageManager = "pacman"
	default:
		info.Family = "unknown"
		info.PackageManager = "unknown"
	}

	if _, err := exec.LookPath("flatpak"); err == nil {
		info.HasFlatpak = true
	}
	if _, err := exec.LookPath("snap"); err == nil {
		info.HasSnap = true
	}

	return info
}
