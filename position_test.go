package main

import (
	"log"
	"testing"
)

func Test_int2Fen(t *testing.T) {

	tests := []struct {
		name string
		p12  int
		want string
	}{
		{"", wP, "P"},
		{"", bP, "p"},
		{"", wK, "K"},
		{"", bK, "k"},
		{"", wN, "N"},
		{"", bN, "n"},
		{"", empty, " "},
		{"", wQ, "Q"},
		{"", bQ, "q"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := int2Fen(tt.p12); got != tt.want {
				t.Errorf("int2Fen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boardStruct_allBB(t *testing.T) {
	tests := []struct {
		name string
		wBB  bitBoard
		bBB  bitBoard
	}{
		{"1", 0xff, 0x00FF},
		{"2", 0x0, 0xAF00},
		{"3", 0xF, 0x00FF},
		{"4", 0xFFFF000000000000, 0xFFFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board.wbBB[WHITE] = tt.wBB
			board.wbBB[BLACK] = tt.bBB
			correct := tt.wBB | tt.bBB
			if got := board.allBB(); got != correct {
				t.Errorf("%v: should be %v but we got %v", tt.name, correct, got)
			}
		})
	}
}

func Test_move(t *testing.T) {
	tests := []struct {
		name      string
		pos       string
		want      []int // pairs of sq + content
		wantEp    int
		wantCastl uint
	}{
		{"", "position startpos moves a2a4", []int{A2, empty, A4, wP}, A3, shortW | longW | shortB | longB},
		{"", "position startpos moves b1c3 g8f6 a1b1", []int{A1, empty, B1, wR}, 0, shortW | shortB | longB},
		{"", "position startpos moves b1c3 g8f6 a1b1 h8g8", []int{A1, empty, B1, wR}, 0, shortW | longB},
		{"", "position startpos moves b1c3 g8f6 g1f3 b8c6 a1b1 a8b8 h1g1 h8g8", []int{A1, empty, B1, wR}, 0, 0},
		{"", "position startpos moves e2e4 e7e5 e1e2 e8e7", []int{}, 0, 0},
		{"", "position startpos moves d2d4 d7d5 b1c3 b8c6 c1f4 c8f5 e1c1 e8c8", []int{A1, empty, B1, empty, C1, wK, D1, wR, E1, empty, A8, empty, B8, empty, C8, bK, D8, bR, E8, empty},
			0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			for ix := 0; ix < len(tt.want); ix += 2 {
				sq := tt.want[ix]
				p12 := tt.want[ix+1]
				if board.sq[sq] != p12 {
					t.Errorf("%v:  want %v on sq=%v got %v", tt.name, p12, sq, board.sq[sq])
				}
			}
			if board.ep != tt.wantEp {
				t.Errorf("%v:  want ep=%v got %v", tt.name, tt.wantEp, board.ep)
			}
			if board.castlings != castlings(tt.wantCastl) {
				t.Errorf("%v:  want castl=%v got %v", tt.name, tt.wantCastl, board.castlings)
			}
		})
	}
}

func Test_unmove(t *testing.T) {
	type movStr struct{ fr, to, p12, cp, pr int }
	tests := []struct {
		name   string
		pos    string
		mov    movStr // the move to unmove
		wantTo int    // p12 on to-sq
		wantSq [2]int // sqNo and content in sq[sqNo]
	}{
		{"", "position startpos e2e4 e7e5 g1f3 g8f6 f1c4 f8c5 e1g1", movStr{E1, G1, wK, empty, empty}, empty, [2]int{F1, empty}},
		{"", "position startpos e2e4 e7e5 g1f3 g8f6 f1c4 f8c5", movStr{F8, C5, bB, empty, empty}, empty, [2]int{-1, empty}},
	}
	for _, tt := range tests {
		handlePosition(tt.pos)
		var mv move
		mv.packMove(tt.mov.fr, tt.mov.to, tt.mov.p12, tt.mov.cp, tt.mov.pr, board.ep, board.castlings)
		board.unmove(mv)
		if int(tt.mov.p12) != board.sq[tt.mov.fr] {
			t.Errorf("%v: fr_sq should have %v (%v) but have %v", tt.name, int2Fen(int(tt.mov.p12)), tt.mov.p12, board.sq[tt.mov.fr])
		}
		if tt.wantTo != board.sq[tt.mov.to] {
			t.Errorf("%v: to_sq should have %v (%v) but have %v", tt.name, int2Fen(tt.wantTo), tt.wantTo, board.sq[tt.mov.to])
		}
		if tt.wantSq[0] >= 0 && tt.wantSq[1] != board.sq[tt.wantSq[0]] {
			t.Errorf("%v: sq=%v should have %v (%v) but have %v", tt.name, tt.wantSq[0], int2Fen(tt.wantSq[1]), tt.wantSq[1], board.sq[tt.wantSq[0]])
		}
	}
}

func Test_isAttacked(t *testing.T) {
	tests := []struct {
		name string
		pos  string
		sq   int
		sd   colour
		want bool
	}{
		{"", "position startpos", E2, WHITE, true},
		{"", "position startpos", E7, BLACK, true},
		{"", "position startpos moves e2e4 d7d5", D5, WHITE, true},
		{"", "position startpos moves e2e4 d7d5", E4, BLACK, true},
		{"", "position startpos moves e2e4 d7d5", E4, WHITE, false},
		{"", "position startpos moves e2e4 d7d5", B5, WHITE, true},
		{"", "position startpos moves e2e4 d7d5 f1b5", E8, WHITE, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			if got := board.isAttacked(tt.sq, tt.sd); got != tt.want {
				t.Errorf("isAttacked(%v,%v) = %v, want %v", sq2Fen[tt.sq], tt.sd, got, tt.want)
			}
		})
	}
}

func Test_setSq(t *testing.T) {
	board.newGame()
	tests := []struct {
		name string
		p12  int
		sq   int
	}{
		{"Ke4", fen2Int("K"), E4},
		{"pe3", fen2Int("p"), E3},
		{"pe5", fen2Int("p"), E5},
		{"Qb6", fen2Int("Q"), B6},
		{"nh6", fen2Int("n"), H6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := board.count[tt.p12]
			board.setSq(tt.p12, tt.sq)
			if board.sq[tt.sq] != tt.p12 {
				t.Errorf("%v: board should contain %v on sq=%v. Got %v", tt.name, tt.p12, tt.sq, board.sq[tt.sq])
			}
			pc := piece(tt.p12)
			sd := p12Colour(tt.p12)
			if !board.wbBB[sd].test(tt.sq) {
				t.Errorf("%v: wbBB[%v] on sq=%v should be set to 1 but is 0", tt.name, sd, tt.sq)
			}
			if !board.pieceBB[pc].test(tt.sq) {
				t.Errorf("%v: pieceBB[%v] on sq=%v should be set to 1 but is 0", tt.name, pc, tt.sq)
			}
			if board.count[tt.p12] != count+1 {
				t.Errorf("%v: count[%v] should be %v. Got %v", tt.name, tt.p12, count+1, board.count[tt.p12])
			}
		})
	}
}

