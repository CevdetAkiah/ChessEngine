package main

import (
	"fmt"
	"math"
	"time"
)

const (
	maxDepth = 100
	maxPly   = 100
)

var cntNodes uint64

type searchLimits struct {
	depth     int
	nodes     uint64
	moveTime  int // in milliseconds
	infinite  bool
	startTime time.Time
	nextTime  time.Time

	stop bool
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

// principle variation
type pvList []move

func (pv *pvList) new() {
	*pv = make(pvList, 0, maxPly)
}

func (pv *pvList) add(mv move) {
	*pv = append(*pv, mv)
}

func (pv *pvList) clear() {
	*pv = (*pv)[:0]
}
func (pv *pvList) addPV(pv2 *pvList) {
	*pv = append(*pv, *pv2...)
}

func (pv *pvList) catenate(mv move, pv2 *pvList) {
	pv.clear()
	pv.add(mv)
	pv.addPV(pv2)
}

func (pv *pvList) String() string {
	s := ""
	for _, mv := range *pv {
		s += mv.String() + " "
	}
	return s
}

func engine() (toEng chan bool, frEng chan string) {
	tell("Hello from engine")
	toEng = make(chan bool)
	frEng = make(chan string)
	go root(toEng, frEng)
	return
}

func root(toEng chan bool, frEng chan string) {
	var pv pvList
	var childPV pvList
	childPV.new()
	b := &board
	ml := make(moveList, 0, 60)

	for range toEng {
		limits.startTime, limits.nextTime = time.Now(), time.Now()
		alpha, beta := -minEval, maxEval
		bm, bs := noMove, noScore
		depth := limits.depth
		cntNodes = 0
		ml.clear()

		genAndSort(b, &ml)

		for ix := range ml {
			mv := &ml[ix]
			childPV.clear()

			b.move(*mv)
			tell("info currmove ", mv.String())
			score := -search(-beta, -alpha, depth-1, 1, &childPV, b)
			b.unmove(*mv)
			mv.packEval(signEval(b.stm, score))
			if score > bs {
				bs = score
				pv.clear()
				pv.catenate(*mv, &childPV)

				bm = *mv
				alpha = score
				tell(fmt.Sprintf("info score cp %v depth %v nodes %v pv ", bs, depth, cntNodes), pv.String())
			}
		}
		ml.sort()
		tell(fmt.Sprintf("info score cp %v depth %v nodes %v pv ", bm.eval(), depth, cntNodes), pv.String())
		frEng <- fmt.Sprintf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
	}
}

func search(alpha, beta, depth, ply int, pv *pvList, b *boardStruct) int {
	cntNodes++
	if depth <= 0 {
		return signEval(b.stm, evaluate(b))
	}
	pv.clear()
	ml := make(moveList, 0, 60)
	//genAndSort(b, &ml)
	genInOrder(b, &ml, ply)

	bm, bs := noMove, noScore
	childPV := make(pvList, 0, maxPly)
	for _, mv := range ml {
		if !b.move(mv) {
			continue
		}

		childPV.clear()

		score := -search(-beta, -alpha, depth-1, ply+1, &childPV, b)

		b.unmove(mv)

		if score > bs {
			bs = score
			pv.catenate(mv, &childPV)

			if score >= beta { // beta cutoff
				// add killer and update history
				if mv.cp() == empty && mv.pr() == empty {
					killers.add(mv, ply)
				}
				return score
			}
			if score > alpha {
				bm = mv
				_ = bm
				alpha = score
			}

		}
		if time.Since(limits.nextTime) >= time.Duration(time.Second) {
			t1 := time.Since(limits.startTime)
			tell(fmt.Sprintf("info time %v nodes %v nps %v", int(t1.Seconds()*1000), cntNodes, cntNodes/uint64(t1.Seconds())))
			limits.nextTime = time.Now()
		}

		if limits.stop {
			return alpha
		}
	}
	return bs
}

func initQS(ml *moveList, b *boardStruct) {
	ml.clear()
	b.genAllCaptures(ml)
}

func qs(beta int, b *boardStruct) int {
	ev := signEval(b.stm, evaluate(b))
	if ev >= beta {
		// we are good. No need to try captures
		return ev
	}
	bs := ev

	qsList := make(moveList, 0, 60)
	initQS(&qsList, b)
	done := bitBoard(0)

	// move loop
	for _, mv := range qsList {
		fr := mv.fr()
		to := mv.to()

		// This works because we pick lower value pieces first
		if done.test(to) { // Don't do the same to-sq again
			continue
		}

		see := see(fr, to, b)

		if see < 0 {
			continue // we're not interested in equal captures
		}
		if see == 0 && mv.cp() == empty {
			// musts be a promotion that didn't capture and was not captured
			see = pieceVal[wQ] - pieceVal[wP]
		}

		sc := ev + see
		if sc > bs {
			bs = sc
			if sc >= beta {
				return sc
			}
		}
	}
	return bs
}

// see (Static Exchange Evaluation)
// Start with the capture fr-to and find out all the other captures to to-sq
func see(fr, to int, b *boardStruct) int {
	pc := b.sq[fr]
	cp := b.sq[to]
	cnt := 1
	us := p12Colour(pc)
	them := us.opp()

	// All the attackers to the to-sq, but first remove the moving piece and use X-ray to the to-sq
	occ := b.allBB()
	occ.clr(fr)
	attackingBB := mRookTab[to].atks(occ)&(b.pieceBB[Rook]|b.pieceBB[Queen]) |
		mBishopTab[to].atks(occ)&(b.pieceBB[Bishop]|b.pieceBB[Queen]) |
		(atksKnights[to] & b.pieceBB[Knight]) |
		(atksKings[to] & b.pieceBB[King]) |
		(b.wPawnAtksFr(to) & b.pieceBB[Pawn] & b.wbBB[BLACK]) |
		(b.bPawnAtksFr(to) & b.pieceBB[Pawn] & b.wbBB[WHITE])
	attackingBB &= occ

	if attackingBB&b.wbBB[them] == 0 { // 'they' have no attackers - good bye
		return abs(pieceVal[cp]) // always return score from 'our' point of view
	}

	// Now we continue to keep track of the material gain/loss for each capture
	// Always remove the last attacker and use x-ray to find possible new attackers

	lastAtkVal := abs(pieceVal[pc]) // save attacker piece value for later use
	var captureList [32]int
	captureList[0] = abs(pieceVal[cp])
	n := 1
	stm := them // change side to move

	for {
		cnt++

		var pt int // piece
		switch {   // select the least valuable attacker
		case (attackingBB & b.pieceBB[Pawn] & b.wbBB[stm]) != 0:
			pt = Pawn
		case (attackingBB & b.pieceBB[Knight] & b.wbBB[stm]) != 0:
			pt = Knight
		case (attackingBB & b.pieceBB[Bishop] & b.wbBB[stm]) != 0:
			pt = Bishop
		case (attackingBB & b.pieceBB[Rook] & b.wbBB[stm]) != 0:
			pt = Rook
		case (attackingBB & b.pieceBB[Queen] & b.wbBB[stm]) != 0:
			pt = Queen
		case (attackingBB & b.pieceBB[King] & b.wbBB[stm]) != 0:
			pt = King
		default:
			panic("Don't come here in see!")
		}

		// now remove the pt above from the attackingBB and scan for new attackers by possible x-ray
		BB := attackingBB & (attackingBB & b.pieceBB[pt] & b.wbBB[stm])
		occ ^= (BB & -BB) // turn off the rightmost bit from BB in occ

		// pick sliding attacks again (do it from to-sq)
		attackingBB |= mRookTab[to].atks(occ)&(b.pieceBB[Rook]|b.pieceBB[Queen]) |
			mBishopTab[to].atks(occ)&(b.pieceBB[Bishop]|b.pieceBB[Queen])
		attackingBB &= occ // but only attacking pieces

		captureList[n] = -captureList[n-1] + lastAtkVal
		n++

		// save the value of the capturing piece to be used later
		lastAtkVal = pieceVal[pc2P12(pt, WHITE)] // using WHITE always gives positive integer
		stm = stm.opp()                          // next side to move

		if pt == King && (attackingBB&b.wbBB[stm]) != 0 { // NOTE: this is just after changing side to move
			// if king capture and 'they' are attacking we have to stop
			captureList[n] = pieceVal[wK]
			n++
			break
		}
		if attackingBB&b.wbBB[stm] == 0 { // if no more attackers
			break
		}
	}
	// find the optimal capture sequence and 'our' material value will be on top
	for n--; n != 0; n-- {
		captureList[n-1] = min(-captureList[n], captureList[n-1])
	}
	return captureList[0]
}

func genAndSort(b *boardStruct, ml *moveList) {
	b.genAllMoves(ml)

	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		v = signEval(b.stm, v)
		(*ml)[ix].packEval(v)
	}

	ml.sort()

}

