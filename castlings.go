package main

import "strings"

type castlings uint

const (
	shortW = uint(0x1) // white can castle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can castle short
	longB  = uint(0x8) // black can castle long
)

// parse castling rights in fenstring
func parseCastlings(fenCast1 string) castlings {
	c := uint(0)
	// no castling possible
	if fenCast1 == "-" {
		return castlings(0)
	}

	if strings.Index(fenCast1, "K") >= 0 {
		c |= shortW
	}

	if strings.Index(fenCast1, "Q") >= 0 {
		c |= longW
	}

	if strings.Index(fenCast1, "k") >= 0 {
		c |= shortB
	}

	if strings.Index(fenCast1, "q") >= 0 {
		c |= longB
	}
	return castlings(c)
}