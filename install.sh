#!/bin/bash
# ==============================================================================
# install.sh - Instalador Universal do Minha-CLI (Multi-Cloud Context Manager)
# ==============================================================================
# Execução simples:
#   chmod +x install.sh && ./install.sh
# ==============================================================================

set -euo pipefail

# Cores ANSI para o instalador
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
TARGET_BIN="$SCRIPT_DIR/bin/mc.sh"
USER_CONFIG_DIR="$HOME/.config/minha-cli"
USER_LOCAL_BIN="$HOME/.local/bin"

print_banner() {
    echo -e "${C_MAGENTA}${C_BOLD}"
    echo "  ╔════════════════════════════════════════════════════════════════════╗"
    echo "  ║                                                                    ║"
    echo "  ║          ☁️   INSTALADOR UNIVERSAL - MINHA-CLI                    ║"
    echo "  ║        Multi-Cloud (AWS • OCI • GCP • Azure • Kubernetes)          ║"
    echo "  ║                                                                    ║"
    echo "  ╚════════════════════════════════════════════════════════════════════╝"
    echo -e "${C_RESET}"
}

print_banner

echo -e "${C_CYAN}🔍 [1/4] Verificando pré-requisitos do sistema...${C_RESET}"

# 1. Checagem de jq (necessário para JSON)
if ! command -v jq >/dev/null 2>&1; then
    echo -e "${C_YELLOW}⚠️  'jq' não foi encontrado. Ele é necessário para manipular os perfis.${C_RESET}"
    if [ "$(id -u)" -eq 0 ] || command -v sudo >/dev/null 2>&1; then
        echo -e "${C_BLUE}Tentando instalar jq automaticamente...${C_RESET}"
        if command -v apt-get >/dev/null 2>&1; then
            sudo apt-get update -y && sudo apt-get install -y jq
        elif command -v dnf >/dev/null 2>&1; then
            sudo dnf install -y jq
        elif command -v yum >/dev/null 2>&1; then
            sudo yum install -y jq
        elif command -v brew >/dev/null 2>&1; then
            brew install jq
        fi
    else
        echo -e "${C_RED}Por favor, instale o 'jq' manualmente (ex: sudo apt install jq).${C_RESET}"
    fi
else
    echo -e "${C_GREEN}  ✅ jq: detectado ($(jq --version))${C_RESET}"
fi

# 2. Checagem de gum (TUI rica)
if ! command -v gum >/dev/null 2>&1; then
    echo -e "${C_YELLOW}  ⚠️  'gum' não detectado. O Minha-CLI funcionará com menu nativo ANSI.${C_RESET}"
    echo -e "${C_BLUE}     (Opcional) Para ter a interface com estilo Charm, instale: https://github.com/charmbracelet/gum${C_RESET}"
else
    echo -e "${C_GREEN}  ✅ gum: detectado (interface visual rica ativada)${C_RESET}"
fi

# 3. Checagem de rtk (proxy otimizado)
if command -v rtk >/dev/null 2>&1 || [ -x "/root/.local/bin/rtk" ] || [ -x "$HOME/.local/bin/rtk" ]; then
    echo -e "${C_GREEN}  ✅ rtk: detectado (consultas ultra-rápidas ativadas)${C_RESET}"
fi

echo ""
echo -e "${C_CYAN}⚙️  [2/4] Configurando permissões e arquivos de dados...${C_RESET}"

# Garante permissões de execução
chmod +x "$TARGET_BIN" "$SCRIPT_DIR/core/"*.sh "$SCRIPT_DIR/providers/"*.sh 2>/dev/null || true
echo -e "${C_GREEN}  ✅ Permissões de execução aplicadas a todos os scripts.${C_RESET}"

# Inicializa o diretório de dados do usuário ~/.config/minha-cli
mkdir -p "$USER_CONFIG_DIR"
if [ ! -f "$USER_CONFIG_DIR/clients.json" ]; then
    if [ -f "$SCRIPT_DIR/config/clients.json" ]; then
        cp "$SCRIPT_DIR/config/clients.json" "$USER_CONFIG_DIR/clients.json"
        echo -e "${C_GREEN}  ✅ Base inicial de clientes sincronizada em $USER_CONFIG_DIR/clients.json${C_RESET}"
    else
        echo "[]" > "$USER_CONFIG_DIR/clients.json"
        echo -e "${C_GREEN}  ✅ Base de dados criada em $USER_CONFIG_DIR/clients.json${C_RESET}"
    fi
