package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type AzureAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	User struct {
		Name string `json:"name"`
	} `json:"user"`
}

func TestAzure() (name, id, user string, err error) {
	out, err := exec.Command("az", "account", "show", "-o", "json").CombinedOutput()
	if err != nil {
		return "", "", "", fmt.Errorf("azure cli não autenticado: %s", string(out))
	}

	var acc AzureAccount
	if err := json.Unmarshal(out, &acc); err != nil {
		return "", "", "", err
	}

	return acc.Name, acc.ID, acc.User.Name, nil
}

func SwitchAzureSubscription(subID string) error {
	cmd := exec.Command("az", "account", "set", "--subscription", subID)
	return cmd.Run()
}
