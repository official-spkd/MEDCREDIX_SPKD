#!/bin/bash
# Menjalankan backend Golang MedCredix (port 8081).
set -e
cd "$(dirname "$0")/../../backend"
export MEDCREDIX_PORT="${MEDCREDIX_PORT:-8081}"
# Build bila binary belum ada
if [ ! -f medcredix-api ]; then
  export PATH="/home/z/my-project/.tools/go/bin:$PATH"
  go build -o medcredix-api ./cmd/server
fi
exec ./medcredix-api
