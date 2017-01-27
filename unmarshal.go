package main

import "encoding/json"

const candidate = `1, 2, 3, 4, 5`

var byteCandidate = []byte(candidate)

// unmarshallFirst expects a well-formed JSON array first, then
// falls back to surrounding the input with square brackets, to
// attempt to make a well-formed JSON array.
func unmarshalFirst(str string) ([]int, error) {
	output := make([]int, 0)
	if err := json.Unmarshal([]byte(str), &output); err != nil {
		if err = json.Unmarshal([]byte("["+str+"]"), &output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// surroundFirst first surrounds the input with square brackets,
// being optimistic about having received our "CSV-style" repeated
// parameters. If this fails, falls back to unmarshaling the original
// input.
func surroundFirst(str string) ([]int, error) {
	if len(str) == 0 {
		return nil, nil
	}
	output := make([]int, 0)
	if err := json.Unmarshal([]byte("["+str+"]"), &output); err != nil {
		if err = json.Unmarshal([]byte(str), &output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// simpleHeuristic inspects the input and determines whether it
// should be surrounded with square brackets before attempting to
// unmarshal.
func simpleHeuristic(str string) ([]int, error) {
	if len(str) == 0 {
		return nil, nil
	}
	output := make([]int, 0)
	if str != "" && str != "null" && !(str[0] == '[' && str[len(str)-1] == ']') {
		str = "[" + str + "]"
		if err := json.Unmarshal([]byte(str), &output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// bytesUnmarshalFirst does the same thing as unmarshalFirst, but
// converts the input to a byte slice only once, instead of twice,
// as shown in https://github.com/TuneLab/go-truss/pull/136.
func bytesUnmarshalFirst(str string) ([]int, error) {
	if len(str) == 0 {
		return nil, nil
	}
	output := make([]int, 0)
	strb := make([]byte, len(str), len(str)+2)
	copy(strb, []byte(str))
	if err := json.Unmarshal(strb, &output); err != nil {
		strb = surround(strb)
		if err = json.Unmarshal(strb, &output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// bytesSurroundFirst does the same thing as surroundFirst, but
// converts the input to a byte slice only once, instead of twice,
// as shown in https://github.com/TuneLab/go-truss/pull/136.
func bytesSurroundFirst(str string) ([]int, error) {
	if len(str) == 0 {
		return nil, nil
	}
	output := make([]int, 0)
	strb := surround([]byte(str))
	if err := json.Unmarshal(strb, &output); err != nil {
		if err = json.Unmarshal(strb[1:len(strb)-1], &output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

func strip(in []byte) []byte {
	return in[1 : len(in)-1]
}

func surround(in []byte) []byte {
	out := in
	if cap(out) < len(in)+2 {
		out = make([]byte, 0, len(in)+2)
	}
	// Grow len(out) by two, for '[' and ']'
	out = out[:len(in)+2]

	// Always copy first so we don't clobber in, if in == out
	copy(out[1:], in)
	out[0] = '['
	out[len(in)+1] = ']'
	return out
}

func main() {

}
