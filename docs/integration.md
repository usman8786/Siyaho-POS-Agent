# Siyaho Printer Agent — Integration Guide

Connect **any web POS** to network thermal printers (TCP port 9100) or **Windows USB/local printers** via a local agent on the cashier PC.

No license keys or origin approval required.

## Quick start

1. Install **Siyaho Printer Agent** on the Windows PC connected to the printer (one-time; runs in the background, no CMD window).
2. Confirm health: `GET http://127.0.0.1:17890/v1/health`
3. From your web app (same machine), send print jobs to `POST /v1/print` or legacy `POST /print`.

**Network printer:** `{ "printer": { "ip": "192.168.1.50", "port": 9100 }, "payload": "..." }`

**Windows USB printer:** `{ "printer": { "type": "windows", "name": "Generic / Text Only Speedex" }, "payload": "..." }`

List installed Windows printers: `GET http://127.0.0.1:17890/v1/printers`

Logs: `%ProgramData%\Siyaho\PrinterAgent\agent.log`

Use an ESC/POS driver when possible; "Generic / Text Only" may produce garbled receipts.

## JavaScript SDK

```bash
npm install siyaho-printer-agent
# or file:../Siyaho-POS-Agent/sdk/js during development
```

```js
import { PrintBridge } from 'siyaho-printer-agent';

const bridge = new PrintBridge({
  baseUrl: 'http://127.0.0.1:17890',
});

if (await bridge.ping()) {
  const printers = await bridge.listPrinters();
  await bridge.printDantsu({ type: 'windows', name: 'Generic / Text Only Speedex' }, '[L]Hello\n[C]<b>Total</b> 10.00');
  await bridge.printDantsu({ ip: '192.168.1.50', port: 9100 }, '[L]Hello\n[C]<b>Total</b> 10.00');
}
```

## Settings UI pattern

- Show agent status (running / not running)
- If not running: **Download agent** → install → **Check again** (poll `bridge.pollUntilReady()`)
- Download URL: `https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe`

## Configuration

Optional file: `%ProgramData%\Siyaho\PrinterAgent\config.json`

```json
{
  "port": 17890
}
```

Legacy configs with `licenseKey` or `apiBaseUrl` are ignored; only `port` is used.

## CORS

The agent reflects the browser `Origin` header so **any web POS domain** can call the API from the same PC. The service listens on **127.0.0.1 only** — not exposed to the LAN.

## API reference

See [openapi.yaml](./openapi.yaml).

## Security

- Agent binds `127.0.0.1` only — never expose on `0.0.0.0`
- Use HTTPS for your POS web app in production
- See [SECURITY.md](../SECURITY.md) for the full security model

## Source

https://github.com/usman8786/Siyaho-POS-Agent
