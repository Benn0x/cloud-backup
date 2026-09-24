package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "net/http/pprof"

	utils "github.com/benn0x/cloud-backup/utilities"
)

func StartProfiler() {
	go func() {
		utils.LogDebug("CPU profiler available at http://localhost:6060/debug/pprof/")

		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			utils.LogError(err)
		}
	}()
}

func GrabLocalFiles(Config utils.BackupConfig) ([]string, error) {
	files := []string{}

	if info, err := os.Stat(Config.LocalConfig.BackupPath); err == nil {
		if info.IsDir() {
			err := filepath.Walk(
				Config.LocalConfig.BackupPath,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if !info.IsDir() {
						files = append(files, path)
					}

					return nil
				},
			)

			if err != nil {
				return files, err
			}
		} else {
			files = append(files, Config.LocalConfig.BackupPath)
		}
	} else {
		return files, err
	}

	return files, nil
}

func CreateZipFile(Config utils.BackupConfig) (*os.File, error) {
	filename := filepath.Join(
		Config.LocalConfig.SaveLocation,
		time.Now().Format("backup-2006-01-02_15-04-05.zip"),
	)

	archive, err := os.Create(filename)
	if err != nil {
		utils.LogError(err)
	}

	return archive, err
}

func StartCopyToArchive(Config utils.BackupConfig, archive *os.File, files []string) {
	utils.LogDebug("Started backup at: " + time.Now().String())

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// 4 MB buffer.
	buf := make([]byte, 4*1024*1024)

	for _, filePath := range files {
		err := func() error {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			fileArchivePath := strings.TrimPrefix(
				filePath,
				Config.LocalConfig.BackupPath,
			)

			header := &zip.FileHeader{
				Name:   fileArchivePath,
				Method: zip.Store,
			}

			fileArchive, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			for {
				n, readErr := file.Read(buf)

				if n > 0 {
					if _, writeErr := fileArchive.Write(buf[:n]); writeErr != nil {
						return writeErr
					}
				}

				if readErr == io.EOF {
					break
				}

				if readErr != nil {
					return readErr
				}
			}

			return nil
		}()

		if err != nil {
			utils.LogError(err)
		}
	}

	utils.LogDebug("Finished backup at: " + time.Now().String())
}

func StartLocalBackup(Config utils.BackupConfig) {
	files, err := GrabLocalFiles(Config)
	if err != nil {
		utils.LogError(err)
		return
	}

	fmt.Println(len(files))
	utils.LogDebug("Grabbed " + fmt.Sprint(len(files)) + " files to backup.")

	archive, err := CreateZipFile(Config)
	if err != nil {
		return
	}
	defer archive.Close()

	StartCopyToArchive(Config, archive, files)
}
