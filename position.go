package main

import (
	"fmt"
	"strconv"
	"strings"
)

// various consts
const (
	nP12     = 12
	nP       = 6
	WHITE    = colour(0)
	BLACK    = colour(1)
	startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - "
	row1     = bitBoard(0x00000000000000FF)
	row2     = bitBoard(0x000000000000FF00)
	row3     = bitBoard(0x0000000000FF0000)
	row4     = bitBoard(0x00000000FF000000)
	row5     = bitBoard(0x000000FF00000000)
	row6     = bitBoard(0x0000FF0000000000)
	row7     = bitBoard(0x00FF000000000000)
	row8     = bitBoard(0xFF00000000000000)
	fileA    = bitBoard(0x0101010101010101)
	fileB    = bitBoard(0x0202020202020202)
	fileG    = bitBoard(0x4040404040404040)
	fileH    = bitBoard(0x8080808080808080)
)

var atksKnights [64]bitBoard
var atksKings [64]bitBoard

// initialize all possible knight attacks
func initAtksKnights() {
	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		// NNE  2,1
		if rk+2 < 8 && fl+1 < 8 {
			to := (rk+2)*8 + fl + 1
			toBB.set(to)
		}

		// ENE  1,2
		if rk+1 < 8 && fl+2 < 8 {
			to := (rk+1)*8 + fl + 2
			toBB.set(to)
		}

		// ESE  -1,2
		if rk-1 >= 0 && fl+2 < 8 {
			to := (rk-1)*8 + fl + 2
			toBB.set(to)
		}

		// SSE  -2,+1
		if rk-2 >= 0 && fl+1 < 8 {
			to := (rk-2)*8 + fl + 1
			toBB.set(to)
		}

		// NNW  2,-1
		if rk+2 < 8 && fl-1 >= 0 {
			to := (rk+2)*8 + fl - 1
			toBB.set(to)
		}

		// WNW  1,-2
		if rk+1 < 8 && fl-2 >= 0 {
			to := (rk+1)*8 + fl - 2
			toBB.set(to)
		}

		// WSW  -1,-2
		if rk-1 >= 0 && fl-2 >= 0 {
			to := (rk-1)*8 + fl - 2
			toBB.set(to)
		}

		// SSW  -2,-1
		if rk-2 >= 0 && fl-1 >= 0 {
			to := (rk-2)*8 + fl - 1
			toBB.set(to)
		}
		atksKnights[fr] = toBB
	}
}

// initialize all possible King attacks
func initAtksKings() {
	fmt.Println("init atksKings")

	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		//N 1,0
		if rk+1 < 8 {
			to := (rk+1)*8 + fl
			toBB.set(to)
		}

		//NE 1,1
		if rk+1 < 8 && fl+1 < 8 {
			to := (rk+1)*8 + fl + 1
			toBB.set(to)
		}

		//E   0,1
		if fl+1 < 8 {
			to := (rk)*8 + fl + 1
			toBB.set(to)
		}

		//SE -1,1
		if rk-1 >= 0 && fl+1 < 8 {
			to := (rk-1)*8 + fl + 1
			toBB.set(to)
		}

		//S  -1,0
		if rk-1 >= 0 {
			to := (rk-1)*8 + fl
			toBB.set(to)
		}

		//SW -1,-1
		if rk-1 >= 0 && fl-1 >= 0 {
			to := (rk-1)*8 + fl - 1
			toBB.set(to)
		}

		//W   0,-1
		if fl-1 >= 0 {
			to := (rk)*8 + fl - 1
			toBB.set(to)
		}

		//NW  1,-1
		if rk+1 < 8 && fl-1 >= 0 {
			to := (rk+1)*8 + fl - 1
			toBB.set(to)
		}
		atksKings[fr] = toBB
	}
}

type boardStruct struct {
	key     uint64
	sq      [64]int
	wbBB    [2]bitBoard
	pieceBB [nP]bitBoard
	King    [2]int
	ep      int
	castlings
	stm    colour
	count  [12]int
	rule50 int //set to 0 if a pawn or capt move otherwise increment
}
type colour int

var board = boardStruct{}

// generates all legal moves
func (b *boardStruct) genAllLegals(ml *moveList) {
	b.genAllMoves(ml)
	b.filterLegals(ml)
}

func (b *boardStruct) allBB() bitBoard {
	return b.wbBB[0] | b.wbBB[1]
}

