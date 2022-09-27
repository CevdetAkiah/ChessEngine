package main

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	initFenSq2Int()
}

type boardStruct struct {
	sq        [64]int      // number of squares on the board
	wbBB      [2]bitBoard  // white black bitBoard. One for white pieces, one for black pieces
	pieceBB   [nP]bitBoard // one bitBoard for each piece type
	King      [2]int       // one int for each king position, one white one black
	ep        int          // en-passant suqare
	castlings              // castling state
	stm       colour       // side to move, white or black
	count     [12]int      // 12 counters to count how many pieces we have
	rule50    int          // set to 0 if pawn moves or capture occurrs. See 50 move rule
}

type colour int

var board = boardStruct{}

// allBB returns a bitBoard with all pieces showing
func (b *boardStruct) allBB() bitBoard {
	return b.wbBB[0] | b.wbBB[1]
}

// clear the board, flags, bitBoards etc
func (b *boardStruct) clear() {
	b.stm = WHITE
	b.rule50 = 0
	b.sq = [64]int{}
	b.King = [2]int{}
	b.ep = 0
	b.castlings = 0

	for ix := A1; ix <= H8; ix++ {
		b.sq[ix] = empty
	}

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

// make a pseudomove
func (b *boardStruct) move(fr, to, pr int) bool {
	newEp := 0
	p12 := b.sq[fr]
	switch {
	// if wK moves turn off castling priveleges for white for future moves
	case p12 == wK:
		b.castlings.off(shortW | longW)
		// castling is occurring
		if abs(to-fr) == 2 {
			if fr == E1 {
				if to == G1 {
					b.setSq(wR, F1)
					b.setSq(empty, H1)
				} else {
					b.setSq(wR, D1)
					b.setSq(empty, H1)
				}
			}
		}
	// if bK moves turn off castling priveleges for white for future moves
	case p12 == bK:
		b.castlings.off(shortB | longB)
		// castling is occurring
		if abs(to-fr) == 2 {
			if fr == F8 {
				if to == G8 {
					b.setSq(bR, F8)
					b.setSq(empty, H8)
				} else {
					b.setSq(bR, D8)
					b.setSq(empty, H8)
				}
			}
		}
	case p12 == wR:
		// if bishop moves from A1 or H1 turn off long castling or short castling respectively
		if fr == A1 {
			b.off(longW)
		}
		if fr == H1 {
			b.off(shortW)
		}
	case p12 == bR:
		// if bishop moves from A1 or H1 turn off long castling or short castling respectively
		if fr == A1 {
			b.off(longB)
		}
		if fr == H1 {
			b.off(shortB)
		}
	case p12 == wP && b.sq[to] == empty: // ep move or set ep
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep square as empty and pawn is attacking
			b.setSq(empty, to-8) // takes enemy pawn off board
		} else if to-fr == 9 { // must be ep square as empty and pawn is attacking
			b.setSq(empty, to-8) // takes enemy pawn off board
		}
	case p12 == bP && b.sq[to] == empty: // ep move or set ep
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep square as empty and pawn is attacking
			b.setSq(empty, to-8) // takes enemy pawn off board
		} else if to-fr == 9 { // must be ep square as empty and pawn is attacking
			b.setSq(empty, to-8) // takes enemy pawn off board
		}
	}
	b.ep = newEp
	// from sq is always empty after the move
	b.setSq(empty, fr)
	// if promotion is not empty set to square to the promotion piece. Else, set the to square to the moving piece
	if pr != empty {
		b.setSq(pr, to)
	} else {
		b.setSq(p12, to)
	}

	// TODO: b.isInCheck needs to be made
	// if b.isInCheck(b.stm) {
	// 	b.stm = b.stm ^ 0x1
	// 	return false
	// }
	// change side to move turn
	b.stm = b.stm ^ 0x1

	return true
}

