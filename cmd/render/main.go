package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	Width  = 10
	Height = 10
)

type Point struct{ X, Y int }

var (
	robot  = Point{0, 0}
	target = Point{9, 9}
	walls  = []Point{{3, 0}, {3, 1}, {3, 2}, {7, 7}, {7, 8}, {8, 7}}
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	drawMap()
	fmt.Println("Waiting for the REPL to connect...")

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := strings.ToLower(strings.TrimSpace(scanner.Text()))
		moveRobot(cmd)
		drawMap()

		if robot == target {
			fmt.Println("🎉 YOU WIN! The robot found the star! 🎉")
			os.Exit(0)
		}
	}
}

func moveRobot(cmd string) {
	next := robot
	switch cmd {
	case "up":
		next.Y--
	case "down":
		next.Y++
	case "left":
		next.X--
	case "right":
		next.X++
	}

	// Boundary and wall checks
	if next.X >= 0 && next.X < Width && next.Y >= 0 && next.Y < Height && !isWall(next) {
		robot = next
	}
}

func isWall(p Point) bool {
	for _, w := range walls {
		if p == w {
			return true
		}
	}
	return false
}

func drawMap() {
	// ANSI escape codes to clear the terminal screen
	fmt.Print("\033[H\033[2J")
	fmt.Println("🤖 Robot Command Center 🤖\n")

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			p := Point{x, y}
			if p == robot {
				fmt.Print("🤖 ")
			} else if p == target {
				fmt.Print("⭐ ")
			} else if isWall(p) {
				fmt.Print("🧱 ")
			} else {
				fmt.Print(".  ")
			}
		}
		fmt.Println()
	}
	fmt.Println("\nWaiting for commands...")
}
