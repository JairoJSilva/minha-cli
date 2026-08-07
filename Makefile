# ==============================================================================
# Makefile para Minha-CLI (Go Edition)
# ==============================================================================

BINARY_NAME=mc
BUILD_DIR=bin
GO=go

.PHONY: all build install clean test

all: build

build:
	@mkdir -p $(BUILD_DIR)
	@echo "🔨 Compilando Minha-CLI em Go..."
	$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "✅ Binário compilado com sucesso em $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@chmod +x install.sh
	@./install.sh

clean:
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	@echo "🧹 Limpeza concluída."

test:
	@$(GO) test -v ./...
