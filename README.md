<div align="center">

# 🐑 BAA 🐑

**The universal Linux package manager updater.**

One command. One password. Every package manager updated.

<br>
</div>

> [!IMPORTANT]
> **YOU HAVE TO UPDATE MANUALLY!** I broke the update script, so you have to update manually this update. In this version it's functioning again. Sorry!

<div align="center">
<br>
</div>

> [!IMPORTANT]
> **Shepherd is now BAA!** The project has been renamed. Please update your software to the latest version: all commands now start with `baa` instead of `shepherd`.

<div align="center">
<br>

[![Release](https://img.shields.io/github/v/release/jstreitb/baa?style=flat-square&color=b7bdf8&label=Release)](https://github.com/jstreitb/baa/releases)
[![License](https://img.shields.io/badge/License-MIT-a6da95?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-8aadf4?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Stars](https://img.shields.io/github/stars/jstreitb/baa?style=flat-square&color=eed49f)](https://github.com/jstreitb/baa/stargazers)

<br>

<img width="1200" height="600" alt="demo" src="https://github.com/user-attachments/assets/12710854-d638-436b-9de4-d6d9ba3cfab7" />

<br>
</div>

---

## Quick Install

```bash
curl -sSfL https://raw.githubusercontent.com/jstreitb/baa/main/install.sh | bash
```

> **Security Note:** You can verify the installer script before running it by downloading it first, reading it, or checking its checksum against the release assets.

To update BAA itself to the latest version, simply run:
```bash
baa --update
```

To uninstall BAA from your system, run:
```bash
baa --uninstall
```

Or build from source:

```bash
git clone https://github.com/jstreitb/baa.git
cd baa && make build
sudo mv baa /usr/local/bin/
```

## How It Works

```
$ baa
```

1. **Detects** installed package managers automatically
2. **Asks** for your sudo password once (masked, secure)
3. **Updates** everything sequentially with a live animated TUI
4. **Shows** a summary of what succeeded and failed

That's it. No config files, no complicated setup.

## Command-Line Flags

| Flag | Description |
|------|-------------|
| `--update` | Updates baa to the latest version |
| `--uninstall` | Uninstalls baa from your system |
| `--detect` | Shows detected package managers and exit |
| `--credits` | Shows the credits and exit |
| `--version` | Prints the current version |
| `--help` | Shows the help message |

## Supported Package Managers

| Manager | Detected via | Sudo | Non-Interactive Flags |
|---------|-------------|------|-----------------------|
| **apt** | `apt-get` | ✓ | `DEBIAN_FRONTEND=noninteractive`, `--force-confold` |
| **dnf** | `dnf` | ✓ | `-y` |
| **zypper** | `zypper` | ✓ | `--non-interactive`, `up` |
| **pacman** | `pacman` | ✓ | `--noconfirm` |
| **nix** | `nix-env` | ✗ | `-u` |
| **brew** | `brew` | ✗ | `upgrade` |
| **flatpak** | `flatpak` | ✗ | `-y --noninteractive` |
| **snap** | `snap` | ✓ | — |

> **Adding a new manager?** It's a single file. See [CONTRIBUTING.md](CONTRIBUTING.md#adding-a-new-package-manager).

## Features

<table>
<tr>
<td width="50%">

### 🔍 Auto-Detection
Probes `$PATH` for every supported package manager — no configuration needed.

### 🔒 Secure by Design
Password stored as `[]byte`, zeroed after use, delivered via `stdin` pipe. Never appears in process args or logs.

### 🤖 Fully Autonomous
Non-interactive flags handle everything. No "Press Y to continue" interruptions.

</td>
<td width="50%">

### 🛟 Interactive Fallback
If a manager *needs* user input (e.g. a dpkg config prompt), BAA suspends the TUI and drops you into the raw terminal. Resume is automatic.

### 🎨 Beautiful TUI
Catppuccin Macchiato theme, ASCII sheep animation, live command output, styled summary screen.

### ⚡ Fast & Lightweight
Single static binary. No runtime dependencies. ~3 MB.

</td>
</tr>
</table>

## Security

BAA takes security seriously:

| Measure | Implementation |
|---------|---------------|
| **Password storage** | `[]byte` — never a Go `string` |
| **Memory cleanup** | Explicit zeroing with `for i := range pw { pw[i] = 0 }` |
| **Password delivery** | Piped via `stdin` to `sudo -S` |
| **Process args** | Password never appears in CLI arguments |
| **Environment** | Env vars passed through `sudo -- env` to reach child processes |
| **Goroutine copies** | Each copy is independently zeroed after use |

## Architecture

```
baa/
├── cmd/baa/main.go           # Entry point
├── internal/
│   ├── ui/                        # Bubbletea TUI (model, views, styles, components)
│   ├── pkgmanager/                # PackageManager interface + implementations
│   └── utils/                     # Secure sudo execution, error sanitization
├── install.sh                     # One-line installer
├── .goreleaser.yml                # Release automation
└── Makefile
```

Built with the [Charmbracelet](https://charm.sh/) ecosystem:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — Components

<div align="center">

### 🐑 The Flock

BAA is a community-driven project. A huge thank you to the amazing people who help guide the herd!

<br>

<a href = "https://github.com/jstreitb/baa/graphs/contributors">
  <img src = "https://contrib.rocks/image?repo=jstreitb/baa"/>
</a>
<br>

## Contributing

Contributions are welcome! Please read the [Contributing Guide](CONTRIBUTING.md) before getting started.
 
## AI Assistance

<i>"If you can't beat them, join them."</i><br>
<p>This project was architected and developed with the support of <b>AI Assistance</b>.
</p>

## License

[MIT](LICENSE) — go build something great.

<br>
<hr>

<div align="center">
<br>
<sub>Made with 🐑 by <a href="https://github.com/jstreitb">jstreitb</a> and the flock.</sub>
<br>
</div>
