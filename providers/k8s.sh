#!/bin/bash
# ==============================================================================
# providers/k8s.sh - Módulo de Operações e Contexto Kubernetes (com suporte RTK)
# ==============================================================================

run_k8s_cmd() {
    if command -v rtk >/dev/null 2>&1; then
        rtk kubectl "$@"
    elif [ -x "/root/.local/bin/rtk" ]; then
        /root/.local/bin/rtk kubectl "$@"
    else
        kubectl "$@"
    fi
}

test_k8s_connection() {
    if ! command -v kubectl >/dev/null 2>&1; then
        ui_warn "kubectl não está instalado no PATH."
        return 1
    fi

    local current_ctx
    current_ctx=$(kubectl config current-context 2>/dev/null || echo "<nenhum>")

    ui_section "Verificando Cluster Kubernetes (Contexto: $current_ctx)"

    if [ "$current_ctx" = "<nenhum>" ]; then
        ui_warn "Nenhum contexto Kubernetes selecionado."
        return 1
    fi

    local output
    if [ "$HAS_GUM" -eq 1 ]; then
        output=$(gum spin --spinner dot --title "Consultando API Server K8s (RTK)..." -- bash -c "kubectl cluster-info --request-timeout='3s' 2>&1" 2>&1)
    else
        echo "⏳ Consultando API Server K8s (RTK)..."
        output=$(kubectl cluster-info --request-timeout='3s' 2>&1)
    fi

    if [ $? -eq 0 ]; then
        local node_count
        node_count=$(kubectl get nodes --no-headers 2>/dev/null | wc -l || echo "N/D")

        local details=" Contexto    : $current_ctx
 Nodes Ativos: $node_count
 RTK Proxy   : DISPONÍVEL
 Status      : CLUSTER ONLINE & RESPONSIVO"

        ui_card "☸️  KUBERNETES - CLUSTER ATIVO" "$details"
        return 0
    else
        ui_warn "Não foi possível conectar ao cluster Kubernetes ativo ($current_ctx)."
        return 1
    fi
}

list_k8s_contexts() {
    if command -v kubectl >/dev/null 2>&1; then
        kubectl config get-contexts --no-headers -o name 2>/dev/null
    fi
}

switch_k8s_context() {
    local target_ctx="$1"
    if [ -n "$target_ctx" ] && command -v kubectl >/dev/null 2>&1; then
        kubectl config use-context "$target_ctx"
        ui_success "Contexto Kubernetes alterado para: $target_ctx"
    fi
}
