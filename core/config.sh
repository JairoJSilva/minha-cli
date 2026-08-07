#!/bin/bash
# ==============================================================================
# core/config.sh - Gerenciamento de Perfis de Clientes (CRUD & Persistência)
# ==============================================================================

get_config_file() {
    local primary="$CLI_DIR/config/clients.json"
    local user_dir="$HOME/.config/minha-cli"
    local fallback="$user_dir/clients.json"

    if [ -f "$primary" ]; then
        echo "$primary"
    else
        mkdir -p "$user_dir"
        if [ ! -f "$fallback" ]; then
            echo "[]" > "$fallback"
        fi
        echo "$fallback"
    fi
}

list_client_names() {
    local cfg
    cfg=$(get_config_file)
    if [ -f "$cfg" ] && command -v jq >/dev/null 2>&1; then
        jq -r '.[].name' "$cfg" 2>/dev/null
    fi
}

apply_profile() {
    local target="$1"
    local cfg
    cfg=$(get_config_file)

    if [ ! -f "$cfg" ] || ! command -v jq >/dev/null 2>&1; then
        # Fallback legado se jq ou json não estiverem presentes
        case "$target" in
            *"Maida"*|"maida")
                export AWS_PROFILE="maida"
                export CLOUDSDK_ACTIVE_CONFIG_NAME="maida"
                unset OCI_CLI_PROFILE
                ui_success "Contexto ativado: MAIDA"
                ;;
            *"Dentalis"*|"dentalis")
                export AWS_PROFILE="dentalis"
                unset CLOUDSDK_ACTIVE_CONFIG_NAME
                unset OCI_CLI_PROFILE
                ui_success "Contexto ativado: DENTALIS"
                ;;
            *"Farmacia"*|"farmacia")
                export AWS_PROFILE="farmacia"
                export CLOUDSDK_ACTIVE_CONFIG_NAME="farmacia"
                unset OCI_CLI_PROFILE
                ui_success "Contexto ativado: FARMÁCIA DIGITAL"
                ;;
            *"Flowti"*|*"Pessoal"*|"flowti"|"pessoal")
                export AWS_PROFILE="flowti"
                export OCI_CLI_PROFILE="pessoal"
                unset CLOUDSDK_ACTIVE_CONFIG_NAME
                ui_success "Contexto ativado: PESSOAL/FLOWTI"
                ;;
            *)
                export AWS_PROFILE="$target"
                export OCI_CLI_PROFILE="$target"
                ui_warn "Contexto customizado aplicado: $target"
                ;;
        esac
        return 0
    fi

    # Busca no JSON por ID ou por substring do Nome
    local client_data
    client_data=$(jq -c --arg q "$target" '.[] | select(.id == $q or (.name | ascii_downcase | contains($q | ascii_downcase)))' "$cfg" | head -n 1)

    if [ -n "$client_data" ] && [ "$client_data" != "null" ]; then
        local c_name c_aws c_oci c_gcp c_az c_k8s
        c_name=$(echo "$client_data" | jq -r '.name')
        c_aws=$(echo "$client_data" | jq -r '.aws_profile // empty')
        c_oci=$(echo "$client_data" | jq -r '.oci_profile // empty')
        c_gcp=$(echo "$client_data" | jq -r '.gcp_config // empty')
        c_az=$(echo "$client_data" | jq -r '.azure_sub // empty')
        c_k8s=$(echo "$client_data" | jq -r '.k8s_context // empty')

        # AWS
        if [ -n "$c_aws" ]; then
            export AWS_PROFILE="$c_aws"
        else
            unset AWS_PROFILE
        fi

        # OCI
        if [ -n "$c_oci" ]; then
            export OCI_CLI_PROFILE="$c_oci"
        else
            unset OCI_CLI_PROFILE
        fi

        # GCP
        if [ -n "$c_gcp" ]; then
            export CLOUDSDK_ACTIVE_CONFIG_NAME="$c_gcp"
        else
            unset CLOUDSDK_ACTIVE_CONFIG_NAME
        fi

        # Azure
        if [ -n "$c_az" ]; then
            export AZURE_SUBSCRIPTION="$c_az"
            if command -v az >/dev/null 2>&1; then
                az account set --subscription "$c_az" >/dev/null 2>&1 || true
            fi
        else
            unset AZURE_SUBSCRIPTION
        fi

        # Kubernetes
        if [ -n "$c_k8s" ] && command -v kubectl >/dev/null 2>&1; then
            kubectl config use-context "$c_k8s" >/dev/null 2>&1 || true
        fi

        ui_success "Contexto ativado: $c_name"
        echo -e "  \033[2mAWS: ${c_aws:-<nenhum>} | OCI: ${c_oci:-<nenhum>} | GCP: ${c_gcp:-<nenhum>} | Azure: ${c_az:-<nenhum>} | K8s: ${c_k8s:-<nenhum>}\033[0m"
    else
        # Se não achou no JSON, aplica como profile AWS/OCI genérico
        export AWS_PROFILE="$target"
        export OCI_CLI_PROFILE="$target"
        ui_warn "Contexto customizado aplicado: $target"
    fi
}

