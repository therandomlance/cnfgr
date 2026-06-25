#!/usr/bin/env bash
set -euo pipefail
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
rm -f "$INSTALL_DIR/cnfgr"
echo "Removed $INSTALL_DIR/cnfgr"
read -rp "Also remove cnfgr config (~/.config/cnfgr)? [y/N] " ans
if [[ "$ans" =~ ^[Yy]$ ]]; then
    rm -rf "$HOME/.config/cnfgr"
    echo "Removed config."
fi
