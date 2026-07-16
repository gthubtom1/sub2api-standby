# Sub2API

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

<a href="https://trendshift.io/repositories/21823" target="_blank"><img src="https://trendshift.io/api/badge/repositories/21823" alt="Wei-Shaw%2Fsub2api | Trendshift" width="250" height="55"/></a>

**AI API Gateway Platform for Subscription Quota Distribution**

English | [中文](README_CN.md) | [日本語](README_JA.md)

</div>

## ⚠️ Important Notice

Please read the following carefully before using this project:

- **🚨 Terms of Service Risk**: Using this project may violate the terms of service of Anthropic and other upstream providers. Please review the relevant providers' user agreements before use; all risks arising from such use are borne solely by the user.
- **⚖️ Compliant Use**: Use this project only in compliance with the laws and regulations of your country or region. Any unlawful use is strictly prohibited.
- **📖 Disclaimer**: This project is provided for technical learning and research purposes only. The authors assume no liability for account bans, service interruptions, data loss, or any other direct or indirect damages resulting from the use of this project.
- **🚫 No Commercial Authorization**: The developers of this project have never authorized any individual or organization to conduct any form of commercial operation based on this project. Any commercial activity conducted in the name of or based on this project is unrelated to this project and its developers, and all resulting disputes, losses, and legal liabilities shall be borne solely by the party conducting such activity.

## ❤️ Sponsors

> [Want to appear here?](mailto:support@sub2api.org)

<table>

<tr>
<td width="180"><a href="https://cctk.ai/register?aff=SUB2API"><img src="assets/partners/logos/cctk.jpg" alt="CCTK.AI" width="150"></a></td>
<td>Thanks to CCTK.AI for sponsoring this project! <a href="https://cctk.ai/register?aff=SUB2API">CCTK.AI</a> is an AI API gateway focused on stability and cost-effectiveness, offering fast relay services for Claude, OpenAI, Gemini, and other popular models. It works seamlessly with Claude Code, Codex, and other mainstream coding tools, delivering the same model capabilities at a fraction of the official cost. Register via <a href="https://cctk.ai/register?aff=SUB2API">this link</a> for faster, more stable, and more affordable AI API access.</td>
</tr>

<tr>
<td width="180"><a href="https://www.openmodel.ai?ref=sub2api"><img src="assets/partners/logos/openmodel.jpg" alt="openmodel" width="150"></a></td>
<td>One API, every top model! <a href="https://www.openmodel.ai?ref=sub2api">OpenModel</a> is a production-grade, high-availability AI API gateway that makes your applications truly fast and stable: automatic failover, smart routing to the best-performing channel, and a production-grade SLA. An SLA that far surpasses any single provider — making stability your core competitive advantage. Works directly with Claude Code, Codex, and Gemini CLI. Register via this link to get started.</td>
</tr>

<tr>
<td width="180"><a href="https://etok.ai"><img src="assets/partners/logos/etok.png" alt="ETok" width="150"></a></td>
<td>Thanks to ETok.ai for sponsoring this project! ETok.ai is dedicated to building a one-stop AI programming tool service platform. We offer professional Claude Code packages and technical community services, with support for Google Gemini and OpenAI Codex. Through carefully designed plans and a professional tech community, we provide developers with reliable service guarantees and continuous technical support, making AI-assisted programming a true productivity tool. Click <a href="https://etok.ai">here</a> to register!</td>
</tr>

<tr>
<td width="180"><a href="https://apikey.fun/register?aff=SUB2API"><img src="assets/partners/logos/apikey-fun.png" alt="APIKEY.FUN" width="150"></a></td>
<td>Thanks to APIKEY.FUN for sponsoring this project! <a href="https://apikey.fun/register?aff=SUB2API">APIKEY.FUN</a> is one of the core contributors to the sub2api open-source project, dedicated to providing open, stable, and cost-effective AI API access. The platform supports API relay services for Claude, OpenAI, Gemini, and other popular models, with pricing starting from as low as 7% of the original rate. Register via the exclusive link: <a href="https://apikey.fun/register?aff=SUB2API">APIKEY</a> to enjoy a permanent 5% discount on all recharges.</td>
</tr>

