# Poor Man's exe.dev - Setup Guide

Follow these steps to deploy the platform on a single Linux VM.

## 1. Prerequisites
- A Linux VM (Ubuntu 22.04+ recommended)
- Docker installed
- Caddy installed (for auto-HTTPS)
- A domain name with wildcard DNS support (`*.yourdomain.com` pointing to the VM IP)

## 2. Platform Installation

### Bootstrap
Run the setup script provided in the repository:
```bash
./bootstrap/setup.sh
```

### Build the Gateway
```bash
go build -o poor-exe-gateway ./cmd/ssh-gateway
```

## 3. Configuration

Set environment variables in a `.env` file or your systemd unit:

- `DOMAIN`: Your base domain (e.g., `rnzlive.com`)
- `SSH_PORT`: Port for the gateway (default: 2222)
- `CADDY_URL`: Caddy Admin API URL (default: http://localhost:2019)
- `DB_PATH`: Path to SQLite database

## 4. Wildcard DNS & Caddy
To support `appname.yourdomain.com`, you need a wildcard Caddy configuration.

Example `Caddyfile`:
```caddyfile
{
    email your@email.com
}

# The gateway will dynamically add routes here via the Admin API
:80, :443 {
    # Default handler if no dynamic route matches
    header Content-Type text/plain
    respond "Poor Man's exe.dev - No app found here." 404
}
```

## 5. Systemd Service
Copy the service file from `deploy/systemd/poor-exe.service` to `/etc/systemd/system/`.

```bash
sudo cp deploy/systemd/poor-exe.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now poor-exe
```

## 6. Security Hardening
- **Firewall**: Ensure ports 22, 80, 443, and 2222 are open.
- **Fail2Ban**: Pre-configured by setup script to protect port 2222.
- **SSH Keys**: The gateway ONLY supports public key authentication. Add your key to the `public_keys` table in SQLite.
