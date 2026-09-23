package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	Width    = 10
	Height   = 10
	NumWalls = 15 // Adjust this to make the game harder or easier

	// ANSI Color Codes
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"

	// Terminal Bell for sound
	Bell = "\a"
)

type Point struct{ X, Y int }

var (
	robot       Point
	target      Point
	trampolines [2]Point // Linked pair: reach one, pop out of the other
	walls       []Point
	lastAction  string
)

func main() {
	generateMap()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	drawMap()
	fmt.Println(Cyan + "Waiting for the REPL to connect..." + Reset)

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
			fmt.Print(Bell, Bell, Bell)
			fmt.Println(Yellow + "\n🎉 YAY! YOU WIN! The robot found the star! 🎉" + Reset)
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}
	}
}

func generateMap() {
	// Place robot
	robot = Point{X: rand.IntN(Width), Y: rand.IntN(Height)}

	// Place target, ensuring it isn't on the robot
	for {
		target = Point{X: rand.IntN(Width), Y: rand.IntN(Height)}
		if target != robot {
			break
		}
	}

	// Place trampolines, ensuring they don't cover the robot, target, or each other
	for i := range trampolines {
		for {
			p := Point{X: rand.IntN(Width), Y: rand.IntN(Height)}
			if p != robot && p != target && !slices.Contains(trampolines[:i], p) {
				trampolines[i] = p
				break
			}
		}
	}

	// Place walls, ensuring they don't cover the robot, target, trampolines, or each other
	walls = make([]Point, 0, NumWalls)
	for len(walls) < NumWalls {
		p := Point{X: rand.IntN(Width), Y: rand.IntN(Height)}
		if p != robot && p != target && !isTrampoline(p) && !isWall(p) {
			walls = append(walls, p)
		}
	}
}

func moveRobot(cmd string) {
	next := robot
	switch cmd {
	case "up", "^":
		next.Y--
	case "down", "v":
		next.Y++
	case "left", "<":
		next.X--
	case "right", ">":
		next.X++
	default:
		lastAction = ""
		return
	}

	// Boundary check
	if !inBounds(next) {
		lastAction = Red + "Ouch! You hit the edge of the world!" + Reset + Bell
		return
	}

	// Wall check
	if isWall(next) {
		lastAction = Red + "BONK! You hit a brick wall!" + Reset + Bell
		return
	}

	// Success!
	robot = next
	lastAction = Green + "Beep boop! Moved " + cmd + "!" + Reset

	// Trampoline: jump into one, pop out of the other
	if i := slices.Index(trampolines[:], robot); i >= 0 {
		robot = trampolines[1-i]
		lastAction = Purple + "BOING! The robot jumped through the trampoline!" + Reset + Bell
	}
}

func inBounds(p Point) bool {
	return p.X >= 0 && p.X < Width && p.Y >= 0 && p.Y < Height
}

func isTrampoline(p Point) bool {
	return slices.Contains(trampolines[:], p)
}

func isWall(p Point) bool {
	return slices.Contains(walls, p)
}

func drawMap() {
	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Colorful Title
	fmt.Println(Cyan + "🤖 " + Purple + "Robot Command Center" + Cyan + " 🤖\n" + Reset)

	for y := range Height {
		for x := range Width {
			p := Point{x, y}
			if p == robot {
				fmt.Print("🤖 ")
			} else if p == target {
				fmt.Print("⭐ ")
			} else if isTrampoline(p) {
				fmt.Print("⭕ ")
			} else if isWall(p) {
				fmt.Print("🧱 ")
			} else {
				fmt.Print(Cyan + ".  " + Reset)
			}
		}
		fmt.Println()
	}

	fmt.Println("\n" + lastAction)
	fmt.Println(Yellow + "Waiting for commands..." + Reset)
}
