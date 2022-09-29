package main

import (
	"fmt"
	"math/bits"
	"strings"
)

var wPieces uint64
var bPieces uint64

// bitboard has 64 bits
// wPieces squares have 1 bit while other squares have 0 bit
// bPieces squares have 1 bit while other squares have 0 bit
// notWhite turns all none white square bits to 1 while other squares have 0 bit
// notBlack turns all none black square bits to 1 while other squares have 0 bit
// emptySquare turns all empty square bits to 1 while squares with pieces have 0 bit
type bitBoard uint64

// func bitManipulation() {
// 	notWhite := ^wPieces
// 	notBlack := ^bPieces
// 	allPieces := wPieces | bPieces
// 	emptySquare := ^allPieces
// }

// count returns number of bits set to 1
func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

// clr sets the bit on given position to 0
func (b *bitBoard) clr(pos uint) {
	*b &= bitBoard(^(uint64(1) << pos))
}

// set board to the given position
func (b *bitBoard) set(pos uint) {
	*b |= bitBoard(uint64(1) << pos)
}

// test is a bit is one or zero
func (b bitBoard) test(pos uint) bool {
	return (b & bitBoard(uint64(1)<<pos)) != 0
}

// firstOne finds first position bit set to one
func (b *bitBoard) firstOne() int {
	bit := bits.TrailingZeros64(uint64(*b))
	// if we hit 64 without finding a 1 return 64
	if bit == 64 {
		return 64
	}
	// update bitBoard, turn off the bit. Can loop through the  board using firstOne to turn off all bits.
	*b = (*b >> uint(bit+1)) << uint(bit+1)
	return bit
}

// returns the full bitstring (with leading zeroes) of the bitboard
func (b bitBoard) String() string {
	zeroes := ""
	for ix := 0; ix < 64; ix++ {
		zeroes = zeroes + "0"
	}

	bits := zeroes + fmt.Sprintf("%b", b)

	return bits[len(bits)-64:]
}

func (b bitBoard) Stringln() string {
	s := b.String()
	row := [8]string{}
	row[0] = s[0:8]
	row[1] = s[8:16]
	row[2] = s[16:24]
	row[3] = s[24:32]
	row[4] = s[32:40]
	row[5] = s[40:48]
	row[6] = s[48:56]
	row[7] = s[56:]
	for ix, r := range row {
		row[ix] = fmt.Sprintf("%v%v%v%v%v%v%v%v\n", r[7:8], r[6:7], r[5:6], r[4:5], r[3:4], r[2:3], r[1:2], r[0:1])
	}

	s = strings.Join(row[:], "")
	s = strings.Replace(s, "1", "1 ", -1)
	s = strings.Replace(s, "0", "0 ", -1)
	return s
}
