
$ErrorActionPreference = "Stop"

# Configuration
$Repo = "C:\Users\bests\Documents\Coding\Exercise\Golang\umamusume-notifier"
$Runtime = "C:\Apps\umamusume-notifier"
$Branch = "lpUpdate"

$Exe = Join-Path $Runtime "server.exe"
$PidFile = Join-Path $Runtime "server.pid"
$ShutdownFlag = Join-Path $Runtime "shutdown.flag"

Write-Host "=== Umamusume Notifier: START ===" -ForegroundColor Cyan

# Check if already running
if (Test-Path $PidFile) {
    $ExistingPID = [int](Get-Content $PidFile -Raw)
    if (Get-Process -Id $ExistingPID -ErrorAction SilentlyContinue) {
        throw "Notifier is already running (PID: $ExistingPID)."
    }
    Remove-Item $PidFile -Force
}

# Pull latest repository changes
Set-Location $Repo

Write-Host "[1/4] Pulling latest Git changes..."
git pull --ff-only origin $Branch

if ($LASTEXITCODE -ne 0) {
    throw "Git pull failed."
}

# Check database
if (!(Test-Path "$Repo\data.db")) {
    throw "Repository database not found."
}

# Synchronize database
Write-Host "[2/4] Synchronizing database..."
Copy-Item "$Repo\data.db" "$Runtime\data.db" -Force

# Build application
Write-Host "[3/4] Building Go application..."
go build -o $Exe ./cmd/server

if ($LASTEXITCODE -ne 0) {
    throw "Go build failed."
}

# Remove previous shutdown signal
if (Test-Path $ShutdownFlag) {
    Remove-Item $ShutdownFlag -Force
}

# Start application
Write-Host "[4/4] Starting application..."

$Process = Start-Process `
    -FilePath $Exe `
    -WorkingDirectory $Runtime `
    -WindowStyle Hidden `
    -RedirectStandardOutput "$Runtime\app.log" `
    -RedirectStandardError "$Runtime\error.log" `
    -PassThru

Set-Content -Path $PidFile -Value $Process.Id

Write-Host "Notifier started successfully. PID: $($Process.Id)" -ForegroundColor Green
