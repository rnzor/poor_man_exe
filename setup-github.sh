#!/bin/bash
set -e

echo "🚀 Creating GitHub repository and pushing..."

# Check if gh is installed
if command -v gh &> /dev/null; then
    echo "📦 Using GitHub CLI..."
    gh repo create poor_man_exe --public --description "Self-hosted deployment platform inspired by exe.dev - Docker + Caddy + Auto-HTTPS" --source=. --push
else
    echo "📦 GitHub CLI not found. Opening browser..."
    echo "Please create the repository at:"
    echo "👉 https://github.com/new?name=poor_man_exe&description=Self-hosted+deployment+platform"
    echo ""
    echo "Then run:"
    echo "  git remote set-url origin https://github.com/rnzor/poor_man_exe.git"
    echo "  git push -u origin main"
    echo ""
    read -p "Press Enter after creating the repository..."
    git remote set-url origin https://github.com/rnzor/poor_man_exe.git
    git push -u origin main
fi

echo "✅ Done! Repository: https://github.com/rnzor/poor_man_exe"
