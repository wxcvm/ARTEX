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
# A 100%-DPI build can pass even when fixed-size Chinese labels are cut off at
# large system text scales. Assert every action group uses measured button sizes.
Require 'private static void ConfigureActionButton(Button button, int minimumWidth)' 'Shared autosize helper must exist'
Require 'button.AutoSize = true;' 'Action width must be measured from button text'
Require 'button.Padding = new Padding(10, 3, 10, 3);' 'Text needs space around the rendered label'
Require 'TextRenderer.MeasureText(button.Text, button.Font)' 'Button labels must fit at large font scale'
Require 'button.FontChanged +=' 'Hidden advanced buttons must resize when DPI text changes'
Require 'button.TextChanged +=' 'Primary button must resize when the setup mode changes'
foreach ($button in @('detectLocal','startLocal','installLocal','create','provisionApp','save','test','cancel')) {
    Require "ConfigureActionButton($button," "Missing DPI-aware sizing for $button"
}
foreach ($fixed in @('Width = 94, Height = 30', 'Width = 112, Height = 30', 'Width = 140, Height = 30', 'Width = 135, Height = 30')) {
    if ($source.Contains($fixed)) { throw "Database setup regression: fixed-size action label still exists ($fixed)" }
}
Require 'ActiveControl = host;' 'Wizard must focus the first form field'
Require 'scroll.AutoScrollPosition = Point.Empty;' 'Wizard must begin with the introduction visible'
Require 'BeginInvoke((Action)(() =>' 'Restore top scroll after the initial focus/layout pass'


Write-Host "[ARTEX] Database setup DPI layout and password key-repeat regressions checked."
