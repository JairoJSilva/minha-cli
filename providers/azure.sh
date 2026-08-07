#!/bin/bash
# ==============================================================================
# providers/azure.sh - Módulo de Operações e Diagnóstico Microsoft Azure
# ==============================================================================

test_azure_connection() {
    if ! command -v az >/dev/null 2>&1; then
        ui_warn "Azure CLI (az) não está instalado no sistema."
        return 1
    fi

    ui_section "Testando Conexão Azure"
    local output
    if [ "$HAS_GUM" -eq 1 ]; then
        output=$(gum spin --spinner dot --title "Consultando Azure Account..." -- az account show --output json 2>&1)
    else
        echo "⏳ Consultando Azure Account..."
        output=$(az account show --output json 2>&1)
    fi

    if [ $? -eq 0 ]; then
        local sub_name sub_id user_name
        sub_name=$(echo "$output" | grep -o '"name": "[^"]*' | head -n 1 | cut -d'"' -f4)
        sub_id=$(echo "$output" | grep -o '"id": "[^"]*' | head -n 1 | cut -d'"' -f4)
        user_name=$(echo "$output" | grep -o '"name": "[^"]*' | tail -n 1 | cut -d'"' -f4)

        local details=" Subscription : $sub_name
 Sub ID       : $sub_id
 Usuário      : $user_name
 Status       : AUTENTICADO"

        ui_card "🔷 MICROSOFT AZURE - IDENTIDADE" "$details"
        return 0
    else
        ui_warn "Sessão Azure não autenticada (rode 'az login')."
        return 1
    fi
}