// is the move legal (except from inCheck)
func (b *boardStruct) isLegal(mv move) bool {
	fr := mv.fr()
	pc := mv.pc()
	if b.sq[fr] != pc || pc == empty {
		return false
	}
	if b.stm != p12Colour(pc) {
		return false
	}

	to := mv.to()
	cp := mv.cp()
	if !((pc == wP || pc == bP) && to == b.ep && b.ep != 0) {
		if b.sq[to] != cp {
			return false
		}
		if cp != empty && p12Colour(cp) == p12Colour(pc) {
			return false
		}
	}

	switch {
	case pc == wP:
		if to-fr == 8 { // wP one step
			if b.sq[to] == empty {
				return true
			}
		} else if to-fr == 16 {
			if b.sq[fr+8] == empty && b.sq[fr+16] == empty { // wP two step
				return true
			}
		} else if b.ep == mv.ep(b.stm) && b.sq[to-8] == bP { // wP ep
			return true
		} else if to-fr == 7 && cp != empty { // wP capture left
			return true
		} else if to-fr == 9 && cp != empty { // wp capture right
			return true
		}

		return false
	case pc == bP:
		if fr-to == 8 { // bP one step
			if b.sq[to] == empty {
				return true
			}
		} else if fr-to == 16 {
			if b.sq[fr-8] == empty && b.sq[fr-16] == empty { // bP two step
				return true
			}
		} else if b.ep == mv.ep(b.stm) && b.sq[to+8] == wP { // bP ep
			return true
		} else if fr-to == 7 && cp != empty { // bP capture right
			return true
		} else if fr-to == 9 && cp != empty { // bp capture left
			return true
		}

		return false
	case pc == wB, pc == bB:
		toBB := bitBoard(1) << uint(to)
		if mBishopTab[fr].atks(b.allBB())&toBB != 0 {
			return true
		}
		return false
	case pc == wR, pc == bR:
		toBB := bitBoard(1) << uint(to)
		if mRookTab[fr].atks(b.allBB())&toBB != 0 {
			return true
		}
		return false
	case pc == wQ, pc == bQ:
		toBB := bitBoard(1) << uint(to)
		if mBishopTab[fr].atks(b.allBB())&toBB != 0 {
			return true
		}
		if mRookTab[fr].atks(b.allBB())&toBB != 0 {
			return true
		}
		return false
	case pc == wK:
		if abs(int(to)-int(fr)) == 2 { //castlings
			if to == G1 {
				if b.sq[H1] != wR || b.sq[E1] != wK {
					return false
				}

				if b.sq[F1] != empty || b.sq[G1] != empty {
					return false
				}

				if !b.isShortOk(b.stm) {
					return false
				}
			} else {
				if b.sq[A1] != wR || b.sq[E1] != wK {
					return false
				}
				if to != C1 {
					return false
				}
				if b.sq[B1] != empty || b.sq[C1] != empty || b.sq[D1] != empty {
					return false
				}
				if !b.isLongOk(b.stm) {
					return false
				}
			}
		}
		return true
	case pc == bK:
		if abs(int(to)-int(fr)) == 2 { //castlings
			if to == G8 {
				if b.sq[H8] != bR || b.sq[E8] != bK {
					return false
				}
				if b.sq[F8] != empty || b.sq[G8] != empty {
					return false
				}
				if !b.isShortOk(b.stm) {
					return false
				}
			} else {
				if b.sq[A8] != bR || b.sq[E8] != bK {
					return false
				}
				if to != C8 {
					return false
				}
				if b.sq[B8] != empty || b.sq[C8] != empty || b.sq[D8] != empty {
					return false
				}
				if !b.isLongOk(b.stm) {
					return false
				}
			}
		}
		return true
	}

	return true
}

// clear the board, flags, bitboards etc
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

// make a move
func (b *boardStruct) move(mv move) bool {
	newEp := 0
	// we assume that the move is legally correct (except inChekc())
	fr := mv.fr()
	to := mv.to()
	pr := int(mv.pr())
	p12 := b.sq[fr]
	switch {
	case p12 == wK:
		b.castlings.off(shortW | longW)
		if abs(int(to)-int(fr)) == 2 {
			if to == G1 {
				b.setSq(wR, F1)
				b.setSq(empty, H1)
			} else {
				b.setSq(wR, D1)
				b.setSq(empty, A1)
			}
		}
	case p12 == bK:
		b.castlings.off(shortB | longB)
		if abs(int(to)-int(fr)) == 2 {
			if to == G8 {
				b.setSq(bR, F8)
				b.setSq(empty, H8)
			} else {
				b.setSq(bR, D8)
				b.setSq(empty, A8)
			}
		}
	case p12 == wR:
		if fr == A1 {
			b.off(longW)
		} else if fr == H1 {
			b.off(shortW)
		}
	case p12 == bR:
		if fr == A8 {
			b.off(longB)
		} else if fr == H8 {
			b.off(shortB)
		}

	case p12 == wP && b.sq[to] == empty: // ep move or set ep
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep
			b.setSq(empty, to-8)
		} else if to-fr == 9 { // must be ep
			b.setSq(empty, to-8)
		}
	case p12 == bP && b.sq[to] == empty: //  ep move or set ep
		if fr-to == 16 {
			newEp = to + 8
		} else if fr-to == 7 { // must be ep
			b.setSq(empty, to+8)
		} else if fr-to == 9 { // must be ep
			b.setSq(empty, to+8)
		}
	}
	b.ep = newEp
	b.setSq(empty, fr)

	if pr != empty {
		b.setSq(pr, to)
	} else {
		b.setSq(p12, to)
	}

	b.stm = b.stm ^ 0x1
	if b.isAttacked(b.King[b.stm^0x1], b.stm) {
		b.unmove(mv)
		return false
	}

	return true
}

