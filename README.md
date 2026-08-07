# ☁️ Minha-CLI — Multi-Cloud & SRE Context Manager (Go Edition)

[![Go](https://img.shields.io/badge/Language-Go%201.25+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![CLI Framework](https://img.shields.io/badge/CLI-Cobra%20Framework-blue)](https://github.com/spf13/cobra)
[![TUI](https://img.shields.io/badge/TUI-Bubbletea%20%26%20Lipgloss-FF79C6)](https://github.com/charmbracelet/bubbletea)
[![Multi-Cloud](https://img.shields.io/badge/Cloud-AWS%20%7C%20OCI%20%7C%20GCP%20%7C%20Azure%20%7C%20K8s-blue)](https://aws.amazon.com/)
[![SRE Ready](https://img.shields.io/badge/Architecture-Modular%20SRE-orange)](#-arquitetura-do-projeto-em-go)

O **`minha-cli`** é uma ferramenta CLI/TUI compilada nativamente em **Go** desenvolvida para engenheiros de DevOps, SREs e desenvolvedores que precisam gerenciar e alternar rapidamente entre múltiplos clientes, contas e contextos de nuvem (**AWS**, **Oracle Cloud OCI**, **Google Cloud Platform**, **Microsoft Azure** e **Kubernetes**), garantindo isolamento de sessão, validação em paralelo com Goroutines e persistência de variáveis de ambiente no terminal.

---

## 📑 Sumário

- [Por que em Go? Destaques e Vantagens](#-por-que-em-go-destaques-e-vantagens)
- [🧠 Como Funciona o `mc switch` (Orquestração Simultânea)](#-como-funciona-o-mc-switch-orquestração-simultânea)
- [📡 Auto-Descoberta e Proteção de Dados (`mc scan`)](#-auto-descoberta-e-proteção-de-dados-mc-scan--mc-leitura)
- [Instalação e Compilação Rápida (1 Comando)](#-instalação-e-compilação-rápida-1-comando)
- [Tabela Rápida de Comandos](#-tabela-rápida-de-comandos)
- [Guia de Uso Passo a Passo](#-guia-de-uso-passo-a-passo)
  - [1. Menu Interativo Visual (`mc`)](#1-menu-interativo-visual-mc)
  - [2. Escanear e Importar Configurações (`mc scan`)](#2-escanear-e-importar-configurações-mc-scan)
  - [3. Cadastrar Novo Cliente (`mc add`)](#3-cadastrar-novo-cliente-mc-add)
  - [4. Alternar de Contexto (`mc switch`)](#4-alternar-de-contexto-mc-switch)
  - [5. Visualizar Contexto Ativo (`mc status`)](#5-visualizar-contexto-ativo-mc-status)
  - [6. Testar Credenciais nas APIs em Paralelo (`mc test`)](#6-testar-credenciais-nas-apis-em-paralelo-mc-test)
  - [7. Editar Cliente Existente (`mc edit`)](#7-editar-cliente-existente-mc-edit)
  - [8. Apagar Cliente (`mc delete`)](#8-apagar-cliente-mc-delete)
  - [9. Limpar Variáveis de Ambiente (`mc clear`)](#9-limpar-variáveis-de-ambiente-mc-clear)
- [Arquitetura do Projeto em Go](#-arquitetura-do-projeto-em-go)
- [Pré-requisitos](#-pré-requisitos)

---

## 🚀 Por que em Go? Destaques e Vantagens

- 📦 **Binário Único Standalone**: Distribuição de um executável compilado de alta performance sem dependências externas de script.
- ⚡ **Goroutines e Testes em Paralelo (`mc test`)**: O diagnóstico de WhoAmI consulta AWS STS, Oracle OCI e Kubernetes simultaneamente em paralelo, retornando o resultado em menos de 200ms.
- 🎨 **TUI Moderna com Bubbletea, Lipgloss e Huh**: Menus selecionáveis com teclado, formulários com validação em tempo real e cards com bordas arredondadas.
- 🔄 **Orquestração Multi-Cloud & K8s**: Exporta e atualiza simultaneamente `AWS_PROFILE`, `OCI_CLI_PROFILE`, `CLOUDSDK_ACTIVE_CONFIG_NAME`, `AZURE_SUBSCRIPTION` e o contexto do `kubectl`.
- 📡 **Auto-Descoberta (`mc scan` / `mc leitura`)**: Varre o terminal (`~/.aws`, `~/.oci`, `~/.kube`, GCP, Azure) e importa tudo automaticamente sem perda de dados.
- ⚡ **Otimização com RTK Proxy**: Suporte nativo a consultas ultra-rápidas e condensadas quando o `rtk` estiver presente.
- 🧹 **Reset de Segurança**: Limpa todas as variáveis de ambiente das nuvens com um único comando (`mc clear`), prevenindo execuções acidentais em contas de clientes.

---

## 🧠 Como Funciona o `mc switch` (Orquestração Simultânea)

O **`mc switch`** funciona como o seu **orquestrador central de contexto**. Em vez de você ter que rodar múltiplos comandos manuais para cada nuvem e para o cluster, ele sincroniza **todas as ferramentas do terminal ao mesmo tempo em 1 segundo**.

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
2. **☁️ Conexão Simultânea a Todas as Nuvens**:
   - **AWS**: Exporta `export AWS_PROFILE="cliente"`.
   - **Oracle Cloud (OCI)**: Exporta `export OCI_CLI_PROFILE="cliente"`.
   - **Google Cloud**: Exporta `export CLOUDSDK_ACTIVE_CONFIG_NAME="cliente"`.
   - **Microsoft Azure**: Configura `AZURE_SUBSCRIPTION` e executa `az account set --subscription "ID"`.
3. **☸️ Chaveamento Automático do Kubernetes (`kubectl`)**:
   - Se o cliente possuir um cluster associado (ex: `oci-mv-devops`, `eks-prod`), ele executa `kubectl config use-context <cluster>`.

---

## 📦 Instalação e Compilação Rápida (1 Comando)

```bash
# 1. Clone ou acesse o diretório do projeto
cd minha-cli

# 2. Execute o instalador universal (compila o binário Go e configura o shell)
chmod +x install.sh && ./install.sh

# 3. Recarregue o seu terminal
source ~/.bashrc
```

---

## ⚡ Tabela Rápida de Comandos

| Comando | Atalho | Descrição |
|---|---|---|
| `mc` | `nuvem` | Abre o **Menu Interativo TUI (Bubbletea)** |
| `mc scan` | `mc leitura` | **Escaneia e importa configurações locais** sem perder nada |
| `mc add` | `mc novo` | **Cadastra um novo cliente/conta** via formulário interativo |
| `mc edit` | `mc editar` | **Edita as configurações** de um cliente cadastrado |
| `mc delete` | `mc apagar` | **Remove um cliente** com confirmação de segurança |
| `mc list` | `mc ls` | **Lista todas as contas** e seus mapeamentos de nuvens |
| `mc switch <nome>` | `mc s <nome>` | **Alterna de contexto** direto (ex: `mc switch maida`) |
| `mc status` | `mc st` | Exibe o card com os **contextos e perfis ativos** |
| `mc test` | `mc whoami` | **Testa a autenticação em paralelo com Goroutines** |
| `mc clear` | `mc reset` | **Dá `unset` em todas as variáveis** de ambiente de nuvem |
| `mc k8s` | `mc kube` | Verifica a conectividade com o **cluster Kubernetes ativo** |
| `mc help` | `mc -h` | Exibe o guia de ajuda |

---

## 🎓 Guia de Uso Passo a Passo

### 1. Menu Interativo Visual (`mc`)
Basta digitar `mc` para navegar com as setas do teclado:

```bash
mc
```

---

### 2. Escanear e Importar Configurações (`mc scan`)
Para quem acabou de clonar o projeto e já tem credenciais configuradas na máquina:

```bash
mc scan
# ou
mc leitura
```
O Minha-CLI detecta seus arquivos de configuração da AWS, OCI, K8s, GCP e Azure e importa automaticamente sem sobrescrever dados.

---

### 3. Cadastrar Novo Cliente (`mc add`)
Para adicionar um novo cliente via formulário guiado:

```bash
mc add
```

---

### 4. Alternar de Contexto (`mc switch`)
Você pode alternar pelo menu ou passando o nome direto:

```bash
mc switch maida
mc switch dentalis
mc switch farmacia
mc switch flowti
```

---

### 5. Visualizar Contexto Ativo (`mc status`)
Veja instantaneamente qual conta e cluster estão selecionados no seu terminal:

```bash
mc status
```

---

### 6. Testar Credenciais nas APIs em Paralelo (`mc test`)
Faz chamadas concorrentes com Goroutines às APIs da AWS (STS), Oracle OCI e Kubernetes para verificar a validade das chaves em milissegundos:

```bash
mc test
```

---

### 7. Editar Cliente Existente (`mc edit`)
```bash
mc edit
```

---

### 8. Apagar Cliente (`mc delete`)
```bash
mc delete
```

---

### 9. Limpar Variáveis de Ambiente (`mc clear`)
```bash
mc clear
```

---

## 🏛️ Arquitetura do Projeto em Go

```
minha-cli/ (Go Edition)
├── cmd/
│   ├── root.go                    # Cobra root command e launcher do Bubbletea
│   ├── switch.go                  # mc switch
│   ├── status.go                  # mc status
│   ├── test.go                    # mc test (Goroutines paralelas)
│   ├── scan.go                    # mc scan / mc leitura
│   ├── add.go                     # mc add (Formulário Huh)
│   ├── edit.go                    # mc edit
│   ├── delete.go                  # mc delete
│   ├── list.go                    # mc list (Cards Lipgloss)
│   └── clear.go                   # mc clear
├── internal/
│   ├── config/
│   │   ├── client.go              # CRUD e persistência de dados JSON
│   │   └── scanner.go             # Auto-descoberta não-destrutiva
│   ├── env/
│   │   └── export.go              # Gerador de scripts de export/unset
│   ├── providers/
│   │   ├── aws.go                 # Integração AWS STS WhoAmI & perfis
│   │   ├── oci.go                 # Integração Oracle Cloud OCI
│   │   ├── gcp.go                 # Integração Google Cloud
│   │   ├── azure.go               # Integração Microsoft Azure
│   │   └── k8s.go                 # Integração Kubernetes
│   └── tui/
│       ├── styles.go              # Estilos e temas Lipgloss
│       └── menu.go                # Modelo interativo Bubbletea
├── bin/
│   └── mc                         # Binário Go compilado standalone
├── config/
│   └── clients.json               # Base de clientes
├── install.sh                     # Instalador e compilador universal
├── Makefile                       # Targets: build, install, test, clean
├── go.mod                         # Módulo Go
├── go.sum                         # Checksums das dependências
└── README.md                      # Documentação completa
```

---

## 📋 Pré-requisitos

- **Go**: Versão 1.21+ (compilação).
- **Sistema Operacional**: Linux (Ubuntu, Debian, RedHat, CentOS, Fedora, Arch) ou WSL2 / macOS.
- **Shell**: Bash 4.4+ ou Zsh.

---

<div align="center">
  <sub>Construído com maestria em Go, SRE e Engenharia de Plataforma. 🚀</sub>
</div>
