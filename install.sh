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
    # Turn off pipefail locally so grep/curl on 404 doesn't panic the script
    set +o pipefail
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
        | grep -Po '"tag_name": "\K.*?(?=")' \
        | head -n 1
    set -o pipefail
}

# ─── Animated Sheep Spinner ──────────────────────────────────────────────────

animate() {
    local pid=$1
    local status_msg="$2"

    # Hide cursor and clear a large block
    printf '\033[?25l'
    printf '\n\n\n\n\n\n\n\n\n\n\n\n\n'

    local frames=(
        # Frame 0
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .    +    .   .        \n                                                          \n                                                          \n                                                          \n   ,@@@.                                                  \n  ( o.o )                                                 \n   /|  |\\                    |-|                          \n    d  b                     | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 1
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +    .  *   .  .          \n                                                          \n                                                          \n                                                          \n         ,@@@.                                            \n        ( o.o )                                           \n         /|  |\\              |-|                          \n          d  b               | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 2
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .    +    .   .        \n                                                          \n                                                          \n                                                          \n               ,@@@.                                      \n              ( o.o )                                     \n               /|  |\\        |-|                          \n                d  b         | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 3
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +    .  *   .  .          \n                                                          \n                                                          \n                    ,@@@.                                 \n                   ( o^o )                                \n                     /  \\                                 \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 4
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .    +    .   .        \n                                                          \n                       ,@@@.                              \n                      ( o^o )                             \n                        /  \\                              \n                                                          \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 5
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +*   .  *   .  .          \n                          ,@@@.                           \n                         ( >w< )                          \n                           /  \\                           \n                                                          \n                                                          \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 6
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .*   +    .   .        \n                             ,@@@.                        \n                            ( >w< )                       \n                              /  \\                        \n                                                          \n                                                          \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 7
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +    .  *   .  .          \n                                                          \n                                ,@@@.                     \n                               ( o^o )                    \n                                 /  \\                     \n                                                          \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 8
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .    +    .   .        \n                                                          \n                                                          \n                                    ,@@@.                 \n                                   ( o^o )                \n                                     /  \\                 \n                             |-|                          \n                             | |          ~               \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 9
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +    .  *   .  .          \n                                                          \n                                                          \n                                                          \n                                         ,@@@.            \n                                        ( o.o )           \n                             |-|         /|  |\\           \n                             | |          d  b            \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 10
        "   .  *      .     +      C    .    +   .  *   .          \n     +   .    *    .      .     *  .    +    .   .        \n                                                          \n                                                          \n                                                          \n                                                ,@@@.     \n                                               ( o.o )    \n                             |-|                /|  |\\    \n                             | |                 d  b     \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
        # Frame 11
        "     .    *    .  +    C     .   .  * .    .    +         \n   +    .   *      .    .       +    .  *   .  .          \n                                                          \n                                                          \n                                                          \n                                                          \n                                                          \n                             |-|                          \n                             | |                          \n---v---*--v------v---*--v------------v---*--v-------v--*--\n"
    )

    local frame=0
    while kill -0 "$pid" 2>/dev/null; do
        # Move cursor up 13 lines
        printf '\033[13A'
        
        # Print frame
        printf '%b' "\033[35m${frames[$((frame % 12))]}\033[0m"
        
        # Calculate dots for status message
        local dots=""
        for ((d=0; d <= frame % 3; d++)); do dots="${dots}."; done
        
        # Empty line then status message centered (padding to 58)
        # Using simple printing (left aligned to match sheep)
        printf '\n\033[K\n'
        printf '\033[K\033[2m%-58s\033[0m\n' "  ${status_msg}${dots}"
        
        frame=$((frame + 1))
        sleep 0.15
    done
    printf '\033[?25h'
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
    if [ "$version" = "latest" ]; then
        url="https://github.com/${REPO}/releases/latest/download/shepherd_${os}_${arch}.tar.gz"
    else
        url="https://github.com/${REPO}/releases/download/${version}/shepherd_${os}_${arch}.tar.gz"
    fi

    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT

    # Download in background with animated sheep
    curl -sSfL "$url" -o "$tmp/shepherd.tar.gz" &
    local dl_pid=$!

    animate "$dl_pid" "Downloading from GitHub Releases"
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
