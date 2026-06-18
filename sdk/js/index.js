import { AgentNotRunningError, PrintAgentError } from './errors.js';

const DEFAULT_BASE_URL = 'http://127.0.0.1:17890';
const DEFAULT_RELEASES_URL = 'https://pos.siyaho.com/downloads/print-agent/releases.json';
const DEFAULT_DOWNLOAD_URL =
  'https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe';

/**
 * @typedef {object} PrintBridgeOptions
 * @property {string} [baseUrl]
 * @property {number} [timeoutMs]
 * @property {string} [releasesUrl]
 * @property {string} [downloadUrl]
 */

export class PrintBridge {
  /**
   * @param {PrintBridgeOptions} [options]
   */
  constructor(options = {}) {
    this.baseUrl = String(options.baseUrl || DEFAULT_BASE_URL).replace(/\/+$/, '');
    this.timeoutMs = Number(options.timeoutMs) || 800;
    this.releasesUrl = options.releasesUrl || DEFAULT_RELEASES_URL;
    this.downloadUrl = options.downloadUrl || DEFAULT_DOWNLOAD_URL;
  }

  async ping() {
    try {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), this.timeoutMs);
      const res = await fetch(`${this.baseUrl}/v1/health`, { signal: controller.signal });
      clearTimeout(timer);
      if (!res.ok) return false;
      const json = await res.json();
      return Boolean(json?.ok);
    } catch {
      try {
        const controller = new AbortController();
        const timer = setTimeout(() => controller.abort(), this.timeoutMs);
        const res = await fetch(`${this.baseUrl}/health`, { signal: controller.signal });
        clearTimeout(timer);
        return res.ok;
      } catch {
        return false;
      }
    }
  }

  async getHealth() {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);
    try {
      const res = await fetch(`${this.baseUrl}/v1/health`, { signal: controller.signal });
      clearTimeout(timer);
      if (!res.ok) throw new AgentNotRunningError();
      return res.json();
    } catch (e) {
      clearTimeout(timer);
      if (e instanceof AgentNotRunningError) throw e;
      throw new AgentNotRunningError();
    }
  }

  /**
   * @param {{ target: { host?: string, ip?: string, port?: number }, data: { format: string, encoding?: string, content: string } }} job
   */
  async print(job) {
    const up = await this.ping();
    if (!up) throw new AgentNotRunningError();

    const res = await fetch(`${this.baseUrl}/v1/print`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(job),
    });

    if (!res.ok) {
      let msg = 'Print agent request failed';
      try {
        const json = await res.json();
        msg = json.error || msg;
      } catch {
        /* ignore */
      }
      throw new PrintAgentError(msg, res.status);
    }
    return res.json();
  }

  /**
   * Legacy DantSu payload used by Siyaho POS Web.
   * @param {{ ip: string, port?: number }} printer
   * @param {string} payload
   */
  async printDantsu(printer, payload) {
    const up = await this.ping();
    if (!up) throw new AgentNotRunningError();

    const res = await fetch(`${this.baseUrl}/print`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        printer: { ip: printer.ip, port: printer.port || 9100 },
        payload,
      }),
    });

    if (!res.ok) {
      let msg = 'Print agent request failed';
      try {
        const json = await res.json();
        msg = json.error || msg;
      } catch {
        /* ignore */
      }
      throw new PrintAgentError(msg, res.status);
    }
    return res.json();
  }

  /**
   * @param {'windows'|'mac'|'linux'|string} [platform]
   */
  getDownloadUrl(platform) {
    if (platform === 'windows' || !platform) return this.downloadUrl;
    return this.downloadUrl;
  }

  async fetchLatestRelease() {
    try {
      const res = await fetch(this.releasesUrl, { cache: 'no-store' });
      if (!res.ok) return null;
      return res.json();
    } catch {
      return null;
    }
  }

  async openDownload(platform) {
    const release = await this.fetchLatestRelease();
    const url = release?.windows || this.getDownloadUrl(platform);
    if (typeof window !== 'undefined') {
      window.open(url, '_blank', 'noopener,noreferrer');
    }
    return url;
  }

  /**
   * Poll until agent responds or timeout.
   * @param {{ intervalMs?: number, timeoutMs?: number, onTick?: () => void }} [options]
   */
  async pollUntilReady(options = {}) {
    const intervalMs = options.intervalMs ?? 2000;
    const timeoutMs = options.timeoutMs ?? 120000;
    const started = Date.now();

    while (Date.now() - started < timeoutMs) {
      if (options.onTick) options.onTick();
      if (await this.ping()) return true;
      await new Promise((r) => setTimeout(r, intervalMs));
    }
    return false;
  }
}

export function createPrintBridge(options) {
  return new PrintBridge(options);
}

export { AgentNotRunningError, PrintAgentError };
