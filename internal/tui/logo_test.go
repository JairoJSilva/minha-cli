package tui_test

import (
	"testing"

	"github.com/JairoJSilva/minha-cli/internal/tui"
)

func TestGetLogo(t *testing.T) {
	logo := tui.GetLogo()
	if logo == "" {
		t.Fatal("Expected rendered logo string, got empty string")
	}
	if len(logo) < 100 {
		t.Fatalf("Rendered logo string too short: %d bytes", len(logo))
	}
	t.Logf("Successfully rendered embedded logo (%d bytes)", len(logo))
}

func TestPrintBanner(t *testing.T) {
	// Ensure PrintBanner runs without errors or panics
	tui.PrintBanner()
}
