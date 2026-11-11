# Regression guard for the native database wizard at Windows high DPI.
# The normal desktop package build is responsible for compiling this WinForms code.
$ErrorActionPreference = "Stop"
$source = Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot "DatabaseSetupForm.cs")

function Require([string]$needle, [string]$reason) {
    if (-not $source.Contains($needle)) { throw "Database setup regression: $reason" }
}

Require 'AutoScaleMode = AutoScaleMode.Dpi;' 'DPI scaling must be enabled'
Require 'FormBorderStyle = FormBorderStyle.Sizable;' 'Wizard must allow resizing'
Require 'Screen.FromControl(this).WorkingArea' 'Scaled wizard must fit the current monitor'
Require 'Size = new Size(Math.Min(Width, maxWidth), Math.Min(Height, maxHeight));' 'Dialog size must be clamped on small displays'
Require 'AutoScroll = true' 'Settings must remain scrollable'
Require 'root.Controls.Add(scroll, 0, 0);' 'Scrollable controls must occupy the main row'
Require 'root.Controls.Add(footer, 0, 1);' 'Primary action buttons must have a dedicated footer row'
Require 'footer.Controls.Add(button);' 'Primary buttons must be inside the pinned footer'
Require 'secondary.Controls.Add(create);' 'Create database action must remain accessible'
Require 'secondary.Controls.Add(provisionApp);' 'Dedicated account setup must remain accessible'
Require 'WrapContents = true' 'Action groups must be allowed to wrap'
Require 'intro.MaximumSize = new Size(width, 0);' 'Long introduction must wrap'
Require 'discoveryStatus.MaximumSize = new Size(width, 0);' 'Diagnostic status must wrap'
Require 'password.KeyDown +=' 'Password key-repeat guard must exist'
Require '!passwordKeysHeld.Add(e.KeyCode)' 'Repeated held password keystrokes must be suppressed'
Require 'password.KeyUp +=' 'Held key must be released'
Require 'password.Leave +=' 'Held key state must clear on focus loss'
Require 'e.KeyCode == Keys.V' 'Repeated Ctrl+V must be covered'

Write-Host "[ARTEX] Database setup DPI layout and password key-repeat regressions checked."
