#!/bin/bash
# ==============================================================================
# bin/mc.sh - Entrypoint Principal do Minha-CLI (Multi-Cloud Context Manager)
# ==============================================================================

# Identifica o diretório base do projeto
if [ -n "${BASH_SOURCE[0]}" ]; then
    CLI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
else
    CLI_DIR="$(pwd)"
fi

# Carrega todos os módulos do core e providers
source "$CLI_DIR/core/ui.sh"
source "$CLI_DIR/core/state.sh"
source "$CLI_DIR/core/config.sh"
source "$CLI_DIR/providers/aws.sh"
source "$CLI_DIR/providers/oci.sh"
source "$CLI_DIR/providers/gcp.sh"
source "$CLI_DIR/providers/azure.sh"
source "$CLI_DIR/providers/k8s.sh"

# Verifica se está rodando via source ou subshell
IS_SOURCED=0
if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
    IS_SOURCED=1
fi

show_help() {
    ui_banner
    echo -e "\n\033[1mUso:\033[0m"
    echo -e "  \033[36mmc\033[0m                           Abre o menu interativo TUI"
    echo -e "  \033[36mmc switch <cliente>\033[0m          Troca para o perfil informado (ex: maida, dentalis, farmacia, flowti)"
    echo -e "  \033[36mmc add | novo\033[0m                Cadastra uma nova conta/cliente interativamente"
    echo -e "  \033[36mmc edit | editar\033[0m             Edita as configurações de uma conta existente"
    echo -e "  \033[36mmc delete | apagar\033[0m           Remove uma conta/cliente cadastrado"
    echo -e "  \033[36mmc list | listar\033[0m             Lista todos os clientes e perfis cadastrados"
    echo -e "  \033[36mmc status\033[0m                    Exibe o status do contexto ativo no terminal"
    echo -e "  \033[36mmc test | whoami\033[0m             Testa as credenciais ativas nas APIs das nuvens"
    echo -e "  \033[36mmc clear\033[0m                     Limpa as variáveis de ambiente de todas as nuvens"
    echo -e "  \033[36mmc k8s\033[0m                       Verifica status do cluster Kubernetes ativo"
    echo -e "  \033[36mmc help\033[0m                      Mostra esta mensagem de ajuda\n"
}

action_switch() {
    local target="$1"
    if [ -z "$target" ]; then
        if [ "$HAS_GUM" -eq 1 ]; then
            local options=()
            while IFS= read -r line; do
                [ -n "$line" ] && options+=("$line")
            done < <(list_client_names)
            options+=("Limpar Contexto (Default)" "Voltar")

            target=$(printf "%s\n" "${options[@]}" | gum choose --header "Selecione o Cliente / Contexto:")
        else
            echo "Selecione o perfil:"
            local client_list
            client_list=$(list_client_names)
            select target in $client_list "Limpar Contexto (Default)" "Voltar"; do
                break
            done
        fi
    fi

    if [ -z "$target" ] || [ "$target" = "Voltar" ]; then
        ui_info "Operação cancelada."
        return 0
    fi

    if [ "$target" = "Limpar Contexto (Default)" ] || [ "$target" = "clear" ]; then
        clear_all_contexts
        return 0
    fi

    apply_profile "$target"
}

action_test_all() {
    ui_banner
    ui_section "Executando Testes de Identidade e Conexão (WhoAmI)..."

    # AWS
    test_aws_connection || true
    echo ""

    # OCI
    test_oci_connection || true
    echo ""

    # K8s
    test_k8s_connection || true
}

# Menu Interativo Principal
interactive_menu() {
    ui_banner

    if [ "$IS_SOURCED" -eq 0 ]; then
        ui_warn "Dica SRE: Execute 'source bin/mc.sh' ou use o alias 'mc' para persistir variáveis no terminal."
    fi

    local choice
    if [ "$HAS_GUM" -eq 1 ]; then
        choice=$(gum choose \
            "☁️  1. Trocar Contexto de Nuvem (Switch Profile)" \
            "📊 2. Status do Contexto Ativo" \
            "🔍 3. Testar Conexão / WhoAmI (AWS & OCI & K8s)" \
            "➕ 4. Configurar Nova Conta / Cliente (Add)" \
            "✏️  5. Editar Conta / Cliente Existente (Edit)" \
            "🗑️  6. Apagar Conta / Cliente (Delete)" \
            "📁 7. Mapeamento de Perfis Cadastrados" \
            "☸️  8. Kubernetes (Status do Cluster)" \
            "🧹 9. Limpar Contexto (Reset de Variáveis)" \
            "🚪 10. Sair")
    else
        echo "1) Trocar Contexto  2) Status  3) Testar WhoAmI  4) Novo Cliente  5) Editar  6) Apagar  7) Listar  8) K8s  9) Limpar  10) Sair"
        read -rp "Opção: " choice
    fi

    case "$choice" in
        *"1. Trocar Contexto"*)
            action_switch ""
            ;;
        *"2. Status do Contexto"*)
            show_context_status
            ;;
        *"3. Testar Conexão"*)
            action_test_all
            ;;
        *"4. Configurar Nova Conta"*)
            add_client_interactive
            ;;
        *"5. Editar Conta"*)
            edit_client_interactive
            ;;
        *"6. Apagar Conta"*)
            delete_client_interactive
            ;;
        *"7. Mapeamento"*)
            ui_section "Perfis Configurados e Detectados"
            discover_local_profiles
            ;;
        *"8. Kubernetes"*)
            test_k8s_connection || true
            ;;
        *"9. Limpar Contexto"*)
            clear_all_contexts
            ;;
        *)
            ui_info "Até logo!"
            return 0
            ;;
    esac
}

# Processamento de Parâmetros de Linha de Comando
main() {
    local cmd="${1:-}"

    case "$cmd" in
        "switch"|"s")
            action_switch "${2:-}"
            ;;
        "status"|"st")
            show_context_status
            ;;
        "test"|"whoami"|"t")
            action_test_all
            ;;
        "add"|"novo"|"criar"|"new")
            add_client_interactive
            ;;
        "edit"|"editar"|"update")
            edit_client_interactive
            ;;
        "delete"|"apagar"|"remover"|"rm"|"del")
            delete_client_interactive
            ;;
        "list"|"listar"|"ls")
            discover_local_profiles
            ;;
        "clear"|"reset"|"c")
            clear_all_contexts
            ;;
        "k8s"|"kube")
            test_k8s_connection || true
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        "")
            interactive_menu
            ;;
        *)
            # Se passou um nome direto de perfil (ex: mc maida)
            action_switch "$cmd"
            ;;
    esac
}

main "$@"
