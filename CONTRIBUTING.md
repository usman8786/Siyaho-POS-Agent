# Contributing to Siyaho Printer Agent

Thank you for contributing. This project is MIT-licensed and open to everyone.

## Getting started

### Prerequisites

- **Go 1.22+** — [go.dev/dl](https://go.dev/dl/)
- **Windows** — required to build and test Windows spooler printing (`agent/print/spooler_windows.go`)
- **Inno Setup 6** — optional, for building the Windows installer (`installer/windows/setup.iss`)

### Clone and run

```bash
git clone https://github.com/usman8786/Siyaho-POS-Agent.git
cd Siyaho-POS-Agent
go run ./agent
```

Open `http://127.0.0.1:17890/v1/health` — you should see JSON with `ok: true`.

### Build binary

```powershell
.\scripts\build.ps1
```

Output: `dist/SiyahoPrinterAgent.exe`

### JavaScript SDK (local)

```bash
cd sdk/js
npm pack
```

In a consuming app: `"siyaho-printer-agent": "file:../Siyaho-POS-Agent/sdk/js"`

## Project layout

```
agent/           Go HTTP server, print targets, CORS
sdk/js/          npm package (PrintBridge)
installer/       Inno Setup script and default config
docs/            Integration guide, OpenAPI, release notes
scripts/         build.ps1
.github/         CI and release workflows
```

## Making changes

1. **Fork** the repo and create a branch from `main`.
2. **Keep diffs focused** — one logical change per PR when possible.
3. **Match existing style** — standard Go formatting (`gofmt`), minimal comments unless logic is non-obvious.
4. **Update docs** if you change HTTP API, config, or install behavior.
5. **Bump version** only when preparing a release (see [docs/RELEASE.md](docs/RELEASE.md)); do not bump in unrelated PRs.

### Windows-only code

Print spooler code lives in `agent/print/spooler_windows.go` with a stub for other platforms. Non-Windows CI still compiles the agent; spooler endpoints return 501 on Linux/macOS.

## Pull requests

- Describe **what** changed and **why**.
- Link related issues if any.
- Confirm you ran `go build ./agent` (or CI will run on push).
- For API changes, update `docs/openapi.yaml` and `docs/integration.md`.

## Reporting bugs

Open a [GitHub issue](https://github.com/usman8786/Siyaho-POS-Agent/issues) with:

- Agent version (`GET /v1/health`)
- Windows version
- Printer type (network IP or Windows printer name)
- Steps to reproduce and relevant log lines from `%ProgramData%\Siyaho\PrinterAgent\agent.log`

## Security

See [SECURITY.md](SECURITY.md) for our security model and how to report vulnerabilities.

## Code of conduct

Be respectful and constructive. We welcome contributors of all experience levels.
