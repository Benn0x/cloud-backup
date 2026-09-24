package utilities

import (
	"encoding/json"
	"errors"
	"os"
)

var ConfigDirectory = "/etc/backups/"
var ConfigFile = "config.json"

type LocalConfig struct {
	SaveLocation string `json:"save_path"`
	BackupPath   string `json:"path_location"`
}

type CloudflareConfig struct {
	CloudflareAccountID   string `json:"account_id"`
	CloudflareAccessKeyId string `json:"access_key_id"`
	CloudflareAccessKey   string `json:"access_key"`
}

type BackupConfig struct {
	Type             string           `json:"type"`
	LocalConfig      LocalConfig      `json:"local"`
	CloudflareConfig CloudflareConfig `json:"cloudflare"`
}

func LoadBackupConfiguration() BackupConfig {
	Path := ConfigDirectory
	var config BackupConfig

	if _, err := os.Stat(Path + ConfigFile); err == nil {
		// If the file does exist
		if data, err := os.ReadFile(Path + ConfigFile); err == nil {
			err = json.Unmarshal(data, &config)
			if err != nil {
				// Handle error and print to LogErrors
				LogError(err)
				panic(err)
			}
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// If the file doesn't exist
		if _, err := os.Stat(Path); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(Path, 0755)
			// Handle error and print to LogErrors
			LogError(err)
			panic(err)
		} else if err != nil {
			// Handle error and print to LogErrors
			LogError(err)
			panic(err)
		}

		DefaultConfig := BackupConfig{
			Type: "local",
			LocalConfig: LocalConfig{
				SaveLocation: "/var/backups/",
				BackupPath:   "path_to_directory",
			},
			CloudflareConfig: CloudflareConfig{
				CloudflareAccountID:   "account_id",
				CloudflareAccessKeyId: "access_key_id",
				CloudflareAccessKey:   "access_key",
			},
		}
		DefaultConfigString, err := json.MarshalIndent(DefaultConfig, "", "    ")
		if err != nil {
			// Handle the error and print to LogErrors
			LogError(err)
			panic(err)
		}

		err = os.WriteFile(ConfigDirectory+ConfigFile, []byte(DefaultConfigString), 0644)
		if err != nil {
			// Handle the error and print to LogErrors
			LogError(err)
			panic(err)
		}

		config = DefaultConfig
	} else {
		// Handle the error and print to LogErrors
		LogError(err)
		panic(err)
	}

	LogDebug("Successfully loaded configuration")
	return config
}