<tr>
<td width="180"><a href="https://aigocode.com/invite/SUB2API"><img src="assets/partners/logos/aigocode.png" alt="AIGoCode" width="150"></a></td>
<td>Thanks to AIGoCode for sponsoring this project! AIGoCode is an all-in-one platform that integrates Claude Code, Codex, and the latest Gemini models, providing you with stable, efficient, and highly cost-effective AI coding services. The platform offers flexible subscription plans, zero risk of account suspension, direct access with no VPN required, and lightning-fast responses. AIGoCode has prepared a special benefit for sub2api users: if you register via <a href="https://aigocode.com/invite/SUB2API">this link</a>, you'll receive an extra 10% bonus credit on your first top-up!</td>
</tr>

<tr>
<td width="180"><a href="https://www.aicodemirror.com/register?invitecode=KMVZQM"><img src="assets/partners/logos/AICodeMirror.jpg" alt="AICodeMirror" width="150"></a></td>
<td>Thanks to AICodeMirror for sponsoring this project! AICodeMirror provides official high-stability relay services for Claude Code / Codex / Gemini CLI, with enterprise-grade concurrency, fast invoicing, and 24/7 dedicated technical support. Claude Code / Codex / Gemini official channels at 38% / 2% / 9% of original price, with extra discounts on top-ups! AICodeMirror offers special benefits for sub2api users: register via <a href="https://www.aicodemirror.com/register?invitecode=KMVZQM">this link</a> to enjoy 20% off your first top-up, and enterprise customers can get up to 25% off!</td>
</tr>

<tr>
<td width="180"><a href="https://shop.bmoplus.com/?utm_source=github"><img src="assets/partners/logos/bmoplus.jpg" alt="bmoplus" width="150"></a></td>
<td>Huge thanks to BmoPlus for sponsoring this project! BmoPlus is a highly reliable AI account provider built strictly for heavy AI users and developers. They offer rock-solid, ready-to-use accounts and official top-up services for ChatGPT Plus / ChatGPT Pro (Full Warranty) / Claude Pro / Super Grok / Gemini Pro. By registering and ordering through <a href="https://shop.bmoplus.com/?utm_source=github">BmoPlus - Premium AI Accounts & Top-ups</a>, users can unlock the mind-blowing rate of 10% of the official GPT subscription price (90% OFF)</td>
</tr>

<tr>
<td width="180"><a href="https://bestproxy.com/?keyword=a2e8iuol"><img src="assets/partners/logos/bestproxy.png" alt="bestproxy" width="150"></a></td>
<td>Thanks to Bestproxy for sponsoring this project! <a href="https://bestproxy.com/?keyword=a2e8iuol">Bestproxy</a> provides high-purity residential IPs with dedicated one-IP-per-account support. By combining real home networks with fingerprint isolation, it enables link environment isolation and reduces the probability of association-based risk control.</td>
</tr>

<tr>
<td width="180"><a href="https://pateway.ai/?ch=1tsfr51"><img src="assets/partners/logos/pateway.png" alt="pateway" width="150"></a></td>
<td>Thanks to PatewayAI for sponsoring this project! <a href="https://pateway.ai/?ch=1tsfr51">PatewayAI</a> is a premium API relay built for heavy AI developers, offering the full Claude and Codex series sourced 100% from official providers, with transparent token-level billing. Enterprise plans include high concurrency, dedicated management, contracts, and invoicing. Register now to get $3 in trial credits, top-ups from 60% off, and referral bonuses up to $150.</td>
</tr>

