package main

import (
	"fmt"
	"math"
)

type searchLimits struct {
	depth    int
	nodes    uint64
	moveTime int // in milliseconds
	infinite bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999 // 3 years
	s.infinite = false
}

func (s *searchLimits) setDepth(d int) {
	s.depth = d
}

func (s *searchLimits) setMoveTime(m int) {
	s.moveTime = m
}

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
			case "go":
				tell("info string I'm thinking")
			}
		}
	}()
	return
}
