#!/bin/bash
# ==============================================================================
# install.sh - Instalador Universal do Minha-CLI (Go Edition)
# ==============================================================================

set -euo pipefail

C_RESET="\033[0m"
C_BOLD="\033[1m"
C_DIM="\033[2m"
C_CYAN="\033[36m"
C_GREEN="\033[32m"
C_YELLOW="\033[33m"
C_RED="\033[31m"
C_MAGENTA="\033[35m"
C_BLUE="\033[34m"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"
TARGET_BIN="$BIN_DIR/mc"
USER_CONFIG_DIR="$HOME/.config/minha-cli"
USER_LOCAL_BIN="$HOME/.local/bin"

print_banner() {
    echo -e "${C_MAGENTA}${C_BOLD}"
    echo "  ╔════════════════════════════════════════════════════════════════════╗"
    echo "  ║                                                                    ║"
    echo "  ║        ☁️   INSTALADOR MINHA-CLI (GO EDITION - SRE)               ║"
    echo "  ║        Multi-Cloud (AWS • OCI • GCP • Azure • Kubernetes)          ║"
    echo "  ║                                                                    ║"
    echo "  ╚════════════════════════════════════════════════════════════════════╝"
    echo -e "${C_RESET}"
}

print_banner

echo -e "${C_CYAN}🔨 [1/4] Compilando o binário Go de alta performance...${C_RESET}"
mkdir -p "$BIN_DIR"

if command -v go >/dev/null 2>&1; then
    (cd "$SCRIPT_DIR" && go build -ldflags="-s -w" -o "$TARGET_BIN" main.go)
    echo -e "${C_GREEN}  ✅ Binário compilado com sucesso: $TARGET_BIN${C_RESET}"
elif [ -f "$TARGET_BIN" ]; then
    echo -e "${C_GREEN}  ✅ Binário pré-compilado detectado: $TARGET_BIN${C_RESET}"
else
    echo -e "${C_RED}❌ Go não foi detectado e não há binário compilado. Instale o Go: https://golang.org${C_RESET}"
    exit 1
fi

chmod +x "$TARGET_BIN"

echo ""
echo -e "${C_CYAN}⚙️  [2/4] Configurando dados locais e permissões...${C_RESET}"
mkdir -p "$USER_CONFIG_DIR"
if [ ! -f "$USER_CONFIG_DIR/clients.json" ]; then
    if [ -f "$SCRIPT_DIR/config/clients.json" ]; then
        cp "$SCRIPT_DIR/config/clients.json" "$USER_CONFIG_DIR/clients.json"
    else
        echo "[]" > "$USER_CONFIG_DIR/clients.json"
    fi
fi
echo -e "${C_GREEN}  ✅ Base de dados ativa em $USER_CONFIG_DIR/clients.json${C_RESET}"

echo ""
echo -e "${C_CYAN}🔗 [3/4] Instalando binário no PATH global (~/.local/bin)...${C_RESET}"
mkdir -p "$USER_LOCAL_BIN"
cp "$TARGET_BIN" "$USER_LOCAL_BIN/mc"
chmod +x "$USER_LOCAL_BIN/mc"
echo -e "${C_GREEN}  ✅ Binário standalone instalado em $USER_LOCAL_BIN/mc${C_RESET}"

echo ""
echo -e "${C_CYAN}🐚 [4/4] Injetando wrapper de contexto de ambiente no Shell...${C_RESET}"

SHELL_FUNCTION="
# --- Minha-CLI Multi-Cloud Context Manager (Go Edition) ---
mc() {
    local target_bin=\"$TARGET_BIN\"
    if [ ! -f \"\$target_bin\" ]; then
        target_bin=\"$USER_LOCAL_BIN/mc\"
    fi

    case \"\${1:-}\" in
        \"switch\"|\"s\"|\"clear\"|\"reset\"|\"c\"|\"\")
            local env_script
            env_script=\$(\"\$target_bin\" env \"\${2:-}\" 2>/dev/null || true)
            if [ -n \"\$env_script\" ]; then
                eval \"\$env_script\"
            fi
            \"\$target_bin\" \"\$@\"
            ;;
        *)
            \"\$target_bin\" \"\$@\"
            ;;
    esac
}
alias nuvem='mc'
"

inject_shell_rc() {
    local rc_file="$1"
    if [ -f "$rc_file" ]; then
        if grep -Fq "Minha-CLI Multi-Cloud" "$rc_file"; then
            sed -i '/Minha-CLI Multi-Cloud/,/alias nuvem=/d' "$rc_file" 2>/dev/null || true
        fi
        echo "$SHELL_FUNCTION" >> "$rc_file"
        echo -e "${C_GREEN}  ✅ Shell wrapper configurado em $rc_file${C_RESET}"
    fi
}

inject_shell_rc "$HOME/.bashrc"
inject_shell_rc "$HOME/.zshrc"
inject_shell_rc "$HOME/.bash_profile"

echo ""
echo -e "${C_GREEN}${C_BOLD}════════════════════════════════════════════════════════════════════${C_RESET}"
echo -e "${C_GREEN}${C_BOLD}  🎉 MINHA-CLI (GO EDITION) INSTALADA COM SUCESSO!                  ${C_RESET}"
echo -e "${C_GREEN}${C_BOLD}════════════════════════════════════════════════════════════════════${C_RESET}"
echo ""
echo -e "${C_BOLD}Para começar a usar imediatamente no terminal atual, execute:${C_RESET}"
echo -e "  ${C_CYAN}${C_BOLD}source ~/.bashrc${C_RESET}  ${C_DIM}(ou source ~/.zshrc)${C_RESET}"
echo ""
echo -e "${C_BOLD}Comandos disponíveis:${C_RESET}"
echo -e "  ${C_GREEN}mc${C_RESET}               Abre o Menu Interativo TUI (Bubbletea)"
echo -e "  ${C_GREEN}mc scan${C_RESET}          Escaneia e importa configurações locais sem perder nada"
echo -e "  ${C_GREEN}mc add${C_RESET}           Cadastra nova conta / cliente (Formulário interativo)"
echo -e "  ${C_GREEN}mc edit${C_RESET}          Edita configurações de um cliente existente"
echo -e "  ${C_GREEN}mc delete${C_RESET}        Apaga uma conta / cliente com confirmação segura"
echo -e "  ${C_GREEN}mc list${C_RESET}          Lista todas as contas cadastradas com Lipgloss"
echo -e "  ${C_GREEN}mc switch <nome>${C_RESET} Alterna instantaneamente o contexto multi-cloud e K8s"
echo -e "  ${C_GREEN}mc status${C_RESET}        Mostra o card com o contexto e cluster K8s ativos"
echo -e "  ${C_GREEN}mc test${C_RESET}          Testa credenciais em paralelo ultra-rápido (Goroutines)"
echo -e "  ${C_GREEN}mc clear${C_RESET}         Limpa todas as variáveis de ambiente das nuvens"
echo ""

# Reinicia o shell automaticamente para carregar as novas funções no terminal atual
echo -e "${C_YELLOW}🔄 Reiniciando o terminal para aplicar as alterações automaticamente...${C_RESET}"
exec "$SHELL"
