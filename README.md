# Atualiza GO

[![Snap Store](https://snapcraft.io/pt/dark/install.svg)](https://snapcraft.io/atualiza-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Interface para atualização e manutenção de sistemas Linux desenvolvida com [Wails](https://wails.io) (Go + Webview).

O **Atualiza GO** detecta automaticamente a distribuição e oferece opções de atualização para pacotes do sistema, Flatpak e Snap.

## Funcionalidades

- **Detecção de Sistema**: Identifica a distro via `/etc/os-release` com fallback para binários do gestor de pacotes.
- **Suporte Multi-distro**: Compatível com famílias Debian, Arch, Fedora, openSUSE, Alpine, Void e Solus.
- **Interface Bilíngue**: Troca dinâmica entre Português (PT-BR) e Inglês (EN).
- **Gerenciamento de Pacotes**: Suporte para APT, Pacman, DNF, Zypper, APK, XBPS e EOPKG.
- **Flatpak e Snap**: Detecta e atualiza aplicações em ambos os ecossistemas.
- **Monitoramento**: Acompanhamento em tempo real de uso de RAM e Disco.
- **Limpeza de Sistema**: Remoção segura de caches temporários de pacotes.
- **Energia**: Opções de Reiniciar e Desligar com confirmação.
- **Segurança**: Autorização gráfica via `pkexec` (Polkit).

## Instalação

### Pacotes Nativos (.deb / .rpm) (Recomendado)
Baixe o instalador `.deb` ou `.rpm` na página de [Releases](https://github.com/erascardsilva/atualiza_go/releases/latest). Esta versão possui acesso total ao sistema e integração com o menu.

```bash
# Debian / Ubuntu
sudo dpkg -i atualiza-go_linux_amd64.deb

# Fedora / openSUSE
sudo rpm -i atualiza-go_linux_amd64.rpm
```

### Executável Portátil
Dê permissão de execução ao binário baixado:
```bash
chmod +x atualiza_go
./atualiza_go
```

### Snap Store (Sandbox / Read-Only)

[![Disponível na Snap Store](https://snapcraft.io/pt/dark/install.svg)](https://snapcraft.io/atualiza-go)

Instale via Snap. Devido ao confinamento `strict`, esta versão funciona apenas como **Painel de Monitoramento (Read-Only)**, sem permissão para modificar o sistema:
```bash
sudo snap install atualiza-go
```

### Flatpak
Instalação via arquivo:
```bash
flatpak install io.github.erascardsilva.atualizago.flatpak
```

### Build do Código

#### Pré-requisitos
- Go 1.21+
- Wails CLI v2
- Bibliotecas GTK (ex: `libwebkit2gtk-4.0-dev`)
- Polkit (`pkexec`)

#### Compilação
```bash
# Instalar Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Compilar
wails build
```

## Arquitetura

```
atualiza_go/
├── main.go          # Ponto de entrada (janela, bindings)
├── app.go           # Estrutura principal e ciclo de vida
├── distro.go        # Lógica de detecção de distro
├── sysinfo.go       # Telemetria (RAM/Disco)
├── updater.go       # Motor de atualização e streaming de logs
├── frontend/
    ├── index.html   # Layout UI
    └── src/
        ├── translations.js # Dicionário i18n
        ├── style.css       # Estilização CSS
        └── main.js         # Lógica do frontend e bindings
```

## Distribuições Suportadas

| Família | Exemplos | Gestor de Pacotes |
|:---:|:---|:---:|
| **Debian** | Debian, Ubuntu, Mint, Pop!_OS, Zorin, Kali, MX | `apt` |
| **Arch** | Arch, Manjaro, EndeavourOS, Garuda, CachyOS | `pacman` |
| **Fedora** | Fedora, RHEL, CentOS, AlmaLinux, Rocky | `dnf` |
| **openSUSE**| Tumbleweed, Leap, SLES | `zypper` |
| **Alpine** | Alpine Linux | `apk` |
| **Void** | Void Linux | `xbps` |
| **Solus** | Solus OS | `eopkg` |

## Uso

1. Inicie a aplicação.
2. O dashboard detecta os recursos e a distribuição automaticamente.
3. Na aba **Update**, selecione as tarefas e clique em **Start Update**.
4. Autentique via Polkit quando solicitado.

---

**Erasmo Cardoso - Dev**