func (c colour) String() string {
	if c == WHITE {
		return "W"
	}
	return "B"
}

func (b *boardStruct) unmove(mv move) {
	b.ep = int(mv.ep(b.stm))
	b.castlings = mv.castl()
	p12 := int(mv.p12())
	fr := int(mv.fr())
	to := int(mv.to())
	b.setSq(int(mv.cp()), to)
	b.setSq(p12, fr)

	if piece(p12) == Pawn {
		if to == b.ep { // ep move
			b.setSq(empty, to)
			switch to - fr {
			case NW, NE:
				b.setSq(bP, to-N)
			case SW, SE:
				b.setSq(wP, to-S)
			}
		}
	} else if piece(p12) == King {
		sd := p12Colour(p12)
		if fr-to == 2 { // long castling
			b.setSq(castl[sd].rook, int(castl[sd].rookL))
			b.setSq(empty, fr-1)
		} else if fr-to == -2 { // short castling
			b.setSq(castl[sd].rook, int(castl[sd].rookSh))
			b.setSq(empty, fr+1)
		}
	}
	b.stm = b.stm ^ 0x1
}

func (b *boardStruct) setSq(p12, sq int) {
	p := piece(p12)
	sd := p12Colour(p12)

	if b.sq[sq] != empty { // capture
		cp := b.sq[sq]
		b.count[cp]--
		b.wbBB[sd^0x1].clr(sq)
		b.pieceBB[piece(cp)].clr(sq)
	}
	b.sq[sq] = p12

	if p12 == empty {
		b.wbBB[WHITE].clr(sq)
		b.wbBB[BLACK].clr(sq)
		for p := 0; p < nP; p++ {
			b.pieceBB[p].clr(sq)
		}
		return
	}

	b.count[p12]++

	if p == King {
		b.King[sd] = sq
	}

	b.wbBB[sd].set(sq)
	b.pieceBB[p].set(sq)
}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}

