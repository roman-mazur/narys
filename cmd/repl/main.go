package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ensure the renderer is running first!")
		panic(err)
	}
	defer conn.Close()

	fmt.Println("🎮 Connected! Type 'up', 'down', 'left', or 'right'.")
	fmt.Println("💡 Tip: Try 'right 4' to move 4 steps at once!\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Command > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		for cmdInput := range strings.SplitSeq(input, ";") {
			parts := strings.Fields(cmdInput)
			if len(parts) == 0 {
				continue
			}
			cmd := parts[0]
			count := 1

			// Parse loop if a number is provided (e.g., "down 5")
			if len(parts) > 1 {
				if n, err := strconv.Atoi(parts[1]); err == nil {
					count = n
				}
			}

			// Send commands with a delay for visual "stepping"
			for i := 0; i < count; i++ {
				fmt.Fprintf(conn, "%s\n", cmd)
				time.Sleep(500 * time.Millisecond) // Half-second delay
			}

		}
	}
}
