#!/bin/bash
###############################################################################
# View logs for an app
# Usage: ./logs.sh <app-name> [lines]
###############################################################################

APPS_DIR="/apps"

if [[ -z "$1" ]]; then
    echo "Usage: $0 <app-name> [lines]"
    echo ""
    echo "Examples:"
    echo "  $0 demo-ai-chat         # View last 100 lines"
    echo "  $0 demo-ai-chat -f      # Follow logs"
    echo "  $0 demo-ai-chat 50      # View last 50 lines"
    exit 1
fi

app_name="$1"
lines="${2:-100}"

if [[ ! -d "${APPS_DIR}/${app_name}" ]]; then
    echo "❌ App '${app_name}' not found"
    exit 1
fi

cd "${APPS_DIR}/${app_name}"

if [[ "$lines" == "-f" ]]; then
    echo "📋 Following logs for ${app_name}..."
    docker compose logs -f --tail=100
else
    echo "📋 Last ${lines} log lines for ${app_name}..."
    docker compose logs --tail="$lines"
fi
