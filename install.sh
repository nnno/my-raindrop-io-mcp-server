#!/bin/sh
set -eu

REPO="nnno/my-raindrop-io-mcp-server"
BINARY="my-raindrop-io-mcp-server"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux) ;;
  darwin) ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# darwin/amd64 is not provided
if [ "$OS" = "darwin" ] && [ "$ARCH" = "amd64" ]; then
  echo "Intel Mac is not supported. Use Apple Silicon Mac or Linux." >&2
  exit 1
fi

# Get latest version
if command -v curl >/dev/null 2>&1; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
elif command -v wget >/dev/null 2>&1; then
  VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
else
  echo "curl or wget is required" >&2
  exit 1
fi

if [ -z "$VERSION" ]; then
  echo "Failed to get latest version" >&2
  exit 1
fi

URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"

echo "Downloading ${BINARY} v${VERSION} (${OS}/${ARCH})..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" | tar xz -C "$TMPDIR" "$BINARY"
else
  wget -qO- "$URL" | tar xz -C "$TMPDIR" "$BINARY"
fi

mkdir -p "$INSTALL_DIR"
mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

echo "Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
