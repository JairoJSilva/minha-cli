package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
)

// VaultSecret armazena todas as credenciais de um cliente
type VaultSecret struct {
	// AWS
	AWSAccessKeyID     string `json:"aws_access_key_id,omitempty"`
	AWSSecretAccessKey string `json:"aws_secret_access_key,omitempty"`
	AWSRegion          string `json:"aws_region,omitempty"`
	AWSSessionToken    string `json:"aws_session_token,omitempty"` // para STS
	AWSRoleARN         string `json:"aws_role_arn,omitempty"`      // para AssumeRole

	// Oracle OCI
	OCIUserOCID        string `json:"oci_user_ocid,omitempty"`
	OCITenancyOCID     string `json:"oci_tenancy_ocid,omitempty"`
	OCIFingerprint     string `json:"oci_fingerprint,omitempty"`
	OCIPrivateKeyPath  string `json:"oci_private_key_path,omitempty"`
	OCIRegion          string `json:"oci_region,omitempty"`

	// GCP
	GCPServiceAccountJSON string `json:"gcp_service_account_json,omitempty"`
	GCPProject            string `json:"gcp_project,omitempty"`

	// Azure
	AzureTenantID     string `json:"azure_tenant_id,omitempty"`
	AzureClientID     string `json:"azure_client_id,omitempty"`
	AzureClientSecret string `json:"azure_client_secret,omitempty"`

	// Kubernetes
	KubeconfigData string `json:"kubeconfig_data,omitempty"`
}

type vaultStore map[string]VaultSecret

// vaultPath retorna o caminho do arquivo vault criptografado
func vaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "minha-cli")
	_ = os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "vault.enc"), nil
}

// keyPath retorna o caminho da chave de criptografia do vault
func keyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "minha-cli")
	_ = os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "vault.key"), nil
}

// initKey cria ou lê a chave mestra do vault (32 bytes aleatórios)
func initKey() ([]byte, error) {
	kp, err := keyPath()
	if err != nil {
		return nil, err
	}

	// Se a chave já existe, lê
	if data, err := os.ReadFile(kp); err == nil && len(data) == 32 {
		return data, nil
	}

	// Gera nova chave aleatória
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("falha ao gerar chave do vault: %w", err)
	}

	if err := os.WriteFile(kp, key, 0600); err != nil {
		return nil, fmt.Errorf("falha ao salvar chave do vault: %w", err)
	}

	return key, nil
}

// deriveKey deriva uma chave AES-256 a partir da chave mestra usando scrypt
func deriveKey(masterKey []byte) ([]byte, error) {
	salt := sha256.Sum256([]byte("minha-cli-vault-v1"))
	return scrypt.Key(masterKey, salt[:], 32768, 8, 1, 32)
}

// encrypt criptografa dados com AES-256-GCM
func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt descriptografa dados com AES-256-GCM
func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("dados criptografados inválidos")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// loadStore lê e descriptografa o vault completo
func loadStore() (vaultStore, []byte, error) {
	masterKey, err := initKey()
	if err != nil {
		return nil, nil, err
	}

	aesKey, err := deriveKey(masterKey)
	if err != nil {
		return nil, nil, err
	}

	vp, err := vaultPath()
	if err != nil {
		return nil, nil, err
	}

	store := make(vaultStore)

	encData, err := os.ReadFile(vp)
	if err != nil {
		if os.IsNotExist(err) {
			// Vault ainda não existe — retorna store vazio
			return store, aesKey, nil
		}
		return nil, nil, fmt.Errorf("falha ao ler vault: %w", err)
	}

	plaintext, err := decrypt(encData, aesKey)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao descriptografar vault: %w", err)
	}

	if err := json.Unmarshal(plaintext, &store); err != nil {
		return nil, nil, fmt.Errorf("vault corrompido: %w", err)
	}

	return store, aesKey, nil
}

// saveStore criptografa e salva o vault completo
func saveStore(store vaultStore, aesKey []byte) error {
	plaintext, err := json.Marshal(store)
	if err != nil {
		return err
	}

	ciphertext, err := encrypt(plaintext, aesKey)
	if err != nil {
		return fmt.Errorf("falha ao criptografar vault: %w", err)
	}

	vp, err := vaultPath()
	if err != nil {
		return err
	}

	return os.WriteFile(vp, ciphertext, 0600)
}

// ─── API Pública ──────────────────────────────────────────────────────────────

// Store salva (ou atualiza) as credenciais de um cliente no vault
func Store(clientID string, secret VaultSecret) error {
	store, aesKey, err := loadStore()
	if err != nil {
		return err
	}
	store[clientID] = secret
	return saveStore(store, aesKey)
}

// Get retorna as credenciais de um cliente do vault
func Get(clientID string) (*VaultSecret, error) {
	store, _, err := loadStore()
	if err != nil {
		return nil, err
	}
	s, ok := store[clientID]
	if !ok {
		return nil, fmt.Errorf("nenhuma credencial salva no vault para '%s'", clientID)
	}
	return &s, nil
}

// Delete remove as credenciais de um cliente do vault
func Delete(clientID string) error {
	store, aesKey, err := loadStore()
	if err != nil {
		return err
	}
	delete(store, clientID)
	return saveStore(store, aesKey)
}

// Has verifica se um cliente tem credenciais no vault
func Has(clientID string) bool {
	store, _, err := loadStore()
	if err != nil {
		return false
	}
	_, ok := store[clientID]
	return ok
}

// List retorna todos os IDs com credenciais no vault
func List() ([]string, error) {
	store, _, err := loadStore()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(store))
	for id := range store {
		ids = append(ids, id)
	}
	return ids, nil
}

// VaultFilePath retorna o caminho do arquivo vault para exibição
func VaultFilePath() string {
	vp, _ := vaultPath()
	return vp
}
