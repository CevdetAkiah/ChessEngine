package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	tell = mainTell // set default tell
	trim = strings.TrimSpace
	low  = strings.ToLower

	saveBm = ""
)

func uci(frGUI chan string) {
	tell("info string Hello from uci")
	frEng, toEng := engine() // what is sent from the engine and what is sent to the engine
	bInfinite := false
	quit := false // when true command stream stops
	cmd := ""
	words := []string{}
	bm := "" // best move
	for !quit {
		select {
		case cmd = <-frGUI:
			words = strings.Split(cmd, " ") // command received from gui
		case bm = <-frEng:
			handleBm(bm, bInfinite)
			continue
		}
		words[0] = trim(low(words[0]))
		switch words[0] {
		case "uci":
			handleUci()
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleNewgame()
		case "position":
			handlePosition(words)
		case "debug":
			handleDebug(words)
		case "register":
			handleRegister(words)
		case "go":
			handleGo(words)
		case "ponderhit":
			handlePonderhit()
		case "setoption":
			handleSetoption(words)
		case "stop":
			handleStop(toEng, &bInfinite)
		case "quit", "q":
			handleQuit(toEng)
			quit = true
			continue
		}
	}
}

func handleUci() {
	tell("id name Bingo")
	tell("id author Cevdet")

	tell("option name Hash type spin default 32 min 1 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}
func handleSetoption(option []string) {
	tell("info string set option", strings.Join(option, " "))
	tell("info string not implemented yet")
}

func handleNewgame() {
	tell("info string ucinewgame not implemented")
}
func handlePosition(words []string) {
	// position [fen <fenstring> | startpos ] moves <move1> .... <movei>
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "startpos":
			tell("info string position startpos not impletmented")
		case "fen":
			tell("info string position fen not implemented")
		default:
			tell("info string position", words[1], " not implemented")
		}
	} else {
		tell("info string position not implemented")
	}
}

func handleGo(words []string) {
	// go searchmoves <move1-moveii>/ponder/wtime <ms>/ btime <ms>/winc/bi
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
			tell("info string go movetime not implemented")
		case "mate":
			tell("info string go mate not implemented")
		case "infinite":
			tell("info string go infinite not implemented")
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
func handleBm(bm string, bInfinite bool) {
	if bInfinite {
		saveBm = bm
		return
	}
	tell(bm)
}

// handleBm handles best move provided from the engine
func handleStop(toEng chan string, bInfinite *bool) {
	// if bInfinite the engine is thnking of a best move
	// if we have a saved best move the engine has done it's job, and can be told to stop
	// the gui is then told the best move
	if *bInfinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}
	}
	toEng <- "stop"
	*bInfinite = false
}

// not really necessary
func handleQuit(toEng chan string) {
	toEng <- "stop"
}

func input() chan string {
	line := make(chan string)
	go func() { // wait for input from gui and sent cmds to uci
		var reader *bufio.Reader
		reader = bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n') // reads each line of input
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 { // if an error occurs part way through input we still get whatever was typed before the error occurred
				line <- text
			}
		}
	}()
	return line
}

func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}
