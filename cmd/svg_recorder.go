package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GenerateSVG() {
	dataFile := "data/session.json"
	outputFile := "output/animation.svg"

	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Error reading session file:", err)
		return
	}
	defer file.Close()

	var sessionData []SessionEntry
	json.NewDecoder(file).Decode(&sessionData)

	os.MkdirAll("output", os.ModePerm)
	svgFile, _ := os.Create(outputFile)
	defer svgFile.Close()

	var svgContent strings.Builder

	terminalWidth := 760 // Adjust as needed, accounting for padding
	xStart := 20
	lineHeight := 18
	fontSize := 14
	charWidth := 8 // Approximate width of a monospace character

	// Helper function to wrap text
	wrapText := func(text string) []string {
		var wrappedLines []string
		words := strings.Split(text, " ")
		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word)+1 <= terminalWidth/charWidth {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			} else {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = word
			}
		}
		wrappedLines = append(wrappedLines, currentLine)
		return wrappedLines
	}

	// Calculate content height
	y := 50
	for _, entry := range sessionData {
		commandLines := 1
		outputLines := 0
		for _, line := range strings.Split(entry.Output, "\n") {
			outputLines += len(wrapText(line))
		}
		y += (commandLines + outputLines) * lineHeight
	}

	svgHeight := y

	svgContent.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="800" height="%d">`, svgHeight))
	svgContent.WriteString(`<rect width="100%" height="100%" fill="#303030" rx="10" ry="10"/>`)
	svgContent.WriteString(`<circle cx="20" cy="20" r="6" fill="#ff5f56"/>`)
	svgContent.WriteString(`<circle cx="40" cy="20" r="6" fill="#ffbd2e"/>`)
	svgContent.WriteString(`<circle cx="60" cy="20" r="6" fill="#27c93f"/>`)
	svgContent.WriteString(`<style>
		.char {
			opacity: 0;
			animation-fill-mode: forwards;
			animation-timing-function: steps(1);
		}
		.cursor {
			width: 2px;
			fill: white;
			animation: blink 1s steps(1) infinite;
		}
		.output {
			opacity: 0;
			animation: fadein 0.5s forwards;
		}
		@keyframes fadein {
			to { opacity: 1; }
		}
		@keyframes blink {
			50% { opacity: 0; }
		}
	</style>`)

	// Positioning below terminal buttons
	x, y := xStart, 55 // Start after buttons
	charDelay := 0.05  // Delay per character in seconds
	lineDelay := 0.5   // Delay before output after typing a command

	totalDelay := 0.0

	for _, entry := range sessionData {
		command := "$ " + entry.Command
		// Typing animation for the command
		currentCommandX := xStart
		for i, char := range command {
			delay := totalDelay + float64(i)*charDelay

			// Character appears with animation
			svgContent.WriteString(fmt.Sprintf(
				`<text x="%d" y="%d" fill="white" font-family="monospace" font-size="%d" class="char" style="animation-delay:%.2fs; animation-name: fadein%d;">%s</text>`,
				currentCommandX, y, fontSize, delay, i, escapeSVG(string(char)),
			))
			svgContent.WriteString(fmt.Sprintf(`<style>
				@keyframes fadein%d {
					to { opacity: 1; }
				}
			</style>`, i))

			currentCommandX += charWidth // Move to the next character's position
		}

		// Update total delay after typing the command
		totalDelay += float64(len(command)) * charDelay

		// Reset x and move y for the next line
		y += lineHeight

		// Output text appears after the command
		for _, line := range strings.Split(entry.Output, "\n") {
			x = xStart // Reset x position for each line of output
			for _, wrappedLine := range wrapText(line) {
				svgContent.WriteString(fmt.Sprintf(
					`<text x="%d" y="%d" fill="white" font-family="monospace" font-size="%d" class="output" style="animation-delay:%.2fs;">%s</text>`,
					x, y, fontSize, totalDelay, escapeSVG(wrappedLine),
				))
				y += lineHeight
			}
		}

		// Update total delay after output
		totalDelay += lineDelay
	}

	svgContent.WriteString(`</svg>`)
	svgFile.WriteString(svgContent.String())

	fmt.Printf("SVG animation saved to %s\n", outputFile)
}

func escapeSVG(text string) string {
	var escaped strings.Builder
	for _, r := range text {
		if r >= ' ' && r != '\b' { // Keep only printable characters and not backspace
			switch r {
			case '&':
				escaped.WriteString("&amp;")
			case '<':
				escaped.WriteString("<")
			case '>':
				escaped.WriteString(">")
			default:
				escaped.WriteRune(r)
			}
		}
	}
	return escaped.String()
}
