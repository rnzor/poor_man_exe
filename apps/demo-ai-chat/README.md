# AI Chat Demo

A real-time AI chat application showcasing Poor Man's exe.dev deployment platform.

## Features

- 🤖 **Real-time AI Chat** - WebSocket-based streaming responses
- 🚀 **FastAPI Backend** - Async Python with WebSocket support
- ⚛️ **React Frontend** - Beautiful, responsive UI
- 🐳 **Docker Native** - Fully containerized
- 🔄 **Auto-updates** - Watchtower integration
- 🔒 **Production Ready** - Health checks, resource limits, restart policies

## Quick Start

### Local Development

```bash
cd apps/demo-ai-chat

# Start with docker compose
docker compose up -d

# Or for development with hot reload
docker compose -f docker-compose.dev.yml up
```

Visit `http://localhost:8080`

### Deploy to Server

```bash
# On your VPS
cd /apps
git clone <your-repo>
cd apps/demo-ai-chat
docker compose up -d
```

### Configure Domain

1. Create `/apps/demo-ai-chat/Caddyfile`:
   ```caddyfile
   chat.yourdomain.com {
       reverse_proxy 127.0.0.1:8080
   }
   ```

2. Reload Caddy:
   ```bash
   caddy reload --config /apps/demo-ai-chat/Caddyfile
   ```

## Architecture

```
┌─────────────┐       ┌──────────────────┐       ┌─────────────────┐
│   Caddy     │──────▶│  Frontend (8080) │──────▶│  React UI       │
│  (HTTPS)    │       │  Nginx           │       │                 │
└─────────────┘       └──────────────────┘       └─────────────────┘
                            │
                            │ WebSocket
                            ▼
                     ┌──────────────────┐       ┌─────────────────┐
                     │  Backend (8000)  │──────▶│  FastAPI        │
                     │  FastAPI         │       │  WebSocket      │
                     └──────────────────┘       └─────────────────┘
```

## Endpoints

- **Frontend**: http://localhost:8080
- **Backend API**: http://localhost:8000
- **Health Check**: http://localhost:8000/health
- **WebSocket**: ws://localhost:8000/ws/{conversation_id}

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8000 | Backend port |
| `PYTHONUNBUFFERED` | 1 | Disable Python output buffering |
| `WS_ENDPOINT` | ws://localhost:8000 | WebSocket endpoint (frontend) |

## Deployment with Poor Man's exe.dev

```bash
# 1. SSH into your server
ssh user@yourdomain.com

# 2. Clone and deploy
cd /apps
git clone <your-repo>
cd apps/demo-ai-chat

# 3. Build and start
docker compose build
docker compose up -d

# 4. Configure Caddy
cat > /etc/caddy/Caddyfile << EOF
chat.yourdomain.com {
    reverse_proxy 127.0.0.1:8080
}
EOF

caddy reload
```

## CI/CD with Watchtower

Watchtower automatically updates containers when new images are pushed:

```bash
# Already included in bootstrap/setup.sh
docker run -d \
    --name watchtower \
    --restart unless-stopped \
    -v /var/run/docker.sock:/var/run/docker.sock \
    containrrr/watchtower:latest \
    --cleanup \
    --interval 300
```

## Performance

- **Backend**: Up to 0.5 CPU, 512MB RAM
- **Frontend**: Up to 0.25 CPU, 256MB RAM
- **Health Checks**: Every 30 seconds
- **Auto-restart**: On failure

## Customization

### Modify AI Responses

Edit `backend/main.py`:
```python
def generate_ai_response(user_message: str) -> str:
    # Your custom logic here
    return "Hello! I'm your custom AI assistant!"
```

### Change UI Colors

Edit `frontend/src/App.css`:
```css
.message.user .message-text {
    background: linear-gradient(135deg, #YOUR_COLOR 0%, #YOUR_COLOR2 100%);
}
```

## Troubleshooting

```bash
# View logs
docker compose logs -f

# Restart services
docker compose restart

# Check health
curl http://localhost:8000/health

# Check resources
docker stats
```

## License

MIT
