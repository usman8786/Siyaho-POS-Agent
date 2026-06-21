# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.2.0] - 2026-06-16

### Removed

- License key requirement and origin allowlist — the agent is fully open source with no approval flow
- `licenseKey` and `apiBaseUrl` from config (config is now `{ "port": 17890 }` only)
- `license` object from `/v1/health` response
- CORS 403 "Origin not allowed" — any web origin is accepted (agent remains localhost-only)

### Changed

- CORS middleware reflects the request `Origin` header for all allowed browser calls
- Documentation: CONTRIBUTING, RELEASE, SECURITY, updated integration guide

## [1.1.0]

- Windows USB/spooler printing (`GET /v1/printers`, `{ type: "windows", name }`)
- npm SDK `listPrinters()` and Windows print targets

## [1.0.x]

- Initial release: network TCP printing, Windows installer, scheduled task, background agent

[1.2.0]: https://github.com/usman8786/Siyaho-POS-Agent/compare/v1.1.0...v1.2.0
