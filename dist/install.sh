#!/bin/sh
set -e

REPO="sverdejot/geemail"
BINARY_NAME="geemail"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS
detect_os() {
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$os" in
        linux*) echo "linux" ;;
        darwin*) echo "darwin" ;;
        mingw*|msys*|cygwin*) echo "windows" ;;
        *) echo "unsupported: $os" >&2; exit 1 ;;
    esac
}

# Detect architecture
detect_arch() {
    arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *) echo "unsupported: $arch" >&2; exit 1 ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
    os="$(detect_os)"
    arch="$(detect_arch)"
    version="${VERSION:-$(get_latest_version)}"

    if [ -z "$version" ]; then
        echo "Error: Could not determine latest version" >&2
        exit 1
    fi

    echo "Installing ${BINARY_NAME} ${version} for ${os}/${arch}..."

    # Determine archive extension
    ext="tar.gz"
    if [ "$os" = "windows" ]; then
        ext="zip"
    fi

    # Build download URL
    archive_name="${BINARY_NAME}_${version#v}_${os}_${arch}.${ext}"
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    # Create temp directory
    tmp_dir="$(mktemp -d)"
    trap 'rm -rf "$tmp_dir"' EXIT

    echo "Downloading ${download_url}..."
    curl -sL "$download_url" -o "${tmp_dir}/${archive_name}"

    # Extract archive
    echo "Extracting..."
    cd "$tmp_dir"
    if [ "$ext" = "zip" ]; then
        unzip -q "$archive_name"
    else
        tar -xzf "$archive_name"
    fi

    # Find and install binary
    binary_file="${BINARY_NAME}_${os}_${arch}"
    if [ "$os" = "windows" ]; then
        binary_file="${binary_file}.exe"
    fi

    if [ ! -f "$binary_file" ]; then
        echo "Error: Binary not found in archive" >&2
        exit 1
    fi

    echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."

    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_file" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo mv "$binary_file" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    echo "Successfully installed ${BINARY_NAME} ${version} to ${INSTALL_DIR}/${BINARY_NAME}"
}

main
