package main

import "fmt"

func engine() (frEng, toEng chan string) {
	tell("Hello from engine")
	frEng = make(chan string)
	toEng = make(chan string)
	go func() {
		for cmd := range toEng {
			switch cmd {
			case "stop":
				fmt.Println("stop from engine")
			case "quit":
				fmt.Println("quit from engine")
			}
		}
	}()
	return
}
