package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hueich/blokus"
)

func getColorAsciiCode(c blokus.Color) int {
	switch c {
	case blokus.Blue:
		return 34
	case blokus.Yellow:
		return 33
	case blokus.Red:
		return 31
	case blokus.Green:
		return 32
	default:
		return 0
	}
}

func getColorTermSymbol(c blokus.Color) string {
	if !c.IsColored() {
		return " "
	}
	return fmt.Sprintf("\033[1;%dm%c\033[0m", getColorAsciiCode(c), strings.ToUpper(c.String())[0])
}

func renderBoard(b *blokus.Board) {
	div := fmt.Sprintf("+%s", strings.Repeat("---+", len(b.Grid())))
	fmt.Println(div)
	for _, r := range b.Grid() {
		fmt.Print("|")
		for _, c := range r {
			fmt.Printf(" %v |", getColorTermSymbol(c))
		}
		fmt.Println("")
		fmt.Println(div)
	}
}

func promptForNewPlayers(g *blokus.Game) error {
	numPlayers := 0
	for numPlayers < 2 || numPlayers > 4 {
		fmt.Printf("How many players? [2-4]: ")
		if _, err := fmt.Scanln(&numPlayers); err != nil {
			fmt.Println("Sorry, I don't know what that number is.")
			continue
		}
		if numPlayers < 2 || numPlayers > 4 {
			fmt.Println("Sorry, this game can only have 2 to 4 players.")
			continue
		}
	}
	fmt.Printf("Setting up a %d player game.\n", numPlayers)

	// Counter-clockwise order from top left.
	startPositions := []blokus.Coord{
		{0,0},
		{0,len(g.Board().Grid()[0])-1},
		{len(g.Board().Grid())-1,len(g.Board().Grid()[0])-1},
		{len(g.Board().Grid())-1,0},
	}
	fmt.Println(startPositions)
	if numPlayers == 2 {
		// Make 2nd player start diagonally across from first player.
		startPositions = append(startPositions[:1], startPositions[2])
	}

	for i := 1; i <= numPlayers; i++ {
		var name string
		for name == "" {
			fmt.Printf("Enter name of player %d: ", i)
			if _, err := fmt.Scanln(&name); err != nil {
				fmt.Println("Sorry, I didn't catch the name.")
				continue
			}
			name = strings.TrimSpace(name)
			if name == "" {
				fmt.Println("Sorry, the name can't be empty.")
				continue
			}
		}
		color := blokus.Color(i)
		startPos := startPositions[i-1]
		if err := g.AddPlayer(name, color, startPos); err != nil {
			return err
		}
		fmt.Printf("Player \033[1m%s\033[0m is color %v and will start at coordinate %v\n", name, color, startPos)
	}
	return nil
}

func main() {
	fmt.Println("Welcome to the game!")

	pieces := make([]*blokus.Piece, 0)
	var p *blokus.Piece
	p, err := blokus.NewPiece([]blokus.Coord{blokus.Coord{0, 0}})
	if err != nil {
		log.Fatalf("Could not create piece: %v", err)
	}
	pieces = append(pieces, p)

	g, err := blokus.NewGame(1, blokus.DefaultBoardSize, pieces)
	if err != nil {
		log.Fatalf("Could not create new game: %v", err)
	}

	if err := promptForNewPlayers(g); err != nil {
		log.Fatal(err.Error())
	}

	renderBoard(g.Board())

	if err := g.PlacePiece(g.Players()[0], 0, blokus.Orientation{}, blokus.Coord{0, 0}); err != nil {
		log.Fatalf("Could not place piece: %v", err)
	}

	renderBoard(g.Board())

	fmt.Println("Done!")
}
