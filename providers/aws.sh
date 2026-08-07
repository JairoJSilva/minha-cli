#!/bin/bash
# ==============================================================================
# providers/aws.sh - Módulo de Operações e Diagnóstico AWS (com suporte RTK)
# ==============================================================================

# Wrapper inteligente que utiliza rtk se disponível
run_aws_cmd() {
    if command -v rtk >/dev/null 2>&1; then
        rtk aws "$@"
    elif [ -x "/root/.local/bin/rtk" ]; then
        /root/.local/bin/rtk aws "$@"
    else
        aws "$@"
    fi
}

test_aws_connection() {
    if ! command -v aws >/dev/null 2>&1; then
        ui_error "AWS CLI não está instalado ou não está no PATH."
        return 1
    fi

    local current_profile="${AWS_PROFILE:-default}"
    ui_section "Testando Conexão AWS (Profile: $current_profile)"

    local output
    if [ "$HAS_GUM" -eq 1 ]; then
        output=$(gum spin --spinner dot --title "Consultando AWS STS via RTK..." -- bash -c "run_aws_cmd sts get-caller-identity --output json 2>&1" 2>&1)
    else
        echo "⏳ Consultando AWS STS via RTK..."
        output=$(run_aws_cmd sts get-caller-identity --output json 2>&1)
    fi

    if [ $? -eq 0 ] && echo "$output" | grep -q '"Account"'; then
        local account_id user_arn
        account_id=$(echo "$output" | grep -o '"Account": "[^"]*' | cut -d'"' -f4)
        user_arn=$(echo "$output" | grep -o '"Arn": "[^"]*' | cut -d'"' -f4)
        local region="${AWS_REGION:-${AWS_DEFAULT_REGION:-$(aws configure get region 2>/dev/null || echo 'us-east-1')}}"

        local details=" Conta ID : $account_id
 Região   : $region
 Identity : $user_arn
 RTK      : ATIVO (Output Otimizado)
 Status   : AUTENTICADO & OPERACIONAL"

        ui_card "☁️  AWS - IDENTIDADE VALIDADA" "$details"
        return 0
    else
        # Tenta fallback caso rtk ou output necessite de tratamento direto
        local direct_out
        direct_out=$(aws sts get-caller-identity --output json 2>&1 || true)
        if echo "$direct_out" | grep -q '"Account"'; then
            local account_id user_arn
            account_id=$(echo "$direct_out" | grep -o '"Account": "[^"]*' | cut -d'"' -f4)
            user_arn=$(echo "$direct_out" | grep -o '"Arn": "[^"]*' | cut -d'"' -f4)
            local region="${AWS_REGION:-${AWS_DEFAULT_REGION:-$(aws configure get region 2>/dev/null || echo 'us-east-1')}}"

            local details=" Conta ID : $account_id
 Região   : $region
 Identity : $user_arn
 Status   : AUTENTICADO & OPERACIONAL"

            ui_card "☁️  AWS - IDENTIDADE VALIDADA" "$details"
            return 0
        fi

        ui_error "Falha de autenticação na AWS com o profile '$current_profile'."
        echo -e "\033[2m$output\033[0m"
        return 1
    fi
}

list_aws_profiles() {
    if [ -f "$HOME/.aws/credentials" ]; then
        grep -E '^\[' "$HOME/.aws/credentials" | tr -d '[]'
    elif [ -f "$HOME/.aws/config" ]; then
        grep -E '^\[' "$HOME/.aws/config" | sed -e 's/\[profile //g' -e 's/\[//g' -e 's/\]//g'
    fi
}
