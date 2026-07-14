#!/usr/bin/env bash

set -euo pipefail

SERVICE="umamusume-notifier"
BINARY="umamusume-notifier"
INSTALL_DIR="/opt/umamusume-notifier"

cleanup() {
    echo "Ensuring service is running..."
    sudo systemctl start "$SERVICE" >/dev/null 2>&1 || true
}
trap cleanup EXIT

git pull
go build -o "$BINARY" ./cmd/server

sudo systemctl stop "$SERVICE"
sudo install -m 755 "$BINARY" "$INSTALL_DIR/$BINARY"

sudo systemctl start "$SERVICE"
systemctl --no-pager --lines=10 status "$SERVICE"

echo "✅ Update completed."