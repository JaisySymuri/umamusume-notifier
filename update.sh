#!/usr/bin/env bash

set -euo pipefail

SERVICE="umamusume-notifier"
BINARY="umamusume-notifier"
INSTALL_DIR="/opt/umamusume-notifier"

echo "==> Pulling latest source..."
git pull

echo "==> Building..."
go build -o "$BINARY" ./cmd/server

echo "==> Stopping service..."
sudo systemctl stop "$SERVICE"

echo "==> Installing binary..."
sudo install -m 755 "$BINARY" "$INSTALL_DIR/$BINARY"

echo "==> Starting service..."
sudo systemctl start "$SERVICE"

echo "==> Checking service status..."
systemctl --no-pager --lines=10 status "$SERVICE"

echo
echo "✅ Update completed successfully."git 