func Test_genQueenMoves(t *testing.T) {
	ml = moveList{}

	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"start", "position startpos", []string{}, 0},
	}
	for _, tt := range tests {
		var mlq moveList
		handlePosition(tt.pos)
		board.genQueenMoves(&mlq, (^board.wbBB[board.stm]))
		if f, mv := findMoves(&mlq, tt.mv); f == false {
			t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
		}
		if tt.cnt >= 0 {
			if tt.cnt != len(mlq) {
				t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(mlq))
			}
		}
	}
}

func Test_genKnightMoves(t *testing.T) {
	ml = moveList{}
	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"startpos", "position startpos", []string{"b1a3", "b1c3", "g1f3", "g1h3"}, 4},
		{"on c3,f3", "position startpos moves b1c3 b8c6 g1f3 g8f6",
			[]string{"c3b1", "c3a4", "c3d5", "c3e4", "f3g1", "f3h4", "f3d4", "f3e5", "f3g5"}, 10},
		{"on c6,f6", "position startpos moves b1c3 b8c6 g1f3 g8f6 a2a3",
			[]string{"c6b8", "c6a5", "c6b4", "c6d4", "c6e5", "f6g4", "f6h5", "f6d5", "f6e4", "f6g4"}, 10},
		{"captures W", "position startpos moves b1c3 b8c6 c3d5 g8f6",
			[]string{"d5c7", "d5e7", "d5f6", "d5f4"}, 10},
		{"Captures B", "startpos moves b1c3 g8f6 g1f3 f6e4 h2h3",
			[]string{"e4c3", "e4d2", "e4d6", "e4f2", "e4g3"}, 10},
	}
	for _, tt := range tests {
		var ml moveList
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			board.genKnightMoves(&ml, (^board.wbBB[board.stm]))
			if f, mv := findMoves(&ml, tt.mv); f == false {
				t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
			}
			if tt.cnt >= 0 {
				if tt.cnt != len(ml) {
					t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(ml))
				}
			}
		})
	}
}
func Test_genBishopMoves(t *testing.T) {
	ml = moveList{}

	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"startpos", "position startpos", []string{}, 0},
		{"A1 blocked", "position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/BN1QKBNR w Kkq - 0 1", []string{}, 0},
		{"A1 open", "position fen rnbqkbnr/pppppp1p/6p1/8/8/1P6/P1PPPPPP/BN1QKBNR w Kkq - 0 2",
			[]string{"a1b2", "a1c3", "a1d4", "a1e5", "a1f6", "a1g7", "a1h8"}, 7},
		{"on E4", "position fen rnbqkbnr/ppp3pp/3ppp2/8/4B3/4P3/PPPP1PPP/RNBQK1NR w KQkq - 0 4",
			[]string{"e4d5", "e4c6", "e4b7", "e4f5", "e4g6", "e4h7", "e4d3", "e4f3"}, 8},
		{"B to b5", "position startpos moves e2e4 d7d5 f1b5",
			[]string{"c8d7"}, 1},
	}
	for _, tt := range tests {
		var ml moveList
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			board.genBishopMoves(&ml, (^board.wbBB[board.stm]))
			board.filterLegals(&ml)

			if f, mv := findMoves(&ml, tt.mv); f == false {
				t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
			}
			if tt.cnt >= 0 {
				if tt.cnt != len(ml) {
					t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(ml))
				}
			}
		})
	}
}