# ==============================================================================
# CRUD de Clientes e Contas
# ==============================================================================

add_client_interactive() {
    ui_section "➕ Cadastrar Nova Conta / Cliente"
    local cfg
    cfg=$(get_config_file)

    local name id aws_p oci_p gcp_c az_s k8s_c

    if [ "$HAS_GUM" -eq 1 ]; then
        name=$(gum input --placeholder "Nome de exibição (ex: Hospital Albert Einstein)" --header "Nome do Cliente:")
        if [ -z "$name" ]; then
            ui_warn "Cadastro cancelado."
            return 0
        fi

        id=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]_-')
        id=$(gum input --value "$id" --placeholder "Identificador único (ex: einstein)" --header "ID / Slug curto:")

        aws_p=$(gum input --placeholder "Nome do AWS Profile (ex: einstein-prod) [Opcional]" --header "AWS Profile:")
        oci_p=$(gum input --placeholder "Nome do OCI Profile (ex: einstein-oci) [Opcional]" --header "Oracle OCI Profile:")
        gcp_c=$(gum input --placeholder "Nome da Config GCP (ex: einstein-gcp) [Opcional]" --header "Google Cloud Config:")
        az_s=$(gum input --placeholder "Subscription ID do Azure [Opcional]" --header "Azure Subscription:")
        k8s_c=$(gum input --placeholder "Nome do Contexto Kubernetes [Opcional]" --header "Contexto K8s:")
    else
        read -rp "Nome do Cliente: " name
        [ -z "$name" ] && return 0
        read -rp "ID curto: " id
        read -rp "AWS Profile (opcional): " aws_p
        read -rp "OCI Profile (opcional): " oci_p
        read -rp "GCP Config (opcional): " gcp_c
        read -rp "Azure Subscription (opcional): " az_s
        read -rp "Contexto K8s (opcional): " k8s_c
    fi

    local new_entry
    new_entry=$(jq -n \
        --arg id "$id" \
        --arg name "$name" \
        --arg aws "${aws_p:-null}" \
        --arg oci "${oci_p:-null}" \
        --arg gcp "${gcp_c:-null}" \
        --arg az "${az_s:-null}" \
        --arg k8s "${k8s_c:-null}" \
        '{
            id: $id,
            name: $name,
            aws_profile: (if $aws == "null" or $aws == "" then null else $aws end),
            oci_profile: (if $oci == "null" or $oci == "" then null else $oci end),
            gcp_config: (if $gcp == "null" or $gcp == "" then null else $gcp end),
            azure_sub: (if $az == "null" or $az == "" then null else $az end),
            k8s_context: (if $k8s == "null" or $k8s == "" then null else $k8s end)
        }')

    # Salva no JSON
    local updated_json
    updated_json=$(jq --argjson item "$new_entry" '. += [$item]' "$cfg")
    echo "$updated_json" > "$cfg"

    ui_success "Cliente '$name' ($id) adicionado com sucesso!"
}

