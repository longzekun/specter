package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/longzekun/specter/client/assets"
)

const (
	AuthConfigDirName        = "auth"
	ClientAuthConfigFilename = "auth.json"
)

func GetClientAuthConfigPath() string {
	appWorkDir := assets.GetRootAppDir()

	serverConfigPath := filepath.Join(appWorkDir, AuthConfigDirName, ClientAuthConfigFilename)

	return serverConfigPath
}

type ClientAuthConfig struct {
	Operator      string `json:"operator"`
	Token         string `json:"token"`
	Lhost         string `json:"lhost"`
	Lport         int    `json:"lport"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	PublicKey     string `json:"public_key"`
}

func (c *ClientAuthConfig) Save() error {
	configPath := GetClientAuthConfigPath()
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetClientAuthConfig() *ClientAuthConfig {
	configPath := GetClientAuthConfigPath()
	config := getDefaultSClientAuthConfig()

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return config
		}

		err = json.Unmarshal(data, config)
		if err != nil {
			return config
		}
	}
	err := config.Save()
	if err != nil {
		panic(err)
	}
	return config
}

func getDefaultSClientAuthConfig() *ClientAuthConfig {
	return &ClientAuthConfig{
		Operator:      "",
		Token:         "",
		Lhost:         "",
		Lport:         0,
		CACertificate: "",
		PrivateKey:    "",
		PublicKey:     "",
	}
}