fi

# Garante que o repositório aponte ou sincronize com o config do usuário
if [ ! -f "$SCRIPT_DIR/config/clients.json" ]; then
    mkdir -p "$SCRIPT_DIR/config"
    cp "$USER_CONFIG_DIR/clients.json" "$SCRIPT_DIR/config/clients.json" 2>/dev/null || true
fi

echo ""
echo -e "${C_CYAN}🔗 [3/4] Criando binários e wrappers no PATH do usuário...${C_RESET}"

# Cria ~/.local/bin se não existir e gera wrapper executável 'mc'
mkdir -p "$USER_LOCAL_BIN"
cat << WRAPPER > "$USER_LOCAL_BIN/mc"
#!/bin/bash
# Wrapper executável do Minha-CLI
exec bash "$TARGET_BIN" "\$@"
WRAPPER
chmod +x "$USER_LOCAL_BIN/mc"
echo -e "${C_GREEN}  ✅ Executável global criado em $USER_LOCAL_BIN/mc${C_RESET}"

echo ""
echo -e "${C_CYAN}🐚 [4/4] Injetando aliases de persistência de ambiente no shell...${C_RESET}"

ALIAS_MC_LINE="alias mc='source \"$TARGET_BIN\"'"
ALIAS_NUVEM_LINE="alias nuvem='source \"$TARGET_BIN\"'"

inject_shell_rc() {
    local rc_file="$1"
    if [ -f "$rc_file" ]; then
        # Remove versão antiga se houver
        if grep -Fq "alias mc=" "$rc_file"; then
            sed -i '/Minha-CLI Multi-Cloud/d' "$rc_file" 2>/dev/null || true
            sed -i '/alias mc=/d' "$rc_file" 2>/dev/null || true
            sed -i '/alias nuvem=/d' "$rc_file" 2>/dev/null || true
        fi

        echo "" >> "$rc_file"
        echo "# --- Minha-CLI Multi-Cloud Context Manager ---" >> "$rc_file"
        echo "$ALIAS_MC_LINE" >> "$rc_file"
        echo "$ALIAS_NUVEM_LINE" >> "$rc_file"
        echo -e "${C_GREEN}  ✅ Aliases 'mc' e 'nuvem' configurados em $rc_file${C_RESET}"
    fi
}

inject_shell_rc "$HOME/.bashrc"
inject_shell_rc "$HOME/.zshrc"
inject_shell_rc "$HOME/.bash_profile"

echo ""
echo -e "${C_GREEN}${C_BOLD}════════════════════════════════════════════════════════════════════${C_RESET}"
echo -e "${C_GREEN}${C_BOLD}  🎉 INSTALAÇÃO FINALIZADA COM SUCESSO!                             ${C_RESET}"
echo -e "${C_GREEN}${C_BOLD}════════════════════════════════════════════════════════════════════${C_RESET}"
echo ""
echo -e "${C_BOLD}Para começar a usar imediatamente no terminal atual, execute:${C_RESET}"
echo -e "  ${C_CYAN}${C_BOLD}source ~/.bashrc${C_RESET}  ${C_DIM}(ou source ~/.zshrc)${C_RESET}"
echo ""
echo -e "${C_BOLD}Comandos disponíveis:${C_RESET}"
echo -e "  ${C_GREEN}mc${C_RESET}               Abre o Menu Interativo (TUI)"
echo -e "  ${C_GREEN}mc add${C_RESET}           Cadastra uma nova conta / cliente"
echo -e "  ${C_GREEN}mc edit${C_RESET}          Edita configurações de um cliente existente"
echo -e "  ${C_GREEN}mc delete${C_RESET}        Apaga uma conta / cliente com confirmação"
echo -e "  ${C_GREEN}mc list${C_RESET}          Lista todas as contas cadastradas"
echo -e "  ${C_GREEN}mc switch${C_RESET}        Alterna rapidamente de contexto"
echo -e "  ${C_GREEN}mc status${C_RESET}        Mostra o contexto ativo e cluster K8s"
echo -e "  ${C_GREEN}mc test${C_RESET}          Testa credenciais em tempo real nas APIs (WhoAmI)"
echo -e "  ${C_GREEN}mc clear${C_RESET}         Limpa todas as variáveis de ambiente das nuvens"
echo ""