edit_client_interactive() {
    ui_section "✏️  Editar Conta / Cliente Existente"
    local cfg
    cfg=$(get_config_file)

    local client_name
    if [ "$HAS_GUM" -eq 1 ]; then
        client_name=$(list_client_names | gum choose --header "Selecione o cliente para editar:" "Voltar")
    else
        read -rp "Digite o nome ou ID do cliente para editar: " client_name
    fi

    if [ -z "$client_name" ] || [ "$client_name" = "Voltar" ]; then
        ui_info "Edição cancelada."
        return 0
    fi

    local client_data
    client_data=$(jq -c --arg n "$client_name" '.[] | select(.name == $n or .id == $n)' "$cfg" | head -n 1)

    if [ -z "$client_data" ]; then
        ui_error "Cliente não encontrado."
        return 1
    fi

    local old_id old_name old_aws old_oci old_gcp old_az old_k8s
    old_id=$(echo "$client_data" | jq -r '.id')
    old_name=$(echo "$client_data" | jq -r '.name')
    old_aws=$(echo "$client_data" | jq -r '.aws_profile // ""')
    old_oci=$(echo "$client_data" | jq -r '.oci_profile // ""')
    old_gcp=$(echo "$client_data" | jq -r '.gcp_config // ""')
    old_az=$(echo "$client_data" | jq -r '.azure_sub // ""')
    old_k8s=$(echo "$client_data" | jq -r '.k8s_context // ""')

    local new_name new_aws new_oci new_gcp new_az new_k8s
    if [ "$HAS_GUM" -eq 1 ]; then
        new_name=$(gum input --value "$old_name" --header "Nome de exibição:")
        new_aws=$(gum input --value "$old_aws" --header "AWS Profile:")
        new_oci=$(gum input --value "$old_oci" --header "Oracle OCI Profile:")
        new_gcp=$(gum input --value "$old_gcp" --header "Google Cloud Config:")
        new_az=$(gum input --value "$old_az" --header "Azure Subscription:")
        new_k8s=$(gum input --value "$old_k8s" --header "Contexto K8s:")
    else
        read -rp "Nome [$old_name]: " new_name
        new_name="${new_name:-$old_name}"
        read -rp "AWS Profile [$old_aws]: " new_aws
        new_aws="${new_aws:-$old_aws}"
        read -rp "OCI Profile [$old_oci]: " new_oci
        new_oci="${new_oci:-$old_oci}"
        read -rp "GCP Config [$old_gcp]: " new_gcp
        new_gcp="${new_gcp:-$old_gcp}"
        read -rp "Azure Sub [$old_az]: " new_az
        new_az="${new_az:-$old_az}"
        read -rp "Contexto K8s [$old_k8s]: " new_k8s
        new_k8s="${new_k8s:-$old_k8s}"
    fi

    local updated_json
    updated_json=$(jq \
        --arg id "$old_id" \
        --arg name "$new_name" \
        --arg aws "$new_aws" \
        --arg oci "$new_oci" \
        --arg gcp "$new_gcp" \
        --arg az "$new_az" \
        --arg k8s "$new_k8s" \
        'map(if .id == $id then {
            id: .id,
            name: $name,
            aws_profile: (if $aws == "" then null else $aws end),
            oci_profile: (if $oci == "" then null else $oci end),
            gcp_config: (if $gcp == "" then null else $gcp end),
            azure_sub: (if $az == "" then null else $az end),
            k8s_context: (if $k8s == "" then null else $k8s end)
        } else . end)' "$cfg")

    echo "$updated_json" > "$cfg"
    ui_success "Cliente '$new_name' atualizado com sucesso!"
}

delete_client_interactive() {
    ui_section "🗑️  Apagar Conta / Cliente"
    local cfg
    cfg=$(get_config_file)

    local client_name
    if [ "$HAS_GUM" -eq 1 ]; then
        client_name=$(list_client_names | gum choose --header "Selecione o cliente para APAGAR:" "Voltar")
    else
        read -rp "Digite o nome ou ID do cliente para apagar: " client_name
    fi

    if [ -z "$client_name" ] || [ "$client_name" = "Voltar" ]; then
        ui_info "Remoção cancelada."
        return 0
    fi

    if [ "$HAS_GUM" -eq 1 ]; then
        if ! gum confirm "Tem certeza que deseja APAGAR permanentemente '$client_name'?"; then
            ui_info "Operação cancelada pelo usuário."
            return 0
        fi
    else
        read -rp "Tem certeza que deseja apagar '$client_name'? (s/N): " confirm
        [[ ! "$confirm" =~ ^[sS]$ ]] && return 0
    fi

    local updated_json
    updated_json=$(jq --arg n "$client_name" 'map(select(.name != $n and .id != $n))' "$cfg")
    echo "$updated_json" > "$cfg"

    ui_success "Cliente '$client_name' removido com sucesso!"
}

discover_local_profiles() {
    local cfg
    cfg=$(get_config_file)
    ui_card "📁 PERFIS CADASTRADOS" "$(jq -r '.[] | " • \(.name) [ID: \(.id)]\n   └─ AWS: \(.aws_profile // "-") | OCI: \(.oci_profile // "-") | GCP: \(.gcp_config // "-") | K8s: \(.k8s_context // "-")"' "$cfg" 2>/dev/null || echo "Nenhum cliente cadastrado.")"
}