func Test_genKingMoves(t *testing.T) {
	ml = moveList{}

	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"startpos", "position startpos", []string{}, 0},
		{"e1d1", "position startpos moves d2d3 d7d6 d1d2 d8d7", []string{"e1d1"}, 1},
		{"short W k moved", "position startpos moves e2e4 e7e5 g1f3 g8f6 f1c4 f8c5 e1f1 c5f8 f1e1 f8c5",
			[]string{"e1e2", "e1f1"}, 2},
		{"short W", "position startpos moves e2e4 e7e5 g1f3 g8f6 f1c4 f8c5 ",
			[]string{"e1g1 e1e2", "e1f1"}, 3},
		{"short B", "position startpos moves e2e4 e7e5 g1f3 g8f6 f1c4 f8c5 b1c3",
			[]string{"e8g8", "e8e7", "e8f8"}, 3},
		{"Short check W", "position startpos moves e2e4 e7e5 g1f3 d8g5 f1c4 g5g2",
			[]string{"e1e2"}, 1},
		{"Short check B", "position startpos moves e2e4 e7e5 d1g4 g8f6 g4g7 f8c5 f1c4",
			[]string{"e8e7"}, 1},
		{"Long check B", "position startpos moves d2d4 d7d5 e2e4 b8a6 g1f3 c8g4 b1c3 e7e6 f1d3 d8d6 d3b5",
			[]string{"e8d8", "e8e7"}, 2},
		{"Long check W", "position startpos moves d2d4 d7d5 c1f4 e7e5 b1a3 f8c5 d1d3 c5b4",
			[]string{"e1d1"}, 1},
		{"long W", "position startpos moves d2d4 d7d5 b1c3 b8c6 c1f4 c8f5 d1d2 d8d7",
			[]string{"e1d1", "e1c1"}, 2},
		{"long B", "position startpos moves d2d4 d7d5 b1c3 b8c6 c1f4 c8f5 d1d2 d8d7 a2a3",
			[]string{"e8c8", "e8d8"}, 2},
	}
	for _, tt := range tests {
		var ml moveList
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			board.genKingMoves(&ml, (^board.wbBB[board.stm]))
			board.filterLegals(&ml)
			if f, mv := findMoves(&ml, tt.mv); f == false {
				board.Print()
				t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
			}
			if tt.cnt >= 0 {
				if tt.cnt != len(ml) {
					t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(ml))
				}
			}
		})
	}
}

