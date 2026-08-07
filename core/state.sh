#!/bin/bash
# ==============================================================================
# core/state.sh - Leitura, Gerenciamento e Reset de Contextos de Ambiente
# ==============================================================================

get_active_aws() {
    echo "${AWS_PROFILE:-<não definido>}"
}

get_active_oci() {
    echo "${OCI_CLI_PROFILE:-<não definido>}"
}

get_active_gcp() {
    echo "${CLOUDSDK_ACTIVE_CONFIG_NAME:-<não definido>}"
}

get_active_azure() {
    if [ -n "$AZURE_SUBSCRIPTION_NAME" ]; then
        echo "$AZURE_SUBSCRIPTION_NAME"
    elif [ -n "$AZURE_SUBSCRIPTION" ]; then
        echo "$AZURE_SUBSCRIPTION"
    else
        echo "<padrão/sessão>"
    fi
}

get_active_k8s() {
    if command -v kubectl >/dev/null 2>&1; then
        local ctx
        ctx=$(kubectl config current-context 2>/dev/null || true)
        if [ -n "$ctx" ]; then
            echo "$ctx"
        else
            echo "<não conectado>"
        fi
    else
        echo "<kubectl não instalado>"
    fi
}

show_context_status() {
    local aws_p oci_p gcp_p az_p k8s_p
    aws_p=$(get_active_aws)
    oci_p=$(get_active_oci)
    gcp_p=$(get_active_gcp)
    az_p=$(get_active_azure)
    k8s_p=$(get_active_k8s)

    local body=" AWS Profile  : $aws_p
 OCI Profile  : $oci_p
 GCP Config   : $gcp_p
 Azure Context: $az_p
 Kubernetes   : $k8s_p"

    ui_card "📊 STATUS DO CONTEXTO ATIVO" "$body"
}

clear_all_contexts() {
    # AWS
    unset AWS_PROFILE
    unset AWS_ACCESS_KEY_ID
    unset AWS_SECRET_ACCESS_KEY
    unset AWS_SESSION_TOKEN
    unset AWS_DEFAULT_REGION

    # OCI
    unset OCI_CLI_PROFILE
    unset OCI_CLI_REGION
    unset OCI_CLI_KEY_FILE

    # GCP
    unset CLOUDSDK_ACTIVE_CONFIG_NAME
    unset GOOGLE_APPLICATION_CREDENTIALS

    # Azure
    unset AZURE_SUBSCRIPTION
    unset AZURE_SUBSCRIPTION_NAME
    unset AZURE_TENANT_ID
    unset AZURE_CLIENT_ID
    unset AZURE_CLIENT_SECRET

    # Kubernetes (opcional, reseta KUBECONFIG se houver customizado)
    unset KUBECONFIG

    ui_success "🧹 Todas as variáveis de ambiente das nuvens foram limpas com sucesso."
}
