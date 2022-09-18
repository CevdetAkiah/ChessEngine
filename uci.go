package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var tell func(text ...string)

func uci(frGUI chan string, myTell func(text ...string)) {
	tell = myTell
	tell("info string Hello from uci")

	frEng, toEng := engine() // what is sent from the engine and what is sent to the engine

	quit := false // when true command stream stops
	cmd := ""
	bm :="" // best move
	for quit == false {
		select {
		case cmd = <-frGUI: // command received from gui
		case bm = <-frEng:
			 handleBm(bm)
			 continue
		}
		switch cmd {
		case "uci":
		case "stop": 
			handleStop(toEng)
		case "quit", "q":
			quit = true
			continue
		}
	}
}

// handleBm handles best move provided from the engine
func handleBm(bm string){
	tell(bm)
}


// handleBm handles best move provided from the engine
func handleStop(toEng chan string){
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
