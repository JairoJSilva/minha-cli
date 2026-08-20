# 🔄 Menu com Loop — Voltar ao Menu Após Cada Ação

**Data:** 2026-08-19  
**Tipo:** Nova Feature — UX / Navegação  
**Arquivos alterados:** `internal/tui/menu.go`, `cmd/root.go`

---

## Contexto / Motivação

Anteriormente, ao executar `mc` sem argumentos:

1. O banner era exibido
2. O menu abria uma única vez
3. O usuário escolhia uma ação (ex: Switch, Status, Add...)
4. A ação executava
5. **O programa encerrava**

Isso forçava o usuário a redigitar `mc` toda vez que quisesse executar outra ação, tornando o fluxo de trabalho fragmentado. Não havia forma de "voltar ao menu" após concluir uma operação, nem dentro dos submenus (pressionar `ESC` encerrava a CLI completamente).

---

## O que foi alterado

### `internal/tui/menu.go`

| Elemento | Mudança |
|---|---|
| `RunMenuLoop()` | **Nova função** — executa o menu sem imprimir o banner (para uso no loop) |
| `RunMenu()` | Refatorada para chamar `PrintBanner()` + `RunMenuLoop()` (mantida para compatibilidade) |

### `cmd/root.go`

| Elemento | Mudança |
|---|---|
| `tui.PrintBanner()` | Movido para fora do loop — exibido **uma única vez** ao abrir |
| Loop `for { ... }` | **Novo loop principal** — mantém o menu ativo entre ações |
| Condição de saída | `selected == ""` (q/ESC) ou `selected == "exit"` encerram o programa |

---

## Como funciona

### Antes
```
mc → banner → menu → [ação] → fim
```

### Depois
```
mc → banner → loop:
                ├─ menu abre
                ├─ usuário escolhe ação
                ├─ ação executa
                ├─ menu reabre automaticamente  ← novo
                └─ repete até: q / ESC / opção "13. Sair"
```

### Condições de saída do loop

| Evento | Resultado |
|---|---|
| Usuário aperta `q` ou `ESC` no menu | `selected == ""` → `return` (encerra) |
| Usuário seleciona opção **13. Sair** | `selected == "exit"` → `return` (encerra) |
| Qualquer outra seleção | Ação executa → loop continua → menu reabre |
| `Ctrl+C` | Interrompe o processo normalmente |

---

## Impacto para o usuário

- ✅ Após executar **Switch**, o menu reabre automaticamente
- ✅ Após ver o **Status**, o menu reabre automaticamente
- ✅ Após **Add**, **Edit**, **Delete** ou qualquer ação — mesmo comportamento
- ✅ Cancelar um submenu com `ESC` (ex: dentro do `huh` form) retorna ao menu principal
- ✅ O banner aparece **apenas uma vez**, sem poluir a tela a cada iteração
- ✅ Para sair: `q`, `ESC` no menu ou digitar `13` (opção Sair)
