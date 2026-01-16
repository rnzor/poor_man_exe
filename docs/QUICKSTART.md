# Quick Start Guide

Get your Poor Man's exe.dev server running in 5 minutes!

## Prerequisites

- A fresh Ubuntu 22.04 or 24.04 VPS
- SSH access as root
- A domain name (optional, but recommended)

## Step 1: Bootstrap the Server

SSH into your server as root:

```bash
ssh root@your-server-ip
```

Download and run the bootstrap script:

```bash
cd /root
curl -LO https://raw.githubusercontent.com/yourusername/poor_man_exe/main/bootstrap/setup.sh
chmod +x setup.sh
./setup.sh
```

This will install:
- Docker + Docker Compose
- Caddy web server
- UFW firewall
- Fail2ban
- Watchtower (auto-updates)

## Step 2: Clone the Repository

```bash
cd /apps
git clone https://github.com/yourusername/poor_man_exe.git
```

## Step 3: Deploy the Demo App

```bash
cd poor_man_exe/apps/demo-ai-chat
docker compose up -d
```

Check if it's running:

```bash
curl http://localhost:8000/health
```

## Step 4: Configure a Domain (Optional)

Edit the Caddyfile:

```bash
cat > /etc/caddy/Caddyfile << EOF
chat.yourdomain.com {
    reverse_proxy 127.0.0.1:8080
}
EOF

caddy reload
```

Make sure your DNS points to your server IP!

## Step 5: Verify Everything

```bash
# Check containers
docker compose ps

# Check health
curl http://localhost:8000/health

# View logs
docker compose logs -f
```

## Common Commands

```bash
# Deploy an app
cd /apps/yourapp
docker compose up -d

# View logs
docker compose logs -f

# Restart
docker compose restart

# Update
docker compose pull
docker compose up -d --build

# Stop
docker compose down
```

## Using the Deployment Script

```bash
# Make the script executable
chmod +x scripts/deploy.sh

# Deploy an app
./scripts/deploy.sh demo-ai-chat up

# View logs
./scripts/deploy.sh demo-ai-chat logs -f

# Update
./scripts/deploy.sh demo-ai-chat update

# Check status
./scripts/deploy.sh demo-ai-chat status
```

## What's Next?

- [Add more apps](DEPLOYMENT.md)
- [Configure custom domains](DOMAINS.md)
- [Set up CI/CD](.github/workflows/deploy.yml)
- [Troubleshoot issues](TROUBLESHOOTING.md)

## Need Help?

- Check [Troubleshooting](TROUBLESHOOTING.md)
- Search existing issues
- Open a new issue on GitHub
