#!/usr/bin/env bash
# Sub2API Standby - one-click pull & deploy (NO source build on server)
# Usage:
#   curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/quick-pull-deploy.sh | bash
set -euo pipefail

REPO="gthubtom1/sub2api-standby"
IMAGE="ghcr.io/gthubtom1/sub2api-standby:latest"
RAW="https://raw.githubusercontent.com/${REPO}/main/deploy"
DIR="${SUB2API_DIR:-$HOME/sub2api-standby}"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'
info() { echo -e "${GREEN}[INFO]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
err()  { echo -e "${RED}[ERROR]${NC} $*" >&2; }

if ! command -v docker >/dev/null 2>&1; then
  err "Docker is required. Install Docker first."
  exit 1
fi
if ! docker compose version >/dev/null 2>&1; then
  err "Docker Compose plugin is required (docker compose)."
  exit 1
fi

mkdir -p "$DIR"
cd "$DIR"
info "Deploy directory: $DIR"

if [[ ! -f docker-compose.yml ]]; then
  info "Downloading compose files from ${REPO}..."
  curl -fsSL "$RAW/docker-compose.local.yml" -o docker-compose.yml
  curl -fsSL "$RAW/.env.example" -o .env.example
fi

if [[ ! -f .env ]]; then
  info "Generating .env secrets..."
  cp .env.example .env
  chmod 600 .env
  gen() { openssl rand -hex 32; }
  if command -v openssl >/dev/null 2>&1; then
    grep -q '^POSTGRES_PASSWORD=' .env || echo "POSTGRES_PASSWORD=$(gen)" >> .env
    grep -q '^JWT_SECRET=' .env || echo "JWT_SECRET=$(gen)" >> .env
    grep -q '^TOTP_ENCRYPTION_KEY=' .env || echo "TOTP_ENCRYPTION_KEY=$(gen)" >> .env
    sed -i "s/^POSTGRES_PASSWORD=$/POSTGRES_PASSWORD=$(gen)/" .env || true
    sed -i "s/^POSTGRES_PASSWORD=your_.*/POSTGRES_PASSWORD=$(gen)/" .env || true
    sed -i "s/^JWT_SECRET=$/JWT_SECRET=$(gen)/" .env || true
    sed -i "s/^JWT_SECRET=your_.*/JWT_SECRET=$(gen)/" .env || true
    sed -i "s/^TOTP_ENCRYPTION_KEY=$/TOTP_ENCRYPTION_KEY=$(gen)/" .env || true
    sed -i "s/^TOTP_ENCRYPTION_KEY=your_.*/TOTP_ENCRYPTION_KEY=$(gen)/" .env || true
  else
    warn "openssl missing; edit .env manually"
  fi
fi

mkdir -p data postgres_data redis_data

info "Pulling image $IMAGE (prebuilt, no local Go compile)..."
docker pull "$IMAGE"

if grep -q 'image:' docker-compose.yml; then
  sed -i "0,/image:.*/s|image:.*|image: ${IMAGE}|" docker-compose.yml || true
fi

info "Starting stack..."
docker compose -f docker-compose.yml --env-file .env up -d

info "Waiting for health..."
for i in $(seq 1 40); do
  if curl -fsS http://127.0.0.1:8080/health >/dev/null 2>&1; then
    echo
    info "HEALTH OK"
    break
  fi
  sleep 3
done

docker compose -f docker-compose.yml ps
docker exec sub2api /app/sub2api --version 2>/dev/null || true

cat <<EOF

========================================
Sub2API Standby is up (pull-only deploy)
URL:  http://SERVER_IP:8080
Dir:  $DIR
Image: $IMAGE

Do NOT use Admin UI "Update" (pulls official binaries).
Upgrade later:
  cd $DIR && docker pull $IMAGE && docker compose up -d
========================================
EOF
