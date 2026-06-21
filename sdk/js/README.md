# siyaho-printer-agent

JavaScript client for **[Siyaho Printer Agent](https://github.com/usman8786/Siyaho-POS-Agent)** — a local Windows bridge that lets web POS apps print ESC/POS receipts to **network thermal printers** (TCP port 9100) or **Windows USB/local printers** without the browser print dialog.

This npm package is the **browser SDK**. The actual print service is a small **Go agent** that runs on the cashier PC. You need **both**: install the agent `.exe` on each Windows machine, then use this package (or plain `fetch`) from your web app.

## Download the Windows agent (required)

Install on every PC that prints to a LAN or USB thermal printer:

| Resource | URL |
|----------|-----|
| **Latest installer** | https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe |
| **Version manifest** | https://pos.siyaho.com/downloads/print-agent/releases.json |
| **Downloads page** | https://pos.siyaho.com/downloads |

### Install steps

1. Download and run `SiyahoPrinterAgent-Setup-latest.exe` on the cashier PC.
2. If Windows SmartScreen shows **“Unknown publisher”**, click **More info** → **Run anyway** (installer is not code-signed yet).
3. Approve the administrator (UAC) prompt.
4. The agent runs in the **background** (no CMD window) and auto-starts on user logon.
5. Confirm it is running: open `http://127.0.0.1:17890/v1/health` in a browser on that PC.

Logs: `%ProgramData%\Siyaho\PrinterAgent\agent.log`

## Install the npm package

```bash
npm install siyaho-printer-agent
```

## Quick start

```js
import { PrintBridge } from 'siyaho-printer-agent';

const bridge = new PrintBridge({
  baseUrl: 'http://127.0.0.1:17890', // default
});

// Check agent on this PC
if (await bridge.ping()) {
  const health = await bridge.getHealth();
  console.log(health.version); // e.g. "1.2.0"
}

// List Windows printers (USB, etc.)
const { printers } = await bridge.listPrinters();

// Network printer
await bridge.printDantsu(
  { ip: '192.168.1.50', port: 9100 },
  '[L]Hello\n[C]<b>Total</b> 10.00',
);

// Windows USB / installed printer (exact name from listPrinters)
await bridge.printDantsu(
  { type: 'windows', name: 'Generic / Text Only Speedex' },
  '[L]Hello\n[C]<b>Total</b> 10.00',
);

// v1 print API
await bridge.print({
  target: { type: 'tcp', host: '192.168.1.50', port: 9100 },
  data: { format: 'escpos', encoding: 'utf-8', content: '\x1b@\nHello\n' },
});
```

## Typical POS integration flow

1. **Ping** the agent (`bridge.ping()`).
2. If not running, show a **Download agent** button → `bridge.openDownload()` or link to the installer URL above.
3. After install, **poll** until ready: `bridge.pollUntilReady({ timeoutMs: 120000 })`.
4. Send print jobs to the printer IP configured in your POS.

```js
import { PrintBridge, AgentNotRunningError } from 'siyaho-printer-agent';

const bridge = new PrintBridge();

async function ensureAgent() {
  if (await bridge.ping()) return true;

  const url = await bridge.openDownload(); // opens installer in new tab
  return bridge.pollUntilReady({ intervalMs: 2000, timeoutMs: 120000 });
}

try {
  await ensureAgent();
  await bridge.printDantsu({ ip: printerIp, port: 9100 }, receiptPayload);
} catch (e) {
  if (e instanceof AgentNotRunningError) {
    alert('Install Siyaho Printer Agent on this PC first.');
  }
}
```

## API overview

### `PrintBridge(options?)`

| Option | Default | Description |
|--------|---------|-------------|
| `baseUrl` | `http://127.0.0.1:17890` | Local agent URL |
| `timeoutMs` | `800` | Health/ping timeout |
| `releasesUrl` | `https://pos.siyaho.com/downloads/print-agent/releases.json` | Version manifest |
| `downloadUrl` | `https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe` | Windows installer |

### Methods

| Method | Description |
|--------|-------------|
| `ping()` | `true` if agent responds on `/v1/health` or legacy `/health` |
| `getHealth()` | Full health JSON (`version`, `service`, `platform`, …) |
| `listPrinters()` | `GET /v1/printers` — Windows installed printers |
| `print(job)` | `POST /v1/print` |
| `printDantsu(printer, payload)` | `POST /print` — network `{ ip, port }` or Windows `{ type: 'windows', name }` |
| `getDownloadUrl(platform?)` | Installer URL for Windows |
| `fetchLatestRelease()` | Parse `releases.json` |
| `openDownload(platform?)` | Open installer in browser; returns URL |
| `pollUntilReady(options?)` | Poll `ping()` until agent is up or timeout |

### Errors

- `AgentNotRunningError` — agent not reachable on localhost
- `PrintAgentError` — agent returned a non-OK print response

### Helper

```js
import { createPrintBridge } from 'siyaho-printer-agent';

const bridge = createPrintBridge({ baseUrl: 'http://127.0.0.1:17890' });
```

## Agent HTTP API (localhost)

The Go agent listens on **`127.0.0.1:17890`** only (not exposed to the network).

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/printers` | GET | List Windows installed printers |
| `/v1/health` | GET | Status and version |
| `/v1/print` | POST | v1 print job |
| `/print` | POST | Legacy `{ printer: { ip, port }, payload }` |
| `/health` | GET | Legacy health (compat) |

OpenAPI and full docs: [github.com/usman8786/Siyaho-POS-Agent](https://github.com/usman8786/Siyaho-POS-Agent/tree/main/docs)

## CORS

The agent reflects your browser `Origin` header — **any web POS domain** works without registration. The service listens on **127.0.0.1 only** on the same PC as your browser.

## Security

- Agent binds **127.0.0.1** only — do not expose it on `0.0.0.0`
- Serve your POS web app over **HTTPS** in production
- Download the installer only from **pos.siyaho.com**

## Repository

- **Source & agent binary builds:** https://github.com/usman8786/Siyaho-POS-Agent
- **Integration guide:** https://github.com/usman8786/Siyaho-POS-Agent/blob/main/docs/integration.md
- **Used by:** [Siyaho POS Web](https://pos.siyaho.com)

## License

MIT
