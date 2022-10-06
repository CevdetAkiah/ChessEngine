package main

import "strings"

type castlings uint

const (
	shortW = uint(0x1) // white can castle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can castle short
	longB  = uint(0x8) // black can castle long
)

type castlOptions struct {
	short    uint
	long     uint
	rook     int
	king     int
	rookSh   uint
	rookLong uint
}

var castl = [2]castlOptions{{shortW, longW, wR, E1, H1, A1}, {shortB, longB, bR, E8, H8, A8}}

// only castling privileges (not if it is legal on board)
func (c castlings) canCastle(sd colour) bool {
	return c.canCastleShort(sd) || c.canCastleLong(sd)
}
func (c castlings) canCastleShort(sd colour) bool {
	return (castl[sd].short & uint(c)) != 0
}
func (c castlings) canCastleLong(sd colour) bool {
	return (castl[sd].long & uint(c)) != 0
}

func (c *castlings) on(val uint) {
	(*c) |= castlings(val)
}

func (c *castlings) off(val uint) {
	(*c) &= castlings(^val)
}

func (c castlings) String() string {
	flags := ""
	if uint(c)&shortW != 0 {
		flags = "K"
	}
	if uint(c)&longW != 0 {
		flags += "Q"
	}
	if uint(c)&shortB != 0 {
		flags += "k"
	}
	if uint(c)&longB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}

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
