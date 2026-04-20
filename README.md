# mirakuta

[![Release](https://img.shields.io/github/v/release/mirakuta-dev/mirakuta?color=blue)](https://github.com/mirakuta-dev/mirakuta/releases)
[![License](https://img.shields.io/github/license/mirakuta-dev/mirakuta)](./LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/mirakuta-dev/mirakuta)](./go.mod)
[![Platform](https://img.shields.io/badge/platform-windows%20amd64%20%7C%20arm64-0078D6)](https://github.com/mirakuta-dev/mirakuta/releases/latest)

**Modern dev environment for Windows — one line away.**

`mirakuta` is a CLI that installs a complete developer setup on Windows in
a single command. Pick a preset, answer a short wizard, or point it at a
YAML file — it runs `winget`, configures Git, and installs VS Code
extensions so you can start coding, not scripting.

## Install

```powershell
irm https://raw.githubusercontent.com/mirakuta-dev/mirakuta/main/install.ps1 | iex
```

The installer:

- Detects your architecture (AMD64 / ARM64)
- Downloads the signed release binary and verifies its SHA-256
- Installs to `%LOCALAPPDATA%\Programs\mirakuta` (no admin required)
- Adds the install directory to your user `PATH`

Verify: open a new terminal and run `mirakuta version`.

## Quick start

Three ways to install a dev environment:

```powershell
# 1. Interactive wizard — guided setup, saves your choices to ~/.mirakuta/profile.yaml
mirakuta install

# 2. Built-in preset — pick a role and go
mirakuta install backend

# 3. External YAML — share a team setup
mirakuta install --file ./team-setup.yaml
```

Useful flags:

- `--dry-run` — print the commands without executing them
- `--verbose` — stream output from each step

## Presets

| Preset        | Extends               | Installs |
|---------------|-----------------------|----------|
| `minimal`     | —                     | Git, Windows Terminal, PowerShell 7, VS Code |
| `backend`     | `minimal`             | + Go, Node.js, Docker Desktop, Postman |
| `frontend`    | `minimal`             | + Node.js, pnpm, Chrome, Figma |
| `fullstack`   | `backend`, `frontend` | everything from backend and frontend |
| `data`        | `minimal`             | + Python, uv, Jupyter extensions |
| `all-rounder` | `fullstack`, `data`   | everything — good default for generalists |

All presets live in [`presets/`](./presets) and are embedded into the
binary via `go:embed`. They share the same YAML schema you can use for
your own `--file` setups.

## Commands

| Command               | What it does |
|-----------------------|--------------|
| `mirakuta install`    | Run a preset (see above) |
| `mirakuta check`      | Diagnose your environment — see below |
| `mirakuta version`    | Print the installed version |

Run `mirakuta <command> --help` for details.

### `mirakuta check`

Scans your system for common developer tooling and prints a pass/fail
table. Useful before running `install`, or when something broke and you
want to know what's actually on the box.

Probed entries:

- **Git**, **Go**, **Node.js**, **npm**, **winget** — detected via each
  tool's own `--version` output
- **WSL2** — enabled / WSL1 / not installed
- **`GOPATH`** — required when Go is installed
- **`GOROOT`** — optional, reported if set

Example output:

```
Mirakuta Environment Check
────────────────────────────────────────
  [✓] Git        git version 2.46.0.windows.1
  [✓] Go         go version go1.26.2 windows/amd64
  [✓] Node.js    v22.11.0
  [✓] npm        10.9.0
  [✓] winget     v1.9.25200
  [✓] WSL2       enabled
  [✓] GOPATH     C:\Users\you\go
  [✓] GOROOT     (not set, optional)
────────────────────────────────────────
All checks passed.
```

Exits with code `1` if any required check fails.

## How it works

`mirakuta install` resolves a preset (with `extends:` merged), converts it
into an ordered plan (package installs → Git config → VS Code extensions),
and runs each step via `winget` / `git` / `code`. Failing steps don't
abort the run — you get a summary of what succeeded, failed, and skipped.

The schema is intentionally small: **data only, no arbitrary scripts**.
That makes `--file` safe to share across a team.

## Build from source

```powershell
git clone https://github.com/mirakuta-dev/mirakuta
cd mirakuta
make build       # produces mirakuta.exe
make test        # run unit tests
make cross       # build amd64 and arm64 into dist/
```

Requires Go (see `go.mod` for the minimum version).

## Requirements

- Windows 10 (21H2+) or Windows 11
- PowerShell 5.1 or later (PowerShell 7+ recommended)
- `winget` (ships with modern Windows; update via Microsoft Store if
  missing)

## License

MIT — see [LICENSE](./LICENSE).
