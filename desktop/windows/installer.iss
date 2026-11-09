; ARTEX Windows installer, built from the verified portable desktop stage.
; AppId is stable across versions for in-place upgrade and standard uninstallation.
#ifndef StageDir
  #error StageDir is required
#endif
#ifndef OutputDir
  #error OutputDir is required
#endif
#ifndef AppVersion
  #define AppVersion "0.0.0-preview"
#endif

[Setup]
AppId={{E6C60BA8-3625-4CA6-BF7E-A74D9D578CCA}
AppName=ARTEX
AppVersion={#AppVersion}
AppVerName=ARTEX {#AppVersion}
AppPublisher=wxcvm
AppPublisherURL=https://github.com/wxcvm/ARTEX
AppSupportURL=https://github.com/wxcvm/ARTEX/issues
AppUpdatesURL=https://github.com/wxcvm/ARTEX/releases
DefaultDirName={localappdata}\Programs\ARTEX
DefaultGroupName=ARTEX
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
WizardStyle=modern
DisableDirPage=no
OutputDir={#OutputDir}
OutputBaseFilename=ARTEX-Windows-Setup-amd64
SetupIconFile={#StageDir}\ARTEX.ico
UninstallDisplayIcon={app}\ARTEX.Desktop.exe
Compression=lzma2
SolidCompression=yes
CloseApplications=yes
RestartApplications=no

[Tasks]
Name: "desktopicon"; Description: "创建桌面快捷方式"; GroupDescription: "附加选项："; Flags: unchecked

[Files]
; The application binaries update in place. config.json, data/ and generated
; session files are NOT in the staged package and are never installed over.
Source: "{#StageDir}\*"; DestDir: "{app}"; Excludes: "\skills\*"; Flags: ignoreversion recursesubdirs createallsubdirs
; Existing skills may be modified through ARTEX. Preserve them on upgrade and
; uninstall; only add missing bundled skills from newer releases.
Source: "{#StageDir}\skills\*"; DestDir: "{app}\skills"; Flags: ignoreversion recursesubdirs createallsubdirs onlyifdoesntexist uninsneveruninstall

[Icons]
Name: "{group}\ARTEX"; Filename: "{app}\ARTEX.Desktop.exe"; WorkingDir: "{app}"
Name: "{autodesktop}\ARTEX"; Filename: "{app}\ARTEX.Desktop.exe"; WorkingDir: "{app}"; Tasks: desktopicon

[Run]
Filename: "{app}\ARTEX.Desktop.exe"; Description: "启动 ARTEX"; Flags: nowait postinstall skipifsilent unchecked
