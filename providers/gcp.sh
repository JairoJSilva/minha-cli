#!/bin/bash
# ==============================================================================
# providers/gcp.sh - Módulo de Operações e Diagnóstico Google Cloud
# ==============================================================================

test_gcp_connection() {
    if ! command -v gcloud >/dev/null 2>&1; then
        ui_warn "Google Cloud SDK (gcloud) não está instalado no sistema."
        return 1
    fi

    local current_config="${CLOUDSDK_ACTIVE_CONFIG_NAME:-default}"
    ui_section "Testando Conexão GCP (Config: $current_config)"

    local project account
    project=$(gcloud config get-value project 2>/dev/null || echo "<não definido>")
    account=$(gcloud config get-value account 2>/dev/null || echo "<não autenticado>")

    local details=" Projeto : $project
 Conta   : $account
 Config  : $current_config"

    ui_card "🌐 GOOGLE CLOUD - IDENTIDADE" "$details"
}
