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

var fedoraIDs = map[string]bool{
	"fedora": true, "rhel": true, "centos": true,
	"almalinux": true, "rocky": true, "nixos": true,
}

var suseIDs = map[string]bool{
	"opensuse": true, "opensuse-tumbleweed": true,
	"opensuse-leap": true, "sles": true,
}

func detectDistro() DistroInfo {
	info := DistroInfo{}

	// Se estiver em Flatpak, tenta ler o os-release do host
	osReleasePath := "/etc/os-release"
	inSandbox := false
	if _, err := os.Stat("/.flatpak-info"); err == nil {
		inSandbox = true
		if _, err := os.Stat("/run/host/etc/os-release"); err == nil {
			osReleasePath = "/run/host/etc/os-release"
		}
	}

	file, err := os.Open(osReleasePath)
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

	// Detecção baseada no ID e ID_LIKE
	switch {
	case debianIDs[info.ID] || strings.Contains(idLike, "debian") || strings.Contains(idLike, "ubuntu"):
		info.Family = "debian"
		info.PackageManager = "apt"
	case archIDs[info.ID] || strings.Contains(idLike, "arch"):
		info.Family = "arch"
		info.PackageManager = "pacman"
	case fedoraIDs[info.ID] || strings.Contains(idLike, "fedora") || strings.Contains(idLike, "rhel") || strings.Contains(idLike, "centos"):
		info.Family = "fedora"
		info.PackageManager = "dnf"
	case suseIDs[info.ID] || strings.Contains(idLike, "suse"):
		info.Family = "suse"
		info.PackageManager = "zypper"
	case info.ID == "alpine":
		info.Family = "alpine"
		info.PackageManager = "apk"
	case info.ID == "void":
		info.Family = "void"
		info.PackageManager = "xbps"
	case info.ID == "solus":
		info.Family = "solus"
		info.PackageManager = "eopkg"
	default:
		// Fallback: Tenta detectar pelo binário do gerenciador de pacotes
		info.Family = "unknown"
		info.PackageManager = "unknown"

		managers := map[string]string{
			"apt":    "debian",
			"pacman": "arch",
			"dnf":    "fedora",
			"zypper": "suse",
			"apk":    "alpine",
			"xbps-install": "void",
			"eopkg":  "solus",
		}

		for cmd, fam := range managers {
			if inSandbox {
				if _, err := os.Stat("/run/host/usr/bin/" + cmd); err == nil {
					info.PackageManager = cmd
					info.Family = fam
					break
				}
			} else {
				if _, err := exec.LookPath(cmd); err == nil {
					info.PackageManager = cmd
					info.Family = fam
					break
				}
			}
		}
	}

	if inSandbox {
		if _, err := os.Stat("/run/host/usr/bin/flatpak"); err == nil {
			info.HasFlatpak = true
		}
		if _, err := os.Stat("/run/host/usr/bin/snap"); err == nil {
			info.HasSnap = true
		}
	} else {
		if _, err := exec.LookPath("flatpak"); err == nil {
			info.HasFlatpak = true
		}
		if _, err := exec.LookPath("snap"); err == nil {
			info.HasSnap = true
		}
	}

	return info
}
