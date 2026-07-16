# Sub2API Standby Deployment

> Repo: `gthubtom1/sub2api-standby`  
> Image: `ghcr.io/gthubtom1/sub2api-standby:latest` (GitHub Actions)  
> **No Go / source build on the server.**

## One-click

```bash
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/quick-pull-deploy.sh | bash
```

## Upgrade

```bash
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d
```

## Never

- In-app Update button
- `weishaw/sub2api`
- Source build on tiny VPS

## CI

`.github/workflows/ghcr-image.yml` publishes to GHCR on `main` / tags / manual dispatch. Make the package Public after first push if anonymous pulls fail.
