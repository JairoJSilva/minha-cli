#!/usr/bin/env python3
"""
generate_screenshots.py - Gerador de Prints de Alta Fidelidade para o Minha-CLI
Gera mockups visuais profissionais com estética moderna de terminal (macOS / Tokyo Night / Charm Gum TUI).
"""

import os
from PIL import Image, ImageDraw, ImageFont

SCREENSHOTS_DIR = "/mnt/c/Users/jairo.sjunior/Documents/Jairo/Pessoal/minha-cli/assets/screenshots"
os.makedirs(SCREENSHOTS_DIR, exist_ok=True)

# Fontes do sistema
FONT_MONO = "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
FONT_MONO_BOLD = "/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Bold.ttf"

# Paleta Tokyo Night / Terminal Moderno
BG_COLOR = (15, 17, 26)          # Fundo geral / canvas
WINDOW_BG = (24, 26, 40)         # Fundo do terminal
TITLEBAR_BG = (30, 32, 48)       # Barra de título
BORDER_COLOR = (49, 54, 82)      # Borda sutil

# Cores de Sintaxe / ANSI
C_WHITE = (220, 226, 245)
C_MUTED = (110, 115, 148)
C_CYAN = (125, 207, 255)
C_MAGENTA = (187, 154, 247)
C_GREEN = (115, 218, 202)
C_YELLOW = (224, 175, 104)
C_RED = (247, 118, 142)
C_BLUE = (122, 162, 247)
C_ORANGE = (255, 158, 100)
C_SELECTION = (45, 52, 85)
C_CARD_BG = (29, 32, 49)

