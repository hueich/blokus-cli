package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hueich/blokus"
)

func getColorSymbol(c blokus.Color) string {
	switch c {
	case blokus.Blue:
		return "\033[1;34mB\033[0m"
	case blokus.Yellow:
		return "\033[1;33mY\033[0m"
	case blokus.Red:
		return "\033[1;31mR\033[0m"
	case blokus.Green:
		return "\033[1;32mG\033[0m"
	default:
		return " "
	}
}

func renderBoard(b *blokus.Board) {
	div := fmt.Sprintf("+%s", strings.Repeat("---+", len(b.Grid())))
	fmt.Println(div)
	for _, r := range b.Grid() {
		fmt.Print("|")
		for _, c := range r {
			fmt.Printf(" %v |", getColorSymbol(c))
		}
		fmt.Println("")
		fmt.Println(div)
	}
}

func main() {
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

	if err := g.AddPlayer("Bob", blokus.Blue, blokus.Coord{0, 0}); err != nil {
		log.Fatalf("Could not add player: %v", err)
	}
	if err := g.AddPlayer("Yeti", blokus.Yellow, blokus.Coord{19, 19}); err != nil {
		log.Fatalf("Could not add player: %v", err)
	}

	renderBoard(g.Board())

	if err := g.PlacePiece(g.Players()[0], 0, blokus.Orientation{}, blokus.Coord{0, 0}); err != nil {
		log.Fatalf("Could not place piece: %v", err)
	}

	renderBoard(g.Board())

	fmt.Println("Done!")
}
