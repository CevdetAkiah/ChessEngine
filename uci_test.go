package main

import (
	"testing"
	"time"
)

var all2GUI []string

func testTell(text ...string) {
	theCMD := ""
	for ix, txt := range text {
		_ = ix
		theCMD += txt
	}

	all2GUI = append(all2GUI, theCMD)
}

func Test_Uci(t *testing.T) {
	tell = testTell
	input := make(chan string)

	go uci(input) // if not go we will block here

	tests := []struct {
		name   string
		cmd    string
		wanted []string
	}{
		{"uci", "uci", []string{"id name GoBit", "id author Cevdet", "option name Hash type spin default", "option name Threads type spin default", "uciok"}},
		{"isready", "isready", []string{"readyok"}},
		{"set Hash", "setoption name Hash value 256", []string{"info string setoption not implemented"}},
		{"skit", "skit", []string{"info string unknown cmd skit"}},
		{"pos skit", "position skit", []string{"info string Error\"skit\" must be \"fen\" or \"startpos\""}},
		{"position no cmd", "position", []string{"info string Error[] wrong length=1"}},
		{"ponderhit", "ponderhit", []string{"info string ponderhit not implemented"}},
		{"debug", "debug on", []string{"info string debug not implemented"}},
		{"go movetime", "go movetime 1000", []string{"info string go movetime not implemented"}},
		{"go movestogo", "go movestogo 20", []string{"info string go movestogo not implemented"}},
		{"go wtime", "go wtime 10000", []string{"info string go wtime not implemented"}},
		{"go btime", "go btime 11000", []string{"info string go btime not implemented"}},
		{"go winc", "go winc 500", []string{"info string go winc not implemented"}},
		{"go binc", "go binc 500", []string{"info string go binc not implemented"}},
		{"go depth", "go depth 7", []string{"info string go depth not implemented"}},
		{"go nodes", "go nodes 11000", []string{"info string go nodes not implemented"}},
		{"go mate", "go mate 11000", []string{"info string go mate not implemented"}},
		{"go ponder", "go ponder", []string{"info string go ponder not implemented"}},
		{"go infinte", "go infinite", []string{"info string go infinite not implemented"}},
		{"stop", "stop", []string{"info string stop not implemented"}},
		{"wrong cmd", "skit", []string{"info string unknown cmd"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all2GUI = []string{}
			input <- tt.cmd
			time.Sleep(10 * time.Millisecond)
			for ix, want := range tt.wanted {
				if len(all2GUI) <= ix {
					t.Errorf("%v: we want %#v in ix=%v but got nothing", tt.name, want, ix)
					continue
				}
				if len(want) > len(all2GUI[ix]) {
					t.Errorf("%v: we want %#v (in index %v) but we got %#v", tt.name, want, ix, all2GUI[ix])
					continue
				}
				if all2GUI[ix][:len(want)] != want {
					t.Errorf("%v: Error. Should be %#v but we got %#v", tt.name, want, all2GUI[ix])
				}
			}

		})
	}
}
