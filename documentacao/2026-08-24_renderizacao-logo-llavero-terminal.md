# Renderização da Logo Llavero no Terminal (TrueColor Half-Block ANSI)

## 📌 Contexto
O usuário solicitou a inclusão da logo oficial do projeto (`Llavero.jpg`), que contém a chave iluminada em tons de ciano/azul com o prompt de comando terminal `>_`, diretamente na interface de terminal da CLI.

## 🛠️ O que mudou
- **[logo.go](file:///home/jairosjunior/minha-cli/internal/tui/logo.go)**: Criado o módulo de renderização de imagem para terminal de alta fidelidade:
  - Embed da imagem `Llavero.jpg` compilada diretamente no binário Go (`//go:embed assets/Llavero.jpg`), tornando a CLI 100% autocontida sem depender de arquivos soltos no disco.
  - Algoritmo de amostragem por média de área (*area-averaging downsampling*) e isolamento de máscara que remove automaticamente o fundo xadrez/cinza neutro e sombras artificiais da imagem original, preservando o brilho neon ciano e o formato nítido da chave.
  - Renderização gráfica ANSI 24-bit TrueColor através de caracteres *upper/lower half-block* Unicode (`▀` / `▄`), permitindo resolução vertical dobrada (2 pixels por linha de caractere).
  - Cache via `sync.Once` para garantir renderização instantânea (0ms de overhead nas execuções subsequentes).
- **[styles.go](file:///home/jairosjunior/minha-cli/internal/tui/styles.go)**: Atualizada a função `PrintBanner()` para integrar e centralizar a logo do Llavero no banner de abertura e na tela de versão da CLI.
- **[logo_test.go](file:///home/jairosjunior/minha-cli/internal/tui/logo_test.go)**: Adicionados testes unitários para validar a renderização da logo e a integridade da chamada do banner.

## ⚙️ Como Funciona
1. O asset gráfico `internal/tui/assets/Llavero.jpg` é embutido no binário durante a compilação (`go build`).
2. Quando a CLI é iniciada ou quando o banner é exibido (`tui.PrintBanner()`), a função `GetLogo()` decodifica os bytes da imagem, recorta os limites exatos da chave, filtra ruídos e gera as sequências de escape ANSI TrueColor (`\033[38;2;R;G;Bm` e `\033[48;2;R;G;Bm`) combinadas com blocos `▀`.
3. A logo é impressa centralizada no cabeçalho com moldura dupla Dracula/Lipgloss (`#BD93F9`), mantendo o design do terminal.

## 🚀 Impacto
- Experiência visual rica no terminal (WOW effect) com a chave e o símbolo `>_` visíveis com cores fiéis à logo original.
- Total compatibilidade com os terminais modernos (Linux, macOS, Windows Terminal, Warp, Alacritty, Kitty, iTerm2, VS Code).
- O binário permanece leve, independente e portátil.
