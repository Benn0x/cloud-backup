package utilities

import (
	"fmt"
	"os"
	"time"
)

var LogDirectory = "/var/backups/logs/"
var LogFile = ""
var (
	logchan = make(chan string)
	logDone = make(chan struct{})
)

func CreateLogFile() {
	time := time.Now()
	LogFile = time.Format("log-2006-01-02_15:04:05") + ".txt"

	if err := os.MkdirAll(LogDirectory, 0755); err != nil {
		panic(err)
	}

	if err := os.WriteFile(LogDirectory+LogFile, []byte(""), 0644); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(LogDirectory+LogFile, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	// go func() {
	// 	defer f.Close()
	// 	for msg := range logchan {
	// 		fmt.Println("LOGGER RECEIVED:", msg)
	// 		if _, err = f.WriteString(msg); err != nil {
	// 			fmt.Println("Failed to print", err)
	// 		}
	// 	}
	// }()
	go func() {
		defer f.Close()
		defer close(logDone)

		for msg := range logchan {
			if _, err := f.WriteString(msg); err != nil {
				fmt.Println("Failed to write:", err)
			}
		}
	}()

}

func CloseLogFile() {
	close(logchan)
	<-logDone
}

func LogError(err error) {
	logchan <- time.Now().Format("[15:04:05]") + " " + err.Error() + "\n"
}

func LogDebug(text string) {
	logchan <- time.Now().Format("[15:04:05]") + " " + text + "\n"
}
