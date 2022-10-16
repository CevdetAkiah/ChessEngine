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
	stop     bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999 // 3 years
	s.infinite = false
}

func (s *searchLimits) setStop(st bool) {
	s.stop = st
}

func (s *searchLimits) setDepth(d int) {
	s.depth = d
}

func (s *searchLimits) setMoveTime(m int) {
	s.moveTime = m
}

func (s *searchLimits) setInfinite(b bool) {
	s.infinite = b
}

func engine() (toEng chan bool, frEng chan string) {
	tell("Hello from engine")
	toEng = make(chan bool)
	frEng = make(chan string)
	go root(toEng, frEng)
	return
}

func root(toEng chan bool, frEng chan string) {
	b := &board
	ml := moveList{}

	for range toEng {
		tell("info string engine got go! X")
		ml = moveList{}
		genAndSort(b, &ml)

		for _, mv := range ml {
			b.move(mv)
			score := search(b)
			b.unmove(mv)

			mv.packEval(adjEval(b, score))
		}
		ml.sort()
		tell("info score cp ", fmt.Sprintf("%v", ml[0].eval()), " depth 1 pv ", ml[0].String())
		frEng <- fmt.Sprintf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
	}
}

func search(b *boardStruct) int {
	return evaluate(b)
}

func genAndSort(b *boardStruct, ml *moveList) {
	b.genAllMoves(ml)

	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		v = adjEval(b, v)
		(*ml)[ix].packEval(v)
	}

	ml.sort()

}

func adjEval(b *boardStruct, ev int) int {
	if b.stm == BLACK {
		return -ev
	}
	return ev
}