func (b *boardStruct) genRookMoves(ml *moveList, targetBB bitBoard) {
	sd := b.stm
	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := pc2P12(Rook, colour(sd))
	var mv move
	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		toBB := mRookTab[fr].atks(b.allBB()) & targetBB
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(fr, to, p12, b.sq[to], empty, b.ep, b.castlings)
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genBishopMoves(ml *moveList, targetBB bitBoard) {
	sd := b.stm
	allBBB := b.pieceBB[Bishop] & b.wbBB[sd]
	p12 := pc2P12(Bishop, colour(sd))
	ep := b.ep
	castlings := b.castlings
	var mv move

	for fr := allBBB.firstOne(); fr != 64; fr = allBBB.firstOne() {
		toBB := mBishopTab[fr].atks(b.allBB()) & targetBB
		for to := toBB.lastOne(); to != 64; to = toBB.lastOne() {
			mv.packMove(fr, to, p12, b.sq[to], empty, ep, castlings)
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genQueenMoves(mlq *moveList, targetBB bitBoard) {
	sd := b.stm
	allQBB := b.pieceBB[Queen] & b.wbBB[sd]
	p12 := int(pc2P12(Queen, colour(sd)))
	ep := b.ep
	castlings := b.castlings
	var mv move

	for fr := allQBB.firstOne(); fr != 64; fr = allQBB.firstOne() {
		toBB := mBishopTab[fr].atks(b.allBB()) & targetBB
		toBB |= mRookTab[fr].atks(b.allBB()) & targetBB
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(fr, to, p12, b.sq[to], empty, ep, castlings)
			mlq.add(mv)
		}
	}
}

func (b *boardStruct) genKnightMoves(ml *moveList, targetBB bitBoard) {
	sd := b.stm
	allNBB := b.pieceBB[Knight] & b.wbBB[sd]
	p12 := int(pc2P12(Knight, colour(sd)))
	ep := b.ep
	castlings := b.castlings
	var mv move
	for fr := allNBB.firstOne(); fr != 64; fr = allNBB.firstOne() {
		toBB := atksKnights[fr] & targetBB
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(fr, to, p12, b.sq[to], empty, ep, castlings)
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genKingMoves(ml *moveList, targetBB bitBoard) {
	sd := b.stm
	// 'normal' moves
	p12 := int(pc2P12(King, colour(sd)))
	ep := b.ep
	castlings := b.castlings
	var mv move

	toBB := atksKings[b.King[sd]] & targetBB
	for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
		mv.packMove(b.King[sd], to, p12, b.sq[to], empty, ep, castlings)
		ml.add(mv)
	}

	// castlings
	if b.King[sd] == castl[sd].kingPos { // NOTE: Maybe not needed. We should know that the king is there if the flags are ok
		// short castling
		if targetBB.test(b.King[sd]+2) && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
			(castl[sd].betweenSh&b.allBB()) == 0 {
			if b.isShortOk(sd) {
				mv.packMove(b.King[sd], b.King[sd]+2, b.sq[b.King[sd]], empty, empty, b.ep, b.castlings)
				ml.add(mv)
			}
		}

		// long castling
		if targetBB.test(b.King[sd]-2) && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
			(castl[sd].betweenL&b.allBB()) == 0 {
			if b.isLongOk(sd) {
				mv.packMove(b.King[sd], b.King[sd]-2, b.sq[b.King[sd]], empty, empty, b.ep, b.castlings)
				ml.add(mv)
			}
		}
	}
}

// check if short castlings is legal
func (b *boardStruct) isShortOk(sd colour) bool {
	if !b.shortFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castl[sd].pawnsSh&b.pieceBB[Pawn]&b.wbBB[opp] != 0 { // stopped by pawns?
		return false
	}
	if castl[sd].pawnsSh&b.pieceBB[King]&b.wbBB[opp] != 0 { // stopped by king?
		return false
	}
	if castl[sd].knightsSh&b.pieceBB[Knight]&b.wbBB[opp] != 0 { // stopped by Knights?
		return false
	}

	// sliding to e1/e8	//NOTE: Maybe not needed during search because we know if we are in check
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// slidings to f1/f8
	if (mBishopTab[sq+1].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+1].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// slidings to g1/g8		//NOTE: Maybe not needed because we always make isAttacked() after a move
	if (mBishopTab[sq+2].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+2].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

// check if long castlings is legal
func (b *boardStruct) isLongOk(sd colour) bool {
	if !b.longFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castl[sd].pawnsL&b.pieceBB[Pawn]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].pawnsL&b.pieceBB[King]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].knightsL&b.pieceBB[Knight]&b.wbBB[opp] != 0 {
		return false
	}

	// sliding e1/e8
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding d1/d8
	if (mBishopTab[sq-1].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-1].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding c1/c8	//NOTE: Maybe not needed because we always make inCheck() before a move
	if (mBishopTab[sq-2].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-2].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

var genPawns = [2]func(*boardStruct, *moveList){(*boardStruct).genWPawnMoves, (*boardStruct).genBPawnMoves}
var genPawnCapt = [2]func(*boardStruct, *moveList){(*boardStruct).genWPawnCapt, (*boardStruct).genBPawnCapt}
var genPawnNonCapt = [2]func(*boardStruct, *moveList){(*boardStruct).genWPawnNonCapt, (*boardStruct).genBPawnNonCapt}

func (b *boardStruct) genPawnCapt(ml *moveList) {
	genPawnCapt[b.stm](b, ml)
}

func (b *boardStruct) genPawnNonCapt(ml *moveList) {
	genPawnNonCapt[b.stm](b, ml)
}

func (b *boardStruct) genPawnMoves(ml *moveList) {
	genPawns[b.stm](b, ml)
}
func (b *boardStruct) genWPawnMoves(ml *moveList) {
	var mv move
	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// one step
	to1Step := (wPawns << N) & ^b.allBB()
	// two steps,
	to2Step := ((to1Step & row3) << N) & ^b.allBB()
	// captures
	toCapL := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCapR := ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]
	// prom
	prom := (to1Step | toCapL | toCapR) & row8

	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := empty
			if b.sq[to] != empty {
				cp = b.sq[to]
				if toCapL.test(to) {
					fr := to - NW
					mv.packMove(fr, to, wP, cp, wQ, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wR, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wN, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wB, b.ep, b.castlings)
					ml.add(mv)
				}
				if toCapR.test(to) {
					fr := to - NE
					mv.packMove(fr, to, wP, cp, wQ, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wR, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wN, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, wP, cp, wB, b.ep, b.castlings)
					ml.add(mv)
				}
			} else {
				fr := to - N
				mv.packMove(fr, to, wP, cp, wQ, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wR, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wN, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wB, b.ep, b.castlings)
				ml.add(mv)
			}
		}
		to1Step &= ^row8
		toCapL &= ^row8
		toCapR &= ^row8

	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((wPawns & ^fileA) << NW) & epBB
		if epToL != 0 {
			mv.packMove(b.ep-NW, b.ep, wP, bP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
		epToR := ((wPawns & ^fileH) << NE) & epBB
		if epToR != 0 {
			mv.packMove(b.ep-NE, b.ep, wP, bP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(to-N, to, wP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}
	// Add two steps forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(to-2*N, to, wP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(to-NW, to, wP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(to-NE, to, wP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

// W pawns  captures or promotions alt 2
func (b *boardStruct) genWPawnCapt(ml *moveList) {
	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// captures
	toCapL := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCapR := ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]
	// prom
	prom := row8 & ((toCapL | toCapR) | ((wPawns << N) & ^b.allBB()))

	var mv move
	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := b.sq[to]
			frTab := make([]int, 0, 3)
			if b.sq[to] == empty {
				frTab = append(frTab, to-N) // not capture
			} else {
				if toCapL.test(to) { // capture left
					frTab = append(frTab, to-NW)
				}
				if toCapR.test(to) { // capture right
					frTab = append(frTab, to-NE)
				}
			}
			for _, fr := range frTab {
				mv.packMove(fr, to, wP, cp, wQ, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wR, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wN, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, wP, cp, wB, b.ep, b.castlings)
				ml.add(mv)
			}
		}
		toCapL &= ^row8
		toCapR &= ^row8
	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((wPawns & ^fileA) << NW) & epBB
		if epToL != 0 {
			mv.packMove(b.ep-NW, b.ep, wP, bP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
		epToR := ((wPawns & ^fileH) << NE) & epBB
		if epToR != 0 {
			mv.packMove(b.ep-NE, b.ep, wP, bP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(to-NW, to, wP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(to-NE, to, wP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

// B pawn captures or promotions alternativ 2
func (b *boardStruct) genBPawnCapt(ml *moveList) {
	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// captures
	toCapL := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCapR := ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]

	var mv move

	// prom
	prom := row1 & ((toCapL | toCapR) | ((bPawns >> (-S)) & ^b.allBB()))
	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := b.sq[to]
			frTab := make([]int, 0, 3)
			if b.sq[to] == empty {
				frTab = append(frTab, to-S) // not capture
			} else {
				if toCapL.test(to) { // capture left
					frTab = append(frTab, to-SW)
				}
				if toCapR.test(to) { // capture right
					frTab = append(frTab, to-SE)
				}
			}

			for _, fr := range frTab {
				mv.packMove(fr, to, bP, cp, bQ, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bR, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bN, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bB, b.ep, b.castlings)
				ml.add(mv)
			}
		}
		toCapL &= ^row1
		toCapR &= ^row1
	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((bPawns & ^fileA) >> (-SW)) & epBB
		if epToL != 0 {
			mv.packMove(b.ep-SW, b.ep, bP, wP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
		epToR := ((bPawns & ^fileH) >> (-SE)) & epBB
		if epToR != 0 {
			mv.packMove(b.ep-SE, b.ep, bP, wP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(to-SW, to, bP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(to-SE, to, bP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

// W pawns moves that don't capture and not promotions
func (b *boardStruct) genWPawnNonCapt(ml *moveList) {
	var mv move
	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// one step
	to1Step := (wPawns << N) & ^b.allBB()
	//two steps
	to2Step := ((wPawns & row3) << N) & ^b.allBB()
	to1Step &= ^row8

	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(to-N, to, wP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// Add two step forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(to-2*N, to, wP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

//B pawns moves that doesn't capture aand not promotions
func (b *boardStruct) genBPawnNonCapt(ml *moveList) {
	var mv move
	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// one step
	to1Step := (bPawns >> (-S)) & ^b.allBB()
	// two steps,
	to2Step := ((to1Step & row6) >> (-S)) & ^b.allBB()
	to1Step &= ^row1

	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(to-S, to, bP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}
	// Add two steps forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(to-2*S, to, bP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

func (b *boardStruct) genBPawnMoves(ml *moveList) {
	var mv move
	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// one step
	to1Step := (bPawns >> (-S)) & ^b.allBB()
	// two steps,
	to2Step := ((to1Step & row6) >> (-S)) & ^b.allBB()
	// captures
	toCapL := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCapR := ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]
	// prom
	prom := (to1Step | toCapL | toCapR) & row1
	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := empty
			if b.sq[to] != empty {
				cp = b.sq[to]
				if toCapL.test(to) {
					fr := to - SW
					mv.packMove(fr, to, bP, cp, bQ, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bR, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bN, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bB, b.ep, b.castlings)
					ml.add(mv)
				}
				if toCapR.test(to) {
					fr := to - SE
					mv.packMove(fr, to, bP, cp, bQ, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bR, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bN, b.ep, b.castlings)
					ml.add(mv)
					mv.packMove(fr, to, bP, cp, bB, b.ep, b.castlings)
					ml.add(mv)
				}
			} else {
				fr := to - S
				mv.packMove(fr, to, bP, cp, bQ, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bR, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bN, b.ep, b.castlings)
				ml.add(mv)
				mv.packMove(fr, to, bP, cp, bB, b.ep, b.castlings)
				ml.add(mv)
			}
		}
		to1Step &= ^row1
		toCapL &= ^row1
		toCapR &= ^row1
	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((bPawns & ^fileA) >> (-SW)) & epBB
		if epToL != 0 {
			mv.packMove(b.ep-SW, b.ep, bP, wP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
		epToR := ((bPawns & ^fileH) >> (-SE)) & epBB
		if epToR != 0 {
			mv.packMove(b.ep-SE, b.ep, bP, wP, empty, b.ep, b.castlings)
			ml.add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(to-S, to, bP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}
	// Add two steps forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(to-2*S, to, bP, empty, empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(to-SW, to, bP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(to-SE, to, bP, b.sq[to], empty, b.ep, b.castlings)
		ml.add(mv)
	}
}

// generates all pseudomoves
func (b *boardStruct) genAllMoves(ml *moveList) {
	b.genPawnMoves(ml)
	b.genKnightMoves(ml, ^b.wbBB[b.stm])
	b.genBishopMoves(ml, ^b.wbBB[b.stm])
	b.genRookMoves(ml, ^b.wbBB[b.stm])
	b.genQueenMoves(ml, ^b.wbBB[b.stm])
	b.genKingMoves(ml, ^b.wbBB[b.stm])
}

func (b *boardStruct) genAllCaptures(ml *moveList) {
	oppBB := b.wbBB[b.stm.opp()]
	b.genPawnCapt(ml)
	b.genKnightMoves(ml, oppBB)
	b.genBishopMoves(ml, oppBB)
	b.genRookMoves(ml, oppBB)
	b.genQueenMoves(ml, oppBB)
	b.genKingMoves(ml, oppBB)
}

func (b *boardStruct) genAllNonCaptures(ml *moveList) {
	emptyBB := ^b.allBB()
	b.genPawnNonCapt(ml)
	b.genKnightMoves(ml, emptyBB)
	b.genBishopMoves(ml, emptyBB)
	b.genRookMoves(ml, emptyBB)
	b.genQueenMoves(ml, emptyBB)
	b.genKingMoves(ml, emptyBB)
}

// generate all legal moves
func (b *boardStruct) filterLegals(ml *moveList) {
	for ix := len(*ml) - 1; ix >= 0; ix-- {
		mov := (*ml)[ix]
		if b.move(mov) {
			b.unmove(mov)
		} else {
			ml.remove(ix)
		}
	}
}

func (b *boardStruct) genFrMoves(p12 int, toBB bitBoard, ml *moveList) {

}

// is sq attacked by the sd color side
func (b *boardStruct) isAttacked(sq int, sd colour) bool {
	if isPawnAtkingSq[sd](b, sq) {
		return true
	}

	if atksKnights[sq]&b.pieceBB[Knight]&b.wbBB[sd] != 0 {
		return true
	}
	if atksKings[sq]&b.pieceBB[King]&b.wbBB[sd] != 0 {
		return true
	}
	if (mBishopTab[sq].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}
	if (mRookTab[sq].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}

	return false
}

//var pawnAtks = [2]func(*boardStruct, int) bool{(*boardStruct).wPawnAtks, (*boardStruct).bPawnAtks}

func (b *boardStruct) attacksBB(us colour) bitBoard {
	allSq := ^bitBoard(0) // all squares
	atkBB := atksKings[b.King[us]]

	atkBB |= allPawnAtksBB[us](b)

	frBB := b.pieceBB[Knight] & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= atksKnights[fr]
	}

	frBB = (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= mBishopTab[fr].atks(allSq)
	}

	frBB = (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= mRookTab[fr].atks(allSq)
	}

	return atkBB

}

var isPawnAtkingSq = [2]func(*boardStruct, int) bool{(*boardStruct).iswPawnAtkingSq, (*boardStruct).isbPawnAtkingSq}
var allPawnAtksBB = [2]func(*boardStruct) bitBoard{(*boardStruct).wPawnAtksBB, (*boardStruct).bPawnAtksBB}
var pawnAtksFr = [2]func(*boardStruct, int) bitBoard{(*boardStruct).wPawnAtksFr, (*boardStruct).bPawnAtksFr}
var pawnAtkers = [2]func(*boardStruct) bitBoard{(*boardStruct).wPawnAtkers, (*boardStruct).bPawnAtkers}

// Returns true or false if to-sq is attacked by white pawn
func (b *boardStruct) iswPawnAtkingSq(to int) bool {
	sqBB := bitBoard(1) << uint(to)

	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// Attacks left and right
	toCap := ((wPawns & ^fileA) >> NW) & b.wbBB[WHITE]
	toCap |= ((wPawns & ^fileH) >> NE) & b.wbBB[WHITE]

	return (toCap & sqBB) != 0
}

// Returns true or false if to-sq is attacked by white pawn
func (b *boardStruct) isbPawnAtkingSq(to int) bool {
	sqBB := bitBoard(1) << uint(to)

	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// Attacks left and right
	toCap := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCap |= ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]

	return (toCap & sqBB) != 0
}

// returns all w pawns that attack black pieces
func (b *boardStruct) wPawnAtkers() bitBoard {
	BB := b.wbBB[BLACK]

	// pretend that all their pieces are pawns
	// Get pawn Attacks left and right from their pieces into our pawns that now are all our pawn attackers
	ourPawnAttackers := ((BB & ^fileA) >> (-SW)) & b.wbBB[WHITE] & b.pieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^fileH) >> (-SE)) & b.wbBB[WHITE] & b.pieceBB[Pawn]

	return ourPawnAttackers
}

// returns all bl pawns that attack white pieces
func (b *boardStruct) bPawnAtkers() bitBoard {
	BB := b.wbBB[WHITE]

	ourPawnAttackers := ((BB & ^fileA) << NW) & b.wbBB[BLACK] & b.pieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^fileA) << NE) & b.wbBB[BLACK] & b.pieceBB[Pawn]

	return ourPawnAttackers
}

// returns captures from fr-sq
func (b *boardStruct) wPawnAtksFr(fr int) bitBoard {
	frBB := bitBoard(1) << uint(fr)

	//Attacks left and right
	toCap := ((frBB & ^fileA) << NW) & b.wbBB[BLACK]
	toCap |= ((frBB & ^fileH) << NE) & b.wbBB[BLACK]
	return toCap
}

// returns captures from fr-sq
func (b *boardStruct) bPawnAtksFr(fr int) bitBoard {
	frBB := bitBoard(1) << uint(fr)

	//Attacks left and right
	toCap := ((frBB & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCap |= ((frBB & ^fileH) << (-SE)) & b.wbBB[WHITE]
	return toCap
}

// returns bitBoard with all attacks, empty or note, from all white pawns
func (b *boardStruct) wPawnAtksBB() bitBoard {
	frBB := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// Attacks left and right
	toCap := ((frBB & ^fileA) << NW)
	toCap |= ((frBB & ^fileA) << NE)

	return toCap
}

// returns bitBoard with all attacks, empty or not, from all black Pawns
func (b *boardStruct) bPawnAtksBB() bitBoard {
	frBB := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// Attacks left and right
	toCap := ((frBB & ^fileA) >> (-SW))
	toCap |= ((frBB & ^fileH) >> (-SE))
	return toCap
}

//////////////////////////////////// my own commands - NOT UCI /////////////////////////////////////

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// print all legal moves
func (b *boardStruct) printAllLegals() {
	var ml moveList
	b.genAllMoves(&ml)
	fmt.Println(ml.String())
}

func (b *boardStruct) Print() {
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}

	fmt.Printf("%v to move; ep: %v  castling:%v\n", txtStm, txtEp, b.castlings.String())

	fmt.Println("  +------+------+------+------+------+------+------+------+")
	for lines := 8; lines > 0; lines-- {
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Printf("%v |", lines)
		for ix := (lines - 1) * 8; ix < lines*8; ix++ {
			if b.sq[ix] == bP {
				fmt.Printf("   o  |")
			} else {
				fmt.Printf("   %v  |", int2Fen(b.sq[ix]))
			}
		}
		fmt.Println()
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Println("  +------+------+------+------+------+------+------+------+")
	}

	fmt.Printf("       A      B      C      D      E      F      G      H\n")
}

func (b *boardStruct) printAllBB() {
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}
	fmt.Printf("%v to move; ep: %v   castling:%v\n", txtStm, txtEp, b.castlings.String())

	fmt.Println("white pieces")
	fmt.Println(b.wbBB[WHITE].Stringln())
	fmt.Println("black pieces")
	fmt.Println(b.wbBB[BLACK].Stringln())

	fmt.Println("wP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wK")
	fmt.Println((b.pieceBB[King] & b.wbBB[WHITE]).Stringln())

	fmt.Println("bP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bK")
	fmt.Println((b.pieceBB[King] & b.wbBB[BLACK]).Stringln())
}

// parse a FEN string and setup that position
func parseFEN(FEN string) {
	board.clear()
	fenIx := 0
	sq := 0
	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {

			char := string(FEN[fenIx])
			fenIx++
			if char == "/" {
				continue
			}

			if i, err := strconv.Atoi(char); err == nil { //numeriskt
				for j := 0; j < i; j++ {
					board.setSq(empty, sq)
					sq++
				}
				continue
			}

			if strings.IndexAny(p12ToFen, char) == -1 {
				tell("info string invalid piece ", char, " try next one")
				continue
			}

			board.setSq(fen2Int(char), sq)

			sq++
		}
	}

	remaining := strings.Split(trim(FEN[fenIx:]), " ")

	// stm
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			board.stm = WHITE
		} else if remaining[0] == "b" {
			board.stm = BLACK
		} else {
			r := fmt.Sprintf("%v; sq=%v;  fenIx=%v", strings.Join(remaining, " "), sq, fenIx)

			tell("info string remaining=", r, ";")
			tell("info string ", remaining[0], " invalid stm color")
			board.stm = WHITE
		}
	}

	// castling
	board.castlings = 0
	if len(remaining) > 1 {
		board.castlings = parseCastlings(remaining[1])
	}

	// ep square
	board.ep = 0
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			board.ep = fenSq2Int[remaining[2]]
		}
	}

	// 50-move
	board.rule50 = 0
	if len(remaining) > 3 {
		board.rule50 = parse50(remaining[3])
	}
}

// parse 50 move rue in fenstring
func parse50(fen50 string) int {
	r50, err := strconv.Atoi(fen50)
	if err != nil || r50 < 0 {
		tell("info string 50 move rule in fenstring ", fen50, " is not a valid number >= 0 ")
		return 0
	}
	return r50
}

// parse and make the moves in position command from GUI
func parseMvs(mvstr string) {

	mvs := strings.Fields(low(mvstr))

	for _, mv := range mvs {
		mv = trim(mv)
		if len(mv) < 4 || len(mv) > 5 {
			tell("info string ", mv, " in the position command is not a correct move")
			return
		}
		// is fr square ok?
		fr, ok := fenSq2Int[mv[:2]]
		if !ok {
			tell("info string ", mv, " in the position command is not a correct fr square")
			return
		}

		p12 := board.sq[fr]
		if p12 == empty {
			tell("info string ", mv, " in the position command. fr_sq is an empty square")
			return
		}
		pCol := p12Colour(p12)
		if pCol != board.stm {
			tell("info string ", mv, " in the position command. fr piece has the wrong color")
			return
		}

		// is to square ok?
		to, ok := fenSq2Int[mv[2:4]]
		if !ok {
			tell("info string ", mv, " in the position has an incorrect to square")
			return
		}

		// is the prom piece ok?
		pr := empty
		if len(mv) == 5 { //prom
			if !strings.ContainsAny(mv[4:5], "QRNBqrnb") {
				tell("info string promotion piece in ", mv, " in the position command is not correct")
				return
			}
			pr = fen2Int(mv[4:5])
			pr = pc2P12(pr, board.stm)
		}
		cp := board.sq[to]

		var intMv move // internal move format
		intMv.packMove(fr, to, p12, cp, pr, board.ep, board.castlings)

		if !board.move(intMv) {
			tell(fmt.Sprintf("tell info string %v-%v is an illegal move", sq2Fen[fr], sq2Fen[to]))
		}
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

// int2fen convert p12 to fenString
func int2Fen(p12 int) string {
	if p12 == empty {
		return " "
	}
	return string(p12ToFen[p12])
}

// piece returns the pc from p12
func piece(p12 int) int {
	return p12 >> 1
}

// p12Colour returns the color of a p12 form
func p12Colour(p12 int) colour {
	return colour(p12 & 0x1)
}

func (c colour) opp() colour {
	return c ^ 0x1
}

// pc2P12 returns p12 from pc and sd
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

// 6 piece types - no color (P)
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// 12 pieces with color (P12)
const (
	wP = iota
	bP
	wN
	bN
	wB
	bB
	wR
	bR
	wQ
	bQ
	wK
	bK
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

//////////////////////////////// TODO: remove this after benchmarking ////////////////////////////////////////
func (b *boardStruct) genSimpleRookMoves(ml *moveList, sd colour) {
	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := int(pc2P12(Rook, colour(sd)))
	ep := b.ep
	castlings := b.castlings
	var mv move
	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		rk := fr / 8
		fl := fr % 8
		//N
		for r := rk + 1; r < 8; r++ {
			to := r*8 + fl
			cp := b.sq[to]
			if cp != empty && p12Colour(int(cp)) == sd {
				break
			}
			mv.packMove(fr, to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//S
		for r := rk - 1; r >= 0; r-- {
			to := r*8 + fl
			cp := b.sq[to]
			if cp != empty && p12Colour(int(cp)) == sd {
				break
			}
			mv.packMove(fr, to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//E
		for f := fl + 1; f < 8; f++ {
			to := rk*8 + f
			cp := b.sq[to]
			if cp != empty && p12Colour(int(cp)) == sd {
				break
			}
			mv.packMove(fr, to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//W
		for f := fl - 1; f >= 0; f-- {
			to := rk*8 + f
			cp := b.sq[to]
			if cp != empty && p12Colour(int(cp)) == sd {
				break
			}
			mv.packMove(fr, to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
	}
}
