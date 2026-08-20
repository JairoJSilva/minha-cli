# 🐛 Fase 1 — Bug Fix: Perfis Falsos + Fundação para Vault

**Data:** 2026-08-19  
**Tipo:** Bug Fix + Refatoração Estrutural  
**Fase:** 1 de 3 (Vault Embutido + Sessões Efêmeras)

---

## Contexto / Motivação

O comando `mc list` e `mc scan` exibiam perfis como "configurados" mesmo quando não existia nenhuma credencial real na máquina. A causa raiz eram **4 clientes hardcoded** gravados automaticamente no `clients.json` na primeira execução da CLI, independente da configuração real do ambiente.

Além disso, a função `checkAWSProfileExists` estava duplicada em vários arquivos (`switch.go`, `show.go`, `add.go`, `edit.go`), violando o princípio DRY.

---

## O que foi alterado

### `internal/config/client.go`
- **Removidos** os 4 clientes hardcoded (`maida`, `dentalis`, `farmacia`, `flowti`) que eram gravados na 1ª execução
- Base de dados começa **vazia** — o usuário cadastra apenas o que realmente usa
- **Adicionados** campos preparatórios para a Fase 2:
  - `AWSAuthMode string` → `"profile"` (padrão) ou `"sts"` (token temporário)
  - `AWSRoleARN *string` → ARN para AssumeRole (opcional)
  - `HasVaultSecret bool` → indica se a credencial está salva no vault (Fase 2)

### `internal/config/credentials.go` ← **NOVO ARQUIVO**
Centraliza todas as verificações de credencial local:
- `AWSProfileLocalExists(profile)` → verifica `~/.aws/credentials` e `~/.aws/config`
- `OCIProfileLocalExists(profile)` → verifica `~/.oci/config`
- `GCPConfigLocalExists(configName)` → verifica `~/.config/gcloud/configurations/`
- `K8sContextLocalExists(ctxName)` → usa `kubectl config get-contexts`

### `cmd/list.go`
- Reescrito para mostrar status real por provider com indicadores visuais:
  - `✅ nome` → credencial encontrada localmente
  - `📋 nome (sem credencial)` → perfil cadastrado, mas sem credencial local
  - `🔐 nome` → credencial no vault (Fase 2)
  - `-` → provider não configurado para este cliente
- Legenda exibida no rodapé do card

### `cmd/switch.go`, `cmd/show.go`, `cmd/add.go`, `cmd/edit.go`
- Todas as chamadas a `checkAWSProfileExists()` substituídas por `config.AWSProfileLocalExists()`
- Função privada duplicada **removida** do `switch.go`

### `config/clients.json`
- Resetado para `[]` (lista vazia) — dados fictícios removidos

---

## Como funciona agora

### `mc list` — Antes
```
• Maida (AWS, GCP, Azure) [ID: maida]
  └─ AWS: maida | OCI: - | GCP: maida | K8s: -
```
*(Exibia como "configurado" mesmo sem credencial real)*

### `mc list` — Depois
```
• Maida (AWS, GCP, Azure) [ID: maida]
  ├─ AWS  : 📋 maida (sem credencial)
  ├─ GCP  : 📋 maida (sem credencial)
  └─ Azure: ID-SUBSCRIPTION-MAIDA

  ✅ credencial local OK  📋 só referência  🔐 no vault
```

### `mc list` — Após Fase 2 (com vault)
```
• Maida (AWS, GCP, Azure) [ID: maida] 🔐
  ├─ AWS  : 🔐 maida
  ├─ GCP  : 🔐 maida
  └─ Azure: ID-SUBSCRIPTION-MAIDA
```

---

## Impacto para o usuário

- ✅ Nenhum perfil fantasma na base — base começa limpa
- ✅ `mc list` é honesto: mostra o que realmente está configurado
- ✅ Usuário consegue distinguir "só cadastrei o nome" vs "tenho credencial aqui"
- ✅ Estrutura preparada para receber o vault na Fase 2
- ✅ Zero duplicação de código de verificação de credencial

---

## Próxima etapa

**Fase 2:** Vault embutido com `age` — armazenar credenciais reais criptografadas em `~/.config/minha-cli/vault.age`
