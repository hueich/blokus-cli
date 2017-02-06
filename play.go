package main

import (
	// "bytes"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hueich/blokus"
	"github.com/jroimartin/gocui"
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

func highlightString(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func renderBoard(out io.Writer, b *blokus.Board) {
	div := fmt.Sprintf(" +%s", strings.Repeat("---+", len(b.Grid())))
	fmt.Fprintln(out, div)
	for _, r := range b.Grid() {
		fmt.Fprint(out, " |")
		for _, c := range r {
			fmt.Fprintf(out, " %v |", getColorTermSymbol(c))
		}
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, div)
	}
}

func fscanln(r *bufio.Reader, a ...interface{}) error {
	r.Discard(r.Buffered())
	_, err := fmt.Fscanln(r, a...)
	return err
}

func promptForNewPlayers(g *blokus.Game) error {
	stdin := bufio.NewReader(os.Stdin)
	numPlayers := 0
	for numPlayers < 2 || numPlayers > 4 {
		fmt.Printf("How many players? [2-4]: ")
		if err := fscanln(stdin, &numPlayers); err != nil {
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
		{0, 0},
		{0, len(g.Board().Grid()[0]) - 1},
		{len(g.Board().Grid()) - 1, len(g.Board().Grid()[0]) - 1},
		{len(g.Board().Grid()) - 1, 0},
	}
	if numPlayers == 2 {
		// Make 2nd player start diagonally across from first player.
		startPositions = append(startPositions[:1], startPositions[2])
	}

	for i := 1; i <= numPlayers; i++ {
		var name string
		for true {
			fmt.Printf("Enter name of player %d: ", i)
			if err := fscanln(stdin, &name); err != nil {
				fmt.Println("Sorry, I didn't catch the name.")
				continue
			}
			name = strings.TrimSpace(name)
			if name == "" {
				fmt.Println("Sorry, the name can't be empty.")
				continue
			}
			color := blokus.Color(i)
			startPos := startPositions[i-1]
			if err := g.AddPlayer(name, color, startPos); err != nil {
				fmt.Printf("Sorry, I couldn't add the player. %v\n", err)
				continue
			}
			fmt.Printf("Player %s is color %v and will start at coordinate %v\n", highlightString(name), color, startPos)
			break
		}
	}
	return nil
}

func promptForNextMove(g *blokus.Game) error {
	stdin := bufio.NewReader(os.Stdin)
	player := g.CurrentPlayer()

	// TODO: Render available blocks

	var input string
	for true {
		fmt.Printf("It's player %s's turn. Which piece do you want to play? (Type 'pass' to pass your turn): ", highlightString(player.Name()))
		if err := fscanln(stdin, &input); err != nil {
			fmt.Println("Sorry, I didn't understand that.")
			continue
		}

		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "pass" {
			if err := g.PassTurn(player); err != nil {
				fmt.Printf("Sorry, I couldn't pass the player's turn. %v\n", err)
				continue
			}
			fmt.Printf("Passing %s's turn.\n", highlightString(player.Name()))
			break
		}

		i, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Sorry, I couldn't understand the piece number.")
			continue
		}
		if err := player.CheckPiecePlaceability(i); err != nil {
			fmt.Printf("Sorry, I can't place that piece. %v\n", err)
			continue
		}

		var o blokus.Orientation
		fmt.Print("How to orient the piece? Enter number of times to rotate 90° clockwise, followed by 'true' or 'false' for flipping horizontally. E.g. '2 false': ")
		if err := fscanln(stdin, &o.Rot, &o.Flip); err != nil {
			fmt.Println("Sorry, I couldn't understand the input.")
			continue
		}
		o.Rot = blokus.Normalize(o.Rot)
		fmt.Printf("Rotating the piece %d times and ", o.Rot)
		if !o.Flip {
			fmt.Print("not ")
		}
		fmt.Print("flipping horizontally.\n")

		var c blokus.Coord
		fmt.Print("Where would you like to place the piece? Enter coordinates as 'row column': ")
		if err := fscanln(stdin, &c.X, &c.Y); err != nil {
			fmt.Println("Sorry, I couldn't understand the coordinates.")
			continue
		}

		if err := g.PlacePiece(player, i, o, c); err != nil {
			fmt.Printf("Sorry, I couldn't place that piece. %v\n", err)
			continue
		}
		fmt.Printf("Player %s has placed piece %d.\n", highlightString(player.Name()), i)
		break
	}
	return nil
}

// func setupHandlers(g *blokus.Game) {
//     // handle key q pressing
//     termui.Handle("/sys/kbd/q", func(termui.Event) {
//         // press q to quit
//         termui.StopLoop()
//         termui.Close()
//     })
// }

// func render(g *blokus.Game) {
// 	var buf bytes.Buffer
// 	renderBoard(&buf, g.Board())
// 	b := termui.NewPar(strings.TrimSpace(buf.String()))
// 	b.Height = 2*(len(g.Board().Grid())+1)+1
// 	b.Width = 4*(len(g.Board().Grid()[0])+1)-1
// 	b.BorderLabel = "Board"

// 	termui.Render(b)
// }

type boardLayoutManager struct {
	game *blokus.Game
}

func (m *boardLayoutManager) Layout(g *gocui.Gui) error {
    if v, err := g.SetView("Board", 0, 0, 4*(len(m.game.Board().Grid()[0])+1), 2*(len(m.game.Board().Grid())+1)); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Clear()
        renderBoard(v, m.game.Board())
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func main() {
	fmt.Println("Welcome to the game!")

	g, err := blokus.NewGame(1, blokus.DefaultBoardSize, blokus.DefaultPieces())
	if err != nil {
		log.Fatalf("Could not create new game: %v\n", err)
	}

	if err := promptForNewPlayers(g); err != nil {
		log.Fatal(err.Error())
	}

    gui, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
		log.Fatalf("Could not initialize GUI: %v\n", err)
    }
    defer gui.Close()

    m := &boardLayoutManager{g}
	gui.SetManager(m)

    if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Fatalf("Could not set key binding for ctrl+c: %v\n", err)
    }
    if err := gui.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
        log.Fatalf("Could not set key binding for q: %v\n", err)
    }

    if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatalf("Did not exit normally: %v\n", err)
    }

	// for !g.IsGameEnd() {
	// 	renderBoard(g.Board())
	// 	if err := promptForNextMove(g); err != nil {
	// 		log.Fatalf("Could not process next move: %v\n", err)
	// 	}
	// 	if err := g.AdvanceTurn(); err != nil {
	// 		log.Fatalf("Could not advance turn: %v\n", err)
	// 	}
	// }

	// TODO: Display score

	fmt.Println("Someone won! Yay!")
}
