# Release process

This project uses [Semantic Versioning](https://semver.org/). Tags look like `v1.2.0`.

## Before you release

1. Merge all changes to `main`.
2. Update version in:
   - [`agent/meta/meta.go`](../agent/meta/meta.go) — `const Version`
   - [`sdk/js/package.json`](../sdk/js/package.json) — `"version"`
   - [`docs/openapi.yaml`](openapi.yaml) — `info.version` (if API version changed)
3. Add an entry to [`CHANGELOG.md`](../CHANGELOG.md).
4. Commit: `chore: release v1.2.0` (example).

## GitHub release (Windows installer)

Push a tag to trigger [`.github/workflows/release.yml`](../.github/workflows/release.yml):

```bash
git tag v1.2.0
git push origin v1.2.0
```

The workflow:

1. Builds `SiyahoPrinterAgent.exe` (Windows, `-H=windowsgui`)
2. Builds Inno Setup installer `SiyahoPrinterAgent-Setup-{version}.exe`
3. Optionally code-signs if `CODESIGN_PFX_*` secrets are set
4. Writes `releases.json` and creates a GitHub Release with artifacts
5. Optionally deploys to CDN if `SELFHOSTED_SSH_PRIVATE_KEY` is configured

### Manual workflow dispatch

In GitHub Actions → **Release** → **Run workflow**, enter the version (e.g. `1.2.0`).

## npm SDK publish

After the GitHub release:

```bash
cd sdk/js
npm publish --access public
```

Requires npm login and publish rights on the `siyaho-printer-agent` package.

## CDN artifacts

Production download URLs (hosted separately):

- `https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe`
- `https://pos.siyaho.com/downloads/print-agent/releases.json`

See [deploy-cdn.md](deploy-cdn.md) for server paths and manual upload steps if CI deploy is skipped.

## Post-release

1. Verify `GET http://127.0.0.1:17890/v1/health` reports the new version on a test PC.
2. Update [Siyaho-POS-Web](https://github.com/usman8786/Siyaho-POS-Web) dependency if needed: `siyaho-printer-agent@^1.2.0`.
3. Announce in GitHub Release notes (auto-generated + edit as needed).

## Hotfix

For urgent fixes on the latest minor version:

1. Branch from `main`, fix, bump patch (`1.2.1`).
2. Tag and push — same workflow as above.
