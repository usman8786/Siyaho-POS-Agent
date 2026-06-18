# Code signing (Windows)

Unsigned installers trigger Microsoft SmartScreen warnings. For production rollout, sign both the agent `.exe` and the setup `.exe`.

## Certificate types

| Type | Cost (approx) | Notes |
|------|----------------|-------|
| Standard OV | $200–400/yr | Reputation builds over time |
| EV | $400–700/yr | Immediate SmartScreen trust; USB token required |

Providers: DigiCert, Sectigo, SSL.com, Certum.

## Sign locally

Install [Windows SDK](https://developer.microsoft.com/windows/downloads/windows-sdk/) for `signtool`.

```powershell
signtool sign /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 `
  /f C:\certs\siyaho-codesign.pfx /p "YOUR_PASSWORD" `
  dist\SiyahoPrinterAgent.exe

signtool sign /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 `
  /f C:\certs\siyaho-codesign.pfx /p "YOUR_PASSWORD" `
  dist\SiyahoPrinterAgent-Setup-1.0.0.exe
```

## GitHub Actions secrets

| Secret | Value |
|--------|--------|
| `CODESIGN_PFX_BASE64` | Base64-encoded `.pfx` file |
| `CODESIGN_PFX_PASSWORD` | PFX password |

The release workflow signs when both secrets are set; otherwise it publishes unsigned artifacts.

## EV hardware token

EV certificates require a USB token on the signing machine. Use a dedicated Windows self-hosted runner with the token plugged in for release builds.
