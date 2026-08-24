package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Cores do Tema SRE / DevOps
	ColorPrimary   = lipgloss.Color("#FF79C6") // Magenta vibrante
	ColorSecondary = lipgloss.Color("#8BE9FD") // Ciano
	ColorSuccess   = lipgloss.Color("#50FA7B") // Verde neon
	ColorWarning   = lipgloss.Color("#FFB86C") // Laranja
	ColorError     = lipgloss.Color("#FF5555") // Vermelho
	ColorMuted     = lipgloss.Color("#6272A4") // Cinza/Azul escuro
	ColorBg        = lipgloss.Color("#282A36") // Fundo escuro

	// Estilos Lipgloss
	BannerBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#BD93F9")).
			Padding(0, 2).
			Margin(0, 1).
			Align(lipgloss.Center).
			Width(68)

	BannerTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	BannerSub = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	CardBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Padding(0, 1).
		Margin(0, 1)

	CardTitle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	SuccessTag = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	WarningTag = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	ErrorTag = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	MutedText = lipgloss.NewStyle().
			Foreground(ColorMuted)
)

func PrintBanner() {
	logo := GetLogo()
	title := BannerTitle.Render("☁️  LLAVERO / MINHA-CLI (GO EDITION) — MULTI-CLOUD VAULT")
	sub := BannerSub.Render("AWS • Oracle OCI • Google Cloud • Azure • Kubernetes (SRE Powered)")
	var content string
	if logo != "" {
		content = fmt.Sprintf("%s\n%s\n%s", logo, title, sub)
	} else {
		content = fmt.Sprintf("%s\n%s", title, sub)
	}
	fmt.Println(BannerBox.Render(content))
}

func PrintCard(title, body string) {
	header := CardTitle.Render("📊 " + title)
	content := fmt.Sprintf("%s\n%s", header, body)
	fmt.Println(CardBox.Render(content))
}

func Success(msg string) {
	fmt.Println(SuccessTag.Render("✅ " + msg))
}

func Warn(msg string) {
	fmt.Println(WarningTag.Render("⚠️  " + msg))
}

func Error(msg string) {
	fmt.Println(ErrorTag.Render("❌ " + msg))
}

func Info(msg string) {
	fmt.Println(BannerSub.Render("ℹ️  " + msg))
}
