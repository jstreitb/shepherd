#!/usr/bin/env bash
# Shepherd Installer — https://github.com/jstreitb/shepherd
# Usage: curl -sSfL https://raw.githubusercontent.com/jstreitb/shepherd/main/install.sh | bash
set -euo pipefail

# ─── Colors ──────────────────────────────────────────────────────────────────
BOLD='\033[1m'
DIM='\033[2m'
GREEN='\033[32m'
RED='\033[31m'
CYAN='\033[36m'
MAGENTA='\033[35m'
RESET='\033[0m'

REPO="jstreitb/shepherd"
INSTALL_DIR="/usr/local/bin"

# ─── Helpers ─────────────────────────────────────────────────────────────────

info()  { printf "  ${CYAN}%s${RESET}\n" "$1"; }
ok()    { printf "  ${GREEN}${BOLD}✓ %s${RESET}\n" "$1"; }
fail()  { printf "  ${RED}${BOLD}✗ %s${RESET}\n" "$1"; exit 1; }

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              fail "Unsupported architecture: $(uname -m)" ;;
    esac
}

detect_os() {
    case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
        linux)  echo "linux" ;;
        *)      fail "Shepherd only supports Linux." ;;
    esac
}

get_latest_version() {
    curl -sI "https://github.com/${REPO}/releases/latest" 2>/dev/null \
        | grep -i "^location:" \
        | sed 's#.*/tag/##' \
        | tr -d '\r\n'
}

# ─── Animated Sheep Spinner ──────────────────────────────────────────────────

animate() {
    local pid=$1

    # Reserve space for animation (7 lines)
    printf '\n\n\n\n\n\n\n'
    printf '\033[?25l' # hide cursor

    local frame=0
    while kill -0 "$pid" 2>/dev/null; do
        printf '\033[7A' # move up 7 lines

        case $((frame % 4)) in
            0)
                printf "\033[K     ${MAGENTA},@@@.${RESET}\n"
                printf "\033[K    ${MAGENTA}( o.o )${RESET}\n"
                printf "\033[K     ${MAGENTA}/|  |\\\\${RESET}\n"
                printf "\033[K      ${MAGENTA}d  b${RESET}\n"
                ;;
            1)
                printf "\033[K     ${MAGENTA},@@@.${RESET}\n"
                printf "\033[K    ${MAGENTA}( o.o )${RESET}\n"
                printf "\033[K     ${MAGENTA}/|  |\\\\${RESET}\n"
                printf "\033[K     ${MAGENTA}d    b${RESET}\n"
                ;;
            2)
                printf "\033[K     ${MAGENTA},@@@.${RESET}\n"
                printf "\033[K    ${MAGENTA}( o^o )${RESET}\n"
                printf "\033[K     ${MAGENTA}/|  |\\\\${RESET}\n"
                printf "\033[K      ${MAGENTA}d  b${RESET}\n"
                ;;
            3)
                printf "\033[K     ${MAGENTA},@@@.${RESET}\n"
                printf "\033[K    ${MAGENTA}( o.- )${RESET}\n"
                printf "\033[K     ${MAGENTA}/|  |\\\\${RESET}\n"
                printf "\033[K     ${MAGENTA}d    b${RESET}\n"
                ;;
        esac

        local dots=""
        for ((d=0; d <= frame % 3; d++)); do dots="${dots}."; done
        printf "\033[K\n"
        printf "\033[K  ${DIM}Installing shepherd${dots}${RESET}\n"
        printf "\033[K\n"

        frame=$((frame + 1))
        sleep 0.25
    done

    printf '\033[?25h' # show cursor
}

# ─── Main ────────────────────────────────────────────────────────────────────

main() {
    local os arch version url tmp

    os="$(detect_os)"
    arch="$(detect_arch)"
    version="$(get_latest_version)"

    if [ -z "$version" ]; then
        version="latest"
    fi

    # Banner
    echo ""
    printf "  ${BOLD}${MAGENTA}  ,@@@.${RESET}\n"
    printf "  ${BOLD}${MAGENTA} ( o.o )${RESET}\n"
    printf "  ${BOLD}${MAGENTA}  /|  |\\\\${RESET}\n"
    printf "  ${BOLD}${MAGENTA}   d  b${RESET}\n"
    echo ""
    printf "  ${BOLD}shepherd installer${RESET}\n"
    echo ""
    printf "  ${DIM}OS:       ${os}${RESET}\n"
    printf "  ${DIM}Arch:     ${arch}${RESET}\n"
    printf "  ${DIM}Version:  ${version}${RESET}\n"
    echo ""

    # Check current version
    local current_version="none"
    if command -v shepherd >/dev/null; then
        current_version="$(shepherd --version 2>/dev/null || echo "unknown")"
        
        # If it says 'dev', it's a locally built dev binary, let's treat it as outdated
        # to allow forcing a release installation, unless we want to keep dev. 
        # But actually, checking strict equality is fine.
        if [ "$current_version" = "$version" ] || [ "$current_version" = "v$version" ]; then
            echo ""
            ok "Shepherd is already up to date ($current_version)"
            exit 0
        fi
    fi

    # Build download URL
    url="https://github.com/${REPO}/releases/download/${version}/shepherd_${os}_${arch}.tar.gz"

    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT

    # Download in background with animated sheep
    info "Downloading from GitHub Releases..."
    curl -sSfL "$url" -o "$tmp/shepherd.tar.gz" &
    local dl_pid=$!

    animate "$dl_pid"
    wait "$dl_pid" || fail "Download failed. Check the URL: ${url}"

    # Extract
    info "Extracting..."
    tar -xzf "$tmp/shepherd.tar.gz" -C "$tmp" 2>/dev/null || {
        # Might be a raw binary (not tarball)
        mv "$tmp/shepherd.tar.gz" "$tmp/shepherd"
    }

    # Determine where to install it (overwrite existing if present)
    local install_dest="${INSTALL_DIR}/shepherd"
    if command -v shepherd >/dev/null; then
        install_dest="$(command -v shepherd)"
    fi

    # Install
    info "Installing to ${install_dest}..."
    if [ -w "$(dirname "$install_dest")" ] && [ ! -d "$install_dest" ]; then
        mv "$tmp/shepherd" "$install_dest"
        chmod +x "$install_dest"
    else
        sudo mv "$tmp/shepherd" "$install_dest"
        sudo chmod +x "$install_dest"
    fi

    echo ""
    ok "shepherd ${version} installed successfully!"
    echo ""
    printf "  ${DIM}Run ${BOLD}shepherd${RESET}${DIM} to get started.${RESET}\n"
    echo ""
}

main "$@"
