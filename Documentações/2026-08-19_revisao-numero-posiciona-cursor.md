# ⌨️ Revisão: Número Posiciona Cursor, Enter Confirma

**Data:** 2026-08-19  
**Tipo:** Revisão de Feature — UX / TUI  
**Arquivo principal:** `internal/tui/menu.go`

---

## Contexto / Motivação

A implementação anterior de seleção por número (`2026-08-19_selecao-por-numero-no-menu.md`) utilizava auto-seleção imediata com timer de 500ms. Isso causou um problema de UX:

- Ao executar uma ação (ex: Switch de contexto), o terminal exibia o resultado mas o programa encerrava antes do usuário conseguir ler o feedback com clareza.
- A abordagem com timer era complexa e frágil.

**Solução adotada:** número apenas **posiciona o cursor** no item correspondente — a ação só é executada quando o usuário pressiona `Enter` (ou `Space`).

---

## O que foi alterado

### `internal/tui/menu.go`

| Elemento | Mudança |
|---|---|
| Import `time` | **Removido** — não é mais necessário |
| `MenuModel.numTimer` | **Removido** — não é mais necessário |
| `numberTimerMsg` | **Removido** — não é mais necessário |
| `numBuffer` | Mantido, mas agora apenas acumula dígitos para calcular o índice |
| `Update()` | Número → move `Cursor`; `Enter`/`Space` → confirma seleção |
| `View()` rodapé | Texto atualizado: `"Digite o número para posicionar • Enter para confirmar"` |

---

## Como funciona agora

```
Usuário digita "1" → cursor move para item 1 (Switch)
Usuário digita "1" + "0" → cursor move para item 10 (Kubernetes)
Usuário pressiona Enter → ação do item em destaque é executada
```

### Suporte a números de dois dígitos (10–13)
O `numBuffer` acumula dígitos consecutivos. Ao digitar `1` seguido de `0`, o buffer vira `"10"` e o cursor vai para o item 10. Qualquer tecla não numérica limpa o buffer.

### Condições de saída
- `q`, `ESC`, `Ctrl+C` → encerram o programa
- Número inválido (> 13) → buffer descartado, cursor permanece

---

## Impacto para o usuário

- ✅ Digitar `2` + `Enter` executa **Status** (com retorno visível no terminal)
- ✅ Digitar `1` + `0` + `Enter` executa **Kubernetes**
- ✅ O terminal exibe o resultado da ação com tempo para leitura
- ✅ Toda navegação anterior (setas, j/k, Enter) continua funcionando normalmente
- ✅ Código mais simples, sem timers ou goroutines desnecessárias
