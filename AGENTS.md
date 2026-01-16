# Poor Man's exe.dev - Agent Instructions

This document provides instructions for AI agents working with this repository.

## Overview

Poor Man's exe.dev is a self-hosted deployment platform that replicates key features of [exe.dev](https://exe.dev). It enables quick deployment of Docker containers with auto-HTTPS, similar to the "exe experience."

## Key Components

### 1. Bootstrap Script (`bootstrap/setup.sh`)

**Purpose:** Provision a fresh Ubuntu 22.04/24.04 VPS with all necessary components.

**What it installs:**
- Docker CE + Docker Compose plugin
- Caddy web server (auto-HTTPS)
- UFW firewall
- Fail2ban
- Watchtower (auto-updates)
- Creates folder structure (`/apps/`, `/apps/_caddy/`, `/apps/_shared/`)

**Usage:**
```bash
curl -LO https://raw.githubusercontent.com/USER/REPO/main/bootstrap/setup.sh
chmod +x setup.sh
./setup.sh
```

### 2. Templates (`templates/`)

**Purpose:** Reusable configuration templates for common app types.

**Available templates:**
- `docker-compose/nodejs.yml` - Node.js/Express API
- `docker-compose/python.yml` - Python/FastAPI
- `docker-compose/static.yml` - Static sites (nginx)
- `docker-compose/fullstack.yml` - Backend + frontend
- `Caddyfile` - Reverse proxy configurations

**Usage:** Copy to your app folder and customize.

### 3. Demo App (`apps/demo-ai-chat/`)

**Purpose:** Showcase application demonstrating the platform's capabilities.

**Features:**
- FastAPI backend with WebSocket support
- React frontend with real-time chat
- Docker Compose orchestration
- Health checks and resource limits

**Architecture:**
```
Caddy (8080) → Frontend (nginx) → Backend (8000, WebSocket)
```

**Deploy:**
```bash
cd /apps/demo-ai-chat
docker compose up -d
```

### 4. Deployment Scripts (`scripts/deploy.sh`)

**Purpose:** Simplify common deployment operations.

**Usage:**
```bash
./deploy.sh <app-name> <command> [options]

Commands:
  up          Start the app
  down        Stop the app
  restart     Restart the app
  logs        View logs (-f to follow)
  status      Check status
  update      Pull and rebuild
  pull        Pull images only
  build       Build images
  ps          List containers
  stop        Stop containers
  start       Start containers
  clean       Remove everything
```

### 5. GitHub Actions (`.github/workflows/`)

**Workflows:**
- `build-push.yml` - Build and push Docker images to GHCR
- `deploy.yml` - Deploy to server on release

**Triggers:**
- Push to main → Build and push
- Release published → Deploy to server

**Required secrets:**
- `SERVER_HOST` - Server IP
- `SERVER_USER` - SSH username
- `SERVER_SSH_KEY` - Private SSH key

## Common Tasks

### Deploy a New App

1. Create app folder: `mkdir /apps/myapp`
2. Copy template: `cp templates/docker-compose/nodejs.yml docker-compose.yml`
3. Customize: Edit `docker-compose.yml` with your image/port
4. Deploy: `docker compose up -d`
5. Configure domain: Add to `/etc/caddy/Caddyfile`
6. Reload Caddy: `caddy reload`

### Update an App

```bash
cd /apps/myapp
docker compose pull
docker compose up -d --build
```

### View Logs

```bash
docker compose logs -f
# or
./deploy.sh myapp logs -f
```

### Check Health

```bash
curl http://localhost:PORT/health
```

### Configure Domain

Add to `/etc/caddy/Caddyfile`:
```caddyfile
app.yourdomain.com {
    reverse_proxy 127.0.0.1:PORT
}
```

Then: `caddy reload`

## Security Considerations

1. **Never commit secrets** - Use environment variables
2. **Use specific image tags** - Avoid `latest`
3. **Set resource limits** - Prevent container monopolization
4. **Enable firewall** - UFW is pre-configured
5. **Regular updates** - Watchtower handles auto-updates

## Troubleshooting

### Container won't start
```bash
docker compose logs
docker compose ps
```

### App not accessible
```bash
curl http://localhost:PORT/health
ufw status
docker compose port app 80
```

### Caddy issues
```bash
caddy validate
journalctl -u caddy -n 100 -f
```

### Disk full
```bash
docker system prune -a
df -h
```

## Best Practices

1. **One app per folder** in `/apps/`
2. **Use health checks** in all containers
3. **Set resource limits** in docker-compose
4. **Monitor with logs** regularly
5. **Keep images updated** with Watchtower
6. **Backup important data** (not included in this repo)

## File Structure

```
poor_man_exe/
├── bootstrap/
│   └── setup.sh              # Server setup
├── templates/
│   ├── Caddyfile             # Caddy config templates
│   ├── docker-compose/
│   │   ├── nodejs.yml
│   │   ├── python.yml
│   │   ├── static.yml
│   │   └── fullstack.yml
│   └── nginx.conf
├── apps/
│   └── demo-ai-chat/         # Demo app
├── scripts/
│   └── deploy.sh             # Deployment helper
├── .github/workflows/
│   ├── build-push.yml
│   └── deploy.yml
├── docs/
│   ├── QUICKSTART.md
│   ├── DEPLOYMENT.md
│   ├── DOMAINS.md
│   └── TROUBLESHOOTING.md
├── README.md
└── .gitignore
```

## Notes for AI Agents

1. When helping users, start with the quick start guide
2. Always verify commands before suggesting them
3. Remind users about security (don't commit secrets)
4. Use the deployment script when possible for consistency
5. Check existing docs before writing new content
6. When debugging, ask for logs first
7. Suggest specific commands, not general guidance
8. Remember this is self-hosted - users control everything

## Version

Current version: 1.0.0

Last updated: 2024
