# Atualiza GO

[![Snap Store](https://badgen.net/badge/snap/store/blue)](https://snapcraft.io/atualiza-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern graphical user interface for Linux system maintenance and updates, built with [Wails](https://wails.io) (Go + Webview).

**Atualiza GO** automatically detects your Linux distribution and provides the appropriate update options, including system packages, Flatpak, and Snap ecosystems.

## Features

- **Robust Detection**: Identifies host distribution via `/etc/os-release` with falling back to package manager binaries.
- **Multi-distro Support**: Native support for Debian, Arch, Fedora, openSUSE, Alpine, Void, and Solus families.
- **Bilingual UI**: Dynamic language switching between Portuguese (PT-BR) and English (EN).
- **Comprehensive Updates**: Manages APT, Pacman, DNF, Zypper, APK, XBPS, and EOPKG.
- **Ecosystem Management**: Detects and updates Flatpak and Snap applications.
- **Telemetry Dashboard**: Real-time monitoring of RAM and Disk usage.
- **Safe Maintenance**: Built-in system cleanup to remove temporary package caches safely.
- **Power Management**: Reboot and Shutdown integration with confirmation prompts.
- **Security**: Graphical authorization via `pkexec` (Polkit).

## Installation

### Snap Store (Recommended)

Install the official Snap package with classic confinement for full system access:

```bash
sudo snap install atualiza-go --classic
```

### Flatpak

Download the latest release and install via:

```bash
flatpak install io.github.erascardsilva.atualizago.flatpak
```

### Build from Source

#### Prerequisites
- Go 1.21+
- Wails CLI v2
- WebKitGTK development files (e.g., `libwebkit2gtk-4.0-dev`)
- Polkit (for `pkexec`)

#### Compilation
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build the project
wails build
```

## Architecture

```
atualiza_go/
├── main.go          # Wails entry point (window, bindings)
├── app.go           # Core App struct and lifecycle
├── distro.go        # Multi-distro detection logic
├── sysinfo.go       # System telemetry (RAM/Disk)
├── updater.go       # Update engine and progress streaming
├── frontend/
│   ├── index.html   # High-fidelity UI layout
│   └── src/
│       ├── translations.js # i18n dictionary
│       ├── style.css       # Design System (Glassmorphism)
│       └── main.js         # Reactive UI logic and Wails bindings
├── build/
│   └── bin/
│       └── atualiza_go     # Final compiled binary
└── snap/
    └── snapcraft.yaml      # Snap packaging configuration
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
2. The dashboard will automatically detect your system resources and distribution.
3. Go to the **Update** tab, select the desired maintenance tasks, and click **Start Update**.
4. Authenticate via Polkit when prompted.

---

**Erasmo Cardoso - Dev**