<tr>
<td width="180"><a href="https://api.pptoken.org/register?promo=SUB2API"><img src="assets/partners/logos/pptoken.png" alt="pptoken" width="150"></a></td>
<td>Thanks to PPToken.org for sponsoring this project! <a href="https://api.pptoken.org/register?promo=SUB2API">PPToken.org</a> specializes in GPT model API relay services, supporting Codex, Claude Code, OpenAI-compatible clients, and Gemini CLI integration. Top-ups are 1:1 (¥1 = $1 credit); GPT models start at 0.16x rate multiplier, with overall cost at roughly 2.2% of official pricing and first-token latency around 1 second — ideal for developers seeking low-cost, high-speed access to GPT model capabilities. Technical support: 24/7 real human responses (no bots), @tech in the group chat and get a reply within 10 minutes. Sponsor benefit: the first 200 users who register via the <a href="https://api.pptoken.org/register?promo=SUB2API">exclusive registration link</a> and enter promo code `SUB2API` can claim free Codex / Claude Code trial credits — no minimum spend, no card required.
</td>
</tr>

<tr>
<td width="180"><a href="https://unity2.ai/register?source=sub2api"><img src="assets/partners/logos/unity2.png" alt="unity2" width="150"></a></td>
<td>Thanks to Unity2 for sponsoring this project! <a href="https://unity2.ai/register?source=sub2api">Unity2</a> is a high-performance AI model API relay for individuals, teams, and enterprises, handling 30B+ tokens/day with 5000 RPM concurrency. One API Key works across Claude Code, Codex, OpenAI models, IDE plugins, and Agent workflows, with balance billing, bundled subscriptions, enterprise invoicing, and 1-on-1 support. <a href="https://unity2.ai/register?source=sub2api">Register</a> to claim $2 in balance, plus $10 more by joining the official group — up to $12 in free credit.
</td>
</tr>

<tr>
<td width="180"><a href="https://veilx.io/#/hello/SJRBRVDV"><img src="assets/partners/logos/veilx.png" alt="veilx" width="150"></a></td>
<td>Thanks to Veilx for sponsoring this project! <a href="https://veilx.io/#/hello/SJRBRVDV">Veilx</a> CDN is purpose-built for large-scale AI API traffic, deeply optimized for relay services and call chains across OpenAI, Claude, Gemini, and scenarios like chat, image generation, embeddings, and streaming — delivering lower latency and higher stability under heavy concurrency. It also offers China three-network optimized return lines, making it ideal for global AI relay platforms, overseas AI SaaS, and cross-border high-concurrency deployments.
</td>
</tr>

<tr>
<td width="180"><a href="https://roxybrowser.com/invite/bgGKG7"><img src="assets/partners/logos/RoxyBrowser.png" alt="veilx" width="150"></a></td>
<td>Thanks to RoxyBrowser for sponsoring this project! <a href="https://roxybrowser.com/invite/bgGKG7">RoxyBrowser</a> RoxyBrowser is the perfect partner for Sub2API: it features a built-in native Roxy AI Agent and high-quality native residential IPs, supports batch automation via simple commands, and significantly boosts security and efficiency for multi-account management! Click <a href="https://roxybrowser.com/invite/bgGKG7">this link</a> to sign up and receive a free residential IP package plus a 10% lifetime discount.
</td>
</tr>

<tr>
<td width="180"><a href="https://apikl.ai"><img src="assets/partners/logos/apikl.png" alt="apikl" width="150"></a></td>
<td>Thanks to Apikl for sponsoring this project! Built on Sub2API, the platform provides developers with relay services for Codex / Claude series models, focusing on long-term stability, high-speed direct connections, and excellent cost-effectiveness. It offers pay-as-you-go balance billing, enterprise-grade official invoices, and one-on-one dedicated support. <a href="https://apikl.ai">Register now</a> for a 1:1 top-up bonus — double your balance!
</td>
</tr>

<tr>
<td width="180"><a href="https://tokeneum.ai"><img src="assets/partners/logos/tokeneum.png" alt="tokeneum" width="150"></a></td>
<td>Thanks to TokenEum for sponsoring this project! <a href="https://tokeneum.ai">TokenEum</a> is a comprehensive AI model aggregation platform and intelligent agent development company. It brings together top-tier international models — including Claude, Gemini, and OpenAI — alongside leading open-source models such as GLM, Qwen, and Kimi, offering a wide range of options across different quality and price tiers to suit every need. TokenEum also provides access to cutting-edge video generation models like Seedance2.0 and Happy Horse. Committed to transparency and honest business practices, TokenEum ensures all model information is accurate and reliable. Visit <a href="https://tokeneum.ai">tokeneum.ai</a> to get started.
</td>
</tr>

