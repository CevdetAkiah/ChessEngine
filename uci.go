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
)

func uci(frGUI chan string) {
	tell("info string Hello from uci")
	frEng, toEng := engine() // what is sent from the engine and what is sent to the engine
	quit := false            // when true command stream stops
	cmd := ""
	words := []string{}
	bm := "" // best move
	for quit == false {
		select {
		case cmd = <-frGUI:
			words = strings.Split(cmd, " ") // command received from gui
		case bm = <-frEng:
			handleBm(bm)
			continue
		}
		words[0] = trim(low(words[0]))
		switch words[0] {
		case "uci":
			handleUci()
		case "isready":
			handleIsReady()
		case "setoption":
			handleSetoption(words)
		case "stop":
			handleStop(toEng)
		case "quit", "q":
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

// handleBm handles best move provided from the engine
func handleBm(bm string) {
	tell(bm)
}

// handleBm handles best move provided from the engine
func handleStop(toEng chan string) {
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
