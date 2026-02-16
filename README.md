# 🛡️ GlassBox Gateway

> **The Zero-Config Security Gateway for Developers.**
> Stop writing Nginx configs. Protect your internal tools, monitor traffic in real-time, and ban bots—all from a single binary.

![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/status-stable-brightgreen)
![Docker](https://img.shields.io/badge/docker-ready-blue?logo=docker&logoColor=white)

---

## ⚡ What is GlassBox?

**GlassBox** is a lightweight, high-performance **Reverse Proxy** written in Go. It is designed to replace complex Nginx/Traefik configurations for local development and internal infrastructure.

It acts as a **Smart Middleware** between the dangerous internet and your backend, instantly adding:
* **🔐 Authentication:** GitHub/Google OAuth2 (Zero-Trust).
* **🛡️ Active Defense:** Honeypots & WAF lite to trap bots.
* **👀 Observability:** A retro sci-fi Terminal User Interface (TUI).
* **🔒 HTTPS:** Automatic TLS certificates for `localhost`.

## 📸 Terminal Dashboard (TUI)

GlassBox provides real-time traffic monitoring directly in your terminal:

```text
┌──────────────── GLASSBOX GATEWAY v1.0 ────────────────┐
│ STATUS: ● ONLINE   |   UPTIME: 04:20:00               │
├───────────────────────────────────────────────────────┤
│ TRAFFIC:  ▂▃▅▇█ (54 req/sec)   LATENCY: 12ms          │
├───────────────────────────────────────────────────────┤
│ ACTIVE SESSIONS:                                      │
│ > user: alex_dev [Github] ──> /api/v1/data            │
│ > user: bot_1337 [BLOCKED] ──> /.env (Honeypot)       │
├───────────────────────────────────────────────────────┤
│ SECURITY LOG:                                         │
│ SQL Injection attempt from 192.168.1.5        │
│ New Login: admin (GitHub Auth)                │
└───────────────────────────────────────────────────────┘

```

## 🚀 Key Features

* **Zero-Config HTTPS:** Automatically generates valid self-signed certificates for local development. No more browser warnings.
* **Instant OAuth2:** Protect any service (Grafana, Portainer, Jenkins, Swagger) with GitHub Login just by setting one flag.
* **🍯 Honeypots (Active Defense):** Automatically exposes fake routes (e.g., `/wp-admin`, `.git`). If a bot touches them, their IP is banned instantly.
* **🐳 Docker Auto-Discovery:** Works like Traefik. Just add a label `glassbox.enable=true` to your containers.
* **Audit Logging:** Tracks *who* did *what*. "User Alex deleted the database" instead of "Anonymous IP request".

## 🛠️ Tech Stack

* **Core:** [Go](https://golang.org/) (Standard Library + `net/http`)
* **UI:** [Bubbletea](https://github.com/charmbracelet/bubbletea) (TUI Framework)
* **Crypto:** `acme/autocert` (Let's Encrypt integration) & `mkcert` logic.

## 🏁 Quick Start

### Option 1: Run Binary (Local Development)

Protect a local Python server running on port 5000 with GitHub Auth and HTTPS:

```bash
# Build
go build -o glassbox main.go

# Run (Generates localhost certs automatically)
./glassbox --target=http://localhost:5000 --auth=github --https

```

### Option 2: Docker Compose (Infrastructure)

Add GlassBox to your `docker-compose.yml` to protect your stack:

```yaml
version: '3'
services:
  glassbox:
    image: yourname/glassbox:latest
    ports:
      - "80:80"
      - "443:443"
    environment:
      - GITHUB_CLIENT_ID=...
      - GITHUB_SECRET=...
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock # For auto-discovery

  # Your internal app (protected)
  my-app:
    image: python:alpine
    labels:
      - "glassbox.enable=true"
      - "glassbox.auth=github"

```

## 🛡️ How it Works

GlassBox uses a **"Bouncer"** architecture:

1. **Intercept:** All traffic hits the Go middleware first.
2. **Identify:** Checks for valid encrypted Session Cookies.
3. **Trap:** If the request matches a Honeypot pattern -> **BAN IP**.
4. **Forward:** If safe, proxies the request to the upstream service.

## 🤝 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

[MIT](https://choosealicense.com/licenses/mit/)
