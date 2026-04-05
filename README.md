# Atualiza GO

[![Snap Store](https://snapcraft.io/pt/dark/install.svg)](https://snapcraft.io/atualiza-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

System update and maintenance interface for Linux, built with [Wails](https://wails.io) (Go + Webview).

**Atualiza GO** automatically detects your distribution and provides update options for system packages, Flatpak, and Snap.

## Features

- **Distro Detection**: Identifies host via `/etc/os-release` with binary fallback.
- **Multi-distro Support**: Native support for Debian, Arch, Fedora, openSUSE, Alpine, Void, and Solus families.
- **Bilingual UI**: Dynamic switching between Portuguese (PT-BR) and English (EN).
- **Package Management**: Supports APT, Pacman, DNF, Zypper, APK, XBPS, and EOPKG.
- **Flatpak & Snap**: Detects and updates apps in both ecosystems.
- **Resource Monitoring**: Real-time RAM and Disk usage tracking.
- **System Cleanup**: Safely removes temporary package caches.
- **Power Management**: Reboot and Shutdown with confirmation prompts.
- **Security**: Graphical authorization via `pkexec` (Polkit).

## Installation

### Native Packages (.deb / .rpm) (Recommended)
Download the `.deb` or `.rpm` installer from the [Releases](https://github.com/erascardsilva/atualiza_go/releases/latest) page. This version has full system access and menu integration.

```bash
# Debian / Ubuntu
sudo dpkg -i atualiza-go_linux_amd64.deb

# Fedora / openSUSE
sudo rpm -i atualiza-go_linux_amd64.rpm
```

### Portable Executable
Grant execution permissions to the downloaded binary:
```bash
chmod +x atualiza_go
./atualiza_go
```

### Snap Store (Sandbox / Read-Only)

[![Get it from the Snap Store](https://snapcraft.io/pt/dark/install.svg)](https://snapcraft.io/atualiza-go)

Install via Snap. Due to `strict` confinement, this version works as a **Read-Only Monitoring Dashboard** and cannot modify your host system:
```bash
sudo snap install atualiza-go
```

### Flatpak
Install from file:
```bash
flatpak install io.github.erascardsilva.atualizago.flatpak
```

### Build from Source

#### Prerequisites
- Go 1.21+
- Wails CLI v2
- GTK Libraries (e.g., `libwebkit2gtk-4.0-dev`)
- Polkit (`pkexec`)

#### Compilation
```bash
# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build
wails build
```

## Architecture

```
atualiza_go/
├── main.go          # Entry point (window, bindings)
├── app.go           # Core struct and lifecycle
├── distro.go        # Distro detection logic
├── sysinfo.go       # Telemetry (RAM/Disk)
├── updater.go       # Update engine and log streaming
├── frontend/
    ├── index.html   # UI Layout
    └── src/
        ├── translations.js # i18n dictionary
        ├── style.css       # CSS Styling
        └── main.js         # Frontend logic and bindings
```

## Supported Distributions

| Family | Examples | Package Manager |
|:---:|:---|:---:|
| **Debian** | Debian, Ubuntu, Mint, Pop!_OS, Zorin, Kali, MX | `apt` |
| **Arch** | Arch, Manjaro, EndeavourOS, Garuda, CachyOS | `pacman` |
| **Fedora** | Fedora, RHEL, CentOS, AlmaLinux, Rocky | `dnf` |
| **openSUSE**| Tumbleweed, Leap, SLES | `zypper` |
| **Alpine** | Alpine Linux | `apk` |
| **Void** | Void Linux | `xbps` |
| **Solus** | Solus OS | `eopkg` |

## Usage

1. Launch the application.
2. The dashboard automatically detects system resources and distribution.
3. In the **Update** tab, select tasks and click **Start Update**.
4. Authenticate via Polkit when prompted.

---

**Erasmo Cardoso**<br>
**Software Engineer | Electronics Specialist**