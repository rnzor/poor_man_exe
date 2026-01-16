# Troubleshooting Guide

Common issues and their solutions.

## Docker Issues

### Container Won't Start

```bash
# Check logs
docker compose logs

# Check container status
docker compose ps

# Check for port conflicts
ss -tlnp | grep 3000
```

### Out of Memory

```bash
# Check memory usage
docker stats

# Increase swap
swapon -s
fallocate -l 2G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
```

### Disk Full

```bash
# Check disk usage
df -h

# Docker system prune
docker system prune -a

# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune
```

## Caddy Issues

### Caddy Won't Start

```bash
# Check Caddy status
systemctl status caddy

# Check logs
journalctl -u caddy -n 100 -f

# Validate Caddyfile
caddy validate
```

### HTTPS Not Working

```bash
# Check if port 443 is open
ufw status

# Verify DNS
dig yourdomain.com

# Check Caddy logs
journalctl -u caddy | grep -i cert
```

## Network Issues

### Can't Access App

```bash
# Check if container is running
docker compose ps

# Test locally
curl http://localhost:3000/health

# Check firewall
ufw status

# Check port binding
docker port container_name
```

### WebSocket Not Working

```bash
# Check WebSocket endpoint
curl -N -i http://localhost:3000/ws/test

# Verify proxy config in Caddy
curl -I http://localhost:3000/ws
```

## Performance Issues

### Slow Response Times

```bash
# Check container resources
docker stats

# Check server load
top
htop

# Check disk I/O
iostat -x 1
```

### High CPU Usage

```bash
# Find resource-heavy containers
docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Check for infinite loops in app logs
docker compose logs | grep -i error
```

## Application-Specific Issues

### Node.js App Crashes

```bash
# Check logs
docker compose logs app

# Check for memory issues
docker stats

# Increase memory limit in docker-compose.yml
```

### Python App Import Errors

```bash
# Rebuild container
docker compose build --no-cache

# Check requirements.txt
docker exec -it container_name pip list
```

## Security Issues

### SSH Brute Force Attempts

```bash
# Check Fail2ban status
fail2ban-client status

# Check banned IPs
fail2ban-client status sshd

# Unban IP
fail2ban-client set sshd unbanip 1.2.3.4
```

### Suspicious Activity

```bash
# Check auth logs
tail -f /var/log/auth.log

# Check failed login attempts
last | head -20
```

## Log Analysis

### View All Container Logs

```bash
# All containers
docker compose logs

# Specific service
docker compose logs app

# With timestamps
docker compose logs -t

# Last 100 lines
docker compose logs --tail 100
```

### Real-Time Log Monitoring

```bash
# Follow all logs
docker compose logs -f

# Follow specific service
docker compose logs -f app
```

## Recovery Procedures

### Reset All Containers

```bash
cd /apps
for app in */; do
    cd "$app"
    docker compose down
    cd ..
done
```

### Reset Caddy

```bash
systemctl stop caddy
rm -f /var/log/caddy/*
systemctl start caddy
```

### Full Server Reboot

```bash
# Graceful shutdown
shutdown -r now

# After reboot
docker compose -f /apps/*/docker-compose.yml up -d
```

## Getting Help

1. Check logs first: `docker compose logs`
2. Check Caddy: `journalctl -u caddy`
3. Search existing issues
4. Open a new issue with:
   - Error messages
   - Steps to reproduce
   - Server specs
   - Docker/Docker Compose versions
