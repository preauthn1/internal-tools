# Starts meetnote in development mode on Windows (counterpart of dev.sh):
# - models come from spike\models (no download needed)
# - the Claude API key comes from ..\.env, so dev builds never touch the credential store
# - cargo builds into a short directory. The Vulkan build of whisper.cpp runs a nested CMake project whose
#   paths exceed 260 characters under app\src-tauri\target and then fails (C1083 "Invalid argument" in
#   CMake's compiler test). Building into C:\t succeeded. Not measured why exactly: only the directory changed.
# Set CARGO_TARGET_DIR yourself before running to use another short directory.
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

if (-not $env:CARGO_TARGET_DIR) { $env:CARGO_TARGET_DIR = 'C:\t' }
$env:MEETNOTE_MODEL_DIR = (Resolve-Path (Join-Path $PSScriptRoot '..\spike\models')).Path

# whisper-rs-sys reads VULKAN_SDK at build time; the SDK installer sets it machine-wide, but a terminal
# opened before the install doesn't have it yet.
if (-not $env:VULKAN_SDK) {
    $env:VULKAN_SDK = [Environment]::GetEnvironmentVariable('VULKAN_SDK', 'Machine')
}

$envFile = Join-Path $PSScriptRoot '..\.env'
if (Test-Path $envFile) {
    $line = Get-Content $envFile | Where-Object { $_ -match '^ANTHROPIC_API_KEY=' } | Select-Object -First 1
    if ($line) {
        $key = ($line -replace '^ANTHROPIC_API_KEY=', '').Trim().Trim('"', "'")
        if ($key) { $env:ANTHROPIC_API_KEY = [REDACTED] }
    }
}

npm run tauri dev
