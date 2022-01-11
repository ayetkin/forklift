package helper

import (
	"errors"
	"fmt"
	"forklift/pkg/config"
	"github.com/hashicorp/vault/api"
)

func GetVcenterPasswordFromVault(configuration config.Configuration, env string) (string, error) {

	var (
		apiConfig = &api.Config{Address: configuration.Vault.Address}
		path      string
		err       error
	)

	client, err := api.NewClient(apiConfig)
	if err != nil {
		//log.Error(err)
		return "", err
	}

	client.SetToken(configuration.Vault.Secret)

	if env == "test" {
		path = "kv/data/vcenterTEST"
	} else if env == "prod" {
		path = "kv/data/vcenterNSXT"
	} else {
		return "", errors.New(fmt.Sprintf("Vcenter environment not valid."))
	}

	secret, err := client.Logical().Read(path)
	if err != nil {
		return "", err
	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", err
	}
	return m["password"].(string), nil
}
