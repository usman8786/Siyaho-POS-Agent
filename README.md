# Siyaho Printer Agent

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/usman8786/Siyaho-POS-Agent/actions/workflows/ci.yml/badge.svg)](https://github.com/usman8786/Siyaho-POS-Agent/actions/workflows/ci.yml)

Open-source local bridge so **any web POS** can send ESC/POS receipts to **network thermal printers** (TCP port 9100) or **Windows USB/local printers** via the spooler.

No license keys. No origin approval. Install the agent on the cashier PC and call `http://127.0.0.1:17890` from your browser app on the same machine.

## Architecture

```
Browser (any origin)  →  http://127.0.0.1:17890  →  Network printer / Windows spooler
```

The agent binds **127.0.0.1 only** — it is not reachable from other devices on the LAN. See [SECURITY.md](SECURITY.md).

## Run (development)

```bash
go run ./agent
```

Listens on `http://127.0.0.1:17890`.

## Windows install (production)

Download: [SiyahoPrinterAgent-Setup-latest.exe](https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe)

The installer:

- Installs to `Program Files\Siyaho\PrinterAgent`
- Registers a **scheduled task** (`SiyahoPrinterAgent`) to start on user logon
- Starts the agent **in the background** (no CMD window)
- Writes logs to `%ProgramData%\Siyaho\PrinterAgent\agent.log`

Config (optional): `%ProgramData%\Siyaho\PrinterAgent\config.json`

```json
{ "port": 17890 }
```

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /v1/health` | Agent status and version |
| `GET /v1/printers` | List Windows installed printers (Windows only) |
| `POST /v1/print` | v1 print job |
| `POST /print` | Legacy DantSu body `{ printer: { ip, port }, payload }` |

## JavaScript SDK

```bash
npm install siyaho-printer-agent
```

See [sdk/js/README.md](sdk/js/README.md).

## Build

```powershell
.\scripts\build.ps1
```

Requires Go 1.22+ and (for the installer) Inno Setup on Windows.

## Documentation

| Doc | Description |
|-----|-------------|
| [Integration guide](docs/integration.md) | Connect your web POS |
| [OpenAPI](docs/openapi.yaml) | HTTP API reference |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Development and pull requests |
| [docs/RELEASE.md](docs/RELEASE.md) | Versioning and publishing releases |
| [SECURITY.md](SECURITY.md) | Security model and reporting |
| [CHANGELOG.md](CHANGELOG.md) | Release history |
| [Code signing](docs/code-signing.md) | Optional Windows code signing |
| [CDN deploy](docs/deploy-cdn.md) | Hosting installer artifacts |

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.

## License

MIT — see [LICENSE](LICENSE).
