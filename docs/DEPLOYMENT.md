# Deployment Guide

Learn how to deploy new applications on your Poor Man's exe.dev server.

## Folder Structure

```
/apps/
├── demo-ai-chat/          # Your app
│   ├── docker-compose.yml
│   ├── backend/
│   ├── frontend/
│   └── Caddyfile
├── my-new-app/            # Another app
│   ├── docker-compose.yml
│   └── ...
└── _caddy/                # Caddy configurations
```

## Adding a New App

### 1. Create the App Folder

```bash
cd /apps
mkdir my-new-app
cd my-new-app
```

### 2. Add Docker Compose File

Use one of the templates from `templates/docker-compose/`:

```bash
cp ../poor_man_exe/templates/docker-compose/nodejs.yml docker-compose.yml
```

Edit `docker-compose.yml` with your app details:

```yaml
services:
  app:
    image: ghcr.io/yourusername/yourapp:latest
    ports:
      - "3000:3000"
```

### 3. Create Caddyfile

```bash
cp ../poor_man_exe/templates/Caddyfile Caddyfile
```

Edit and add to Caddy config:

```bash
cat >> /etc/caddy/Caddyfile << EOF

myapp.yourdomain.com {
    reverse_proxy 127.0.0.1:3000
}
EOF

caddy reload
```

### 4. Deploy

```bash
docker compose up -d
```

## Using Templates

### Node.js API

```bash
cp templates/docker-compose/nodejs.yml docker-compose.yml
```

### Python/FastAPI

```bash
cp templates/docker-compose/python.yml docker-compose.yml
```

### Static Site

```bash
cp templates/docker-compose/static.yml docker-compose.yml
```

### Fullstack (Backend + Frontend)

```bash
cp templates/docker-compose/fullstack.yml docker-compose.yml
```

## Environment Variables

Create a `.env` file in your app folder:

```bash
# .env
IMAGE=ghcr.io/yourusername/yourapp:latest
PORT=3000
DATABASE_URL=postgresql://user:pass@localhost:5432/db
JWT_SECRET=your-secret-key
```

## Health Checks

All templates include health checks. Verify:

```bash
curl http://localhost:YOUR_PORT/health
```

## Resource Limits

Templates include Docker deploy resources. Adjust as needed:

```yaml
deploy:
  resources:
    limits:
      cpus: '1.0'
      memory: 1G
    reservations:
      cpus: '0.5'
      memory: 512M
```

## Auto-Updates (Watchtower)

Watchtower is already installed and will:
- Poll every 5 minutes
- Auto-restart containers when new images are pushed
- Clean up old images

To manually trigger an update:

```bash
docker restart watchtower
```

Or update a specific container:

```bash
docker compose pull
docker compose up -d
```

## Multi-Container Apps

Example with database:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: app
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d app"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
```

## Best Practices

1. **Use specific image tags** (not `latest`)
2. **Set resource limits** to prevent container monopolization
3. **Add health checks** for monitoring
4. **Use .env files** for secrets
5. **Regularly update** base images
6. **Monitor logs** with `docker compose logs`

## Deploying from GitHub Actions

See [GitHub Actions workflow](../.github/workflows/deploy.yml)

Required secrets:
- `SERVER_HOST` - Server IP
- `SERVER_USER` - SSH username
- `SERVER_SSH_KEY` - Private SSH key
