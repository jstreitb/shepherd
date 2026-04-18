#!/usr/bin/env bash
# Shepherd Installer — https://github.com/jstreitb/shepherd
# Usage: curl -sSfL https://raw.githubusercontent.com/jstreitb/shepherd/main/install.sh | bash
set -euo pipefail

# ─── Styling ─────────────────────────────────────────────────────────────────

BOLD="$(tput bold 2>/dev/null || printf '')"
DIM="$(tput dim 2>/dev/null || printf '\033[2m')"
RESET="$(tput sgr0 2>/dev/null || printf '\033[0m')"
BLUE="$(tput setaf 4 2>/dev/null || printf '\033[34m')"
GREEN="$(tput setaf 2 2>/dev/null || printf '\033[32m')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '\033[33m')"
RED="$(tput setaf 1 2>/dev/null || printf '\033[31m')"

REPO="jstreitb/shepherd"
INSTALL_DIR="/usr/local/bin"

# ─── Helpers ─────────────────────────────────────────────────────────────────

step()  { printf "\n  ${BLUE}${BOLD}==>${RESET} ${BOLD}%s${RESET}\n" "$1"; }
ok()    { printf "      ${GREEN}✔${RESET} %s\n" "$1"; }
warn()  { printf "      ${YELLOW}⚠${RESET} %s\n" "$1"; }
abort() { printf "\n  ${RED}${BOLD}✗ Error:${RESET} %s\n\n" "$1"; exit 1; }

# ─── Main ────────────────────────────────────────────────────────────────────

main() {
    printf "\n  ${BOLD}🐑 shepherd${RESET} ${DIM}installer${RESET}\n"

    # 1. Detect System
    step "Detecting system environment"
    local osarch
    
    if [ "$(uname -s)" != "Linux" ]; then
        abort "Shepherd currently only supports Linux."
    fi
    
    case "$(uname -m)" in
        x86_64|amd64)   osarch="linux_amd64" ;;
        aarch64|arm64)  osarch="linux_arm64" ;;
        *)              abort "Unsupported architecture: $(uname -m)" ;;
    esac
    ok "Detected Linux (${osarch})"

    # 2. Find Latest Version
    step "Resolving latest version"
    local version=""
    
    # Check if a specific version parameter was passed, otherwise default to "latest" API
    set +o pipefail
    version=$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep -Po '"tag_name": "\K.*?(?=")' | head -n 1)
    set -o pipefail

    if [ -z "$version" ]; then
        # Fallback to predefined if API fails
        version="v1.0.0"
        warn "Could not fetch latest release tag, defaulting to ${version}"
    fi
    ok "Target version is ${version}"
    
    # Check if already installed
    if command -v shepherd >/dev/null; then
        local current
        current="$(shepherd --version 2>/dev/null || true)"
        if [ "$current" = "$version" ] || [ "$current" = "${version#v}" ]; then
            printf "\n  ${GREEN}✨ Shepherd is already up to date (${current}).${RESET}\n\n"
            exit 0
        fi
    fi

    # 3. Download
    step "Downloading artifacts"
    local tmp
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT

    local url="https://github.com/${REPO}/releases/download/${version}/shepherd_${osarch}.tar.gz"
    
    # Try downloading
    if curl -sSfL "$url" -o "$tmp/shepherd.tar.gz" 2>/dev/null; then
        ok "Download complete"
    else
        abort "Download failed (Release might not be published yet)"
    fi

    # 4. Extract
    if ! tar -xzf "$tmp/shepherd.tar.gz" -C "$tmp" 2>/dev/null; then
        # Fallback if it was a raw binary instead of tarball
        mv "$tmp/shepherd.tar.gz" "$tmp/shepherd"
    fi

    # 5. Install
    step "Installing binary"
    
    local dest="${INSTALL_DIR}/shepherd"
    if command -v shepherd >/dev/null; then
        dest="$(command -v shepherd)"
    fi
    
    if [ -w "$(dirname "$dest")" ] && [ ! -d "$dest" ]; then
        mv "$tmp/shepherd" "$dest"
        chmod +x "$dest"
    else
        warn "Sudo permissions required to install to ${dest}"
        sudo mv "$tmp/shepherd" "$dest"
        sudo chmod +x "$dest"
    fi
    ok "Installed to ${dest}"

    printf "\n  ${GREEN}${BOLD}✨ Installation successful!${RESET}\n"
    printf "  ${DIM}Run ${BOLD}shepherd${DIM} in your terminal to start updating your system.${RESET}\n\n"
}

main "$@"
