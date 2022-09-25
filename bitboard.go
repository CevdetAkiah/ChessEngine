package main


var wPieces uint64
var bPieces uint64

// bitboard has 64 bits
// wPieces squares have 1 bit while other squares have 0 bit
// bPieces squares have 1 bit while other squares have 0 bit
// notWhite turns all none white square bits to 1 while other squares have 0 bit
// notBlack turns all none black square bits to 1 while other squares have 0 bit
// emptySquare turns all empty square bits to 1 while squares with pieces have 0 bit
type bitBoard uint64

func bitManipulation() {
	notWhite := ^wPieces
	notBlack := ^bPieces
	allPieces := wPieces | bPieces
	emptySquare := ^allPieces
}