// setSq sets a piece on a certain square on the board
func (b *boardStruct) setSq(p12, s int) {
	b.sq[s] = p12

	// if p12 is empty then we clear that bit (s) for white and black on the board
	if p12 == empty {
		b.wbBB[WHITE].clr(uint(s))
		b.wbBB[BLACK].clr(uint(s))
		// we clear the piece type on position s
		for p := 0; p < nP; p++ {
			b.pieceBB[p].clr(uint(s))
		}
		return
	}

	p := piece(p12)
	sd := p12Colour(p12)

	// king is in this square
	if p == King {
		b.King[sd] = s
	}

	// set the colour board to have a piece on position s
	b.wbBB[sd].set(uint(s))
	// set the piece board to have piece p on position s
	b.pieceBB[p].set(uint(s))
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

			board.setSq(fen2Int(char), sq)
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
		mv = trim(low(mvstr))
		if len(mv) < 4 {
			tell("info string ", mv, " in the position command is not a correct move")
			return
		}
		// is from square ok?
		fr, ok := fenSq2Int[mv[:2]]
		if !ok {
			tell("info string ", mv, " in the position command. fr_sq is not a correct move")
			return
		}
		p12 := board.sq[fr]
		if p12 == empty {
			tell("info string ", mv, " in the position command. fr_sq is not a correct move")
			return
		}
		pCol := p12Colour(p12)
		// check if piece colour matches current side to move
		if pCol != board.stm {
			tell("info string ", mv, " in the position command. fr_piece has the wrong colour")
			return
		}
		// is move to square ok?
		to, ok := fenSq2Int[mv[2:4]]
		if !ok {
			tell("info string ", mv, " in the position command. to_sq is not ok")
			return
		}

		// is the promotion piece ok?
		pr := 0
		if len(mv) == 5 { //prom character
			if !strings.ContainsAny(mv[4:5], "qrbn") {
				tell("info string promotion piece in ", mv, " in the position")
				return
			}
			pr = fen2Int(mv[4:5])
			pr = pc2P12(pr, board.stm)
		}
		board.move(fr, to, pr)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// fen2Int convert pieceString to p12 int
func fen2Int(c string) int {
	for p, x := range p12ToFen {
		if string(x) == c {
			return p
		}
	}
	return empty
}

// int2Fen converts p12 to fenString
func int2Fen(p12 int) string {
	if p12 == empty {
		return " "
	}
	return string(p12ToFen[p12])
}

// piece returns the colourless pc from p12.
// bitshifting backwards reveals the base piece, eg bQ >> 1 and wQ >> 1 both equal Queen. 1001 >> 1 & 1000 >> 1 == 100 (4 in decimal)
func piece(p12 int) int {
	return p12 >> 1
}

// p12Colour returns the colour of a p12 form. 0 for white 1 for black.
func p12Colour(p12 int) colour {
	return colour(p12 & 0x1)
}

// pcP12 returns p12 from pc and sd
func pc2P12(pc int, sd colour) int {
	return (pc << 1) | int(sd)
}

// map fen-sq to int
var fenSq2Int = make(map[string]int)

// map int-sq to fen
var sq2Fen = make(map[int]string)

// init the square map from string to int and int to string
func initFenSq2Int() {
	fenSq2Int["a1"] = A1
	fenSq2Int["a2"] = A2
	fenSq2Int["a3"] = A3
	fenSq2Int["a4"] = A4
	fenSq2Int["a5"] = A5
	fenSq2Int["a6"] = A6
	fenSq2Int["a7"] = A7
	fenSq2Int["a8"] = A8

	fenSq2Int["b1"] = B1
	fenSq2Int["b2"] = B2
	fenSq2Int["b3"] = B3
	fenSq2Int["b4"] = B4
	fenSq2Int["b5"] = B5
	fenSq2Int["b6"] = B6
	fenSq2Int["b7"] = B7
	fenSq2Int["b8"] = B8

	fenSq2Int["c1"] = C1
	fenSq2Int["c2"] = C2
	fenSq2Int["c3"] = C3
	fenSq2Int["c4"] = C4
	fenSq2Int["c5"] = C5
	fenSq2Int["c6"] = C6
	fenSq2Int["c7"] = C7
	fenSq2Int["c8"] = C8

	fenSq2Int["d1"] = D1
	fenSq2Int["d2"] = D2
	fenSq2Int["d3"] = D3
	fenSq2Int["d4"] = D4
	fenSq2Int["d5"] = D5
	fenSq2Int["d6"] = D6
	fenSq2Int["d7"] = D7
	fenSq2Int["d8"] = D8

	fenSq2Int["e1"] = E1
	fenSq2Int["e2"] = E2
	fenSq2Int["e3"] = E3
	fenSq2Int["e4"] = E4
	fenSq2Int["e5"] = E5
	fenSq2Int["e6"] = E6
	fenSq2Int["e7"] = E7
	fenSq2Int["e8"] = E8

	fenSq2Int["f1"] = F1
	fenSq2Int["f2"] = F2
	fenSq2Int["f3"] = F3
	fenSq2Int["f4"] = F4
	fenSq2Int["f5"] = F5
	fenSq2Int["f6"] = F6
	fenSq2Int["f7"] = F7
	fenSq2Int["f8"] = F8

	fenSq2Int["g1"] = G1
	fenSq2Int["g2"] = G2
	fenSq2Int["g3"] = G3
	fenSq2Int["g4"] = G4
	fenSq2Int["g5"] = G5
	fenSq2Int["g6"] = G6
	fenSq2Int["g7"] = G7
	fenSq2Int["g8"] = G8

	fenSq2Int["h1"] = H1
	fenSq2Int["h2"] = H2
	fenSq2Int["h3"] = H3
	fenSq2Int["h4"] = H4
	fenSq2Int["h5"] = H5
	fenSq2Int["h6"] = H6
	fenSq2Int["h7"] = H7
	fenSq2Int["h8"] = H8

	// -------------- sq2Fen
	sq2Fen[A1] = "a1"
	sq2Fen[A2] = "a2"
	sq2Fen[A3] = "a3"
	sq2Fen[A4] = "a4"
	sq2Fen[A5] = "a5"
	sq2Fen[A6] = "a6"
	sq2Fen[A7] = "a7"
	sq2Fen[A8] = "a8"

	sq2Fen[B1] = "b1"
	sq2Fen[B2] = "b2"
	sq2Fen[B3] = "b3"
	sq2Fen[B4] = "b4"
	sq2Fen[B5] = "b5"
	sq2Fen[B6] = "b6"
	sq2Fen[B7] = "b7"
	sq2Fen[B8] = "b8"

	sq2Fen[C1] = "c1"
	sq2Fen[C2] = "c2"
	sq2Fen[C3] = "c3"
	sq2Fen[C4] = "c4"
	sq2Fen[C5] = "c5"
	sq2Fen[C6] = "c6"
	sq2Fen[C7] = "c7"
	sq2Fen[C8] = "c8"

	sq2Fen[D1] = "d1"
	sq2Fen[D2] = "d2"
	sq2Fen[D3] = "d3"
	sq2Fen[D4] = "d4"
	sq2Fen[D5] = "d5"
	sq2Fen[D6] = "d6"
	sq2Fen[D7] = "d7"
	sq2Fen[D8] = "d8"

	sq2Fen[E1] = "e1"
	sq2Fen[E2] = "e2"
	sq2Fen[E3] = "e3"
	sq2Fen[E4] = "e4"
	sq2Fen[E5] = "e5"
	sq2Fen[E6] = "e6"
	sq2Fen[E7] = "e7"
	sq2Fen[E8] = "e8"

	sq2Fen[F1] = "f1"
	sq2Fen[F2] = "f2"
	sq2Fen[F3] = "f3"
	sq2Fen[F4] = "f4"
	sq2Fen[F5] = "f5"
	sq2Fen[F6] = "f6"
	sq2Fen[F7] = "f7"
	sq2Fen[F8] = "f8"

	sq2Fen[G1] = "g1"
	sq2Fen[G2] = "g2"
	sq2Fen[G3] = "g3"
	sq2Fen[G4] = "g4"
	sq2Fen[G5] = "g5"
	sq2Fen[G6] = "g6"
	sq2Fen[G7] = "g7"
	sq2Fen[G8] = "g8"

	sq2Fen[H1] = "h1"
	sq2Fen[H2] = "h2"
	sq2Fen[H3] = "h3"
	sq2Fen[H4] = "h4"
	sq2Fen[H5] = "h5"
	sq2Fen[H6] = "h6"
	sq2Fen[H7] = "h7"
	sq2Fen[H8] = "h8"
}

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

// piece char definitions
const (
	pc2Char  = "PNBRQK?"
	p12ToFen = "PpNnBbRrQqKk"
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
