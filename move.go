package main

import "fmt"

// directions
const (
	E  = +1
	W  = -1
	N  = 8
	S  = -8
	NW = +7
	NE = +9
	SW = -NE
	SE = -NW

	frMask     = 0x00000003f                //00 0000 0000  0000 0000  0000 0000  0011 1111
	toMask     = 0x000000fd0                //00 0000 0000  0000 0000  0000 1111  1100 0000
	p12Mask    = 0x00000f000                //00 0000 0000  0000 0000  1111 0000  0000 0000
	cpMask     = 0x0000f0000                //00 0000 0000  0000 1111  0000 0000  0000 0000
	prMask     = 0x000f00000                //00 0000 0000  1111 0000  0000 0000  0000 0000
	epMask     = 0x03f000000                //00 0011 1111  0000 0000  0000 0000  0000 0000
	castlMask  = 0x3c0000000                //11 1100 0000  0000 0000  0000 0000  0000 0000
	evalMask   = uint64(0xffff000000000000) // The 16 first bits in uint64
	toShift    = 6
	p12Shift   = 12 //6+6
	cpShift    = 16 //6+6+4
	prShift    = 20 //6+6+4+4
	epShift    = 24 //6+6+4+4+4
	castlShift = 30 //6+6+4+4+4+6
	evalShift  = 64 - 16
)

var pieceRules [nP][]int // not pawns

func init() {
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], S)
	pieceRules[Rook] = append(pieceRules[Rook], W)
}

type move uint64

func (m move) String() string {
	s := m.StringFull()
	s = s[1:3] + s[5:]
	return s
}

func (m move) StringFull() string {
	fr := sq2Fen[int(m.fr())]
	to := sq2Fen[int(m.to())]
	p := int2Fen(int(m.p12()))
	cp := int2Fen(int(m.cp())) + " "
	pr := int2Fen(int(m.pr()))
	return fmt.Sprintf("%v%v-%v%v%v", p, fr, cp[:1], to, pr)
}

func (m *move) packMove(fr, to, p12, cp, pr, ep uint, castl uint) {
	// 6 bits (fr), 6 bits (2), 4 bits (piece), 4 bits (cp), 4 bits (pr), 6 bits (ep), 4 bits (cast1), x bits value
	*m = move(fr | (to << toShift) | (p12 | p12Shift) | (cp << cpShift) | (pr << prShift) | (ep << epShift) | uint(castl<<castlShift))
}

func (m *move) packEval(score int) {
	(*m) |= move(score+30000) << evalShift
}

func (m move) eval() int {
	return int((uint(m)&uint(evalMask))>>evalShift) - 30000
}

func (m move) fr() int {
	return int(m & frMask)
}

func (m move) to() int {
	return int(m&toMask) >> toShift
}

func (m move) p12() int {
	return int(m&p12Mask) >> p12Shift
}

func (m move) ep() int {
	return int(m&epMask) >> epShift
}

func (m move) cp() int {
	return int(m&cpMask) >> cpShift
}

func (m move) cast1() castlings {
	return castlings(m&castlMask) >> castlShift
}

func (m move) pr() int {
	return int(m&prMask) >> prShift
}

type moveList []move

func (mvs *moveList) add(mv move) {
	*mvs = append(*mvs, mv)
}

func (ml *moveList) remove(ix int) {
	if len(*ml) > ix && ix >= 0 {
		*ml = append((*ml)[:ix], (*ml)[ix+1:]...)
	}
}

func (ml *moveList) sort() {
	bSwap := true
	for bSwap {
		bSwap = false
		// arrange each move in descending order of score value
		for i := 0; i < len(*ml)-1; i++ {
			if (*ml)[i+1].eval() > (*ml)[i].eval() {
				(*ml)[i], (*ml)[i+1] = (*ml)[i+1], (*ml)[i]
				bSwap = true
			}
		}
	}
}

func (ml moveList) String() string {
	theString := ""
	for _, mv := range ml {
		theString += mv.String() + " "
	}
	return theString
}

var ml = moveList{}
