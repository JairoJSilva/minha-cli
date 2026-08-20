package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Title string
	Key   string
}

type MenuModel struct {
	Items     []MenuItem
	Cursor    int
	Selected  string
	Quitting  bool
	numBuffer string
}

func InitialMenuModel() MenuModel {
	return MenuModel{
		Items: []MenuItem{
			{Title: "☁️   1. Trocar Contexto de Nuvem (Switch Profile)", Key: "switch"},
			{Title: "📊  2. Status do Contexto Ativo", Key: "status"},
			{Title: "⚡  3. Testar Conexão / WhoAmI (AWS & OCI & K8s)", Key: "test"},
			{Title: "🔍  4. Detalhes da Conta (Show Info)", Key: "show"},
			{Title: "📡  5. Escanear e Importar Configurações (Scan)", Key: "scan"},
			{Title: "➕  6. Cadastrar Nova Conta / Cliente (Add)", Key: "add"},
			{Title: "✏️   7. Editar Conta / Cliente Existente (Edit)", Key: "edit"},
			{Title: "🗑️   8. Apagar Conta / Cliente (Delete)", Key: "delete"},
			{Title: "📁  9. Mapeamento de Perfis Cadastrados (List)", Key: "list"},
			{Title: "☸️  10. Kubernetes (Status do Cluster)", Key: "k8s"},
			{Title: "🧹  11. Limpar Contexto (Reset de Variáveis)", Key: "clear"},
			{Title: "ℹ️   12. Sobre / Versão da CLI (Version)", Key: "version"},
			{Title: "🚪  13. Sair", Key: "exit"},
		},
		Cursor: 0,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// Teclas numéricas: posicionam o cursor no item correspondente.
		// Suporta 1–9 direto e 10–13 com dois dígitos consecutivos.
		if len(key) == 1 && key >= "0" && key <= "9" {
			m.numBuffer += key
			if idx, err := strconv.Atoi(m.numBuffer); err == nil {
				if idx >= 1 && idx <= len(m.Items) {
					m.Cursor = idx - 1
				} else if idx > len(m.Items) {
					// Número inválido — descarta o buffer
					m.numBuffer = ""
				}
			}
			return m, nil
		}

		// Qualquer outra tecla limpa o buffer numérico
		m.numBuffer = ""

		switch key {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Items) - 1
			}
		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
		case "enter", " ":
			m.Selected = m.Items[m.Cursor].Key
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	if m.Quitting {
		return ""
	}

	s := "\n"
	cursorStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))

	for i, item := range m.Items {
		cursor := "  "
		itemStr := normalStyle.Render(item.Title)

		if m.Cursor == i {
			cursor = cursorStyle.Render("▶ ")
			itemStr = cursorStyle.Render(item.Title)
		}

		s += fmt.Sprintf("%s %s\n", cursor, itemStr)
	}

	s += "\n" + MutedText.Render("  (Navegue com ↑/↓ ou j/k • Digite o número para posicionar • Enter para confirmar • q para sair)") + "\n"
	return s
}

// RunMenuLoop executa o menu sem imprimir o banner (usado no loop do root).
func RunMenuLoop() (string, error) {
	p := tea.NewProgram(InitialMenuModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	model := m.(MenuModel)
	return model.Selected, nil
}

// RunMenu imprime o banner e executa o menu (mantido para compatibilidade).
func RunMenu() (string, error) {
	PrintBanner()
	return RunMenuLoop()
}
