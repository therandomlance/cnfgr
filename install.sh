#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
"$DIR/build.sh"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"
cp "$DIR/cnfgr" "$INSTALL_DIR/cnfgr"
chmod +x "$INSTALL_DIR/cnfgr"
"$INSTALL_DIR/cnfgr" -init
echo "Installed to $INSTALL_DIR/cnfgr"
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *) echo "Note: $INSTALL_DIR is not on your PATH. Add it to your shell rc, e.g.:" \
       && echo "  export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac
