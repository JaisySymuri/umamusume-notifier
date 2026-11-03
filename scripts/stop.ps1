
$ErrorActionPreference = "Stop"

# Configuration
$Repo = "C:\Users\bests\Documents\Coding\Exercise\Golang\umamusume-notifier"
$Runtime = "C:\Apps\umamusume-notifier"
$Branch = "lpUpdate"

$PidFile = Join-Path $Runtime "server.pid"
$ShutdownFlag = Join-Path $Runtime "shutdown.flag"

Write-Host "=== Umamusume Notifier: STOP ===" -ForegroundColor Cyan

# Check whether the application is running
if (!(Test-Path $PidFile)) {
    Write-Host "Notifier is not running."
    exit 0
}

$AppPID = [int](Get-Content $PidFile -Raw)
$Process = Get-Process -Id $AppPID -ErrorAction SilentlyContinue

if ($Process) {
    # Request graceful shutdown
    Write-Host "Requesting graceful shutdown..."
    New-Item -ItemType File -Path $ShutdownFlag -Force | Out-Null

    # Wait up to 30 seconds
    $Process.WaitForExit(30000)

    $Process.Refresh()

    if (!$Process.HasExited) {
        throw "Notifier did not stop within 30 seconds. No forceful termination performed."
    }

    Write-Host "Application stopped successfully." -ForegroundColor Green
}
else {
    Write-Host "Process not found. Cleaning up stale PID file."
}

Remove-Item $PidFile -Force -ErrorAction SilentlyContinue
Remove-Item $ShutdownFlag -Force -ErrorAction SilentlyContinue

# Database synchronization
Write-Host "Backing up database..."

$Source = Join-Path $Runtime "data.db"
$Destination = Join-Path $Repo "data.db"
$WalFile = "$Source-wal"

if (!(Test-Path $Source)) {
    throw "Runtime database not found."
}

# Normally, a clean SQLite shutdown checkpoints and removes WAL files.
# If a WAL file remains, use SQLite backup instead of raw file copying.
if ((Test-Path $WalFile) -and (Get-Item $WalFile).Length -gt 0) {
    Write-Host "WAL file detected. Using SQLite backup..."

    if (!(Get-Command sqlite3 -ErrorAction SilentlyContinue)) {
        throw "SQLite WAL file exists, but sqlite3 CLI is unavailable. Database was not copied."
    }

    & sqlite3 $Source ".backup '$Destination'"

    if ($LASTEXITCODE -ne 0) {
        throw "SQLite backup failed."
    }
}
else {
    Copy-Item $Source $Destination -Force
}

# Verify backup integrity if sqlite3 is installed
if (Get-Command sqlite3 -ErrorAction SilentlyContinue) {
    $Integrity = & sqlite3 $Destination "PRAGMA integrity_check;"

    if ($LASTEXITCODE -ne 0 -or $Integrity -ne "ok") {
        throw "Database integrity check failed."
    }
}

Write-Host "Database backup completed." -ForegroundColor Green

# Git operations
Write-Host "Committing database to Git..."

Set-Location $Repo

git add -- data.db

if ($LASTEXITCODE -ne 0) {
    throw "Git add failed."
}

git diff --cached --quiet -- data.db

if ($LASTEXITCODE -eq 0) {
    Write-Host "No database changes detected. Skipping commit."
    exit 0
}
elseif ($LASTEXITCODE -ne 1) {
    throw "Git diff failed."
}

git commit -m "Weekly database backup: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -- data.db

if ($LASTEXITCODE -ne 0) {
    throw "Git commit failed."
}

git push origin $Branch

if ($LASTEXITCODE -ne 0) {
    throw "Git push failed."
}

Write-Host "Database backup committed and pushed successfully." -ForegroundColor Green