func Test_genRookMoves(t *testing.T) {

	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"startpos", "position startpos", []string{}, 0},
		{"startpos B", "position startpos", []string{}, 0},
		{"A2-A7", "position startpos moves a2a4 b7b5 a4b5 h7h6", []string{"a1a2", "a1a3", "a1a4", "a1a5", "a1a6", "a1a7"}, 6},
		{"A7-A2", "position startpos moves b2b4 a7a5 h2h3 a5b4 e2e3", []string{"a8a2", "a8a3", "a8a4", "a8a5", "a8a6", "a8a7"}, 6},
		{"on E4", "position startpos moves a2a4 b7b5 a4a5 a7a6 a1a4 h7h6 a4e4 a8a7",
			[]string{"e4e3", "e4e5", "e4e6", "e4e7", "e4d4", "e4c4", "e4b4", "e4a4", "e4f4", "e4g4", "e4h4"}, 11},
	}
	for _, tt := range tests {
		var ml moveList
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.pos)
			board.genRookMoves(&ml, (^board.wbBB[board.stm]))
			if f, mv := findMoves(&ml, tt.mv); f == false {
				t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
			}
			if tt.cnt >= 0 {
				if tt.cnt != len(ml) {
					t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(ml))
				}
			}
		})
	}
}

func Test_genPawnMoves(t *testing.T) {

	tests := []struct {
		name string
		pos  string
		mv   []string
		cnt  int
	}{
		{"extra", "position fen 8/8/1k2P3/8/8/6K1/2p5/8 w - - 0 47", []string{}, 1},
		{"startpos", "position startpos", []string{"a2a3", "a2a4", "e2e3", "e2e4", "g2g3", "g2g4", "h2h3", "h2h4"}, 16},

		{"cap L W", "position fen r1bqkbnr/1ppp1p1p/2n5/p3p1p1/1P1P3P/6P1/P1P1PP2/RNBQKBNR w KQkq - 0 5",
			[]string{"a2a3", "a2a4", "e2e3", "e2e4", "b4b5", "b4a5", "g3g4", "d4e5", "h4h5"}, 15},
		{"cap R W", "position fen rnbqkbnr/p1pp1pp1/8/1p2p2p/P2P2P1/8/1PP1PP1P/RNBQKBNR w KQkq - 0 4",
			[]string{"a4b5", "a4a5", "e2e3", "e2e4", "d4e5", "d4d5", "g4h5", "h2h3", "h2h4"}, 16},

		{"ep L W", "position fen rnbqkbnr/1ppp1p2/4p1pp/pP5P/8/8/P1PPPPP1/RNBQKBNR w KQkq a6 0 5",
			[]string{"b5a6", "b5b6", "h5g6", "e2e4", "g2g3", "g2g4"}, 15},
		{"ep R W", "position fen rnbqkbnr/1p2ppp1/p1p5/P2p2Pp/8/8/1PPPPP1P/RNBQKBNR w KQkq h6 0 5",
			[]string{"b2b4", "e2e3", "e2e4", "f2f3", "f2f4", "h2h3", "h2h4", "g5h6"}, 14},

		{"pr W", "position fen r1n5/1PPP2P1/8/8/7p/4k2P/1p2p1P1/1N2K3 w - - 0 47",
			[]string{"b7a8Q", "b7a8B", "b7b8R", "b7c8N", "d7c8Q", "d7d8R", "g7g8Q", "g7g8R", "g7g8N", "g7g8B"}, 26},

		{"startpos B", "position startpos moves a2a3", []string{"a7a6", "e7e5", "h7h6", "h7h5"}, 16},

		{"cap L B", "position fen rnbqkbnr/p1pp1pp1/8/1p2p2p/P2P2P1/8/1PP1PP1P/RNBQKBNR b KQkq - 0 4",
			[]string{"b5a4", "b5b4", "d7d5", "e5d4", "e5e4", "h5g4", "h5h4"}, 16},
		{"cap R B", "position fen r1bqkbnr/1ppp1p1p/2n5/p3p1p1/1P1P3P/6P1/P1P1PP2/RNBQKBNR b KQkq - 0 5",
			[]string{"a5b4", "a5a4", "b7b5", "e5d4", "g5g4", "a5b4", "f7f5", "e5d4", "g5h4", "b7b6"}, 14},

		{"ep L B", "position fen rnbqkbnr/p1ppppp1/8/8/Pp5p/2P1PPP1/1P1P3P/RNBQKBNR b KQkq a3 0 5",
			[]string{"b4a3", "a7a5", "e7e6", "e7e5", "h4g3", "g7g5", "h4h3"}, 17},
		{"ep R B", "position fen rnbqkbnr/1p2pp1p/2p5/3p4/p5pP/P7/1PPPPPP1/RNBQKBNR b KQkq h3 0 5",
			[]string{"g4h3", "d5d4"}, 12},

		{"pr B", "position fen 8/1PPP2P1/8/6K1/8/4k3/pppp2p1/BNR5 b - - 0 47",
			[]string{"a2b1q", "b2a1r", "b2c1n", "c2b1b", "c2b1n", "c2b1q", "d2c1b", "d2c1q", "d2c1r", "d2c1n", "g2g1q", "g2g1r", "g2g1n", "g2g1b"}, 28},
	}
	for _, tt := range tests {
		var ml moveList
		t.Run(tt.name, func(t *testing.T) {
			board.newGame()
			handlePosition(tt.pos)
			ml = moveList{}
			board.genPawnMoves(&ml)
			if f, mv := findMoves(&ml, tt.mv); f == false {
				t.Errorf("%v: %v wasn't generated", tt.name, mv.String())
				log.Print(ml.String())
			}
			if tt.cnt >= 0 {
				if tt.cnt != len(ml) {
					t.Errorf("%v: number of moves should be %v. Got %v", tt.name, tt.cnt, len(ml))
					log.Print(ml.String())
				}
			}
		})
	}
}

