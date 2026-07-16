#!/usr/bin/env bash
# Sub2API Standby - one-click pull & deploy (official Docker style: AUTO_SETUP)
# No Setup Wizard. After up, open http://SERVER_IP:8080 and login.
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
  gen_pass() { openssl rand -base64 18 | tr -d '=+/' | cut -c1-20; }
  if command -v openssl >/dev/null 2>&1; then
    sed -i "s/^POSTGRES_PASSWORD=$/POSTGRES_PASSWORD=$(gen)/" .env || true
    sed -i "s/^POSTGRES_PASSWORD=your_.*/POSTGRES_PASSWORD=$(gen)/" .env || true
    sed -i "s/^JWT_SECRET=$/JWT_SECRET=$(gen)/" .env || true
    sed -i "s/^JWT_SECRET=your_.*/JWT_SECRET=$(gen)/" .env || true
    sed -i "s/^TOTP_ENCRYPTION_KEY=$/TOTP_ENCRYPTION_KEY=$(gen)/" .env || true
    sed -i "s/^TOTP_ENCRYPTION_KEY=your_.*/TOTP_ENCRYPTION_KEY=$(gen)/" .env || true
    # ensure admin password exists
    if ! grep -q '^ADMIN_PASSWORD=' .env || grep -q '^ADMIN_PASSWORD=$' .env || grep -q '^ADMIN_PASSWORD=your_' .env; then
      if grep -q '^ADMIN_PASSWORD=' .env; then
        sed -i "s/^ADMIN_PASSWORD=.*/ADMIN_PASSWORD=$(gen_pass)/" .env
      else
        echo "ADMIN_PASSWORD=$(gen_pass)" >> .env
      fi
    fi
    if ! grep -q '^ADMIN_EMAIL=' .env; then
      echo "ADMIN_EMAIL=admin@sub2api.local" >> .env
    fi
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

# ensure official auto setup
sed -i 's/AUTO_SETUP=false/AUTO_SETUP=true/' docker-compose.yml 2>/dev/null || true

info "Starting stack (AUTO_SETUP, no wizard)..."
docker compose -f docker-compose.yml --env-file .env up -d

info "Waiting for app..."
for i in $(seq 1 40); do
  code=$(curl -sS -m 3 -o /dev/null -w '%{http_code}' http://127.0.0.1:8080/ 2>/dev/null || echo 000)
  if [ "$code" = "200" ] || [ "$code" = "302" ]; then
    echo
    info "APP READY"
    break
  fi
  sleep 3
done

docker compose -f docker-compose.yml ps
ADMIN_EMAIL=$(grep -E '^ADMIN_EMAIL=' .env | head -1 | cut -d= -f2- || echo admin@sub2api.local)
ADMIN_PASSWORD=$(grep -E '^ADMIN_PASSWORD=' .env | head -1 | cut -d= -f2- || echo)

cat <<EOF

========================================
Sub2API Standby is up (official Docker style)
URL:   http://SERVER_IP:8080
Dir:   $DIR
Image: $IMAGE

Login (auto-created, no setup wizard):
  Email:    ${ADMIN_EMAIL}
  Password: ${ADMIN_PASSWORD}

Do NOT use Admin UI "Update" (pulls official binaries).
Upgrade later (keeps data):
  cd $DIR && docker pull $IMAGE && docker compose up -d
========================================
EOF
