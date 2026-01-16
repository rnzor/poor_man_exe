#!/bin/bash
set -euo pipefail

###############################################################################
# Poor Man's exe.dev - Deployment Script
# Usage: ./deploy.sh <app-name> [command]
###############################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

APPS_DIR="/apps"
SCRIPT_VERSION="1.0.0"

log() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

usage() {
    echo "Poor Man's exe.dev - Deployment Script v${SCRIPT_VERSION}"
    echo ""
    echo "Usage: $0 <app-name> [command] [options]"
    echo ""
    echo "Commands:"
    echo "  up          Start the app (default)"
    echo "  down        Stop the app"
    echo "  restart     Restart the app"
    echo "  logs        View logs (use -f to follow)"
    echo "  status      Check app status"
    echo "  update      Pull latest changes and rebuild"
    echo "  pull        Pull Docker images only"
    echo "  build       Build Docker images"
    echo "  ps          List containers"
    echo "  stop        Stop all containers"
    echo "  start       Start all containers"
    echo "  clean       Remove containers and volumes"
    echo ""
    echo "Options:"
    echo "  -h, --help  Show this help"
    echo "  -f          Follow logs (with logs command)"
    echo "  -d          Detached mode (default)"
    echo "  --port      Specify port to check (default: 3000)"
    echo ""
    echo "Examples:"
    echo "  $0 demo-ai-chat up"
    echo "  $0 demo-ai-chat logs -f"
    echo "  $0 demo-ai-chat update"
    echo "  $0 demo-ai-chat status --port 8080"
}

check_app_exists() {
    local app_name=$1
    if [[ ! -d "${APPS_DIR}/${app_name}" ]]; then
        error "App '${app_name}' not found in ${APPS_DIR}"
    fi
}

check_docker_compose() {
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
    fi

    if ! docker compose version &> /dev/null; then
        error "Docker Compose is not installed"
    fi
}

cmd_up() {
    local app_name=$1
    log "Starting ${app_name}..."
    cd "${APPS_DIR}/${app_name}"
    docker compose up -d
    log "✅ ${app_name} started successfully"
}

cmd_down() {
    local app_name=$1
    log "Stopping ${app_name}..."
    cd "${APPS_DIR}/${app_name}"
    docker compose down
    log "✅ ${app_name} stopped"
}

cmd_restart() {
    local app_name=$1
    log "Restarting ${app_name}..."
    cd "${APPS_DIR}/${app_name}"
    docker compose restart
    log "✅ ${app_name} restarted"
}

cmd_logs() {
    local app_name=$1
    local follow=""
    if [[ "${FOLLOW_LOGS:-}" == "true" ]]; then
        follow="-f"
    fi
    cd "${APPS_DIR}/${app_name}"
    docker compose logs $follow --tail=100
}

cmd_status() {
    local app_name=$1
    local port=${PORT:-3000}
    cd "${APPS_DIR}/${app_name}"

    info "Container Status:"
    docker compose ps

    echo ""
    info "Checking health..."
    if curl -sf "http://localhost:${port}/health" > /dev/null 2>&1; then
        log "✅ App is healthy on port ${port}"
    else
        warn "App may not be responding on port ${port}"
        info "Try: curl http://localhost:${port}/health"
    fi
}

cmd_update() {
    local app_name=$1
    log "Updating ${app_name}..."

    cd "${APPS_DIR}/${app_name}"

    if [[ -d ".git" ]]; then
        log "Pulling latest code..."
        git pull origin main 2>/dev/null || warn "Git pull failed or not a git repo"
    fi

    log "Pulling Docker images..."
    docker compose pull

    log "Rebuilding and restarting..."
    docker compose up -d --build

    log "✅ ${app_name} updated successfully"
}

cmd_pull() {
    local app_name=$1
    log "Pulling Docker images for ${app_name}..."
    cd "${APPS_DIR}/${app_name}"
    docker compose pull
    log "✅ Images pulled"
}

cmd_build() {
    local app_name=$1
    log "Building ${app_name}..."
    cd "${APPS_DIR}/${app_name}"
    docker compose build
    log "✅ Build complete"
}

cmd_ps() {
    local app_name=$1
    cd "${APPS_DIR}/${app_name}"
    docker compose ps
}

cmd_stop() {
    local app_name=$1
    log "Stopping ${app_name} containers..."
    cd "${APPS_DIR}/${app_name}"
    docker compose stop
    log "✅ Containers stopped"
}

cmd_start() {
    local app_name=$1
    log "Starting ${app_name} containers..."
    cd "${APPS_DIR}/${app_name}"
    docker compose start
    log "✅ Containers started"
}

cmd_clean() {
    local app_name=$1
    warn "This will remove all containers, networks, and volumes for ${app_name}"
    read -p "Are you sure? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd "${APPS_DIR}/${app_name}"
        docker compose down -v
        log "✅ ${app_name} cleaned"
    else
        info "Cancelled"
    fi
}

main() {
    if [[ $# -lt 1 ]]; then
        usage
        exit 1
    fi

    local app_name="$1"
    local command="${2:-up}"
    local follow_logs="false"

    while [[ $# -gt 0 ]]; do
        case $1 in
            -f)
                follow_logs="true"
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            --port)
                PORT="$2"
                shift 2
                ;;
            *)
                shift
                ;;
        esac
    done

    export FOLLOW_LOGS="$follow_logs"

    check_docker_compose
    check_app_exists "$app_name"

    case $command in
        up)
            cmd_up "$app_name"
            ;;
        down)
            cmd_down "$app_name"
            ;;
        restart)
            cmd_restart "$app_name"
            ;;
        logs)
            cmd_logs "$app_name"
            ;;
        status)
            cmd_status "$app_name"
            ;;
        update)
            cmd_update "$app_name"
            ;;
        pull)
            cmd_pull "$app_name"
            ;;
        build)
            cmd_build "$app_name"
            ;;
        ps)
            cmd_ps "$app_name"
            ;;
        stop)
            cmd_stop "$app_name"
            ;;
        start)
            cmd_start "$app_name"
            ;;
        clean)
            cmd_clean "$app_name"
            ;;
        *)
            error "Unknown command: ${command}. Use -h for help."
            ;;
    esac
}

main "$@"
