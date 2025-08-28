package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Duck struct {
	Name           string
	Position       int
	AnimationFrame int
}

func NewDuck(name string) *Duck {
	return &Duck{
		Name:           name,
		Position:       0,
		AnimationFrame: 0,
	}
}

func (d *Duck) GetSprite() string {
	sprites := []string{"🦆"}
	return sprites[d.AnimationFrame%len(sprites)]
}

func getOrdinal(n int) string {
	suffix := "th"
	switch n % 10 {
	case 1:
		if n%100 != 11 {
			suffix = "st"
		}
	case 2:
		if n%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if n%100 != 13 {
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suffix)
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printUsage() {
	fmt.Println("Usage: ./duckrace -n <name1,name2,name3,...>")
	fmt.Println("Example: ./duckrace -n alice,bob,charlie,diana")
}

type RaceConfig struct {
	Duration int
	Names    []string
}

func parseArgs() (*RaceConfig, error) {
	var namesFlag = flag.String("n", "", "Comma-separated duck names")
	
	flag.Parse()

	if *namesFlag == "" {
		return nil, fmt.Errorf("at least one name must be specified with -n")
	}

	names := strings.Split(*namesFlag, ",")
	var cleanNames []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			cleanNames = append(cleanNames, name)
		}
	}

	if len(cleanNames) == 0 {
		return nil, fmt.Errorf("at least one valid name must be provided")
	}

	return &RaceConfig{
		Names:    cleanNames,
	}, nil
}

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}
	clearScreen()

	var ducks []*Duck
	for _, name := range config.Names {
		ducks = append(ducks, NewDuck(name))
	}

	trackWidth := 60
	finishLine := trackWidth - 15
	frameDuration := 200 * time.Millisecond

	fmt.Println("🏁 DUCK RACE STARTING! 🏁")
	fmt.Printf("Racers: %s\n", strings.Join(config.Names, ", "))
	fmt.Println("\nGet ready...")

	time.Sleep(1 * time.Second)
	fmt.Println("3...")
	time.Sleep(1 * time.Second)
	fmt.Println("2...")
	time.Sleep(1 * time.Second)
	fmt.Println("1...")
	time.Sleep(1 * time.Second)
	fmt.Println("GO! 🏁")

	time.Sleep(500 * time.Millisecond)

	rand.Seed(time.Now().UnixNano())
	var winner string
	startTime := time.Now()
	totalDuration := time.Duration(15) * time.Second

	for {
		elapsed := time.Since(startTime)
		clearScreen()

		fmt.Println("🦆 DUCK RACE IN PROGRESS 🦆")
		fmt.Println(strings.Repeat("=", trackWidth+20))

		for _, duck := range ducks {

			if rand.Float64() < 0.7 {
				speed := rand.Intn(3) + 1
				duck.Position += speed
				if duck.Position > finishLine {
					duck.Position = finishLine
				}
			}
			duck.AnimationFrame++

			if duck.Position >= finishLine && winner == "" {
				winner = duck.Name
			}
		}

		for _, duck := range ducks {
			fmt.Printf("%-10s |", duck.Name)

			for pos := 0; pos < trackWidth; pos++ {
				if pos == duck.Position {
					fmt.Print(duck.GetSprite())
				} else if pos == finishLine {
					fmt.Print("🏁")
				} else if pos < duck.Position {
					fmt.Print("~")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println("|")
		}

		fmt.Println(strings.Repeat("=", trackWidth+20))

		fmt.Println("\nCurrent Positions:")
		
		sortedDucks := make([]*Duck, len(ducks))
		copy(sortedDucks, ducks)
		sort.Slice(sortedDucks, func(i, j int) bool {
			return sortedDucks[i].Position > sortedDucks[j].Position
		})

		medals := []string{"🥇", "🥈", "🥉"}
		for rank, duck := range sortedDucks {
			medal := "  "
			if rank < len(medals) {
				medal = medals[rank]
			}
			fmt.Printf("%s %s %s\n", getOrdinal(rank+1), medal, duck.Name)
		}

		if winner != "" || elapsed >= totalDuration {
			break
		}

		time.Sleep(frameDuration)
	}

	clearScreen()
	fmt.Println("🏁 RACE FINISHED! 🏁")
	fmt.Println(strings.Repeat("=", 60))

	if winner != "" {
		fmt.Printf("🎉 WINNER: %s 🎉\n", winner)	
		fmt.Println("      👑")
		fmt.Println("      🦆")
		fmt.Println("   \\  o  /")
		fmt.Println("    \\   /")
		fmt.Println("   🏆 🏆 🏆")
	} else {

		sort.Slice(ducks, func(i, j int) bool {
			return ducks[i].Position > ducks[j].Position
		})
		
		fmt.Println("Time's up! Final positions:")

		medals := []string{"🥇", "🥈", "🥉"}
		for rank, duck := range ducks {
			medal := "  "
			if rank < len(medals) {
				medal = medals[rank]
			}
			fmt.Printf("%s %s %s\n", getOrdinal(rank+1), medal, duck.Name)
		}

		if len(ducks) > 0 {
			fmt.Printf("\n🎉 WINNER BY DISTANCE: %s 🎉\n", ducks[0].Name)
		}
	}

	fmt.Println("\nThanks for watching the Duck Race! 🦆🏁")
}
