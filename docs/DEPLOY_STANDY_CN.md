# Sub2API Standby 部署教程

> 仓库：`gthubtom1/sub2api-standby`（公开）  
> 镜像：`ghcr.io/gthubtom1/sub2api-standby:latest`（GitHub Actions 自动构建推送）  
> **服务器上不需要装 Go，不需要源码编译。**

## 一键安装（推荐，像官方一样快）

前提：机器已装 Docker + Compose。

```bash
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/quick-pull-deploy.sh | bash
```

做的事：
1. 下载本 fork 的 compose / env  
2. `docker pull ghcr.io/gthubtom1/sub2api-standby:latest`  
3. 启动 postgres + redis + sub2api  

访问：`http://服务器IP:8080`

### 升级

```bash
cd ~/sub2api-standby   # 或你的部署目录
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d
```

### 绝对禁止

- 管理后台「检查更新 / 更新」（会拉官方 `Wei-Shaw/sub2api`）  
- `weishaw/sub2api` 官方镜像  
- 在 2G 小机上 `docker build` 源码（慢且易 OOM）

---

## 首次安装：官方 Docker 自动安装（无需向导）

本 fork 的 deploy/docker-compose.local.yml 默认 **AUTO_SETUP=false**。

1. 打开 http://你的IP:8080 或 /setup
2. 向导里填写（Docker 内网默认值）：
   - PostgreSQL：host=postgres port=5432 user=sub2api db=sub2api password= 看部署目录 .env 里的 POSTGRES_PASSWORD
   - Redis：host=redis port=6379 密码一般为空
3. **自己创建管理员**：邮箱必须是合法 email（不要用 dmin@local），密码自设
4. 安装完成后才会进入后台登录

不要用 AUTO_SETUP=true + 非法邮箱跳过向导，否则会装完却登不进去。

若误开了 AUTO_SETUP 想重来：
`ash
cd /path/to/deploy
docker compose down
# 会清空数据
rm -rf data postgres_data redis_data
mkdir -p data postgres_data redis_data
# 确认 compose 里 AUTO_SETUP=false
docker compose up -d
`


## 镜像从哪来

GitHub Actions 工作流：`.github/workflows/ghcr-image.yml`

- 推送 `main` → 自动构建并推送  
  - `ghcr.io/gthubtom1/sub2api-standby:latest`  
  - `ghcr.io/gthubtom1/sub2api-standby:0.1.157-standby`  
  - `ghcr.io/gthubtom1/sub2api-standby:sha-<commit>`  
- 打 tag `v*` → 同步推送 tag 名  
- 也可在 Actions 页手动 **Run workflow**

公共仓库的 GHCR 包首次推送后，建议在 GitHub Packages 页面把 visibility 设为 **Public**，否则别人 `docker pull` 可能要登录。

---

## 手动 compose（同样只 pull）

```bash
mkdir -p sub2api-standby && cd sub2api-standby
curl -fsSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/docker-compose.local.yml -o docker-compose.yml
curl -fsSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/.env.example -o .env.example
cp .env.example .env && chmod 600 .env
# 编辑 .env 里的密码/密钥
mkdir -p data postgres_data redis_data
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d
```

---

## 误点网页更新了

1. 别删 `data` / `postgres_data`  
2. 重新 pull 本镜像并 up：

```bash
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
docker compose up -d --force-recreate sub2api
```

---

## 开发者：本地改代码后出镜像

不需要在用户 VPS 上编。推到 GitHub 即可：

```bash
git push origin main
# 等 Actions 变绿后
docker pull ghcr.io/gthubtom1/sub2api-standby:latest
```
