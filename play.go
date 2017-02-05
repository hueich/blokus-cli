package main

import (
	"log"

	"github.com/hueich/blokus"
)

func main() {
	pieces := make([]*blokus.Piece, 0)
	var p *blokus.Piece
	p, err := blokus.NewPiece([]blokus.Coord{blokus.Coord{0, 0}})
	if err != nil {
		log.Fatalf("Could not create piece: %v", err)
	}
	pieces = append(pieces, p)

	_, err = blokus.NewGame(1, blokus.DefaultBoardSize, pieces)
	if err != nil {
		log.Fatalf("Could not create new game: %v", err)
	}
	log.Println("Done!")
}