def render_terminal(title, prompt_cmd, lines, width=960):
    font_size = 15
    line_height = 25
    padding_x = 32
    padding_top = 56
    padding_bottom = 28
    
    font_reg = ImageFont.truetype(FONT_MONO, font_size)
    font_bold = ImageFont.truetype(FONT_MONO_BOLD, font_size)
    font_title = ImageFont.truetype(FONT_MONO_BOLD, 13)
    
    total_lines = len(lines) + (2 if prompt_cmd else 0)
    calculated_height = padding_top + (total_lines * line_height) + padding_bottom
    
    margin = 24
    canvas_w = width + (margin * 2)
    canvas_h = calculated_height + (margin * 2)
    
    img = Image.new("RGBA", (canvas_w, canvas_h), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Sombra suave multicamadas
    for i in range(12, 0, -1):
        alpha = int(65 * (1 - i / 12))
        shadow_rect = [
            margin - i + 2, margin - i + 6,
            canvas_w - margin + i - 2, canvas_h - margin + i + 4
        ]
        draw.rounded_rectangle(shadow_rect, radius=14, fill=(5, 6, 12, alpha))
        
    # Janela Principal
    win_rect = [margin, margin, canvas_w - margin, canvas_h - margin]
    draw.rounded_rectangle(win_rect, radius=12, fill=WINDOW_BG, outline=BORDER_COLOR, width=1)
    
    # Barra de Título
    titlebar_rect = [margin, margin, canvas_w - margin, margin + 38]
    draw.rounded_rectangle(titlebar_rect, radius=12, fill=TITLEBAR_BG)
    draw.rectangle([margin, margin + 26, canvas_w - margin, margin + 38], fill=TITLEBAR_BG)
    draw.line([(margin, margin + 38), (canvas_w - margin, margin + 38)], fill=BORDER_COLOR, width=1)
    
    # Botões macOS
    btn_y = margin + 13
    btn_r = 6
    draw.ellipse([margin + 18, btn_y, margin + 18 + btn_r*2, btn_y + btn_r*2], fill=(255, 95, 86))
    draw.ellipse([margin + 38, btn_y, margin + 38 + btn_r*2, btn_y + btn_r*2], fill=(255, 189, 46))
    draw.ellipse([margin + 58, btn_y, margin + 58 + btn_r*2, btn_y + btn_r*2], fill=(39, 201, 63))
    
    # Título central
    title_text = f"💻 {title}"
    draw.text((canvas_w // 2, margin + 19), title_text, fill=C_MUTED, font=font_title, anchor="mm")
    
    curr_y = margin + padding_top
    
    # Linha de Prompt do Shell
    if prompt_cmd:
        draw.text((margin + padding_x, curr_y), "➜ ", fill=C_GREEN, font=font_bold)
        curr_x = margin + padding_x + 22
        draw.text((curr_x, curr_y), "minha-cli ", fill=C_CYAN, font=font_bold)
        curr_x += 88
        draw.text((curr_x, curr_y), "(main) ", fill=C_MAGENTA, font=font_reg)
        curr_x += 65
        draw.text((curr_x, curr_y), prompt_cmd, fill=C_WHITE, font=font_bold)
        curr_y += int(line_height * 1.5)
        
    for item in lines:
        if len(item) == 4:
            f_type, color, text, bg_highlight = item
        elif len(item) == 3:
            f_type, color, text = item
            bg_highlight = None
        else:
            continue
            
        f = font_bold if f_type == 'bold' else font_reg
        
        if bg_highlight:
            highlight_rect = [
                margin + padding_x - 6, curr_y - 2,
                canvas_w - margin - padding_x + 6, curr_y + line_height - 4
            ]
            draw.rounded_rectangle(highlight_rect, radius=6, fill=bg_highlight)
            
        draw.text((margin + padding_x, curr_y), text, fill=color, font=f)
        curr_y += line_height
        
    return img

def generate_all():
    banner_lines = [
        ("bold", C_MAGENTA, " ╔════════════════════════════════════════════════════════════════╗ "),
        ("bold", C_CYAN,    " ║           ☁️  MINHA CLI - MULTI-CLOUD & SRE CONTEXT            ║ "),
        ("bold", C_WHITE,   " ║      AWS • Oracle OCI • Google Cloud • Azure • Kubernetes      ║ "),
        ("bold", C_MAGENTA, " ╚════════════════════════════════════════════════════════════════╝ "),
        ("reg",  C_MUTED,   ""),
    ]

    # 1. 01_mc_menu.png (Menu Interativo Gum TUI)
    lines_menu = list(banner_lines) + [
        ("bold", C_YELLOW, "? Selecione a operação desejada:"),
        ("bold", C_CYAN,   "> 1. Alternar Cliente / Contexto (switch)", C_SELECTION),
        ("reg",  C_WHITE,  "  2. Ver Status Atual do Terminal (status)"),
        ("reg",  C_WHITE,  "  3. Listar Clientes e Perfis Cadastrados (list)"),
        ("reg",  C_WHITE,  "  4. Escanear e Importar Configs Existentes (scan)"),
        ("reg",  C_WHITE,  "  5. Cadastrar Novo Cliente (add)"),
        ("reg",  C_WHITE,  "  6. Editar Configurações de Cliente (edit)"),
        ("reg",  C_WHITE,  "  7. Testar Conexão e Credenciais (test)"),
        ("reg",  C_WHITE,  "  8. Limpar Todas as Variáveis de Nuvem (clear)"),
        ("reg",  C_MUTED,  "  9. Sair"),
    ]
    img1 = render_terminal("mc — Menu Interativo TUI", "mc", lines_menu)
    img1.save(os.path.join(SCREENSHOTS_DIR, "01_mc_menu.png"))
    print("✅ Gerado: 01_mc_menu.png")

    # 2. 02_mc_switch.png (Orquestração Simultânea Multi-Cloud)
    lines_switch = list(banner_lines) + [
        ("bold", C_YELLOW,  "▶ Sincronizando contexto com Flowti / Pessoal..."),
        ("bold", C_GREEN,   "✅ Contexto ativado: Flowti / Pessoal (AWS, Oracle OCI)"),
        ("reg",  C_WHITE,   "   • AWS Profile       : flowti"),
        ("reg",  C_WHITE,   "   • Oracle OCI Profile: pessoal"),
        ("reg",  C_MUTED,   "   • Google Cloud      : <não configurado>"),
        ("reg",  C_MUTED,   "   • Microsoft Azure   : <não configurado>"),
        ("bold", C_CYAN,    "   • Kubernetes (k8s)  : Switched to context \"oci-mv-devops\"."),
        ("reg",  C_MUTED,   ""),
        ("bold", C_BLUE,    "ℹ️  Terminal 100% sincronizado com Flowti em 0.8s (Sem contaminação cruzada)"),
    ]
    img2 = render_terminal("mc switch flowti — Orquestração Multi-Cloud & K8s", "mc switch flowti", lines_switch)
    img2.save(os.path.join(SCREENSHOTS_DIR, "02_mc_switch.png"))
    print("✅ Gerado: 02_mc_switch.png")

    # 3. 03_mc_status.png (Status do Contexto Ativo)
    lines_status = [
        ("bold", C_CYAN,  " ┌────────────────────────────────────────────────────────┐ "),
        ("bold", C_YELLOW," │ 📊 STATUS DO CONTEXTO ATIVO                            │ "),
        ("bold", C_WHITE, " │  AWS Profile  : flowti                                 │ "),
        ("bold", C_WHITE, " │  OCI Profile  : pessoal                                │ "),
        ("reg",  C_MUTED, " │  GCP Config   : <não definido>                         │ "),
        ("reg",  C_MUTED, " │  Azure Context: <padrão/sessão>                        │ "),
        ("bold", C_GREEN, " │  Kubernetes   : oci-mv-devops (Oracle Cloud OKE)       │ "),
        ("bold", C_CYAN,  " └────────────────────────────────────────────────────────┘ "),
    ]
    img3 = render_terminal("mc status — Diagnóstico em Tempo Real", "mc status", lines_status)
    img3.save(os.path.join(SCREENSHOTS_DIR, "03_mc_status.png"))
    print("✅ Gerado: 03_mc_status.png")

    # 4. 04_mc_list.png (Listagem em Árvore)
    lines_list = [
        ("bold", C_CYAN,  " ┌────────────────────────────────────────────────────────────────┐ "),
        ("bold", C_YELLOW," │ 📁 PERFIS E CLIENTES CADASTRADOS                               │ "),
        ("bold", C_WHITE, " │  • Maida (AWS, GCP, Azure) [ID: maida]                         │ "),
        ("reg",  C_MUTED, " │    └─ AWS: maida | OCI: - | GCP: maida | K8s: -                │ "),
        ("bold", C_WHITE, " │  • Dentalis (AWS) [ID: dentalis]                               │ "),
        ("reg",  C_MUTED, " │    └─ AWS: dentalis | OCI: - | GCP: - | K8s: -                 │ "),
        ("bold", C_WHITE, " │  • Farmacia Digital (AWS, GCP, Azure) [ID: farmacia]           │ "),
        ("reg",  C_MUTED, " │    └─ AWS: farmacia | OCI: - | GCP: farmacia | K8s: -          │ "),
        ("bold", C_GREEN, " │  • Flowti / Pessoal (AWS, Oracle OCI) [ID: flowti]             │ "),
        ("reg",  C_CYAN,  " │    └─ AWS: flowti | OCI: pessoal | GCP: - | K8s: oci-mv-devops │ "),
        ("bold", C_CYAN,  " └────────────────────────────────────────────────────────────────┘ "),
    ]
    img4 = render_terminal("mc list — Lista de Perfis e Provedores", "mc list", lines_list)
    img4.save(os.path.join(SCREENSHOTS_DIR, "04_mc_list.png"))
    print("✅ Gerado: 04_mc_list.png")

    # 5. 05_mc_scan.png (Auto-Descoberta e Proteção)
    lines_scan = list(banner_lines) + [
        ("bold", C_YELLOW, "▶ 🔍 Escaneando Configurações Existentes na Máquina..."),
        ("bold", C_CYAN,   " ┌── [ 📊 DIAGNÓSTICO DO TERMINAL / MÁQUINA ] ─────────────────── "),
        ("reg",  C_WHITE,  " │  AWS Profiles     : 4 encontrados (maida, dentalis, farmacia, flowti)"),
        ("reg",  C_WHITE,  " │  OCI Profiles     : 2 encontrados (pessoal, devops-prod)"),
        ("reg",  C_WHITE,  " │  GCP Configs      : 2 encontrados (maida, farmacia)"),
        ("reg",  C_WHITE,  " │  K8s Contextos    : 3 encontrados (oci-mv-devops, eks-prod, gke-cluster)"),
        ("reg",  C_WHITE,  " │  Azure Assinaturas: 2 encontradas (ID-MAIDA, ID-FARMACIA)"),
        ("bold", C_CYAN,   " └─────────────────────────────────────────────────────────────── "),
        ("bold", C_GREEN,  "✅ 4 perfis importados/sincronizados com sucesso em config/clients.json."),
        ("bold", C_BLUE,   "ℹ️  100% seguro: nenhuma configuração preexistente foi sobrescrita ou perdida."),
    ]
    img5 = render_terminal("mc scan — Auto-Descoberta e Proteção de Dados", "mc scan", lines_scan)
    img5.save(os.path.join(SCREENSHOTS_DIR, "05_mc_scan.png"))
    print("✅ Gerado: 05_mc_scan.png")

    # 6. 06_mc_test.png (WhoAmI & Conexão)
    lines_test = list(banner_lines) + [
        ("bold", C_YELLOW, "▶ Executando Testes de Identidade e Conexão (WhoAmI)..."),
        ("bold", C_CYAN,   " [AWS STS]    "),
        ("bold", C_GREEN,  "   ✅ Conectado como arn:aws:iam::123456789012:user/jairo (Profile: flowti)"),
        ("bold", C_CYAN,   " [Oracle OCI] "),
        ("bold", C_GREEN,  "   ✅ Autenticado na Tenancy: ocid1.tenancy.oc1..aaaa (User: jairo@flowti)"),
        ("bold", C_CYAN,   " [Kubernetes] "),
        ("bold", C_GREEN,  "   ✅ Contexto 'oci-mv-devops' ativo - Cluster Kubernetes v1.29 OK"),
        ("reg",  C_MUTED,  ""),
        ("bold", C_GREEN,  "✅ Todas as credenciais validadas com sucesso nas APIs oficiais das nuvens."),
    ]
    img6 = render_terminal("mc test — Testes WhoAmI em Tempo Real", "mc test", lines_test)
    img6.save(os.path.join(SCREENSHOTS_DIR, "06_mc_test.png"))
    print("✅ Gerado: 06_mc_test.png")

    # 7. 07_mc_clear.png (Reset de Segurança)
    lines_clear = [
        ("bold", C_YELLOW, "▶ 🧹 Resetando variáveis de ambiente de todas as nuvens..."),
        ("bold", C_GREEN,  "✅ Todas as variáveis de ambiente das nuvens foram limpas com sucesso."),
        ("reg",  C_MUTED,  ""),
        ("reg",  C_WHITE,  "   • Unset: AWS_PROFILE, AWS_ACCESS_KEY_ID, AWS_SESSION_TOKEN"),
        ("reg",  C_WHITE,  "   • Unset: OCI_CLI_PROFILE, OCI_CLI_REGION"),
        ("reg",  C_WHITE,  "   • Unset: CLOUDSDK_ACTIVE_CONFIG_NAME, GOOGLE_APPLICATION_CREDENTIALS"),
        ("reg",  C_WHITE,  "   • Unset: AZURE_SUBSCRIPTION, AZURE_TENANT_ID"),
        ("reg",  C_WHITE,  "   • Unset: KUBECONFIG"),
        ("reg",  C_MUTED,  ""),
        ("bold", C_BLUE,   "🔒 Sessão segura: nenhum comando acidental será executado em contas de clientes."),
    ]
    img7 = render_terminal("mc clear — Reset de Segurança", "mc clear", lines_clear)
    img7.save(os.path.join(SCREENSHOTS_DIR, "07_mc_clear.png"))
    print("✅ Gerado: 07_mc_clear.png")

if __name__ == "__main__":
    generate_all()
