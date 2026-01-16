#!/bin/bash
set -euo pipefail

###############################################################################
# Poor Man's exe.dev - Bootstrap Script
# Turns a fresh Ubuntu 22.04/24.04 VPS into a deployment platform
###############################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

check_os() {
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        if [[ "$ID" != "ubuntu" ]]; then
            log_warn "This script is designed for Ubuntu. It may work on Debian but is not tested."
        fi
        if [[ "$VERSION_ID" != "22.04" && "$VERSION_ID" != "24.04" ]]; then
            log_warn "This script is tested on Ubuntu 22.04 and 24.04. You have $VERSION_ID"
        fi
    fi
}

update_system() {
    log_info "Updating system packages..."
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq
    apt-get upgrade -y -qq
}

install_basic_deps() {
    log_info "Installing basic dependencies..."
    apt-get install -y -qq \
        ca-certificates \
        curl \
        gnupg \
        ufw \
        fail2ban \
        git \
        vim \
        htop \
        unzip \
        wget \
        apt-transport-https \
        ca-certificates \
        gnupg2 \
        lsb-release
}

install_docker() {
    log_info "Installing Docker..."

    if command -v docker &> /dev/null; then
        log_warn "Docker already installed. Skipping..."
        return 0
    fi

    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    UBUNTU_CODENAME="$(. /etc/os-release && echo "$VERSION_CODENAME")"
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
        https://download.docker.com/linux/ubuntu ${UBUNTU_CODENAME} stable" \
        > /etc/apt/sources.list.d/docker.list

    apt-get update -qq
    apt-get install -y -qq \
        docker-ce \
        docker-ce-cli \
        containerd.io \
        docker-buildx-plugin \
        docker-compose-plugin

    systemctl enable --now docker
    usermod -aG docker $SUDO_USER || true

    log_info "Docker installed successfully"
}

install_caddy() {
    log_info "Installing Caddy..."

    if command -v caddy &> /dev/null; then
        log_warn "Caddy already installed. Skipping..."
        return 0
    fi

    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' > /etc/apt/sources.list.d/caddy-stable.list

    apt-get update -qq
    apt-get install -y -qq caddy

    systemctl enable --now caddy

    log_info "Caddy installed successfully"
}

setup_folders() {
    log_info "Creating folder structure..."

    mkdir -p /apps
    mkdir -p /apps/_caddy
    mkdir -p /apps/_shared
    mkdir -p /var/log/caddy
    mkdir -p /var/log/apps

    chown -R caddy:caddy /var/log/caddy 2>/dev/null || true
    chmod -R 755 /apps

    log_info "Folders created:"
    echo "  /apps          - Your application deployments"
    echo "  /apps/_caddy   - Caddy configuration"
    echo "  /apps/_shared  - Shared data between apps"
    echo "  /var/log/caddy - Caddy logs"
    echo "  /var/log/apps  - Application logs"
}

configure_firewall() {
    log_info "Configuring firewall (UFW)..."

    ufw default deny incoming
    ufw default allow outgoing

    ufw allow OpenSSH
    ufw allow 80/tcp
    ufw allow 443/tcp

    ufw --force enable

    log_info "Firewall configured"
}

configure_fail2ban() {
    log_info "Configuring Fail2ban..."

    cat > /etc/fail2ban/jail.local << 'EOF'
[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
findtime = 600

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 3
bantime = 3600
findtime = 600
EOF

    systemctl enable --now fail2ban

    log_info "Fail2ban configured"
}

install_watchtower() {
    log_info "Installing Watchtower for auto-updates..."

    docker run -d \
        --name watchtower \
        --restart unless-stopped \
        -v /var/run/docker.sock:/var/run/docker.sock \
        containrrr/watchtower:latest \
        --cleanup \
        --interval 300 \
        --include-restarting

    log_info "Watchtower installed (polls every 5 minutes)"
}

create_caddy_base_config() {
    log_info "Creating base Caddy configuration..."

    cat > /etc/caddy/Caddyfile << 'EOF'
{
    admin off
    auto_https off
}

(log) {
    log {
        format json {
            time_format "rfc3339"
        }
        output file /var/log/caddy/access.log {
            roll_size 100MiB
            roll_keep 5
            roll_keep_gzip 5
        }
    }
}

:80 {
    respond "Poor Man's exe.dev - Server is running!" 200
}
EOF

    caddy reload --config /etc/caddy/Caddyfile

    log_info "Base Caddy configuration created"
}

print_summary() {
    echo ""
    echo "=============================================================================="
    echo -e "${GREEN}✅ Bootstrap Complete!${NC}"
    echo "=============================================================================="
    echo ""
    echo "Installed versions:"
    docker --version
    docker compose version
    caddy version
    echo ""
    echo "Next steps:"
    echo "  1. Add your domain to /etc/caddy/Caddyfile"
    echo "  2. Create your app in /apps/"
    echo "  3. Configure DNS to point to this server"
    echo ""
    echo "Useful commands:"
    echo "  docker compose up -d      # Start an app"
    echo "  docker compose logs -f    # View logs"
    echo "  caddy reload              # Reload Caddy config"
    echo ""
    echo "Docs: /apps/README.md"
    echo "=============================================================================="
}

main() {
    log_info "Starting Poor Man's exe.dev bootstrap..."
    log_info "This may take a few minutes..."

    check_root
    check_os
    update_system
    install_basic_deps
    install_docker
    install_caddy
    setup_folders
    configure_firewall
    configure_fail2ban
    install_watchtower
    create_caddy_base_config
    print_summary
}

main "$@"
