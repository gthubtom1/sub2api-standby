# Sub2API Standby Deployment Guide

> Private fork: `gthubtom1/sub2api-standby`  
> Features: background standby probe, auto error quarantine, frontend **Standby Healthy** badge.  
> **Do not use official install scripts, official images, or the in-app Update button.**

## Canonical sources

| Item | Value |
|------|-------|
| Repo | `https://github.com/gthubtom1/sub2api-standby` |
| Image tag | `sub2api-custom:0.1.157-standby` |
| Optional GHCR | `ghcr.io/gthubtom1/sub2api-standby:latest` |
| Build file | `Dockerfile.custom` |

## Recommended: build from source (Docker)

```bash
git clone https://github.com/gthubtom1/sub2api-standby.git
cd sub2api-standby
# private repo: gh auth login or token URL

# low-RAM hosts: enable 2G swap first
export DOCKER_BUILDKIT=1
docker build --memory=2g --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom -t sub2api-custom:0.1.157-standby .

cd deploy
cp .env.example .env && chmod 600 .env
# edit secrets
mkdir -p data postgres_data redis_data
docker compose -f docker-compose.local.yml up -d
```

## Replace existing container (keep data)

Do **not** delete data volumes.

```bash
docker inspect sub2api --format '{{range .Config.Env}}{{println .}}{{end}}' \
  | grep -v '^PATH=' > /tmp/sub2api.env
TS=$(date +%Y%m%d%H%M%S)
docker stop sub2api && docker rename sub2api sub2api-old-$TS
docker run -d --name sub2api --restart unless-stopped \
  --env-file /tmp/sub2api.env -p 127.0.0.1:8080:8080 \
  -v /YOUR/DATA/PATH:/app/data sub2api-custom:0.1.157-standby
```

## Upgrade

```bash
git pull
docker build --memory=2g --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom -t sub2api-custom:0.1.157-standby .
# recreate container / compose up -d
```

## Never do this

- In-app **Check for Updates / Update** (hardcoded to official `Wei-Shaw/sub2api`)
- `weishaw/sub2api` images
- Official `curl ... Wei-Shaw/sub2api ... | bash`

## If you clicked Update by mistake

Rebuild/recreate from this fork image; keep DB volumes. Fork features return; data usually intact.
