# 🤖 Adaptação dos Agentes e Skills para Contexto CLI (Llavero)

**Data:** 2026-08-19  
**Tipo:** Atualização de Diretrizes e Regras  

---

## Contexto / Motivação

O arquivo `AGENTS.md` e a skill `qa_tester` ainda continham diretrizes focadas no antigo portal web "Omnicloud" (UX/UI, React, CSS, Dashboards BI). Como a ferramenta atual evoluiu para ser puramente baseada em terminal (CLI construída em Go com Bubble Tea), os agentes precisavam ter suas habilidades e responsabilidades redefinidas para atuarem corretamente no novo ambiente.

---

## O que foi alterado

### `.agents/AGENTS.md`
Os papéis do time foram reescritos para a identidade **Llavero**:
- **Agente 1:** Llavero Core Developer (Go) — Focado na engine do Llavero, Cobra, Huh, IPC.
- **Agente 2:** Security & Software Architect — Focado em AES-GCM (Go stdlib crypto) e Vault embutido.
- **Agente 3:** Cloud Integrations Engineer — Especializado em injeção de sessão e tokens temporários multi-cloud.
- **Agente 4:** CLI UX & TUI Designer — Especializado no "WOW effect", Bubble Tea e usabilidade de terminal.

### `.agents/skills/qa_tester/SKILL.md`
As validações do QA Tester agora exigem:
- **Testes de TUI:** Quebra de linha no terminal, feedback visual (cores/spinners), navegação por setas.
- **Testes de Sessão Efêmera:** Garantir export e unset correto no Shell, sem vazar credenciais.
- **Testes Go / Pré-Commit:** Build sem falhas (`go build`), detecção de imports quebrados e checagem de formatação.

---

## Como funciona

A partir deste momento, quando o time for acionado (seja para desenvolver uma nova feature, arquitetar uma solução, ou realizar o QA), ele assumirá a identidade de desenvolvedores de CLI em Go, fornecendo código idiomático, respeitando padrões de terminal e buscando a melhor experiência possível para o usuário (TUI).

---

## Impacto

O time não vai mais recomendar uso de HTML/CSS ou bibliotecas JS neste repositório. O escopo está travado no ecossistema Go/CLI.

> **Nota:** A imagem da logo `Llavero.jpg` compartilhada é a nossa identidade visual! Embora não seja possível exibi-la dentro da CLI em terminal de texto, ela é excelente para constar no `README.md` do repositório no Github (ou GitLab).
