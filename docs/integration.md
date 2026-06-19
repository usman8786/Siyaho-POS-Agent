# Siyaho Printer Agent — Integration Guide

Connect your web POS to network thermal printers (TCP port 9100) via a local agent on the cashier PC.

## Quick start

1. Install **Siyaho Printer Agent** on the Windows PC connected to the LAN printer (one-time; runs in the background, no CMD window).
2. Confirm health: `GET http://127.0.0.1:17890/v1/health`
3. From your web app (same machine), send print jobs to `POST /v1/print` or legacy `POST /print`.

Logs: `%ProgramData%\Siyaho\PrinterAgent\agent.log`

## JavaScript SDK

```bash
npm install @siyaho/printer-agent
# or file:../Siyaho-POS-Agent/sdk/js during development
```

```js
import { PrintBridge } from '@siyaho/printer-agent';

const bridge = new PrintBridge({
  baseUrl: 'http://127.0.0.1:17890',
});

if (await bridge.ping()) {
  await bridge.printDantsu({ ip: '192.168.1.50', port: 9100 }, '[L]Hello\n[C]<b>Total</b> 10.00');
}
```

## Settings UI pattern

- Show agent status (running / not running)
- If not running: **Download agent** → install → **Check again** (poll `bridge.pollUntilReady()`)
- Download URL: `https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe`

## Third-party POS vendors

Submit an access request at `https://pos.siyaho.com/print-agent-request`.

After approval you receive:

- `licenseKey` — add to `%ProgramData%\Siyaho\PrinterAgent\config.json`
- `allowedOrigins` — your POS web origin (e.g. `https://app.yourpos.com`)

```json
{
  "port": 17890,
  "licenseKey": "spa_live_...",
  "apiBaseUrl": "https://siyaho.com"
}
```

## CORS

The agent allows:

- `*.siyaho.com` origins (first-party)
- `http://localhost:5173` and `http://localhost:5000` (development)
- Approved third-party origins after license validation

## API reference

See [openapi.yaml](./openapi.yaml).

## Security

- Agent binds `127.0.0.1` only
- Never expose the agent on `0.0.0.0`
- Use HTTPS for your POS web app in production
