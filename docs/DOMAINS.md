# Domain Configuration

Configure custom domains with automatic HTTPS.

## Basic Setup

### 1. Point DNS to Your Server

Add an A record:
- **Type:** A
- **Name:** @ (or your subdomain)
- **Value:** Your server IP
- **TTL:** 3600 (or lower for testing)

### 2. Configure Caddy

Edit `/etc/caddy/Caddyfile`:

```caddyfile
# Simple setup
yourdomain.com {
    reverse_proxy 127.0.0.1:3000
}

# With subdomain
app.yourdomain.com {
    reverse_proxy 127.0.0.1:3001
}
```

Reload Caddy:

```bash
caddy reload
```

## Subdomains

### Multiple Apps on One Domain

```caddyfile
# /etc/caddy/Caddyfile

api.yourdomain.com {
    reverse_proxy 127.0.0.1:3001
}

chat.yourdomain.com {
    reverse_proxy 127.0.0.1:3002

    # WebSocket support
    websocket {
        header_upstream Connection "upgrade"
        header_upstream Upgrade websocket
    }
}

www.yourdomain.com {
    redir https://yourdomain.com{uri}
}

yourdomain.com {
    reverse_proxy 127.0.0.1:3000
}
```

## Cloudflare Setup

If using Cloudflare with proxy enabled, use DNS challenge:

```caddyfile
yourdomain.com {
    reverse_proxy 127.0.0.1:3000
    tls {
        dns cloudflare {env.CLOUDFLARE_API_TOKEN}
    }
}
```

Set the environment variable:

```bash
echo "CLOUDFLARE_API_TOKEN=your-token" >> /etc/environment
```

## Security Headers

Add headers for improved security:

```caddyfile
yourdomain.com {
    reverse_proxy 127.0.0.1:3000

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Frame-Options "SAMEORIGIN"
        X-Content-Type-Options "nosniff"
        X-XSS-Protection "1; mode=block"
        Referrer-Policy "strict-origin-when-cross-origin"
    }
}
```

## HTTP to HTTPS Redirect

Caddy does this automatically!

## Wildcard Certificates

Requires DNS provider setup:

```caddyfile
*.yourdomain.com {
    tls {
        dns cloudflare {env.CLOUDFLARE_API_TOKEN}
    }

    @app host app.yourdomain.com
    reverse_proxy @app 127.0.0.1:3000

    @api host api.yourdomain.com
    reverse_proxy @api 127.0.0.1:3001
}
```

## Per-App Caddyfile

For cleaner config, create per-app Caddyfiles:

```bash
# /apps/yourapp/Caddyfile
yourdomain.com {
    reverse_proxy 127.0.0.1:3000
}
```

Then include in main Caddyfile:

```caddyfile
# /etc/caddy/Caddyfile
(apps) {
    import /apps/yourapp/Caddyfile
    import /apps/otherapp/Caddyfile
}

yourdomain.com {
    import apps
}
```

## Troubleshooting

### Certificate Issues

```bash
# Check Caddy logs
journalctl -u caddy -n 100 -f

# Test manually
curl -I https://yourdomain.com
```

### DNS Not Resolving

```bash
# Check DNS
dig yourdomain.com

# Verify server IP
curl ifconfig.me
```

### Port Not Listening

```bash
# Check if app is running
docker compose ps

# Check port
ss -tlnp | grep 3000
```

## SSL Certificate Info

Caddy automatically:
- Requests certificates from Let's Encrypt
- Renews before expiry
- Serves HTTPS on port 443

View certificate details:

```bash
echo | openssl s_client -servername yourdomain.com -connect yourdomain.com:443 2>/dev/null | openssl x509 -noout -dates
```
