# 🔐 Fase 2 — Vault Embutido com AES-256-GCM + Sessões Efêmeras

**Data:** 2026-08-19  
**Tipo:** Nova Feature — Segurança / Vault / Sessão Efêmera  
**Fase:** 2 de 3 (Vault Embutido + Sessões Efêmeras)

---

## Contexto / Motivação

Após a Fase 1 (remoção de clientes fictícios + indicadores de status reais), a Fase 2 implementa o núcleo da proposta arquitetural aprovada:

> *"As secrets são guardadas no Vault. Ao chamar o contexto, a CLI vai ao vault, utiliza a credencial para se autenticar e entrega a sessão do terminal ao usuário. Ao encerrar a sessão, o terminal padrão não tem nada configurado."*

---

## Arquitetura do Vault

### Armazenamento
```
~/.config/minha-cli/
├── clients.json    ← metadados apenas (sem credenciais)
├── vault.enc       ← credenciais criptografadas (AES-256-GCM)
└── vault.key       ← chave mestra (32 bytes aleatórios, chmod 600)
```

### Criptografia
- **Algoritmo:** AES-256-GCM (autenticado — detecta adulteração)
- **Derivação de chave:** `scrypt` (N=32768, r=8, p=1) — resistente a brute-force
- **Salt:** derivado de `SHA-256("minha-cli-vault-v1")` — fixo por design
- **Nonce:** 12 bytes aleatórios gerados a cada escrita (nunca reutilizados)
- **Biblioteca:** `crypto/aes`, `crypto/cipher` (stdlib Go) + `golang.org/x/crypto/scrypt`

---

## O que foi implementado

### `internal/vault/vault.go` ← **NOVO PACOTE**

API pública do vault:

| Função | Descrição |
|---|---|
| `vault.Store(id, secret)` | Salva/atualiza credenciais criptografadas |
| `vault.Get(id)` | Retorna as credenciais descriptografadas |
| `vault.Delete(id)` | Remove as credenciais do vault |
| `vault.Has(id)` | Verifica se há credencial salva |
| `vault.List()` | Lista todos os IDs com credenciais |
| `vault.VaultFilePath()` | Retorna o caminho do arquivo para exibição |

`VaultSecret` armazena:
- **AWS:** `access_key_id`, `secret_access_key`, `region`, `session_token`, `role_arn`
- **OCI:** `user_ocid`, `tenancy_ocid`, `fingerprint`, `private_key_path`, `region`
- **GCP:** `service_account_json`, `project`
- **Azure:** `tenant_id`, `client_id`, `client_secret`

### `cmd/add.go`
- Nova **Etapa 2** no cadastro: após metadados, pergunta se deseja salvar credenciais no vault
- **AWS:** escolha de modo `profile` ou `sts` (com suporte a ARN de role)
- **OCI:** coleta OCID, tenancy, fingerprint, chave privada, região
- `HasVaultSecret: true` é definido no `clients.json` quando credenciais são salvas

### `cmd/switch.go` — Dual-mode
```
mc switch flowti
  │
  ├─[HasVaultSecret=true]──► vault.Get("flowti")
  │                          GenerateExportScriptFromVault(client, secret)
  │                          → AWS_ACCESS_KEY_ID=...  ← variáveis efêmeras
  │                          → AWS_SECRET_ACCESS_KEY=...
  │                          → AWS_DEFAULT_REGION=...
  │                          ✅ "Contexto ativado via Vault"
  │
  └─[HasVaultSecret=false]─► GenerateExportScript(client)  ← modo legado
                             → AWS_PROFILE=flowti  (depende de ~/.aws/credentials)
                             ⚠️ Aviso + sugestão de migrar para o vault
```

### `internal/env/export.go`
- Nova função `GenerateExportScriptFromVault(client, secret)`:
  - Injeta `AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY` diretamente (sem `~/.aws/credentials`)
  - Injeta variáveis OCI via env vars (sem `~/.oci/config`)
  - Limpa `AWS_PROFILE` para evitar conflito

### `go.mod`
- Adicionado `golang.org/x/crypto v0.38.0` (para `scrypt`)

---

## Fluxo Completo (Fase 2)

```
1. mc add
   → informa: Nome, ID, AWS Profile, OCI Profile, etc.
   → "Deseja salvar credenciais no vault?" → Sim
   → informa: AWS Key, Secret, Region
   → vault.Store("flowti", {AWSAccessKeyID: "AKIA...", ...})
   → clients.json: {..., "has_vault_secret": true}

2. mc switch flowti
   → client.HasVaultSecret == true
   → secret = vault.Get("flowti")
   → script = GenerateExportScriptFromVault(client, secret)
   → export AWS_ACCESS_KEY_ID="AKIA..."
   → export AWS_SECRET_ACCESS_KEY="xxxxx"
   → export AWS_DEFAULT_REGION="us-east-1"
   → Terminal recebe as variáveis

3. Usuário trabalha normalmente com aws cli, kubectl, oci, etc.

4. Terminal fecha
   → Variáveis de ambiente somem
   → Zero rastro de credenciais no sistema
   → ~/.aws/credentials não foi tocado
```

---

## Segurança

| Aspecto | Implementação |
|---|---|
| Credenciais em disco | Nunca em plaintext — sempre AES-256-GCM |
| Chave mestra | 32 bytes aleatórios, `chmod 600` |
| Autenticidade | GCM detecta adulteração (authenticated encryption) |
| Sessão | Variáveis de ambiente efêmeras — somem com o terminal |
| Backward compat | Clientes sem vault continuam funcionando (modo legado) |

---

## Próxima etapa

**Fase 3:** Suporte a STS Token temporário (AssumeRole + MFA) para clientes AWS que exigem autenticação federada.
