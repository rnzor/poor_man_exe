# 🚀 Poor Man's exe.dev

### *Deploy Docker containers like a boss. No Kubernetes required.*

---

<p align="center">
  <img src="https://img.shields.io/github/stars/rnzor/poor_man_exe?style=for-the-badge&logo=github" alt="Stars">
  <img src="https://img.shields.io/github/license/rnzor/poor_man_exe?style=for-the-badge&logo=mit" alt="License">
  <img src="https://img.shields.io/github/last-commit/rnzor/poor_man_exe?style=for-the-badge&logo=git" alt="Last Commit">
</p>

---

## 🎯 What Even Is This?

Remember [exe.dev](https://exe.dev)? That magical place where you spin up VMs with SSH and it just works?

This is that. But **you host it**. On your own VPS. For **free**. With Docker. Because why pay for something you can self-host?

It's like having your own personal app store, but for *your* servers. Deploy anything, anywhere, whenever you want.

### The Vibe

```
🤔 "I need to deploy this API quickly"
💭 "Hmm, should I use Kubernetes? Terraform? ArgoCD?"
🧠 [brain explodes]
🎉 "Nah, just use poor man's exe.dev"
```

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🚀 **Blazing Fast** | Deploy apps in seconds, not hours |
| 🔒 **Auto-HTTPS** | Caddy handles TLS like magic ✨ |
| 🔄 **Auto-Updates** | Watchtower keeps everything fresh |
| 📦 **Resource Guardrails** | No container can eat all your RAM |
| 🛡️ **Security Hardened** | Firewall + Fail2ban pre-configured |
| 📝 **Batteries Included** | Templates for Node, Python, Static, Fullstack |
| 🤖 **CI/CD Ready** | GitHub Actions workflows that actually work |
| 🎮 **Demo App** | A cool AI chat to show off to your friends |
| 🔑 **SSH Gateway CLI** | exe.dev-style CLI to manage apps via SSH |

---

## 🎮 The "It Just Works" Experience

```bash
# SSH into your server
ssh root@yourserver.com

# Run one command
./deploy.sh demo-ai-chat up

# Boom 💥
→ App is live at https://chat.yourdomain.com
```

That's it. No kubectl. No helm. No crying.

---

## 🚀 Get Started in 5 Minutes

### Step 1: Bootstrap Your Server

```bash
# SSH into your fresh Ubuntu 22.04/24.04 VPS as root
ssh root@YOUR_SERVER_IP

# Download and run the magic script
curl -LO https://raw.githubusercontent.com/rnzor/poor_man_exe/main/bootstrap/setup.sh
chmod +x setup.sh
./setup.sh
```

**What this does:**
- 🐳 Installs Docker + Docker Compose
- 🌐 Installs Caddy (auto-HTTPS wizard)
- 🔥 Sets up UFW firewall
- 🚫 Installs Fail2ban (SSH bodyguard)
- 📡 Installs Watchtower (auto-updates)
- 📁 Creates `/apps/` folder structure

### Step 2: Clone & Deploy

```bash
cd /apps
git clone https://github.com/rnzor/poor_man_exe.git
cd poor_man_exe/apps/demo-ai-chat
docker compose up -d
```

Visit `http://localhost:8080` and say hello to your new AI chat buddy! 🤖

### Step 3: (Optional) Add a Custom Domain

```bash
cat >> /etc/caddy/Caddyfile << 'EOF'

chat.yourdomain.com {
    reverse_proxy 127.0.0.1:8080
}
EOF

caddy reload
```

**Caddy will:**
1. Detect the domain
2. Provision an SSL certificate (free!)
3. Redirect HTTP → HTTPS automatically
4. Handle all the TLS nonsense

*Chef's kiss* 👨‍🍳✨

---

## 🎨 Demo App: AI Chat

A real-time, streaming AI chat application that screams "I know what I'm doing."

### ✨ Features

- ⚡ **WebSocket Streaming** - Character-by-character responses like ChatGPT
- 🎨 **Beautiful Dark Theme** - Cyberpunk vibes, smooth animations
- 📱 **Fully Responsive** - Works on your phone, laptop, and toaster
- 🐳 **Docker Native** - Backend + Frontend + Nginx, all containerized
- 💪 **Production Ready** - Health checks, resource limits, restart policies

### 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│                   Caddy                      │
│         (Auto-HTTPS Reverse Proxy)           │
└─────────────────────────────────────────────┘
                    │
         ┌──────────┴──────────┐
         │                     │
         ▼                     ▼
┌──────────────┐      ┌──────────────┐
│   Frontend   │      │   Backend    │
│   :8080      │      │   :8000      │
│  (React +    │      │ (FastAPI +   │
│   Nginx)     │      │  WebSocket)  │
└──────────────┘      └──────────────┘
```

### 🚀 Deploy the Demo

```bash
cd /apps/demo-ai-chat

# Start it up
docker compose up -d

# Check if it's healthy
curl http://localhost:8000/health
# → {"status":"healthy"}

# Watch the magic
docker compose logs -f
```

---

## 📁 Repository Structure

```
poor_man_exe/
├── 📄 README.md                    # You are here
├── 📄 AGENTS.md                    # For AI assistants 🤖
│
├── 📂 bootstrap/
│   └── setup.sh                    # ⭐ The magic bootstrap script
│
├── 📂 templates/
│   ├── Caddyfile                   # Caddy configuration templates
│   ├── nginx.conf                  # Nginx config for static sites
│   └── docker-compose/
│       ├── nodejs.yml              # Node.js/Express API template
│       ├── python.yml              # Python/FastAPI template
│       ├── static.yml              # Static site template
│       └── fullstack.yml           # Backend + Frontend template
│
├── 📂 apps/
│   └── demo-ai-chat/               # ⭐ The cool demo app
│       ├── backend/                # FastAPI + WebSocket
│       │   ├── Dockerfile
│       │   ├── main.py             # The brains
│       │   └── requirements.txt
│       ├── frontend/               # React + TypeScript
│       │   ├── Dockerfile
│       │   ├── src/
│       │   │   ├── App.tsx         # The beauty
│       │   │   └── App.css         # The style
│       │   └── vite.config.ts
│       ├── docker-compose.yml      # Orchestrates everything
│       └── Caddyfile               # Domain config
│
├── 📂 scripts/
│   ├── deploy.sh                   # ⭐ Deployment Swiss Army knife
│   └── logs.sh                     # Quick log viewer
│
├── 📂 .github/workflows/
│   ├── build-push.yml              # Build & push to GHCR
│   └── deploy.yml                  # Auto-deploy on release
│
└── 📂 docs/
    ├── QUICKSTART.md               # Get running fast
    ├── DEPLOYMENT.md               # How to add apps
    ├── DOMAINS.md                  # Domain + HTTPS config
    └── TROUBLESHOOTING.md          # When things go sideways
```

---

## 🛠️ Deploy Your Own App

### 1. Copy a Template

```bash
cd /apps
mkdir my-awesome-app
cd my-awesome-app

# Pick your poison:
cp ../poor_man_exe/templates/docker-compose/nodejs.yml docker-compose.yml
# or
cp ../poor_man_exe/templates/docker-compose/python.yml docker-compose.yml
# or
cp ../poor_man_exe/templates/docker-compose/static.yml docker-compose.yml
```

### 2. Customize It

Edit `docker-compose.yml`:

```yaml
services:
  app:
    image: ghcr.io/rnzor/my-awesome-app:latest  # <- your image
    ports:
      - "3000:3000"                              # <- your port
    environment:
      - NODE_ENV=production
```

### 3. Deploy

```bash
docker compose up -d
```

### 4. Configure Domain

```bash
cat >> /etc/caddy/Caddyfile << 'EOF'

myapp.yourdomain.com {
    reverse_proxy 127.0.0.1:3000
}
EOF

caddy reload
```

**Done!** 🎉 Your app is now live with HTTPS.

---

## 🔑 SSH Gateway CLI

The SSH gateway provides an exe.dev-style CLI experience. Connect via SSH and manage your apps directly:

```bash
# Connect to the gateway
ssh -p 2222 poor-exe@yourserver.com

# Or run commands directly
ssh -p 2222 poor-exe@yourserver.com ls
```

### Available Commands

| Command | Description |
|---------|-------------|
| `ls` | List your apps with Docker status |
| `new --name=X [--image=Y]` | Create a new app container |
| `rm <app>` | Delete an app and its container |
| `share <cmd> <vm>` | Manage sharing (public/private/port) |
| `keys [add\|rm]` | Manage SSH keys |
| `whoami` | Show current user info |
| `help` | Show available commands |
| `exit` | Disconnect |

### Examples

```bash
# Create a new app with custom image
ssh -p 2222 poor-exe@server.com new --name=myapi --image=node:20-alpine

# List all your apps
ssh -p 2222 poor-exe@server.com ls

# Get JSON output (for scripting)
ssh -p 2222 poor-exe@server.com ls --json

# Attach to an app's shell
ssh -p 2222 myapi@server.com

# Set an app to public
ssh -p 2222 poor-exe@server.com share set-public myapi

# Change the HTTP port
ssh -p 2222 poor-exe@server.com share port myapi 3000

# Add a new SSH key
ssh -p 2222 poor-exe@server.com keys add "ssh-ed25519 AAAA... user@host"

# Remove an app
ssh -p 2222 poor-exe@server.com rm myapi
```

### JSON Output

All commands support `--json` flag for automation:

```bash
ssh -p 2222 poor-exe@server.com ls --json
```
```json
{
  "success": true,
  "data": {
    "vms": [
      {"vm_name": "myapi", "image": "node:20-alpine", "status": "running", "created_at": "2026-01-17"}
    ]
  }
}
```

---

## 🎣 Using the Deployment Script

Our Swiss Army knife for deployments:

```bash
./deploy.sh <app-name> <command> [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `up` | Start the app (default) |
| `down` | Stop the app |
| `restart` | Restart the app |
| `logs -f` | View logs (follow mode) |
| `status` | Check app status |
| `update` | Pull latest + rebuild |
| `pull` | Pull images only |
| `build` | Build images |
| `ps` | List containers |
| `stop` | Stop containers |
| `start` | Start containers |
| `clean` | Remove everything |

### Examples

```bash
# Deploy
./deploy.sh demo-ai-chat up

# Watch logs
./deploy.sh demo-ai-chat logs -f

# Update to latest
./deploy.sh demo-ai-chat update

# Check if it's alive
./deploy.sh demo-ai-chat status --port 8080

# Restart everything
./deploy.sh demo-ai-chat restart

# Nuclear option
./deploy.sh demo-ai-chat clean
```

---

## 🔄 Auto-Updates (Watchtower)

Watchtower is pre-installed and runs in the background. It:

- ⏰ Polls every 5 minutes
- 📦 Checks for new images
- 🔄 Restarts containers automatically
- 🧹 Cleans up old images

**Manual trigger:**
```bash
docker restart watchtower
```

---

## 🤖 CI/CD (GitHub Actions)

### Workflow 1: Build & Push

**Trigger:** Push to `main`

**What it does:**
1. Builds both backend and frontend images
2. Pushes to GitHub Container Registry (GHCR)
3. Tags with `latest` and semantic versions

### Workflow 2: Deploy

**Trigger:** Release published OR manual trigger

**What it does:**
1. SSHs into your server
2. Pulls latest code/images
3. Rebuilds and restarts containers
4. Verifies deployment

### Required Secrets

Add these in GitHub Settings → Secrets:

| Secret | Value |
|--------|-------|
| `SERVER_HOST` | Your server IP |
| `SERVER_USER` | `root` (or your SSH user) |
| `SERVER_SSH_KEY` | Your private SSH key |

---

## 📊 Resource Limits

All templates include sane defaults to prevent container chaos:

```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'        # Half a CPU
      memory: 512M      # 512MB RAM
    reservations:
      cpus: '0.1'
      memory: 128M
```

Adjust based on your app's needs. Some apps are hungry. We don't judge.

---

## 🔐 Security Features

| Feature | What It Does |
|---------|--------------|
| **UFW Firewall** | Blocks everything except 22, 80, 443 |
| **Fail2ban** | Bans bots that brute-force your SSH |
| **Non-root Containers** | No container runs as root |
| **Auto-HTTPS** | TLS 1.3 with modern cipher suites |
| **Resource Limits** | Containers can't eat all your RAM |
| **Ed25519 Host Keys** | Modern, high-security SSH host keys |
| **SHA256 Fingerprints** | Proper SSH key fingerprint calculation |
| **Rate Limiting** | Token bucket rate limiting per IP |
| **Session Limits** | Max 10 concurrent sessions per key |
| **Audit Logging** | All operations logged with IP + timestamp |

---

## 🎯 Comparison Table

| Feature | exe.dev | Poor Man's exe.dev |
|---------|---------|-------------------|
| 💸 Free | ❌ | ✅ |
| 🏠 Self-hosted | ❌ | ✅ |
| 🔒 Auto-HTTPS | ✅ | ✅ |
| 🐳 Docker | ✅ | ✅ |
| 🔌 WebSockets | ✅ | ✅ |
| 📊 Resource limits | ✅ | ✅ |
| 🚀 Fast deploy | ✅ | ✅ |
| 🤖 Custom infrastructure | ❌ | ✅ |

We're basically the open-source, free, self-hosted version. You're welcome.

---

## 📚 Documentation

- [📖 Quick Start](docs/QUICKSTART.md) - Get running in 5 minutes
- [🚀 Deployment Guide](docs/DEPLOYMENT.md) - Add new applications
- [🌐 Domain Config](docs/DOMAINS.md) - HTTPS + custom domains
- [🐛 Troubleshooting](docs/TROUBLESHOOTING.md) - Fix stuff

---

## 🧰 Useful Commands

```bash
# Deployment script
./deploy.sh demo-ai-chat up
./deploy.sh demo-ai-chat logs -f
./deploy.sh demo-ai-chat update

# Docker directly
docker compose up -d              # Start
docker compose down               # Stop
docker compose logs -f            # Logs
docker compose pull               # Update images
docker compose up -d --build      # Rebuild
docker stats                      # Check resources

# Caddy
caddy reload                      # Reload config
caddy validate                    # Check config
journalctl -u caddy -f            # Caddy logs
```

---

## 🤝 Contributing

1. 🍴 Fork the repo
2. 🌿 Create a feature branch
3. ✨ Make it awesome
4. 📝 Commit with a fun message
5. 📤 Push to your fork
6. 🔀 Open a PR

All contributions welcome! Especially ones that make this even more "exe-like."

---

## 📝 License

MIT License - Fork it, break it, fix it, deploy it.

---

## 🐛 Support

- 📖 [Docs](docs/QUICKSTART.md)
- 🐛 [Issues](https://github.com/rnzor/poor_man_exe/issues)
- 💬 [Discussions](https://github.com/rnzor/poor_man_exe/discussions)

---

<p align="center">
  Made with ❤️, ☕, and too much Docker
  <br><br>
  <a href="https://github.com/rnzor/poor_man_exe">⭐ Star us on GitHub!</a>
</p>
