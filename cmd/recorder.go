package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type SessionEntry struct {
	Command   string  `json:"command"`
	Output    string  `json:"output"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

func RecordTerminal() {
	// Create a slice to store session entries
	var sessionEntries []SessionEntry

	// Handle signals like Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Interactive terminal-like prompt
	reader := bufio.NewReader(os.Stdin)
	for {
		// Show a terminal-like prompt
		fmt.Print("$ ")

		// Read user input
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
		userInput = strings.TrimSpace(userInput)
		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Exiting the terminal.")
			break
		}

		// Record the start time
		startTime := float64(time.Now().UnixNano()) / 1e9

		// Create a command
		commandArgs := strings.Split(userInput, " ")
		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

		// Capture output and error streams
		outputBuf := new(strings.Builder)
		multiWriter := io.MultiWriter(os.Stdout, outputBuf)
		cmd.Stdout = multiWriter
		cmd.Stderr = multiWriter

		// Start the command
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting command '%s': %v\n", userInput, err)
			continue
		}

		// Wait for the command or handle interruptions
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()

		select {
		case <-signalChan:
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("Failed to kill process:", err)
			}
			fmt.Println()
		case err := <-done:
			if err != nil {
				fmt.Printf("Command '%s' exited with error: %v\n", userInput, err)
			}
		}

		// Record the end time
		endTime := float64(time.Now().UnixNano()) / 1e9

		// Store the session entry
		sessionEntries = append(sessionEntries, SessionEntry{
			Command:   userInput,
			Output:    outputBuf.String(),
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	// Save session entries to a JSON file
	sessionLogFile := "data/session.json"
	file, err := os.Create(sessionLogFile)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sessionEntries); err != nil {
		fmt.Println("Error encoding session entries to JSON:", err)
		return
	}

	fmt.Println("Session log saved to", sessionLogFile)
}
