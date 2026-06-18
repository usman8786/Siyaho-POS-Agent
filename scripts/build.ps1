param(
  [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Dist = Join-Path $Root "dist"
New-Item -ItemType Directory -Force -Path $Dist | Out-Null

Write-Host "Building SiyahoPrinterAgent.exe..."
Push-Location $Root
go build -ldflags "-X main.version=$Version" -o "$Dist\SiyahoPrinterAgent.exe" ./agent
Pop-Location

$Inno = "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe"
if (Test-Path $Inno) {
  Write-Host "Compiling installer..."
  & $Inno "/DMyAppVersion=$Version" (Join-Path $Root "installer\windows\setup.iss")
} else {
  Write-Warning "Inno Setup not found. Binary built at $Dist\SiyahoPrinterAgent.exe"
}

$releases = @{
  version    = $Version
  windows    = "https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-$Version.exe"
  minVersion = $Version
} | ConvertTo-Json
Set-Content -Path (Join-Path $Root "releases.json") -Value $releases -Encoding UTF8

Write-Host "Done."
