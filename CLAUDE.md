# Mirakuta - CLAUDE.md

## Project Overview
- **Name:** Mirakuta (mira)
- **Description:** Modern dev environment for Windows, one line away.
- **Concept:** A CLI tool that makes Windows as comfortable as Mac for developers
- **Module:** github.com/mirakuta-dev/mirakuta

## Tech Stack
- **Language:** Go
- **CLI Framework:** Cobra + Viper
- **Config Format:** YAML
- **Bootstrap:** PowerShell

## Project Structure
```
mirakuta/
├── cmd/mirakuta/main.go
├── internal/
│   ├── cli/        # root.go, install.go, font.go, update.go
│   ├── config/
│   ├── preset/
│   └── runner/
├── scripts/        # PowerShell scripts
├── presets/        # minimal.yaml, all-rounder.yaml, etc.
├── install.ps1
├── go.mod
├── Makefile
└── CLAUDE.md
```

## Target Users
- Developers who moved from Mac to Windows
- Cross-platform developers using both Windows and Mac
- Beginner to intermediate developers who find Windows dev setup tedious

## Branch Strategy (Git Flow)
- **main** — Production releases only (tag required)
- **dev** — Integration branch (default working branch)
- **feature/xxx** — Feature development (branch from dev, merge to dev)
- **release/x.x.x** — Release preparation (QA before merging to main)
- **hotfix/xxx** — Emergency fixes (branch from main, merge to main + dev)

Rules:
- Direct commits to main are prohibited
- Feature branches should be scoped to a single unit of work
- Human must review diff before any merge

## Collaboration Workflow

### Code
- Claude Code handles code writing and `git add`
- Security-sensitive code (permissions, credentials, execution policies) must be reviewed by human before proceeding

### Commit Message Convention (Conventional Commits)
- feat: new feature
- fix: bug fix
- chore: build or config changes
- docs: documentation changes
- refactor: code restructuring
- test: add or update tests

## Coding Conventions
- Follow Go standard format (`gofmt`)
- Handle errors explicitly (no panic)
- Package names: lowercase, singular
- Comments: written in English

## Current Development Phase
Phase 5: `mirakuta check` command (environment diagnostics)
- Goal: check if Go, Git, PowerShell 7, Node.js, WSL are installed
- Print versions in a clean, readable format
- Next: install.ps1 bootstrap script

## Notes
- Windows-first (WSL integration included)
- PowerShell 7 as default shell
- Minimize external dependencies
- README final edit is done by human
- Build and test must be verified by human before merge

## Related
- Shared org-level config: `D:\Projects\mirakuta-dev\CLAUDE.md`
- Website repo: `mirakuta-web` (see parent CLAUDE.md for details)