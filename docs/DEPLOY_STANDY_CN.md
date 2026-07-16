# Sub2API Standby 部署教程

> 本仓库是 **私有 fork**：`gthubtom1/sub2api-standby`  
> 含：后台预备健康探测、失败自动隔离、前端「预备健康」标识。  
> **不要使用任何官方安装脚本 / 官方镜像 / 网页「更新」按钮。**

## 你以后只用这一套

| 项目 | 值 |
|------|----|
| 源码仓库 | `https://github.com/gthubtom1/sub2api-standby` |
| 推荐镜像名 | `sub2api-custom:0.1.157-standby` |
| 可选 GHCR 镜像 | `ghcr.io/gthubtom1/sub2api-standby:latest`（需自己 Actions 推送后才有） |
| 构建文件 | 仓库根目录 `Dockerfile.custom` |

---

## 方式一：Docker 源码构建部署（推荐）

适合：新机器、小内存 VPS（先 build 再 run）、替换现有官方容器。

### 前置

- Docker 20.10+
- 内存建议 ≥ 2G，并加 **2G swap**（低配机编译必做）
- 已有 PostgreSQL + Redis，或本仓库 `deploy/docker-compose*.yml` 一起起

### 1. 拉源码

```bash
git clone https://github.com/gthubtom1/sub2api-standby.git
cd sub2api-standby
```

私有仓需要登录：

```bash
# GitHub CLI
gh auth login
# 或 HTTPS 带 token
git clone https://<TOKEN>@github.com/gthubtom1/sub2api-standby.git
```

### 2. 构建镜像（预编译前端已嵌入时用 custom Dockerfile）

```bash
# 低配机务必先有 swap
sudo fallocate -l 2G /swapfile || sudo dd if=/dev/zero of=/swapfile bs=1M count=2048
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

export DOCKER_BUILDKIT=1
docker build --memory=2g \
  --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom \
  -t sub2api-custom:0.1.157-standby \
  .
```

若 `backend/internal/web/dist` 不完整，先在有 Node 的机器构建前端再拷进仓库后 build。

### 3A. 全新 Compose 部署

```bash
cd deploy
cp .env.example .env
chmod 600 .env
# 编辑 .env：POSTGRES_PASSWORD / JWT_SECRET / TOTP_ENCRYPTION_KEY / ADMIN_*
mkdir -p data postgres_data redis_data

# compose 默认 image 已是本 fork 镜像名
docker compose -f docker-compose.local.yml up -d
docker compose -f docker-compose.local.yml logs -f sub2api
```

访问：`http://服务器IP:8080`

### 3B. 已有 Docker 容器（只换程序，保留数据）

**不要删** `data` / `postgres_data` / `redis_data`。

```bash
# 导出旧容器环境（名称按你实际改）
docker inspect sub2api --format '{{json .Config.Env}}' > /tmp/sub2api.env.json
# 或逐行 env：
docker inspect sub2api --format '{{range .Config.Env}}{{println .}}{{end}}' \
  | grep -v '^PATH=' > /tmp/sub2api.env

TS=$(date +%Y%m%d%H%M%S)
docker stop sub2api
docker rename sub2api sub2api-old-$TS

# 网络/挂载按你原来的改；下面是常见写法
docker run -d --name sub2api --restart unless-stopped \
  --network ai-stack-net \
  --env-file /tmp/sub2api.env \
  -p 127.0.0.1:8080:8080 \
  -v /opt/ai-stack/sub2api/data:/app/data \
  sub2api-custom:0.1.157-standby

curl -fsS http://127.0.0.1:8080/health
docker exec sub2api /app/sub2api --version
# 期望类似：0.1.157-standby
```

---

## 方式二：本仓库一键准备脚本（仅本 fork）

脚本只从 **`gthubtom1/sub2api-standby`** 拉 compose / env，**不碰官方**。

```bash
mkdir -p sub2api-deploy && cd sub2api-deploy
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/docker-deploy.sh | bash
# 先确保本机已有镜像 sub2api-custom:0.1.157-standby（见方式一 build）
docker compose -f docker-compose.local.yml up -d
```

---

## 方式三：二进制 / systemd（仅当你自己发了 Release）

```bash
# 仅当 https://github.com/gthubtom1/sub2api-standby/releases 已有包时
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/install.sh | sudo bash
```

当前若还没有 Release，**请用方式一**，不要硬跑 install。

---

## 升级（正确姿势）

```bash
cd /path/to/sub2api-standby
git pull
docker build --memory=2g --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom -t sub2api-custom:0.1.157-standby .
# 再按 3B 替换容器，或 docker compose up -d
```

### 绝对禁止

1. **管理后台左上角「检查更新 / 更新」** — 代码仍会去拉官方 `Wei-Shaw/sub2api`，会把本 fork 功能冲掉  
2. 任何 `weishaw/sub2api` 镜像  
3. 任何 `Wei-Shaw/sub2api` 的 curl 安装命令  
4. 删除 `postgres_data` / `data` 当“重装”

---

## 误点网页更新了怎么办？

1. **别删数据目录**  
2. 用上面 **3B** 重新换回 `sub2api-custom:0.1.157-standby`  
3. 账号/配置一般还在，只是程序变回官方逻辑了；换回镜像即恢复 fork 功能  

如有旧容器备份：`docker start sub2api-old-时间戳` 可临时顶上。

---

## 功能确认

- 版本：`docker exec sub2api /app/sub2api --version` 含 `standby`  
- 健康：`curl -fsS http://127.0.0.1:8080/health`  
- 前端账号列表出现 **预备健康**  
- 后台按活跃账号自动预探测，失败进错误态  

---

## JSON 导入账号

导入后请检查：分组绑定、可调度、平台类型正确。导入本身不等于健康；依赖 standby 探测把坏号隔开。
