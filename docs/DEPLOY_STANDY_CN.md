# Sub2API Standby 部署教程（本 fork）

> 镜像：`ghcr.io/gthubtom1/sub2api-standby:latest`  
> 仓库：https://github.com/gthubtom1/sub2api-standby  
> 与官方 Docker 一样：**AUTO_SETUP 自动装好，不用手填数据库/安装向导**。

---

## 一键安装（推荐）

```bash
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/quick-pull-deploy.sh | bash
```

脚本会：

1. 下载 compose / `.env`
2. `docker pull ghcr.io/gthubtom1/sub2api-standby:latest`（预编译镜像，服务器上**不装 Go、不编译**）
3. `docker compose up -d`，`AUTO_SETUP=true` 自动初始化
4. 打印管理员邮箱和密码

然后打开：

```text
http://你的服务器IP:8080
```

用脚本打印的账号登录即可。

---

## 首次安装说明（官方 Docker 风格）

| 项 | 说明 |
|---|---|
| 安装方式 | `AUTO_SETUP=true` 自动初始化 |
| 安装向导 | **不会出现**（Docker 路径不需要） |
| 数据库/Redis | compose 已内置，不用你填 host |
| 管理员 | `.env` 里的 `ADMIN_EMAIL` / `ADMIN_PASSWORD` |
| 默认邮箱 | `admin@sub2api.local`（合法 email 格式） |

登录后可在后台改密码。

---

## 手动 compose（同样只 pull）

```bash
mkdir -p ~/sub2api-standby && cd ~/sub2api-standby
curl -fsSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/docker-compose.local.yml -o docker-compose.yml
curl -fsSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/.env.example -o .env.example
cp .env.example .env && chmod 600 .env
# 编辑 .env：至少设置 POSTGRES_PASSWORD / JWT_SECRET / TOTP_ENCRYPTION_KEY / ADMIN_PASSWORD
mkdir -p data postgres_data redis_data
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d
```

---

## 升级（保留数据）

```bash
cd ~/sub2api-standby   # 或你的部署目录
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d
```

## 绝对禁止

- ~~管理后台「立即更新」二进制通道~~ 已改为本仓库 Docker 热更新提示（不会再拉官方）
- `weishaw/sub2api` 官方镜像
- 在小内存机器上源码 `docker build`（慢且易 OOM）

误点了网页更新怎么办：不要用它。只用上面的 `docker pull` 本镜像升级。

---

## 镜像从哪来

GitHub Actions：`.github/workflows/ghcr-image.yml`

- 推送 `main` → 自动构建并推送  
  - `ghcr.io/gthubtom1/sub2api-standby:latest`  
  - `ghcr.io/gthubtom1/sub2api-standby:0.1.157-standby`  
  - `ghcr.io/gthubtom1/sub2api-standby:sha-<commit>`  
- 打 tag `v*` → 同步推送 tag 名  
- 也可在 Actions 页手动 **Run workflow**

公共仓库的 GHCR 包首次推送后，建议在 GitHub Packages 页面把 visibility 设为 **Public**，否则别人 `docker pull` 可能要登录。

---

## 常见问题

**Q: 为什么不用手填 postgres host？**  
A: Docker 官方路径是 `AUTO_SETUP=true`，应用连容器名 `postgres` / `redis`，compose 已写好。

**Q: 登录邮箱是什么？**  
A: 看部署目录 `.env` 的 `ADMIN_EMAIL` / `ADMIN_PASSWORD`，一键脚本结束时也会打印。

**Q: 和官方有什么区别？**  
A: 部署方式一样简单；镜像是本 fork（含预备健康探测等），不要用官方更新按钮。

## 管理后台一键热更新

部署 compose 需挂载：

- `/var/run/docker.sock`
- 部署目录到 `/opt/sub2api-standby`（只读即可）

并设置 `UPDATE_DOCKER_ENABLED=true`。之后在后台点「立即更新」→「立即重启」即可，无需 SSH。

