package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	tell  = mainTell // set default tell
	trim  = strings.TrimSpace
	low   = strings.ToLower
	split = strings.Split

	saveBm = ""
)

func uci(input chan string) {
	tell("info string Hello from uci")
	toEng, frEng := engine() // what is sent from the engine and what is sent to the engine
	quit := false            // when true command stream stops
	cmd := ""
	bm := "" // best move
	for !quit {
		select {
		case cmd = <-input:
		case bm = <-frEng:
			handleBm(bm)
			continue
		}
		words := strings.Split(cmd, " ") // command received from gui
		words[0] = trim(low(words[0]))
		switch words[0] {
		case "uci":
			handleUci()
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleNewgame()
		case "position":
			handlePosition(cmd)
		case "debug":
			handleDebug(words)
		case "register":
			handleRegister(words)
		case "go":
			handleGo(toEng, words)
		case "ponderhit":
			handlePonderhit()
		case "setoption":
			handleSetoption(words)
		case "stop":
			handleStop(toEng)
		case "quit", "q":
			handleQuit()
			quit = true
			continue
		case "pb":
			board.Print()
		case "pbb":
			board.printAllBB()
		case "pm":
			board.printAllLegals()
		}
	}
	tell("info string leaving uci(")
}

func handleUci() {
	tell("id name GoBit")
	tell("id author Carokanns")

	tell("option name Hash type spin default 128 min 16 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}
func handleSetoption(option []string) {
	tell("info string not implemented yet")
}

func handleNewgame() {
	board.newGame()
}
func handlePosition(cmd string) {
	// position [fen <fenstring> | startpos ] moves <move1> .... <movei>
	board.newGame()
	cmd = trim(strings.TrimPrefix(cmd, "position"))
	parts := split(cmd, "moves")

	if len(cmd) == 0 || len(parts) > 2 {
		err := fmt.Errorf("%v wrong length=%v", parts, len(parts))
		tell("info string Error ", fmt.Sprint(err))
		return
	}

	alt := split(parts[0], " ")
	alt[0] = trim(alt[0])

	if alt[0] == "startpos" {
		// black position, then number of empty spaces, then white position, then turn order (w is first in this case), then castling potential, then 50 pawn move rule, then move number
		parts[0] = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	} else if alt[0] == "fen" { // if alt[0] is fen then alt[1] is the moves
		parts[0] = trim(strings.TrimPrefix(parts[0], "fen"))
	} else {
		err := fmt.Errorf("%#v must be %#v or %#v", alt[0], "fen", "startpos")
		tell("info string Error ", err.Error())
		return
	}
	// Now parts[0] is the fen-string only

	// start the parsing

	parseFEN(parts[0])

	if len(parts) == 2 {
		parts[1] = low(trim(parts[1]))
		fmt.Printf("info string parse %#v\n", parts[1])
		parseMvs(parts[1])
	}

}

// handleGo parses the go command. The go command tells us to start thinking about best moves.
func handleGo(toEng chan bool, words []string) {
	// go searchmoves <move1-moveii>/ponder/wtime <ms>/ btime <ms>/winc/bi
	limits.init()

	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "searchmoves":
			tell("info string go searchmoves not implemented")
		case "ponder":
			tell("info string go ponder not implemented")
		case "wtime":
			tell("info string go wtime not implemented")
		case "btime":
			tell("info string go btime not implemented")
		case "winc":
			tell("info string go winc not implemented")
		case "binc":
			tell("info string go binc not implemented")
		case "movestogo":
			tell("info string go movestogo not implemented")
		case "depth":
			tell("info string go depth not implemented")
		case "nodes":
			tell("info string go nodes not implemented")
		case "movetime":
			mt, err := strconv.Atoi(words[2])
			if err != nil {
				tell("info string ", words[2], " not numeric")
				return
			}
			limits.setMoveTime(mt)
			toEng <- true
		case "mate": // mate <x> mate in x moves
			tell("info string go mate not implemented")
		case "infinite":
			limits.setInfinite(true)
			toEng <- true
		case "register":
			tell("info string go register not implemented")
		default:
			tell("info string", words[1], " not implemented")
		}
	} else {
		tell("info string go not implemented")
	}
}

func handlePonderhit() {
	tell("info string ponderhit not implemented")
}

func handleDebug(words []string) {
	// debug [ on | off]
	tell("info string debug not implemented")
}

func handleRegister(words []string) {
	// register later/name <x>/code <y>
	tell("info string register not implemented")
}

// handleBm handles best move provided from the engine
func handleBm(bm string) {
	if limits.infinite {
		saveBm = bm
		return
	}
	tell(bm)
}

// handleBm handles best move provided from the engine
func handleStop(toEng chan bool) {
	// if bInfinite the engine is thnking of a best move
	// if we have a saved best move the engine has done it's job, and can be told to stop
	// the gui is then told the best move
	if limits.infinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}
	}
	limits.setStop(true)
	limits.setInfinite(false)
}

// not really necessary
func handleQuit() {
}

//------------------------------------------------------
func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	var reader *bufio.Reader
	reader = bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}
