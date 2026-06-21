# Security Policy

## Model

Siyaho Printer Agent is a **localhost-only** HTTP service:

- Binds to **127.0.0.1** (not `0.0.0.0`)
- Accepts print requests from browsers and apps on **the same machine**
- Reflects the request `Origin` header for CORS so any web POS domain can integrate without a license

### What this means

Any website open in the user's browser **on that PC** could theoretically send print jobs to the agent if the user visits a malicious page while the agent is running. This is the same trust model as other local helper apps (e.g. localhost dev servers, desktop companion apps).

**Mitigations:**

- Agent is not exposed on the network
- Users should install the agent only from trusted sources ([pos.siyaho.com](https://pos.siyaho.com/downloads) or official GitHub releases)
- POS vendors should serve their web app over **HTTPS**

## Supported versions

Security fixes are applied to the latest release on `main`. Download the latest installer from [GitHub Releases](https://github.com/usman8786/Siyaho-POS-Agent/releases) or [pos.siyaho.com/downloads](https://pos.siyaho.com/downloads).

## Reporting a vulnerability

If you discover a security issue, please **do not** open a public GitHub issue with exploit details.

Instead:

1. Email or contact the maintainers privately (via GitHub Security Advisories if enabled, or the contact on the repository).
2. Include steps to reproduce and impact assessment.
3. Allow reasonable time to fix before public disclosure.

We appreciate responsible disclosure.

## Out of scope

- Issues that require physical access to an unlocked cashier PC
- Social engineering of store staff
- Network attacks against printers on LAN (agent does not route LAN traffic from remote machines)

## Recommendations for integrators

- Do not reconfigure the agent to listen on `0.0.0.0`
- Pin installer downloads to official URLs
- Validate print payloads in your POS before sending to the agent
