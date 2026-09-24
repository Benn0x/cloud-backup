package backup

import (
	utils "github.com/benn0x/cloud-backup/utilities"
)

func StartBackup(Config utils.BackupConfig) {
	if Config.Type == "local" {
		StartLocalBackup(Config)
	}
}
