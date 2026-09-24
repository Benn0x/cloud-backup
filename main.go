package main

import (
	"fmt"
	"os"

	backup "github.com/benn0x/cloud-backup/backup"
	utils "github.com/benn0x/cloud-backup/utilities"
)

func main() {
	if os.Geteuid() == 0 {
		utils.CreateLogFile()
		config := utils.LoadBackupConfiguration()

		backup.StartProfiler()
		backup.StartBackup(config)

		utils.CloseLogFile()
	} else {
		fmt.Println("Please run this programme as sudo.")
	}
}
