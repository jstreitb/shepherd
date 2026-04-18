package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jstreitb/shepherd/internal/ui/components"
)

func main() {
	anim := components.NewAnimation()
	
	var out strings.Builder
	out.WriteString(`#!/usr/bin/env bash
# Shepherd Installer — https://github.com/jstreitb/shepherd
# Usage: curl -sSfL https://raw.githubusercontent.com/jstreitb/shepherd/main/install.sh | bash
set -euo pipefail

REPO="jstreitb/shepherd"
INSTALL_DIR="/usr/local/bin"
STATUS_FILE="$(mktemp)"
trap 'rm -f "$STATUS_FILE"; printf "\033[?25h"' EXIT

# ─── Helpers ─────────────────────────────────────────────────────────────────

fail()  { 
    printf '\033[2J\033[H\033[?25h\n  \033[31m\033[1m✗ %s\033[0m\n\n' "$1"
    exit 1
}

update_status() {
    echo "$1" > "$STATUS_FILE"
}

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
    set +o pipefail
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
        | grep -Po '"tag_name": "\K.*?(?=")' \
        | head -n 1
    set -o pipefail
}

# ─── Animated Sheep Spinner ──────────────────────────────────────────────────

animate() {
    local pid=$1

    # Clear screen and hide cursor
    printf '\033[2J\033[H'
    printf '\033[?25l'

    local frames=(
`)
	
	for i := 0; i < 12; i++ {
		frameLines := strings.Split(anim.Frame(), "\n")
		out.WriteString("        # Frame " + fmt.Sprint(i) + "\n")
		
		out.WriteString("        \"")
		for _, l := range frameLines {
			escaped := strings.ReplaceAll(l, "\\", "\\\\")
			escaped = strings.ReplaceAll(escaped, "`", "\\`")
			escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
			escaped = strings.ReplaceAll(escaped, "$", "\\$")
			out.WriteString(fmt.Sprintf("%s\\n", escaped))
		}
		out.WriteString("\"\n")
		anim.NextFrame()
	}

	out.WriteString(`    )

    local frame=0
    while kill -0 "$pid" 2>/dev/null; do
        # Move cursor to top left
        printf '\033[H'
        
        # Vertical padding
        printf '\n\n\n\n\n\n'
        
        local current_frame="${frames[$((frame % 12))]}"
        
        # Read dynamic status
        local status_msg="Initializing..."
        if [ -s "$STATUS_FILE" ]; then
            status_msg="$(cat "$STATUS_FILE")"
        fi
        
        # Print frame indented to center it
        while IFS= read -r line; do
            [ -z "$line" ] && continue
            printf '           \033[35m%s\033[0m\033[K\n' "$line"
        done <<< "$current_frame"
        
        local dots=""
        for ((d=0; d <= frame % 3; d++)); do dots="${dots}."; done
        
        printf '\n\033[K\n'
        printf '           \033[36m\033[2m%-58s\033[0m\033[K\n' "${status_msg}${dots}"
        printf '\033[J' # clear to end of screen
        
        frame=$((frame + 1))
        sleep 0.15
    done
}

# ─── Main ────────────────────────────────────────────────────────────────────

main() {
    # Start the dummy background process that keeps the animation running
    # until we explicitly kill it or it finishes. We use a subshell.
    (
        update_status "Detecting system..."
        os="$(detect_os)"
        arch="$(detect_arch)"
        
        # Fallback to hardcoded version if no release
        version="$(get_latest_version)"
        if [ -z "$version" ]; then
            version="v1.0.0"
        fi

        update_status "Checking local installation..."
        local current_version="none"
        if command -v shepherd >/dev/null; then
            current_version="$(shepherd --version 2>/dev/null || echo "unknown")"
            if [ "$current_version" = "$version" ] || [ "$current_version" = "v$version" ]; then
                echo "ALREADY_UP_TO_DATE" > "$STATUS_FILE"
                exit 0
            fi
        fi

        update_status "Resolving download URL..."
        if [ "$version" = "latest" ]; then
            url="https://github.com/${REPO}/releases/latest/download/shepherd_${os}_${arch}.tar.gz"
        else
            url="https://github.com/${REPO}/releases/download/${version}/shepherd_${os}_${arch}.tar.gz"
        fi

        tmp="$(mktemp -d)"
        
        update_status "Downloading shepherd ${version}..."
        if ! curl -sSfL "$url" -o "$tmp/shepherd.tar.gz"; then
            echo "ERROR_DOWNLOAD:$url" > "$STATUS_FILE"
            exit 1
        fi

        update_status "Extracting binaries..."
        tar -xzf "$tmp/shepherd.tar.gz" -C "$tmp" 2>/dev/null || {
            mv "$tmp/shepherd.tar.gz" "$tmp/shepherd"
        }

        update_status "Installing to system..."
        local install_dest="${INSTALL_DIR}/shepherd"
        if command -v shepherd >/dev/null; then
            install_dest="$(command -v shepherd)"
        fi

        if [ -w "$(dirname "$install_dest")" ] && [ ! -d "$install_dest" ]; then
            mv "$tmp/shepherd" "$install_dest"
            chmod +x "$install_dest"
        else
            # Request sudo inside the animation
            update_status "Requesting sudo to install to ${install_dest}..."
            sudo mv "$tmp/shepherd" "$install_dest"
            sudo chmod +x "$install_dest"
        fi
        
        echo "SUCCESS:$version" > "$STATUS_FILE"
        rm -rf "$tmp"
    ) &
    local MAIN_PID=$!

    # Start animation which will loop as long as MAIN_PID is running
    animate "$MAIN_PID"
    wait "$MAIN_PID" || true

    # Check the result left in STATUS_FILE
    local result=""
    if [ -s "$STATUS_FILE" ]; then
        result="$(cat "$STATUS_FILE")"
    fi

    printf '\033[2J\033[H\033[?25h'

    if [[ "$result" == ERROR_DOWNLOAD* ]]; then
        local failed_url="${result#ERROR_DOWNLOAD:}"
        printf '\n  \033[31m\033[1m✗ Download failed.\033[0m\n'
        printf '  \033[2mURL: %s\033[0m\n\n' "$failed_url"
        printf '  \033[33m(This usually means the GitHub Release is not uploaded or public yet.)\033[0m\n\n'
        exit 1
    elif [ "$result" = "ALREADY_UP_TO_DATE" ]; then
        printf '\n  \033[32m\033[1m✓ Shepherd is already up to date!\033[0m\n\n'
    elif [[ "$result" == SUCCESS* ]]; then
        local ver="${result#SUCCESS:}"
        printf '\n  \033[32m\033[1m✓ Shepherd %s installed successfully!\033[0m\n' "$ver"
        printf '  \033[2mRun \033[1mshepherd\033[2m to get started.\033[0m\n\n'
    else
        printf '\n  \033[31m\033[1m✗ Installation failed mysteriously.\033[0m\n\n'
        exit 1
    fi
}

main "$@"
`)

	os.WriteFile("install.sh", []byte(out.String()), 0755)
}
