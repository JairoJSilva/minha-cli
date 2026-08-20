# 🏛️ Análise Arquitetural — Vault + Sessões Efêmeras

> **Time**: Arquiteto · Fullstack · DevOps · UX/Data  
> **Data**: 2026-08-19

---

## 1. Bug: Scan Mostra Perfis "Configurados" que Não Existem

### Diagnóstico

O `mc list` e o `mc scan` exibem **o que está registrado no `clients.json`** — que é apenas um mapeamento de nomes (ex: `aws_profile: "maida"`). Eles **não verificam** se essas credenciais realmente existem nos arquivos de configuração locais (`~/.aws/credentials`, `~/.oci/config`, etc.).

**Raiz do problema:** o `client.go` possui clientes padrão (hardcoded) que são gravados automaticamente na primeira execução, mesmo que o usuário não tenha essas credenciais configuradas na máquina.

```go
// client.go — linha 63-90: clientes fictícios gravados na 1ª execução
defaultClients := []Client{
    {ID: "maida",    AWSProfile: StringPtr("maida"), ...},
    {ID: "dentalis", AWSProfile: StringPtr("dentalis"), ...},
    ...
}
```

### Correção proposta (simples, independente do Vault)

- Remover os clientes hardcoded padrão — a base começa vazia
- No `mc list`, adicionar um indicador visual de status real: ✅ credencial existe no sistema | ⚠️ apenas referência cadastrada

---

## 2. Proposta de Sessão Efêmera (Já funciona hoje, mas incompleta)

### Como funciona hoje

O `mc switch` gera **variáveis de ambiente** (`AWS_PROFILE`, `OCI_CLI_PROFILE`, etc.) e as injeta no shell via wrapper. Isso já é uma forma de sessão efêmera — quando o terminal é fechado, as variáveis somem.

**Porém há um problema:** o mecanismo usa apenas o **nome do profile** (ex: `AWS_PROFILE=maida`), o que significa que a AWS CLI ainda precisa ler as chaves de `~/.aws/credentials`. Se esse arquivo não existir ou não tiver o perfil, a autenticação falha.

### O modelo ideal que você descreveu

```
mc switch flowti
   │
   ├── vai ao Vault → busca as credenciais do "flowti"
   ├── usa as chaves para autenticar temporariamente (STS/token)
   └── injeta as variáveis temporárias no shell atual
        → AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN
        → (sem precisar de ~/.aws/credentials)
```

Quando o terminal fecha → as variáveis somem → **zero rastro**.

---

## 3. Análise de Opções de Vault

### Opção A — HashiCorp Vault ⭐ Recomendada para times

| Aspecto | Avaliação |
|---|---|
| Segurança | ✅ Máxima — audit log, políticas granulares, criptografia |
| Custo | ✅ Open source (self-hosted) / Pago (HCP Vault) |
| Complexidade de setup | ⚠️ Precisa de servidor Vault rodando |
| Integração com a CLI | ✅ API REST simples — fácil integrar em Go |
| Ideal para | Times, ambientes corporativos, múltiplos usuários |

**Fluxo:**
```
mc switch flowti
  → POST /v1/secret/data/flowti → retorna { aws_key, aws_secret, oci_key... }
  → injeta como vars temporárias
```

---

### Opção B — age + arquivo local criptografado ⭐ Recomendada para uso pessoal

| Aspecto | Avaliação |
|---|---|
| Segurança | ✅ Boa — criptografia forte com chave local ou SSH key |
| Custo | ✅ Gratuito, sem servidor |
| Complexidade de setup | ✅ Mínima — só instalar `age` |
| Integração com a CLI | ✅ Fácil — decrypt na hora do switch |
| Ideal para | Uso pessoal, SRE solo, ambientes simples |

**Fluxo:**
```
mc switch flowti
  → lê ~/.config/minha-cli/secrets.age
  → decrypt com chave SSH local (sem senha extra)
  → extrai as credenciais do flowti
  → injeta como vars temporárias
```

---

### Opção C — aws-vault (para AWS only)

Ferramenta existente especializada apenas em AWS. Não cobre OCI/GCP/Azure. Descartada por não ser multicloud.

---

### Opção D — Bitwarden Secrets Manager

| Aspecto | Avaliação |
|---|---|
| Segurança | ✅ Boa |
| Custo | 💰 Plano pago para o Secrets Manager |
| Complexidade | ✅ Baixa |
| Ideal para | Quem já usa Bitwarden |

---

## 4. Arquitetura Recomendada pelo Time

### Para o seu cenário (SRE/DevOps pessoal + multicloud):

> **Opção B (age local) como default + suporte a HashiCorp Vault como backend opcional**

**Motivo:**
- Não depende de servidor externo para uso pessoal
- Simples de configurar — integra com sua chave SSH existente
- Quando o contexto for corporativo/time → plug in do Vault
- O usuário pode escolher o backend na configuração da CLI

---

## 5. Fluxo Proposto Completo

```
┌─────────────────────────────────────────────────────┐
│                   mc add / mc edit                   │
│                                                      │
│  1. Usuário informa:                                 │
│     - Nome, ID, AWS Key, AWS Secret, OCI, GCP...    │
│  2. CLI criptografa com age (ou envia ao Vault)     │
│  3. client.json guarda APENAS: id, name, provider   │
│     (sem nenhuma credencial em plaintext)           │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│                   mc switch flowti                   │
│                                                      │
│  1. Busca no Vault/age as credenciais do "flowti"   │
│  2. Para AWS: faz STS GetSessionToken (token temp.) │
│  3. Injeta no shell:                                │
│     AWS_ACCESS_KEY_ID=...                           │
│     AWS_SECRET_ACCESS_KEY=...                       │
│     AWS_SESSION_TOKEN=...  ← expira em X horas     │
│     OCI_CLI_PROFILE=... (com config temp.)          │
│  4. Usuário trabalha normalmente                    │
└─────────────────────────────────────────────────────┘
                          │
                     Terminal fecha
                          │
                          ▼
              ✅ Zero rastro de credenciais
```

---

## 6. Open Questions — Preciso da sua decisão

> [!IMPORTANT]
> **Q1 — Backend de segredos:**
> Prefere começar com **age (local, sem servidor)** ou **HashiCorp Vault (servidor)**?
> A CLI pode suportar os dois com um flag de configuração.

> [!IMPORTANT]
> **Q2 — Tokens temporários AWS (STS):**
> Para AWS, ao fazer switch, prefere:
> - **A)** Injetar apenas `AWS_PROFILE` (comportamento atual) — depende do `~/.aws/credentials`
> - **B)** Gerar token STS temporário — independe de qualquer arquivo local, mais seguro

> [!IMPORTANT]
> **Q3 — Escopo do MVP:**
> Implementar tudo de uma vez ou em fases?
> - **Fase 1**: Corrigir o bug dos perfis falsos + remover clientes hardcoded
> - **Fase 2**: Implementar storage criptografado com `age`
> - **Fase 3**: Adicionar suporte a HashiCorp Vault como backend alternativo

> [!NOTE]
> **Q4 — Plataforma do Vault (se escolher HashiCorp):**
> - Self-hosted na sua infra?
> - HCP Vault (cloud gerenciada pela HashiCorp)?
> - Vault em container local (dev mode para testes)?

---

## 7. O que NÃO muda nessa proposta

- A interface TUI (menu, navegação) — permanece igual
- O `mc scan` — continua escaneando o ambiente local
- O `mc list` — continua listando perfis cadastrados (com status melhorado)
- O conceito de `clients.json` — permanece, mas **sem credenciais**