<tr>
<td width="180"><a href="https://666api.work/sub2api"><img src="assets/partners/logos/666api.jpg" alt="666api" width="150"></a></td>
<td>Thanks to 666api for sponsoring this project! <a href="https://666api.work/sub2api">666api</a> is an all-in-one platform offering:<br>
⚡ API Relay — Pay-as-you-go access to global models sourced 100% from official providers, up to 75% off official pricing<br>
&nbsp;&nbsp;&nbsp;&nbsp;Exclusive: Zhipu GLM 50% off · DeepSeek V4-pro 50% off · Seedance 2.0 8% off (whitelisted) · HappyHorse Overseas 30% off (whitelisted)<br>
🔑 GPT Subscription Accounts (same-origin IP included) · Global Residential IP <br>
💰 Invoices supported
</td>
</tr>

<tr>
<td width="180"><a href="https://dis.chatdesks.cn/chatdesk/hsyqsub2api.html"><img src="assets/partners/logos/byteplus.png" alt="BytePlus" width="150"></a></td>
<td>Thanks to Dola seed for sponsoring this project! Dola Seed 2.0 is a full‑modal general large model independently developed by ByteDance for the global market. Built on a unified multimodal architecture, it supports joint understanding and generation of text, images, audio, and video. It natively enables agent collaboration, with strong reasoning, long‑task execution, tool integration, and coding capabilities. It is widely applicable to smart cockpits, personal assistants, education, customer support, marketing, retail, and other scenarios. It excels in multimodal perception, end‑to‑end complex task delivery, stable interaction, and data security, and is readily accessible and deployable via the ModelArk platform.Register via <a href="https://dis.chatdesks.cn/chatdesk/hsyqsub2api.html">this link</a> to get 500,000 tokens of free inference quota per model.<a href="https://dis.chatdesks.cn/chatdesk/hsyqsub2api.html"> >>中国大陆地区的开发者请点击这里</a></td>
</tr>

<tr>
<td width="180"><a href="https://sui-xiang.com/"><img src="assets/partners/logos/sui-xiang.jpg" alt="sui-xiang" width="150"></a></td>
<td>Thanks to Suixiang AI Gateway for sponsoring this project! <a href="https://sui-xiang.com/">Suixiang AI Gateway</a> is a reliable and efficient API relay service provider offering relay services for Claude, Codex, Gemini, and more. A privacy-focused relay — no data reselling, no model dilution; privacy, transparency, and lightning-fast after-sales support. New accounts get ¥0.5 in trial credit daily by signing in; top-ups are 1:1, no subscription required, pay-as-you-go. Multi-line redundancy, cross-region disaster recovery, automatic failover, and uninterrupted long-link SSE. 99.9% availability — critical calls never fall behind.
</td>
</tr>

<tr>
<td width="180"><a href="https://www.miyaip.com/?invitecode=sub2api"><img src="assets/partners/logos/miyaip.png" alt="miyaip" width="150"></a></td>
<td>Thanks to MiyaIP for sponsoring this project! <a href="https://www.miyaip.com/?invitecode=sub2api">MiyaIP</a> is a platform dedicated to global residential proxy network services, committed to providing high-quality, pure overseas residential IP resources for enterprise developers, cross-border business teams, and AI application users. It delivers stable, independent overseas network environments for AI platforms, overseas SaaS, and other online services, supporting multi-region access testing and project environment isolation. Ideal for development and testing scenarios that require access to overseas AI services, such as: AI model platform access, AI development testing, AI SaaS service usage, AI API debugging, and multi-region network environment validation.
</td>
</tr>