// generate capture moves first, then killers, then non captures
func genInOrder(b *boardStruct, ml *moveList, ply int) {
	ml.clear()
	b.genAllCaptures(ml)
	noCaptIx := len(*ml)
	b.genAllNonCaptures(ml)

	if len(*ml)-noCaptIx > 2 {
		// place killers first among non captures
		for ix := noCaptIx; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			if killers[ply].k1.cmpFrTo(mv) {
				(*ml)[ix], (*ml)[noCaptIx] = (*ml)[noCaptIx], (*ml)[ix]
			} else if killers[ply].k2.cmpFrTo(mv) {
				(*ml)[ix], (*ml)[noCaptIx+1] = (*ml)[noCaptIx+1], (*ml)[ix]
			}
		}
	}
}

func signEval(stm colour, ev int) int {
	if stm == BLACK {
		return -ev
	}
	return ev
}

/////////////////  Killers ///////////////////////////////////////////////

// killerStruct holds the killer movs per ply
type killerStruct [maxPly]struct {
	k1 move
	k2 move
}

// Clear killer moves
func (k *killerStruct) clear() {
	for ply := 0; ply < maxPly; ply++ {
		k[ply].k1 = noMove
		k[ply].k2 = noMove
	}
}

// add killer 1 and 2 (Not inCheck, captures and promotions)
func (k *killerStruct) add(mv move, ply int) {
	if !k[ply].k1.cmp(mv) {
		k[ply].k2 = k[ply].k1
		k[ply].k1 = mv
	}
}

var killers killerStruct

///////////////////////////// history table //////////////////////////////////

type historyStruct [2][64][64]uint

func (h *historyStruct) inc(fr, to int, stm colour, depth int) {
	h[stm][fr][to] += uint(depth * depth)
}

func (h *historyStruct) clear() {
	for fr := 0; fr < 64; fr++ {
		for to := 0; to < 64; to++ {
			h[0][fr][to] = 0
			h[1][fr][to] = 0
		}
	}
}

var history historyStruct
