#define MyAppName "Siyaho Printer Agent"
#ifndef MyAppVersion
#define MyAppVersion "1.2.1"
#endif
#define MyAppPublisher "Siyaho"
#define MyAppExeName "SiyahoPrinterAgent.exe"
#define MyAppTaskName "SiyahoPrinterAgent"

[Setup]
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}}
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
; Stop/kill the agent in [Code] PrepareToInstall before files are replaced.
CloseApplications=no

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\..\dist\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "config.default.json"; DestDir: "{commonappdata}\Siyaho\PrinterAgent"; Flags: onlyifdoesntexist uninsneveruninstall

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\Uninstall {#MyAppName}"; Filename: "{uninstallexe}"

[Run]
Filename: "schtasks"; Parameters: "/Create /TN ""{#MyAppTaskName}"" /TR ""\""{app}\{#MyAppExeName}\"""" /SC ONLOGON /RL LIMITED /F"; Flags: runhidden waituntilterminated; StatusMsg: "Registering auto-start..."
Filename: "schtasks"; Parameters: "/Run /TN ""{#MyAppTaskName}"""; Flags: runhidden waituntilterminated; Description: "Start {#MyAppName} in the background"; StatusMsg: "Starting {#MyAppName}..."

[Code]
procedure StopRunningAgent;
var
  ResultCode: Integer;
  I: Integer;
begin
  for I := 1 to 4 do
  begin
    Exec('schtasks.exe', '/End /TN "{#MyAppTaskName}"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/IM {#MyAppExeName} /F', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Sleep(750);
  end;
end;

function PrepareToInstall(var NeedsRestart: Boolean): String;
begin
  StopRunningAgent;
  Result := '';
end;

function InitializeUninstall(): Boolean;
begin
  StopRunningAgent;
  Result := True;
end;

[UninstallRun]
Filename: "schtasks"; Parameters: "/End /TN ""{#MyAppTaskName}"""; Flags: runhidden
Filename: "schtasks"; Parameters: "/Delete /TN ""{#MyAppTaskName}"" /F"; Flags: runhidden
Filename: "taskkill"; Parameters: "/IM {#MyAppExeName} /F"; Flags: runhidden
