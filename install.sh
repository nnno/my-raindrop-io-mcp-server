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

FILENAME="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt"

echo "Downloading ${BINARY} v${VERSION} (${OS}/${ARCH})..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Download archive and checksums
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" -o "$TMPDIR/$FILENAME"
  curl -fsSL "$CHECKSUMS_URL" -o "$TMPDIR/checksums.txt"
else
  wget -qO "$TMPDIR/$FILENAME" "$URL"
  wget -qO "$TMPDIR/checksums.txt" "$CHECKSUMS_URL"
fi

# Verify checksum
EXPECTED=$(grep "  ${FILENAME}$" "$TMPDIR/checksums.txt" | awk '{print $1}')
if [ -z "$EXPECTED" ]; then
  echo "Checksum entry not found for ${FILENAME}" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "$TMPDIR/$FILENAME" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "$TMPDIR/$FILENAME" | awk '{print $1}')
else
  echo "sha256sum or shasum is required for checksum verification" >&2
  exit 1
fi

if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "Checksum verification failed" >&2
  echo "  expected: $EXPECTED" >&2
  echo "  actual:   $ACTUAL" >&2
  exit 1
fi

tar xzf "$TMPDIR/$FILENAME" -C "$TMPDIR" "$BINARY"

mkdir -p "$INSTALL_DIR"
mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

echo "Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
