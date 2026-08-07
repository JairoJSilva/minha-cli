#!/bin/bash
# ==============================================================================
# core/scanner.sh - Auto-Descoberta e Importação Segura de Configurações Existentes
# ==============================================================================
# Comando: mc scan | mc leitura
# Lê os arquivos existentes no computador do usuário (~/.aws, ~/.oci, ~/.kube, GCP, Azure)
# e importa tudo para o clients.json sem apagar nada que o usuário já tinha.
# ==============================================================================

scan_and_import_existing_configs() {
    ui_banner
    ui_section "🔍 Escaneando Configurações Existentes na Máquina..."

    local cfg
    cfg=$(get_config_file)

    if [ ! -f "$cfg" ]; then
        echo "[]" > "$cfg"
    fi

    local aws_profiles=()
    local oci_profiles=()
    local gcp_configs=()
    local k8s_contexts=()
    local azure_subs=()

    # 1. Escanear AWS (~/.aws/credentials e ~/.aws/config)
    if [ -f "$HOME/.aws/credentials" ]; then
        while IFS= read -r line; do
            if [[ "$line" =~ ^\[(.*)\]$ ]]; then
                local p="${BASH_REMATCH[1]}"
                [[ "$p" != "" ]] && aws_profiles+=("$p")
            fi
        done < "$HOME/.aws/credentials"
    fi
    if [ -f "$HOME/.aws/config" ]; then
        while IFS= read -r line; do
            if [[ "$line" =~ ^\[profile[[:space:]]+(.*)\]$ ]] || [[ "$line" =~ ^\[(.*)\]$ ]]; then
                local p="${BASH_REMATCH[1]}"
                p="${p#profile }"
                if [[ "$p" != "" ]] && [[ ! " ${aws_profiles[*]} " =~ " ${p} " ]]; then
                    aws_profiles+=("$p")
                fi
            fi
        done < "$HOME/.aws/config"
    fi

    # 2. Escanear Oracle Cloud OCI (~/.oci/config)
    if [ -f "$HOME/.oci/config" ]; then
        while IFS= read -r line; do
            if [[ "$line" =~ ^\[(.*)\]$ ]]; then
                local p="${BASH_REMATCH[1]}"
                [[ "$p" != "" ]] && oci_profiles+=("$p")
            fi
        done < "$HOME/.oci/config"
    fi

    # 3. Escanear Google Cloud SDK (~/.config/gcloud/configurations)
    if [ -d "$HOME/.config/gcloud/configurations" ]; then
        for f in "$HOME/.config/gcloud/configurations/config_"*; do
            if [ -f "$f" ]; then
                local bname
                bname=$(basename "$f")
                local cname="${bname#config_}"
                [[ "$cname" != "*" && "$cname" != "" ]] && gcp_configs+=("$cname")
            fi
        done
    fi

    # 4. Escanear Kubernetes (~/.kube/config ou kubectl)
    if [ -f "$HOME/.kube/config" ] && command -v jq >/dev/null 2>&1; then
        if command -v kubectl >/dev/null 2>&1; then
            while IFS= read -r ctx; do
                [ -n "$ctx" ] && k8s_contexts+=("$ctx")
            done < <(kubectl config get-contexts --no-headers -o name 2>/dev/null || true)
        fi
    fi

    # 5. Escanear Azure CLI (~/.azure/azureProfile.json)
    if [ -f "$HOME/.azure/azureProfile.json" ] && command -v jq >/dev/null 2>&1; then
        while IFS= read -r sub; do
            [ -n "$sub" ] && azure_subs+=("$sub")
        done < <(jq -r '.subscriptions[]?.name // empty' "$HOME/.azure/azureProfile.json" 2>/dev/null || true)
    fi

    # Resumo do que foi encontrado
    local summary_body=" AWS Profiles    : ${#aws_profiles[@]} encontrados (${aws_profiles[*]:-<nenhum>})
 OCI Profiles    : ${#oci_profiles[@]} encontrados (${oci_profiles[*]:-<nenhum>})
 GCP Configs     : ${#gcp_configs[@]} encontrados (${gcp_configs[*]:-<nenhum>})
 K8s Contextos   : ${#k8s_contexts[@]} encontrados (${k8s_contexts[*]:-<nenhum>})
 Azure Assinaturas: ${#azure_subs[@]} encontrados (${azure_subs[*]:-<nenhum>})"

    ui_card "📊 DIAGNÓSTICO DO TERMINAL / MÁQUINA" "$summary_body"

    # Processamento e fusão inteligente no clients.json (sem apagar nem duplicar)
    local imported_count=0
    local updated_count=0

    # Unifica todos os identificadores únicos encontrados
    local all_keys=()
    for p in "${aws_profiles[@]}" "${oci_profiles[@]}" "${gcp_configs[@]}"; do
        local slug
        slug=$(echo "$p" | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]_-')
        if [[ -n "$slug" ]] && [[ ! " ${all_keys[*]} " =~ " ${slug} " ]] && [[ "$slug" != "default" ]]; then
            all_keys+=("$slug")
        fi
    done

    for slug in "${all_keys[@]}"; do
        # Verifica se o slug já existe no clients.json
        local exists
        exists=$(jq --arg s "$slug" '[.[] | select(.id == $s or (.name | ascii_downcase | contains($s)))] | length' "$cfg" 2>/dev/null || echo "0")

        # Procura correspondências nas nuvens
        local match_aws=""
        for ap in "${aws_profiles[@]}"; do
            local ap_slug
            ap_slug=$(echo "$ap" | tr '[:upper:]' '[:lower:]')
            if [[ "$ap_slug" == "$slug" || "$ap_slug" == *"$slug"* ]]; then
                match_aws="$ap"
                break
            fi
        done

        local match_oci=""
        for op in "${oci_profiles[@]}"; do
            local op_slug
            op_slug=$(echo "$op" | tr '[:upper:]' '[:lower:]')
            if [[ "$op_slug" == "$slug" || "$op_slug" == *"$slug"* ]]; then
                match_oci="$op"
                break
            fi
        done

        local match_gcp=""
        for gp in "${gcp_configs[@]}"; do
            local gp_slug
            gp_slug=$(echo "$gp" | tr '[:upper:]' '[:lower:]')
            if [[ "$gp_slug" == "$slug" || "$gp_slug" == *"$slug"* ]]; then
                match_gcp="$gp"
                break
            fi
        done

        local match_k8s=""
        for kp in "${k8s_contexts[@]}"; do
            local kp_slug
            kp_slug=$(echo "$kp" | tr '[:upper:]' '[:lower:]')
            if [[ "$kp_slug" == *"$slug"* ]]; then
                match_k8s="$kp"
                break
            fi
        done

        if [ "$exists" -eq 0 ]; then
            # Novo cliente identificado na máquina -> Adiciona
            local display_name
            display_name="$(echo "${slug:0:1}" | tr '[:lower:]' '[:upper:]')${slug:1}"
            
            local new_item
            new_item=$(jq -n \
                --arg id "$slug" \
                --arg name "$display_name (Importado)" \
                --arg aws "${match_aws:-null}" \
                --arg oci "${match_oci:-null}" \
                --arg gcp "${match_gcp:-null}" \
                --arg k8s "${match_k8s:-null}" \
                '{
                    id: $id,
                    name: $name,
                    aws_profile: (if $aws == "null" or $aws == "" then null else $aws end),
                    oci_profile: (if $oci == "null" or $oci == "" then null else $oci end),
                    gcp_config: (if $gcp == "null" or $gcp == "" then null else $gcp end),
                    azure_sub: null,
                    k8s_context: (if $k8s == "null" or $k8s == "" then null else $k8s end)
                }')

            local updated
            updated=$(jq --argjson item "$new_item" '. += [$item]' "$cfg")
            echo "$updated" > "$cfg"
            imported_count=$((imported_count + 1))
            ui_info "📥 Importado novo cliente detectado: $display_name"
        else
            updated_count=$((updated_count + 1))
        fi
    done

    echo ""
    ui_success "Varredura concluída com segurança!"
    echo -e "  \033[32m✔ $imported_count novo(s) cliente(s) importado(s) automaticamente.\033[0m"
    echo -e "  \033[36m✔ $updated_count cliente(s) já existentes foram preservados intactos.\033[0m"
    echo -e "  \033[2mBase ativa: $cfg\033[0m"
}
