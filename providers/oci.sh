#!/bin/bash
# ==============================================================================
# providers/oci.sh - Módulo de Operações e Diagnóstico Oracle Cloud (OCI)
# ==============================================================================

test_oci_connection() {
    if ! command -v oci >/dev/null 2>&1; then
        ui_error "OCI CLI não está instalado ou não está no PATH."
        return 1
    fi

    local current_profile="${OCI_CLI_PROFILE:-DEFAULT}"
    ui_section "Testando Conexão Oracle OCI (Profile: $current_profile)"

    local output
    if [ "$HAS_GUM" -eq 1 ]; then
        output=$(gum spin --spinner dot --title "Consultando Oracle Cloud Namespace..." -- oci os ns get --profile "$current_profile" 2>&1)
    else
        echo "⏳ Consultando Oracle Cloud Namespace..."
        output=$(oci os ns get --profile "$current_profile" 2>&1)
    fi

    if [ $? -eq 0 ]; then
        local ns
        ns=$(echo "$output" | grep -o '"data": "[^"]*' | cut -d'"' -f4)
        local region
        region=$(oci session validate --profile "$current_profile" 2>/dev/null | grep -o '"region": "[^"]*' | cut -d'"' -f4 || echo "sa-saopaulo-1")

        local details=" Namespace : ${ns:-oci-tenancy}
 Profile   : $current_profile
 Região    : ${region:-sa-saopaulo-1}
 Status    : AUTENTICADO & OPERACIONAL"

        ui_card "🏛️  ORACLE CLOUD (OCI) - IDENTIDADE VALIDADA" "$details"
        return 0
    else
        ui_error "Falha de autenticação na Oracle OCI com o profile '$current_profile'."
        echo -e "\033[2m$output\033[0m"
        return 1
    fi
}

list_oci_profiles() {
    if [ -f "$HOME/.oci/config" ]; then
        grep -E '^\[' "$HOME/.oci/config" | tr -d '[]'
    fi
}
