# ☁️ Minha-CLI — Multi-Cloud & SRE Context Manager

[![Bash](https://img.shields.io/badge/Language-Bash%205.0+-4EAA25?logo=gnu-bash&logoColor=white)](https://www.gnu.org/software/bash/)
[![TUI](https://img.shields.io/badge/Interface-Charm%20Gum%20TUI-FF79C6)](https://github.com/charmbracelet/gum)
[![Multi-Cloud](https://img.shields.io/badge/Cloud-AWS%20%7C%20OCI%20%7C%20GCP%20%7C%20Azure%20%7C%20K8s-blue)](https://aws.amazon.com/)
[![SRE Ready](https://img.shields.io/badge/Architecture-Modular%20SRE-orange)](#-arquitetura-modular-do-projeto)

O **`minha-cli`** é uma ferramenta de linha de comando (CLI/TUI) desenvolvida para engenheiros de DevOps, SREs e desenvolvedores que precisam gerenciar e alternar rapidamente entre múltiplos clientes, contas e contextos de nuvem (**AWS**, **Oracle Cloud OCI**, **Google Cloud Platform**, **Microsoft Azure** e **Kubernetes**), garantindo isolamento de sessão, validação em tempo real e persistência de variáveis de ambiente.

---

## 📑 Sumário

- [Destaques e Funcionalidades](#-destaques-e-funcionalidades)
- [🧠 Como Funciona o `mc switch` (Orquestração Simultânea)](#-como-funciona-o-mc-switch-orquestração-simultânea)
- [Instalação Rápida (1 Comando)](#-instalação-rápida-1-comando)
- [Tabela Rápida de Comandos](#-tabela-rápida-de-comandos)
- [Guia de Uso Passo a Passo](#-guia-de-uso-passo-a-passo)
  - [1. Menu Interativo Visual (`mc`)](#1-menu-interativo-visual-mc)
  - [2. Cadastrar Novo Cliente (`mc add`)](#2-cadastrar-novo-cliente-mc-add)
  - [3. Alternar de Contexto (`mc switch`)](#3-alternar-de-contexto-mc-switch)
  - [4. Visualizar Contexto Ativo (`mc status`)](#4-visualizar-contexto-ativo-mc-status)
  - [5. Testar Credenciais nas APIs (`mc test`)](#5-testar-credenciais-nas-apis-mc-test)
  - [6. Editar Cliente Existente (`mc edit`)](#6-editar-cliente-existente-mc-edit)
  - [7. Apagar Cliente (`mc delete`)](#7-apagar-cliente-mc-delete)
  - [8. Limpar Variáveis de Ambiente (`mc clear`)](#8-limpar-variáveis-de-ambiente-mc-clear)
- [Arquitetura Modular do Projeto](#-arquitetura-modular-do-projeto)
- [Pré-requisitos e Dependências](#-pré-requisitos-e-dependências)

---

## ✨ Destaques e Funcionalidades

- 🔄 **Troca de Contexto Multi-Cloud em 1 Segundo**: Exporta e atualiza simultaneamente `AWS_PROFILE`, `OCI_CLI_PROFILE`, `CLOUDSDK_ACTIVE_CONFIG_NAME`, `AZURE_SUBSCRIPTION` e o contexto do `kubectl`.
- 🧠 **Orquestração Inteligente Total**: Conecta você a **TUDO** do cliente em uma única ação (todas as nuvens vinculadas e o cluster Kubernetes correspondente).
- ➕ **CRUD Completo de Clientes**: Cadastre, edite, remova e liste clientes diretamente pelo terminal ou pelo menu visual.
- 🎨 **Interface Rica com Charm Gum**: Menus selecionáveis, banners estilizados, spinners de carregamento e confirmações de segurança com fallback nativo em ANSI.
- 🔍 **Testes de Conexão WhoAmI**: Valida se suas credenciais estão ativas consultando as APIs da AWS (STS) e Oracle OCI em tempo real.
- ⚡ **Otimização com RTK Proxy**: Suporte nativo a comandos ultra-rápidos e condensados quando o `rtk` estiver disponível.
- 🧹 **Reset de Segurança**: Limpa todas as variáveis de ambiente das nuvens com um único comando (`mc clear`), prevenindo execuções acidentais em contas de clientes.

---

## 🧠 Como Funciona o `mc switch` (Orquestração Simultânea)

O **`mc switch`** funciona como o seu **orquestrador único de contexto**. Em vez de você ter que lembrar de rodar múltiplos comandos manuais para cada nuvem e para o cluster, ele sincroniza **todas as ferramentas do terminal ao mesmo tempo em 1 segundo**.

```
                        ┌─────────────────────────────────┐
                        │        mc switch <cliente>      │
                        └────────────────┬────────────────┘
                                         │
               ┌─────────────────────────┼─────────────────────────┐
               ▼                         ▼                         ▼
        [ 1. LIMPEZA ]            [ 2. NUVENS ]            [ 3. KUBERNETES ]
      Dá unset no que não       Exporta todas as          Executa:
     pertence ao novo cliente   variáveis simultâneas     kubectl config use-context
     (Sem contaminação cruzada) (AWS, OCI, GCP, Azure)    (Aponta para o cluster certo)
```

### O Ciclo em 3 Etapas:

1. **🧹 Limpeza de Segurança (Isolamento de Sessão)**:
   - Antes de ativar o novo cliente, o script dá `unset` em qualquer variável ou nuvem que pertencia ao cliente anterior.
   - *Exemplo*: Se você estava na Maida (GCP) e foi para o Dentalis (apenas AWS), o GCP e a Oracle são limpos para você não rodar comandos por engano na nuvem de outro cliente.

2. **☁️ Conexão Simultânea a Todas as Nuvens**:
   - **AWS**: Exporta `export AWS_PROFILE="cliente"` (usado por `aws-cli`, `terraform`, SDKs).
   - **Oracle Cloud (OCI)**: Exporta `export OCI_CLI_PROFILE="cliente"` (usado por `oci-cli` e SDKs).
   - **Google Cloud**: Exporta `export CLOUDSDK_ACTIVE_CONFIG_NAME="cliente"` (`gcloud`).
   - **Microsoft Azure**: Configura `AZURE_SUBSCRIPTION` e executa `az account set --subscription "ID"`.

3. **☸️ Chaveamento Automático do Kubernetes (`kubectl`)**:
   - Se o cliente possuir um cluster associado (ex: `oci-mv-devops`, `eks-prod`), o `mc switch` executa:
     ```bash
     kubectl config use-context <nome-do-cluster>
     ```
   - A partir desse exato milissegundo, comandos como `kubectl get pods`, `helm`, `k9s` ou `lens` já operam no cluster correto do cliente.

---

## 📦 Instalação Rápida (1 Comando)

Qualquer pessoa que baixar ou clonar este repositório pode instalar em sua máquina executando:

```bash
# 1. Entre no diretório do projeto
cd minha-cli

# 2. Execute o instalador universal autônomo
chmod +x install.sh && ./install.sh

# 3. Recarregue o seu terminal (apenas uma vez)
source ~/.bashrc
```

> **O que o instalador faz automaticamente:**
> - Valida e instala dependências necessárias (`jq`).
> - Aplica permissões de execução nos scripts (`chmod +x`).
> - Cria o diretório de dados em `~/.config/minha-cli/clients.json`.
> - Cria o binário no PATH do usuário (`~/.local/bin/mc`).
> - Injeta os aliases `mc` e `nuvem` com suporte a `source` no `~/.bashrc`, `~/.zshrc` e `~/.bash_profile`.

---

## ⚡ Tabela Rápida de Comandos

| Comando | Atalho | Descrição |
|---|---|---|
| `mc` | `nuvem` | Abre o **Menu Interativo (TUI)** com todas as opções |
| `mc add` | `mc novo` | **Cadastra um novo cliente/conta** interativamente |
| `mc edit` | `mc editar` | **Edita as configurações** de um cliente cadastrado |
| `mc delete` | `mc apagar` | **Remove um cliente** com confirmação de segurança |
| `mc list` | `mc ls` | **Lista todas as contas** e seus mapeamentos de nuvens |
| `mc switch <nome>` | `mc s <nome>` | **Alterna de contexto** direto (ex: `mc switch maida`) |
| `mc status` | `mc st` | Exibe o card com os **contextos e perfis ativos** |
| `mc test` | `mc whoami` | **Testa a autenticação real** nas APIs das nuvens |
| `mc clear` | `mc reset` | **Dá `unset` em todas as variáveis** de ambiente de nuvem |
| `mc k8s` | `mc kube` | Verifica a conectividade com o **cluster Kubernetes ativo** |
| `mc help` | `mc -h` | Exibe o guia de ajuda rápido |

---

## 🎓 Guia de Uso Passo a Passo

### 1. Menu Interativo Visual (`mc`)
Basta digitar `mc` ou `nuvem` para navegar com as setas do teclado:

```bash
mc
```

```
  ╔══════════════════════════════════════════════════════════════╗
  ║          ☁️   MINHA CLI - MULTI-CLOUD & SRE CONTEXT           ║
  ║        AWS • Oracle OCI • Google Cloud • Azure • K8s         ║
  ╚══════════════════════════════════════════════════════════════╝

  > ☁️  1. Trocar Contexto de Nuvem (Switch Profile)
    📊 2. Status do Contexto Ativo
    🔍 3. Testar Conexão / WhoAmI (AWS & OCI & K8s)
    ➕ 4. Configurar Nova Conta / Cliente (Add)
    ✏️  5. Editar Conta / Cliente Existente (Edit)
    🗑️  6. Apagar Conta / Cliente (Delete)
    📁 7. Mapeamento de Perfis Cadastrados
    ☸️  8. Kubernetes (Status do Cluster)
    🧹 9. Limpar Contexto (Reset de Variáveis)
    🚪 10. Sair
```

---

### 2. Cadastrar Novo Cliente (`mc add`)
Para adicionar um novo cliente à base de dados da CLI:

```bash
mc add
```
O assistente solicitará os dados de forma guiada:
- **Nome de exibição**: `Hospital Albert Einstein`
- **ID curto / Slug**: `einstein`
- **AWS Profile**: `einstein-prod` *(opcional)*
- **Oracle OCI Profile**: `einstein-oci` *(opcional)*
- **Google Cloud Config**: `einstein-gcp` *(opcional)*
- **Azure Subscription**: `ID-SUBSCRIPTION` *(opcional)*
- **Contexto K8s**: `cluster-einstein` *(opcional)*

---

### 3. Alternar de Contexto (`mc switch`)
Você pode alternar pelo menu ou passando o nome direto:

```bash
# Modo Interativo:
mc switch

# Modo Direto:
mc switch maida
mc switch dentalis
mc switch farmacia
mc switch flowti
```

---

### 4. Visualizar Contexto Ativo (`mc status`)
Veja instantaneamente qual conta e cluster estão selecionados no seu terminal:

```bash
mc status
```

```
 ┌─────────────────────────────────┐ 
 │ 📊 STATUS DO CONTEXTO ATIVO     │ 
 │  AWS Profile  : dentalis        │ 
 │  OCI Profile  : <não definido>  │ 
 │  GCP Config   : <não definido>  │ 
 │  Azure Context: <padrão/sessão> │ 
 │  Kubernetes   : oci-mv-devops   │ 
 └─────────────────────────────────┘ 
```

---

### 5. Testar Credenciais nas APIs (`mc test`)
Faz chamadas com spinner às APIs da AWS (STS) e Oracle OCI para verificar a validade das chaves:

```bash
mc test
```

---

### 6. Editar Cliente Existente (`mc edit`)
Se precisar alterar um profile AWS, subscription Azure ou contexto Kubernetes:

```bash
mc edit
```
Selecione o cliente na lista e atualize os valores desejados.

---

### 7. Apagar Cliente (`mc delete`)
Para remover uma conta que não é mais utilizada:

```bash
mc delete
```
O sistema pedirá uma confirmação (`Y/n`) para evitar exclusões acidentais.

---

### 8. Limpar Variáveis de Ambiente (`mc clear`)
Para resetar o terminal no final do expediente ou antes de trocar de projeto sensível:

```bash
mc clear
```

---

## 🏛️ Arquitetura Modular do Projeto

```
minha-cli/
├── bin/
│   └── mc.sh                      # Entrypoint principal (TUI Interativa e CLI)
├── config/
│   └── clients.json               # Base de dados estruturada de clientes e contas
├── core/
│   ├── ui.sh                      # Componentes visuais, cards, spinners e banners (gum/ansi)
│   ├── config.sh                  # CRUD de clientes (Add, Edit, Delete, List e Apply)
│   └── state.sh                   # Leitura de variáveis ativas e rotina de reset
├── providers/
│   ├── aws.sh                     # Módulo AWS (STS get-caller-identity, perfis)
│   ├── oci.sh                     # Módulo Oracle Cloud (Tenancy, namespaces, profiles)
│   ├── gcp.sh                     # Módulo Google Cloud SDK (gcloud configs)
│   ├── azure.sh                   # Módulo Azure CLI (Subscriptions e accounts)
│   └── k8s.sh                     # Módulo Kubernetes (kubectl contexts e clusters)
├── install.sh                     # Instalador universal autônomo
└── README.md                      # Documentação completa
```

---

## 📋 Pré-requisitos e Dependências

- **Sistema Operacional**: Linux (Ubuntu, Debian, RedHat, CentOS, Fedora, Arch) ou WSL2 no Windows ou macOS.
- **Shell**: Bash 4.4+ ou Zsh.
- **Ferramentas Recomendadas**:
  - `jq` *(instalado automaticamente pelo `install.sh` se disponível)*
  - `gum` *(opcional, fornece a interface gráfica moderna no terminal)*
  - `rtk` *(opcional, proxy de otimização de saída)*
  - `aws-cli`, `oci-cli`, `azure-cli`, `google-cloud-sdk`, `kubectl` *(conforme a sua necessidade de nuvens)*

---

<div align="center">
  <sub>Construído com maestria em SRE e Engenharia de Plataforma. 🚀</sub>
</div>