func Benchmark_genRookMoves(b *testing.B) {
	// run the genRook function b.N times
	var ml moveList
	initFenSq2Int()
	initMagic()
	initAtksKings()
	initAtksKnights()

	handleNewgame()
	handlePosition("position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 e7e6 b1d2 c7c5 c2c3 b8c6 g1f3 f8e7 f1d3 c8d7")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		board.genRookMoves(&ml, (^board.wbBB[board.stm]))
	}
}

func findMoves(mlf *moveList, mvs []string) (bool, move) {
	found := false
	var mv1 move
	for _, mvStr := range mvs {
		fr := fenSq2Int[mvStr[:2]]
		to := fenSq2Int[mvStr[2:4]]
		cp := board.sq[to]
		if cp == empty && (board.sq[fr] == wP || board.sq[fr] == bP) && to == board.ep { // ep
			cp = wP
			if board.sq[fr] == wP {
				cp = bP
			}
		}
		pr := empty
		if len(mvStr) == 5 {
			pr = fen2Int(mvStr[4:5])
		}

		mv1.packMove(fr, to, board.sq[fr], cp, pr, board.ep, board.castlings)
		found = false
		for _, mv2 := range *mlf {
			if mv1 == mv2 {
				found = true
				break
			}
		}
		if found == false {
			return false, mv1
		}
	}
	return true, 0
}
