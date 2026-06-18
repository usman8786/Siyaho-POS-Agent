#define MyAppName "Siyaho Printer Agent"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Siyaho"
#define MyAppExeName "SiyahoPrinterAgent.exe"
#define MyAppId "{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}"

[Setup]
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\Siyaho\PrinterAgent
DefaultGroupName=Siyaho Printer Agent
DisableProgramGroupPage=yes
OutputDir=..\..\dist
OutputBaseFilename=SiyahoPrinterAgent-Setup-{#MyAppVersion}
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\..\dist\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "config.default.json"; DestDir: "{commonappdata}\Siyaho\PrinterAgent"; Flags: onlyifdoesntexist uninsneveruninstall

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\Uninstall {#MyAppName}"; Filename: "{uninstallexe}"

[Run]
Filename: "schtasks"; Parameters: "/Create /TN ""SiyahoPrinterAgent"" /TR ""\""{app}\{#MyAppExeName}\"""" /SC ONLOGON /RL LIMITED /F"; Flags: runhidden; StatusMsg: "Registering auto-start..."
Filename: "{app}\{#MyAppExeName}"; Description: "Start {#MyAppName}"; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "schtasks"; Parameters: "/Delete /TN ""SiyahoPrinterAgent"" /F"; Flags: runhidden
Filename: "taskkill"; Parameters: "/IM {#MyAppExeName} /F"; Flags: runhidden

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    // Agent started via [Run] section.
  end;
end;
