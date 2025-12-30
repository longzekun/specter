package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/longzekun/specter/server/assets"
)

const (
	ConfigDirName        = "config"
	ServerConfigFilename = "server.json"
	LogFilename          = "server.log"
)

func GetServerConfigPath() string {
	appWorkDir := assets.GetRootAppDir()

	serverConfigPath := filepath.Join(appWorkDir, ConfigDirName, ServerConfigFilename)

	return serverConfigPath
}

type ServerConfig struct {
	LogFilename string `json:"log_filename"`
	TmpDir      string `json:"tmp_dir"`
	RunMode     string `json:"run_mode"`
}

func (c *ServerConfig) Save() error {
	configPath := GetServerConfigPath()
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

func GetServerConfig() *ServerConfig {
	configPath := GetServerConfigPath()
	config := getDefaultServerConfig()

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

func getDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		LogFilename: filepath.Join(assets.GetRootAppDir(), "log", LogFilename),
		TmpDir:      filepath.Join(assets.GetRootAppDir(), "tmp"),
		RunMode:     "DEBUG",
	}
}
