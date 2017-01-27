package main

import (
	"fmt"
	"testing"
)

var favorableInputs = []string{
	"[1,2,3,4,5]",
	"[6,7,8,9,10,11,12]",
	"[17,24,19,0,-1,14,53,79]",
	"[144,121,81,225,400,1600,256000]",
	"[22, 19, 17, 23, 34, 17, 44, 46, 29]",
	"[243, -51, 42, 36, 70, 101, 102, 244]",
}

var unfavorableInputs = []string{
	"1,2,3,4,5",
	"6,7,8,9,10,11,12",
	"17,24,19,0,-1,14,53,79",
	"144,121,81,225,400,1600,256000",
	"22, 19, 17, 23, 34, 17, 44, 46, 29",
	"243, -51, 42, 36, 70, 101, 102, 244",
}

var mixedInputs = []string{
	"1,2,3,4,5",
	"[1,2,3,4,5]",
	"17,24,19,0,-1,14,53,79",
	"[144,121,81,225,400,1600,256000]",
	"22, 19, 17, 23, 34, 17, 44, 46, 29",
	"[ 243, -51, 42, 36, 70, 101, 102, 244 ]",
}

var inputs = map[string][]string{
	"favorable":   favorableInputs,
	"unfavorable": unfavorableInputs,
	"mixed":       mixedInputs,
}

var methods = map[string]func(string) ([]int, error){
	"UnmarshalFirst":      unmarshalFirst,
	"SurroundFirst":       surroundFirst,
	"SimpleHeuristic":     simpleHeuristic,
	"BytesUnmarshalFirst": bytesUnmarshalFirst,
	"BytesSurroundFirst":  bytesSurroundFirst,
}

func BenchmarkAll(b *testing.B) {
	for methdesc, bf := range methods {
		fmt.Printf(">>> method: %s\n", methdesc)
		for inpdesc, in := range inputs {
			b.Run(methdesc+"-"+inpdesc, makeRunBenchmark(in, bf))
		}
	}
}

func makeRunBenchmark(inputs []string, f func(string) ([]int, error)) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in := inputs[i%len(inputs)]
			out, err := f(in)
			if err != nil {
				b.Fatalf("error running benchmark: %v, input: %q", err, in)
			}
			_ = out
		}
	}
}

func BenchmarkUnmarshalFirst(b *testing.B) {
	for desc, in := range inputs {
		b.Run(desc+"-UnmarshalFirst", makeRunBenchmark(in, unmarshalFirst))
	}
}

func BenchmarkSurroundFirst(b *testing.B) {
	for desc, in := range inputs {
		b.Run(desc+"-SurroundFirst", makeRunBenchmark(in, surroundFirst))
	}
}

func BenchmarkSimpleHeuristic(b *testing.B) {
	for desc, in := range inputs {
		b.Run(desc+"-SimpleHeuristic", makeRunBenchmark(in, simpleHeuristic))
	}
}

func BenchmarkBytesUnmarshalFirst(b *testing.B) {
	for desc, in := range inputs {
		b.Run(desc+"-BytesUnmarshalFirst", makeRunBenchmark(in, bytesUnmarshalFirst))
	}
}

func BenchmarkBytesSurroundFirst(b *testing.B) {
	for desc, in := range inputs {
		b.Run(desc+"-BytesSurroundFirst", makeRunBenchmark(in, bytesSurroundFirst))
	}
}

func TestDummy(t *testing.T) {
	// This is here just to squelch the "warning: no tests to run" message.
}
