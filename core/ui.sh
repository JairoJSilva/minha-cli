#!/bin/bash
# ==============================================================================
# core/ui.sh - Biblioteca de Interface e Estilização (TUI com gum e fallback)
# ==============================================================================

# Cores e Estilos ANSI (Fallback caso gum não esteja disponível)
CLR_RESET="\033[0m"
CLR_BOLD="\033[1m"
CLR_DIM="\033[2m"
CLR_CYAN="\033[36m"
CLR_GREEN="\033[32m"
CLR_YELLOW="\033[33m"
CLR_RED="\033[31m"
CLR_MAGENTA="\033[35m"
CLR_BLUE="\033[34m"

HAS_GUM=0
if command -v gum >/dev/null 2>&1; then
    HAS_GUM=1
fi

ui_banner() {
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style \
            --foreground 212 --border-foreground 99 --border double \
            --align center --width 64 --margin "0 1" --padding "0 2" \
            "☁️  MINHA CLI - MULTI-CLOUD & SRE CONTEXT" \
            "AWS • Oracle OCI • Google Cloud • Azure • Kubernetes"
    else
        echo -e "${CLR_MAGENTA}${CLR_BOLD}================================================================${CLR_RESET}"
        echo -e "${CLR_CYAN}${CLR_BOLD}       ☁️  MINHA CLI - MULTI-CLOUD & SRE CONTEXT MANAGER       ${CLR_RESET}"
        echo -e "${CLR_MAGENTA}================================================================${CLR_RESET}"
    fi
}

ui_section() {
    local title="$1"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style --foreground 214 --bold --margin "0 0 0 1" "▶ $title"
    else
        echo -e "\n${CLR_YELLOW}${CLR_BOLD}▶ $title${CLR_RESET}"
    fi
}

ui_success() {
    local msg="$1"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style --foreground 82 --bold "✅ $msg"
    else
        echo -e "${CLR_GREEN}${CLR_BOLD}✅ $msg${CLR_RESET}"
    fi
}

ui_info() {
    local msg="$1"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style --foreground 117 "ℹ️  $msg"
    else
        echo -e "${CLR_CYAN}ℹ️  $msg${CLR_RESET}"
    fi
}

ui_warn() {
    local msg="$1"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style --foreground 214 --bold "⚠️  $msg"
    else
        echo -e "${CLR_YELLOW}${CLR_BOLD}⚠️  $msg${CLR_RESET}"
    fi
}

ui_error() {
    local msg="$1"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum style --foreground 196 --bold "❌ $msg"
    else
        echo -e "${CLR_RED}${CLR_BOLD}❌ $msg${CLR_RESET}"
    fi
}

ui_card() {
    local title="$1"
    local content="$2"
    if [ "$HAS_GUM" -eq 1 ]; then
        local formatted_title
        formatted_title=$(gum style --foreground 39 --bold "$title")
        gum style \
            --border normal --border-foreground 39 \
            --padding "0 1" --margin "0 1" \
            "$formatted_title" "$content"
    else
        echo -e "${CLR_CYAN}┌── [ $title ] ────────────────────────────────────────${CLR_RESET}"
        echo -e "$content"
        echo -e "${CLR_CYAN}└─────────────────────────────────────────────────────${CLR_RESET}"
    fi
}

ui_spin() {
    local title="$1"
    local cmd="$2"
    if [ "$HAS_GUM" -eq 1 ]; then
        gum spin --spinner dot --title "$title" -- bash -c "$cmd"
    else
        echo -e "${CLR_DIM}⏳ $title...${CLR_RESET}"
        eval "$cmd"
    fi
}
