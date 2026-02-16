# 🛡️ GlassBox Gateway

> **The Zero-Config Security Gateway for Developers.**
> Protect your internal tools, monitor traffic in real-time, and ban bots—all from a single binary.

![Go Version](https://img.shields.io/github/go-mod/go-version/username/glassbox)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/status-stable-brightgreen)

## ⚡ What is GlassBox?

GlassBox is a lightweight **Reverse Proxy** written in Go. It is designed to replace complex Nginx configurations for local development and internal tools.

It acts as a **Middleware** between the internet and your backend, adding:
* **Authentication:** GitHub/Google OAuth2.
* **Security:** Honeypots & WAF lite.
* **Observability:** A retro sci-fi Terminal User Interface (TUI).

## 📸 Dashboard (TUI)

*(Insert a screenshot of your terminal dashboard here)*
`[ TUI Screenshot Placeholder ]`

## 🚀 Key Features

* **🔒 Zero-Config HTTPS:** Automatically generates valid self-signed certificates for `localhost` development. No more "Insecure Connection" warnings.
* **👤 Instant OAuth2:** Protect any service (Grafana, Portainer, Jenkins) with GitHub Login just by setting one flag.
* **🍯 Active Defense (Honeypots):** Automatically creates fake routes (e.g., `/wp-admin`, `.env`). If a bot touches them, their IP is banned instantly.
* **🐳 Docker Auto-Discovery:** Works like Traefik. Just add a label `glassbox.enable=true` to your containers.
* **📟 TUI Dashboard:** Monitor requests per second (RPS), latency, and active users directly in your terminal.

## 🛠️ Built With

* **Language:** [Go](https://golang.org/) (Standard Library + `net/http`)
* **UI:** [Bubbletea](https://github.com/charmbracelet/bubbletea) (TUI Framework)
* **Crypto:** `acme/autocert` (Let's Encrypt integration)

## 🏁 Quick Start

### Option 1: Run Binary
```bash
# Protect a local Python server running on port 5000
./glassbox --target=http://localhost:5000 --auth=github --https
Option 2: Docker Compose

YAML
version: '3'
services:
  glassbox:
    image: glassbox:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
🛡️ Security Logic
GlassBox uses a "Bouncer" architecture:

Intercept: All traffic passes through the Go middleware.

Verify: Checks for valid Session Cookies.

Trap: If the request matches a Honeypot pattern -> BAN.

Forward: If safe, proxies the request to the upstream service.

🤝 Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

📄 License
MIT
