package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	AWSProfile *string `json:"aws_profile"`
	OCIProfile *string `json:"oci_profile"`
	GCPConfig  *string `json:"gcp_config"`
	AzureSub   *string `json:"azure_sub"`
	K8sContext *string `json:"k8s_context"`
}

// Helper para converter string em ponteiro de string
func StringPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

// Helper para ler string de ponteiro
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Retorna o caminho do arquivo clients.json preferencial
func GetConfigPath() string {
	// 1. Tenta caminho local no projeto config/clients.json
	localPath := "config/clients.json"
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// 2. Caminho global do usuário ~/.config/minha-cli/clients.json
	home, err := os.UserHomeDir()
	if err == nil {
		userDir := filepath.Join(home, ".config", "minha-cli")
		_ = os.MkdirAll(userDir, 0755)
		return filepath.Join(userDir, "clients.json")
	}

	return localPath
}

// Carrega todos os clientes da base JSON
func LoadClients() ([]Client, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Retorna lista padrão se não existir
			defaultClients := []Client{
				{
					ID:         "maida",
					Name:       "Maida (AWS, GCP, Azure)",
					AWSProfile: StringPtr("maida"),
					GCPConfig:  StringPtr("maida"),
					AzureSub:   StringPtr("ID-SUBSCRIPTION-MAIDA"),
				},
				{
					ID:         "dentalis",
					Name:       "Dentalis (AWS)",
					AWSProfile: StringPtr("dentalis"),
				},
				{
					ID:         "farmacia",
					Name:       "Farmacia Digital (AWS, GCP, Azure)",
					AWSProfile: StringPtr("farmacia"),
					GCPConfig:  StringPtr("farmacia"),
					AzureSub:   StringPtr("ID-SUBSCRIPTION-FARMACIA"),
				},
				{
					ID:         "flowti",
					Name:       "Flowti / Pessoal (AWS, Oracle OCI)",
					AWSProfile: StringPtr("flowti"),
					OCIProfile: StringPtr("pessoal"),
					K8sContext: StringPtr("oci-mv-devops"),
				},
			}
			_ = SaveClients(defaultClients)
			return defaultClients, nil
		}
		return nil, err
	}

	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

// Salva a lista de clientes no arquivo JSON
func SaveClients(clients []Client) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Busca um cliente por ID ou por parte do nome
func FindClient(query string) (*Client, error) {
	clients, err := LoadClients()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	for _, c := range clients {
		if strings.ToLower(c.ID) == q || strings.Contains(strings.ToLower(c.Name), q) {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("cliente '%s' não encontrado", query)
}

// Adiciona um novo cliente
func AddClient(client Client) error {
	clients, err := LoadClients()
	if err != nil {
		clients = []Client{}
	}

	// Evita IDs duplicados
	for _, c := range clients {
		if strings.EqualFold(c.ID, client.ID) {
			return fmt.Errorf("já existe um cliente cadastrado com o ID '%s'", client.ID)
		}
	}

	clients = append(clients, client)
	return SaveClients(clients)
}

// Atualiza um cliente existente
func UpdateClient(updated Client) error {
	clients, err := LoadClients()
	if err != nil {
		return err
	}

	found := false
	for i, c := range clients {
		if strings.EqualFold(c.ID, updated.ID) {
			clients[i] = updated
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("cliente '%s' não encontrado para atualização", updated.ID)
	}

	return SaveClients(clients)
}

// Remove um cliente por ID ou nome
func DeleteClient(idOrName string) error {
	clients, err := LoadClients()
	if err != nil {
		return err
	}

	var newList []Client
	removed := false
	q := strings.ToLower(strings.TrimSpace(idOrName))

	for _, c := range clients {
		if strings.ToLower(c.ID) == q || strings.ToLower(c.Name) == q {
			removed = true
			continue
		}
		newList = append(newList, c)
	}

	if !removed {
		return fmt.Errorf("nenhum cliente encontrado com o identificador '%s'", idOrName)
	}

	return SaveClients(newList)
}
