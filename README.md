# Atualiza GO

Ferramenta gráfica para atualização de sistemas Linux, construída com [Wails](https://wails.io) (Go + WebView).

Detecta automaticamente a distribuição e disponibiliza as opções de atualização adequadas para o sistema.

## Funcionalidades

- Detecção automática da distribuição Linux
- Suporte a Debian/Ubuntu e Arch/Manjaro (e derivados)
- Atualização de pacotes do sistema (APT / Pacman)
- Atualização de Flatpak e Snap (quando disponíveis)
- Barras de progresso em tempo real
- Reiniciar e desligar com confirmação
- Autenticação gráfica via `pkexec` (Polkit)

## Arquitetura

```
atualiza_go/
├── main.go          # Entry point Wails (janela, bindings)
├── app.go           # Struct App principal
├── distro.go        # Detecção de distro via /etc/os-release
├── updater.go       # Motor de atualização + streaming de progresso
├── frontend/
│   ├── index.html   # Layout (sidebar + páginas)
│   └── src/
│       ├── style.css  # Tema escuro, animações, barras de progresso
│       └── main.js    # Navegação, bindings Wails, eventos
├── build/
│   └── bin/
│       └── atualiza_go  # Binário compilado
└── wails.json       # Configuração Wails
```

### Fluxo

```
┌──────────────┐     bindings      ┌──────────────┐
│   Frontend   │ ◄══════════════► │   Backend Go  │
│  HTML/CSS/JS │                   │               │
│              │   EventsEmit()    │  distro.go    │
│  Progresso ◄─┤◄─────────────────┤  updater.go   │
│  Log output  │                   │  app.go       │
└──────────────┘                   └───────┬───────┘
                                           │
                                    pkexec │ bash -c
                                           ▼
                                   ┌───────────────┐
                                   │  apt / pacman  │
                                   │  flatpak / snap│
                                   └───────────────┘
```

### Backend

| Arquivo | Responsabilidade |
|---------|-----------------|
| `main.go` | Inicialização Wails, configuração de janela |
| `app.go` | Struct App, expõe `GetDistroInfo()` ao frontend |
| `distro.go` | Lê `/etc/os-release`, identifica família e gerenciador de pacotes |
| `updater.go` | Executa comandos via `pkexec`, streaming de stdout via `EventsEmit` |

### Frontend

| Arquivo | Responsabilidade |
|---------|-----------------|
| `index.html` | Layout com sidebar + 4 páginas (Início, Atualizar, Ações, Sobre) |
| `style.css` | Design system: tema escuro, glassmorphism, animações |
| `main.js` | Navegação SPA, bindings Wails, barras de progresso, modais |

## Distros Suportadas

| Família | Distros | Gerenciador |
|---------|---------|-------------|
| Debian | Debian, Ubuntu, Mint, Pop!_OS, Elementary, Zorin, Kali, MX, Deepin | `apt` |
| Arch | Arch, Manjaro, EndeavourOS, Garuda, ArcoLinux, Artix, CachyOS | `pacman` |

## Requisitos

- Go 1.21+
- Wails CLI v2
- WebKitGTK (libwebkit2gtk-4.0-dev ou equivalente)
- Polkit (para `pkexec`)

## Build

```bash
# Instalar Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Compilar
wails build

# Ou modo desenvolvimento (hot reload)
wails dev
```

## Uso

```bash
./build/bin/atualiza_go
```

O app abre e detecta seu sistema automaticamente. Selecione as atualizações desejadas e clique em "Iniciar Atualização". O `pkexec` pedirá sua senha quando necessário.

---

**Erasmo Cardoso - Dev**
