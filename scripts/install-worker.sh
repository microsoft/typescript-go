#!/usr/bin/env bash
#
# Install a prebuilt luchta-tsc-worker for this OS/arch from the fork's GitHub
# releases into ~/.local/bin (override with LUCHTA_WORKER_BINDIR).
#
# Usage:
#   scripts/install-worker.sh [VERSION]        # VERSION defaults to latest
#   curl -fsSL <raw-url>/install-worker.sh | bash
#   curl -fsSL <raw-url>/install-worker.sh | bash -s -- v2026.06.21
#
# Works without auth on the public fork (curl fallback); uses gh when present.
set -euo pipefail

REPO="${LUCHTA_WORKER_REPO:-dobesv/typescript-go}"
BIN="luchta-tsc-worker"
BINDIR="${LUCHTA_WORKER_BINDIR:-$HOME/.local/bin}"
VERSION="${1:-}"

case "$(uname -s)" in
  Linux)  OS=linux ;;
  Darwin) OS=darwin ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64|amd64)  ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac
ASSET="$BIN-$OS-$ARCH"

mkdir -p "$BINDIR"
DEST="$BINDIR/$BIN"
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

if command -v gh >/dev/null 2>&1; then
  echo ">> Downloading $ASSET (${VERSION:-latest}) via gh..."
  # No tag arg => latest release.
  gh release download ${VERSION:+"$VERSION"} --repo "$REPO" --pattern "$ASSET" --output "$TMP" --clobber
else
  if [ -n "$VERSION" ]; then
    URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
  else
    URL="https://github.com/$REPO/releases/latest/download/$ASSET"
  fi
  echo ">> Downloading $URL via curl..."
  curl -fsSL "$URL" -o "$TMP"
fi

chmod +x "$TMP"
mv -f "$TMP" "$DEST"
trap - EXIT
echo ">> Installed $BIN -> $DEST"

case ":$PATH:" in
  *":$BINDIR:"*) ;;
  *) echo ">> NOTE: $BINDIR is not on your PATH; luchta resolves the worker via PATH." >&2 ;;
esac
