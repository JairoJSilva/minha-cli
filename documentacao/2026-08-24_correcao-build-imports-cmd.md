# Correção de Compilação: Imports em `cmd/list.go` e `cmd/switch.go`

## 📌 Contexto
Ao executar o script de instalação `install.sh` do `minha-cli`, o processo de compilação do binário Go (`go build`) falhava com dois erros:
1. `cmd/list.go:95:25: undefined: providers`
2. `cmd/switch.go:5:2: "strings" imported and not used`

## 🛠️ O que mudou
- **[list.go](file:///home/jairosjunior/minha-cli/cmd/list.go)**: Adicionada a importação do pacote `github.com/JairoJSilva/minha-cli/internal/providers` necessária para invocar `providers.TestK8s("")` na função `runK8s()`.
- **[switch.go](file:///home/jairosjunior/minha-cli/cmd/switch.go)**: Removida a importação não utilizada do pacote `"strings"`.

## ⚙️ Como Funciona
- O Go possui verificação estrita de pacotes compilados: variáveis/funções não declaradas e imports não utilizados barram o `go build`.
- Com a importação correta do pacote `internal/providers` em `list.go` e a limpeza do import não utilizado em `switch.go`, o binário compila perfeitamente tanto de forma standalone quanto via script `install.sh`.

## 🚀 Impacto
- O script `bash install.sh` e o comando `go build` agora executam e concluem sem erros de compilação.
- O binário `bin/mc` é gerado e instalado corretamente no PATH do usuário.
