package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type SessionEntry struct {
	Command   string  `json:"command"`
	Output    string  `json:"output"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

func RecordTerminal() {
	sessionData := []SessionEntry{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Recording terminal session. Type 'exit' to stop.")
	for {
		fmt.Print("$ ")
		scanner.Scan()
		command := scanner.Text()

		if strings.ToLower(command) == "exit" {
			break
		}

		startTime := float64(time.Now().UnixNano()) / 1e9

		cmd := exec.Command("sh", "-c", command)
		outputBytes, _ := cmd.CombinedOutput()
		output := string(outputBytes)

		endTime := float64(time.Now().UnixNano()) / 1e9

		sessionData = append(sessionData, SessionEntry{
			Command:   command,
			Output:    output,
			StartTime: startTime,
			EndTime:   endTime,
		})

		fmt.Println(output)
	}

	dataFile := "data/session.json"
	os.MkdirAll("data", os.ModePerm)
	file, _ := os.Create(dataFile)
	defer file.Close()

	jsonData, _ := json.MarshalIndent(sessionData, "", "  ")
	file.Write(jsonData)

	fmt.Printf("Session recorded in %s\n", dataFile)
}
