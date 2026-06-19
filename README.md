# Siyaho Printer Agent

Local bridge so web POS apps can send ESC/POS receipts to network thermal printers (TCP port 9100).

## Run (development)

```bash
go run ./agent
```

Listens on `http://127.0.0.1:17890`.

## Windows install (production)

Run `SiyahoPrinterAgent-Setup-latest.exe`. The installer:

- Installs to `Program Files\Siyaho\PrinterAgent`
- Registers a **scheduled task** (`SiyahoPrinterAgent`) to start on user logon
- Starts the agent **in the background** (no CMD window)
- Writes logs to `%ProgramData%\Siyaho\PrinterAgent\agent.log`

Reboot or log off/on is not required after install — the agent starts immediately via the scheduled task.

## Endpoints

- `GET /v1/health` — agent status
- `POST /v1/print` — v1 print job
- `POST /print` — legacy DantSu body `{ printer: { ip, port }, payload }`

## Build

```powershell
.\scripts\build.ps1
```

## Docs

- [Integration guide](./docs/integration.md)
- [OpenAPI](./docs/openapi.yaml)
- [Code signing](./docs/code-signing.md)
- [CDN deploy](./docs/deploy-cdn.md)

## License

MIT
