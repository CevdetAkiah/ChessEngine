package main

import "math/bits"

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
func (b bitBoard) count() int{
	return bits.OnesCount64(uint64(b))
}

// clr sets the bit on given position to 0
func (b *bitBoard) clr(pos uint) {
	*b &= bitBoard(^(uint64(1) << pos))
}

// set board to the given position
func (b *bitBoard) set(pos uint){
	*b |= bitBoard(uint64(1) << pos)
}

// test is a bit is one or zero
func (b bitBoard) test(pos uint) bool{
	return (b & bitBoard(uint64(1) << pos)) != 0
}

// firstOne finds first position bit set to one
func (b * bitBoard) firstOne() int{
	bit := bits.TrailingZeros64(uint64(*b))
	// if we hit 64 without finding a 1 return 64
	if bit == 64 {
		return 64
	}
	// update bitBoard, turn off the bit. Can loop through the  board using firstOne to turn off all bits.
	*b = (*b >> uint(bit+1)) << uint(bit+1)
	return bit
}