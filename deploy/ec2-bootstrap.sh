#!/usr/bin/env bash
# ClearRouter — EC2 bootstrap (run as root / sudo)
# Installs Docker + Cloudflare Tunnel and runs the backend.
set -euo pipefail

log() { echo -e "\n[SETUP] $1"; }

log "Updating system..."
apt-get update -y

log "Installing Docker..."
curl -fsSL https://get.docker.com | sh
systemctl enable --now docker

log "Installing cloudflared..."
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | gpg --dearmor -o /usr/share/keyrings/cloudflare-warp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/cloudflare-warp-archive-keyring.gpg] https://pkg.cloudflare.com/cloudflared any main" > /etc/apt/sources.list.d/cloudflared.list
apt-get update -y
apt-get install -y cloudflared

log "Creating app directory..."
mkdir -p /opt/clearrouter

log "Done. Next steps:"
echo "  1. docker login ghcr.io -u hariomop12  (paste PAT)"
echo "  2. Copy .env to /opt/clearrouter/.env"
echo "  3. Copy deploy/docker-compose.prod.yml to /opt/clearrouter/docker-compose.yml"
echo "  4. cd /opt/clearrouter && docker compose up -d"
echo "  5. cloudflared tunnel login   →  pick account"
echo "  6. cloudflared tunnel create clearrouter"
echo "  7. cloudflared tunnel route dns clearrouter api.clearrouter.hariomop.in"
echo "  8. Copy deploy/cloudflared.yml → ~/.cloudflared/config.yml with tunnel UUID"
echo "  9. cloudflared service install"