<tr>
<td width="180"><a href="https://anpin.ai"><img src="assets/partners/logos/anpin.jpg" alt="anpin" width="150"></a></td>
<td>Thanks to <a href="https://anpin.ai">anpin.ai</a> for sponsoring this project! anpin.ai is a premium AI relay service platform dedicated to advancing AI accessibility. With an advanced technical architecture and globally distributed deployment, it provides users with a direct high-speed channel to the world's top-tier large language models.<br>
Self-built primary account pool: 1-3s ultra-fast response, supports channel-partner distribution<br>
Extreme stability: multi-line intelligent routing + redundant backup system, ensuring year-round high-availability operation;<br>
Model authenticity: no content intervention or secondary filtering — experience the purest, most powerful native model capabilities.<br>
1:1 top-up, enterprise-grade service with invoicing available. Anpin AI is not just a relay — it's your secure, reliable, and efficient bridge to the frontier world of intelligence.
</td>
</tr>

<tr>
<td width="180"><a href="https://www.proxy4free.com/?keyword=4yjqecpc"><img src="assets/partners/logos/proxy4free.png" alt="proxy4free" width="150"></a></td>
<td>Thanks to Proxy4Free for sponsoring this project! Proxy4Free is a data proxy service provider for developers and AI applications, offering residential proxies, static residential proxies, ISP proxies, and datacenter proxies for scenarios such as Web Scraping, Browser Automation, and AI Agents. With global IP resources, stable connections, and flexible switching, it helps developers improve data collection success rates and reduce the risk of IP bans. Register via <a href="https://www.proxy4free.com/?keyword=4yjqecpc">this link</a> to get started and easily build more stable and efficient automation workflows.
</td>
</tr>

<tr>
<td width="180"><a href="http://www.fastaitoken.com/register"><img src="assets/partners/logos/fastaitoken.jpg" alt="fastaitoken" width="150"></a></td>
<td>🎉 Thanks to FastAIToken for sponsoring this project! <a href="http://www.fastaitoken.com/register">FastAIToken</a> is an AI API aggregation platform for developers, supporting mainstream large models such as OpenAI, Claude, and Gemini. Top-up at 1:1 — 1 CNY = 1 USD of API credit — letting developers use the world's leading large model services at lower cost and with greater convenience.<br>

🚀 The platform offers a variety of channels to choose from: an ultra-low-price 0.02x OpenAI promotional group (limited time), groups as low as 0.25x OpenAI, 0.7x Claude with 95% fixed cache, and a 1.2x Claude Max channel. It also provides a public status page showing real-time availability, latency, and operating status of each group for transparent and reliable service, plus 7×24 human technical support (not bots) with fast responses to developer needs.
</td>
</tr>

<tr>
<td width="180"><a href="http://aimzoon.com"><img src="assets/partners/logos/aimzoon.jpg" alt="aimzoon" width="150"></a></td>
<td>Thanks to Aimzoon for sponsoring this project! <a href="http://aimzoon.com">Aimzoon</a> provides stable, cost-effective AI API access services, enabling developers to quickly connect popular AI services to coding tools such as Codex, Claude Code, and Gemini CLI. No complex configuration — faster onboarding, more stable calls, and lower costs. Ongoing promotions including discounted Codex rates and special pricing, with free trial credits upon registration, bringing AI coding into your daily workflow. <a href="http://aimzoon.com">Click here</a> to register and try it out!
</td>
</tr>

</table>

## Overview

Sub2API is an AI API gateway platform designed to distribute and manage API quotas from AI product subscriptions. Users can access upstream AI services through platform-generated API Keys, while the platform handles authentication, billing, load balancing, and request forwarding.

## Features

- **Multi-Account Management** - Support multiple upstream account types (OAuth, API Key)
- **API Key Distribution** - Generate and manage API Keys for users
- **Precise Billing** - Token-level usage tracking and cost calculation
- **Smart Scheduling** - Intelligent account selection with sticky sessions
- **Concurrency Control** - Per-user and per-account concurrency limits
- **Rate Limiting** - Configurable request and token rate limits
- **Built-in Payment System** - Supports EasyPay, Alipay, WeChat Pay, and Stripe for user self-service top-up, no separate payment service needed ([Configuration Guide](docs/PAYMENT.md))
- **Admin Dashboard** - Web interface for monitoring and management
- **External System Integration** - Embed external systems (e.g. ticketing) via iframe to extend the admin dashboard

## Ecosystem

Community projects that extend or integrate with Sub2API:

| Project | Description | Features |
|---------|-------------|----------|
| ~~[Sub2ApiPay](https://github.com/touwaeriol/sub2apipay)~~ | ~~Self-service payment system~~ | **Now Built-in** — Payment is now integrated into Sub2API, no separate deployment needed. See [Payment Configuration Guide](docs/PAYMENT.md) |
| [sub2api-mobile](https://github.com/ckken/sub2api-mobile) | Mobile admin console | Cross-platform app (iOS/Android/Web) for user management, account management, monitoring dashboard, and multi-backend switching; built with Expo + React Native |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25.7, Gin, Ent |
| Frontend | Vue 3.4+, Vite 5+, TailwindCSS |
| Database | PostgreSQL 15+ |
| Cache/Queue | Redis 7+ |

---

## Nginx Reverse Proxy Note

When using Nginx as a reverse proxy for Sub2API (or CRS) with Codex CLI, add the following to the `http` block in your Nginx configuration:

```nginx
underscores_in_headers on;
```

Nginx drops headers containing underscores by default (e.g. `session_id`), which breaks sticky session routing in multi-account setups.

---

## Deployment

> **This repo is the fork `gthubtom1/sub2api-standby`.**  
> Full guide: **[docs/DEPLOY_STANDY_EN.md](docs/DEPLOY_STANDY_EN.md)** / **[中文](docs/DEPLOY_STANDY_CN.md)**  
> Commands below are fork-only. **No official install entrypoints.**

### Warnings

- **Do not** use Admin UI **Check for Updates / Update** (pulls official `Wei-Shaw/sub2api` binaries)
- **Do not** use official image `weishaw/sub2api` or official curl installers from `Wei-Shaw/sub2api`
- Upgrade only via this repo + `Dockerfile.custom`

### Recommended: Docker build from source

```bash
git clone https://github.com/gthubtom1/sub2api-standby.git
cd sub2api-standby
export DOCKER_BUILDKIT=1
docker build --memory=2g --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom -t sub2api-custom:0.1.157-standby .
cd deploy
cp .env.example .env && chmod 600 .env
mkdir -p data postgres_data redis_data
docker compose -f docker-compose.local.yml up -d
```

### One-click compose file prep (this fork only)

```bash
mkdir -p sub2api-deploy && cd sub2api-deploy
curl -sSL https://raw.githubusercontent.com/gthubtom1/sub2api-standby/main/deploy/docker-deploy.sh | bash
docker compose -f docker-compose.local.yml up -d
```

### Upgrade

```bash
git pull
docker build --memory=2g --build-arg VERSION=0.1.157-standby \
  -f Dockerfile.custom -t sub2api-custom:0.1.157-standby .
# recreate container / compose up -d
```

---

## Project Structure

```
sub2api/
├── backend/                  # Go backend service
│   ├── cmd/server/           # Application entry
│   ├── internal/             # Internal modules
│   │   ├── config/           # Configuration
│   │   ├── model/            # Data models
│   │   ├── service/          # Business logic
│   │   ├── handler/          # HTTP handlers
│   │   └── gateway/          # API gateway core
│   └── resources/            # Static resources
│
├── frontend/                 # Vue 3 frontend
│   └── src/
│       ├── api/              # API calls
│       ├── stores/           # State management
│       ├── views/            # Page components
│       └── components/       # Reusable components
│
└── deploy/                   # Deployment files
    ├── docker-compose.yml    # Docker Compose configuration
    ├── .env.example          # Environment variables for Docker Compose
    ├── config.example.yaml   # Full config file for binary deployment
    └── install.sh            # One-click installation script
```

## Star History

<a href="https://star-history.com/#gthubtom1/sub2api-standby&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=gthubtom1/sub2api-standby&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=gthubtom1/sub2api-standby&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=gthubtom1/sub2api-standby&type=Date" />
 </picture>
</a>

---

## License

This project is licensed under the [GNU Lesser General Public License v3.0](LICENSE) (or later).

Copyright (c) 2026 Wesley Liddick

---

<div align="center">

**If you find this project useful, please give it a star!**

</div>
