package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	nP12     = 12 // number of pieces total
	nP       = 6  // number of individual pieces
	WHITE    = colour(0)
	BLACK    = colour(1)
	startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

// 6 piece types - no colour (P)
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// 12 pieces with color (P12)
// every white piece is an even number and every black piece is an un even number
const (
	wP    = iota // white piece
	bP           // black piece
	wN           // white knight
	bN           // black knight
	wB           // white bishop
	bB           // black bishop
	wR           // white rook
	bR           // black rook
	wQ           // white queen
	bQ           // black queen
	wK           // white king
	bK           // black king
	empty = 15
)

// square names
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

type boardStruct struct {
	sq        [64]int      // number of squares on the board
	wbBB      [2]bitBoard  // white black bitBoard. One for white pieces, one for black pieces
	pieceBB   [nP]bitBoard // one bitBoard for each piece type
	King      [2]int       // one int for each king position, one white one black
	ep        int          // en-passant suqare
	castlings              // castling state
	stm       colour       // side to move, white or black
	count     [12]int      // 12 counters to count how many pieces we have
}

type colour int
type castlings uint

var board = boardStruct{}

// allBB returns a bitBoard with all pieces showing
func (b *boardStruct) allBB() bitBoard {
	return b.wbBB[0] | b.wbBB[1]
}

// clear the board, flags, bitBoards etc
func (b *boardStruct) clear() {
	b.stm = WHITE
	b.sq = [64]int{}
	for ix := A1; ix <= H8; ix++ {
		b.sq[ix] = empty
	}
	b.ep = 0
	b.castlings = 0

	for ix := 0; ix < nP12; ix++ {
		b.count[ix] = 0
	}

	// bitBoards
	b.wbBB[WHITE], b.wbBB[BLACK] = 0, 0
	for ix := 0; ix < nP; ix++ {
		b.pieceBB[ix] = 0
	}
}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}

func parseFEN(FEN string) {
	fenIx := 0

	for row := 7; row >= 0; row-- {
		for sq := row * 8; sq < row*8+8; { // start drawing characters on white side

			char := string(FEN[fenIx])
			fenIx++
			if char == "/" {
				continue
			}

			if i, err := strconv.Atoi(char); err == nil { // if no error then this is an empty square. Empty squares are represented by numbers.
				fmt.Println(i, "empty from sq", sq)
				sq += i
				continue
			}

			fmt.Println(char, "on sq", sq)
			sq++
		}
	}

	// take care of side to move
	// take care of castling rights
	// set the 50 move rule
	// set number of full moves
}

// parse and make the moves in position command from GUI
func parseMvs(mvstr string) {
	mvs := strings.Split(mvstr, " ")

	for _, mv := range mvs {
		fmt.Println("make move", mv)
	}
}
