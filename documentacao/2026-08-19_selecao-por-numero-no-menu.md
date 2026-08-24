# ⌨️ Seleção por Número no Menu Interativo

**Data:** 2026-08-19  
**Tipo:** Nova Feature — UX / TUI  
**Arquivo principal:** `internal/tui/menu.go`

---

## Contexto / Motivação

O menu interativo da CLI (`mc`) já suportava navegação com setas (`↑`/`↓`) e seleção com `Enter`. No entanto, como o menu exibe os itens numerados (ex: `1. Switch`, `2. Status`), era natural esperar que digitar o número correspondente já selecionasse a opção diretamente — comportamento ausente até então.

---

## O que foi alterado

### `internal/tui/menu.go`

| Elemento | Mudança |
|---|---|
| Imports | Adicionados `strconv` e `time` |
| `MenuModel` | Adicionados os campos `numBuffer string` e `numTimer *time.Timer` |
| `numberTimerMsg` | Novo tipo de mensagem interno para o timer do bubbletea |
| `Update()` | Lógica de captura de teclas numéricas com buffer de 500ms |
| `View()` | Rodapé atualizado com a dica `• Digite o número •` |

---

## Como funciona

### Lógica de seleção imediata (opções 2–9)
Ao pressionar qualquer dígito de `2` a `9`, o item correspondente é selecionado **imediatamente**, sem aguardar `Enter`.

### Lógica de buffer para opções 10–13
O dígito `1` é ambíguo: pode ser a opção `1` (Switch) **ou** o primeiro dígito de `10`, `11`, `12`, `13`. Por isso:

1. Ao pressionar `1`, o sistema inicia um **timer de 500ms** e aguarda um segundo dígito.
2. Se o segundo dígito for `0`, `1`, `2` ou `3` — forma o número completo (10–13) e seleciona.
3. Se o timer expirar sem segundo dígito — seleciona a **opção 1** (Switch).

### Fluxo resumido

```
Tecla pressionada
       │
       ├─ "2"–"9"  ──► Seleciona item[n-1] diretamente
       │
       ├─ "1"      ──► Inicia timer 500ms
       │                    │
       │              Novo dígito "0"–"3" ──► Seleciona item[10|11|12|13]
       │                    │
       │              Timer expira ──────────► Seleciona item[1] (Switch)
       │
       └─ outros   ──► Navegação normal (↑/↓/Enter/q)
```

---

## Impacto para o usuário

- ✅ Digitar `2` abre diretamente o **Status do Contexto Ativo**
- ✅ Digitar `9` abre **Mapeamento de Perfis**
- ✅ Digitar `1` + `0` abre **Kubernetes**
- ✅ Digitar `1` sozinho (aguardar meio segundo) abre **Switch Profile**
- ✅ Toda navegação anterior (setas, j/k, Enter) continua funcionando normalmente
