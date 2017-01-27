# go-json-unmarshal-experiment

This is a dumb little experiment to help determine a sound approach for
unmarshaling parameter values in the
[Truss](https://github.com/TuneLab/go-truss) project.

When a member of a composite type in Truss is a repeated element (i.e., an
array), Truss renders it like a JSON array, but without the surrounding square
brackets, and likewise, when unmarshaling these values from HTTP parameters, it
needs to be able to interpret these comma-separated arrays. This is easy to do
if, given one of these comma-separated multi-values, we wrap it in square
brackets and unmarshal it as JSON. This gives us the added benefit of still
being able to interpret JSON- encoded parameter values, if we're intelligent
about whether we surround the parameter value with square brackets.

## So what's the problem?

There are a couple ways to do this, though, and we weren't sure which made the
most sense. This is going to matter because this ends up generating code that
will be called on every HTTP request that gets handled, so it had better be
efficient. So, we thought we'd use science. Or benchmarks. Whatever.

Any approach requires making some assumptions about the kind of input we're
likely to encounter. For this experiment, I've defined *favorable* and
*unfavorable* kinds of input, which is a bit of a misnomer, but deal with it.
*Unfavorable* means that the parameter values are all mostly-JSON arrays,
meaning they're missing the square brackets. I've called this "unfavorable"
because a generated Truss client still renders its parameters as JSON; this
added functionality is just to support clients or manual requests that send
comma-separated strings as multi-values. It follows that *favorable* means that
the parameter values are actually proper JSON arrays.

## What approaches are you considering?

* Unmarshal first: optimistically expect legitimate JSON array strings. If
  unmarshaling fails, surround the string with square brackets and try
  unmarshaling again.
* Surround first: pessimistically expect comma-separated multi-values strings,
  unlike what the Truss clients would send. Surround the input with square brackets
  and try to unmarshal. If this fails, try unmarshaling the input as provided.
* Simple heuristic: Determine whether the input needs to be surrounded with
  square brackets based on a couple simple rules.
* Single-byte-slice-conversion unnmarshal first: same as "Unmarshal first," but
  avoid converting the input string to a byte slice a second time if the
  initial unmarshal attempt fails.
* Single-byte-slice-conversion surround first: same as "Surround first," but
  avoid converting the input string to a byte slice a second time if the
  initial unmarshal attempt fails.

## So, what'd you find out?

Well, this could be because of the bogus input data, but it's remained
consistent: the simple heuristic is the fastest and requires the fewest
allocations. These are the results from my Late 2013 13" MacBook Pro with 8GB
RAM, running macOS Sierra 10.12.3.

```
 ➜  make
go test -v -bench=BenchmarkAll -benchmem -run=XXX .
>>> method: UnmarshalFirst
BenchmarkAll/UnmarshalFirst-favorable-4         	  500000	      3599 ns/op	     572 B/op	      17 allocs/op
BenchmarkAll/UnmarshalFirst-unfavorable-4       	  300000	      4302 ns/op	     957 B/op	      25 allocs/op
BenchmarkAll/UnmarshalFirst-mixed-4             	  300000	      3730 ns/op	     741 B/op	      20 allocs/op
>>> method: SurroundFirst
BenchmarkAll/SurroundFirst-favorable-4          	  300000	      4969 ns/op	     953 B/op	      23 allocs/op
BenchmarkAll/SurroundFirst-unfavorable-4        	  500000	      3825 ns/op	     588 B/op	      17 allocs/op
BenchmarkAll/SurroundFirst-mixed-4              	  300000	      4204 ns/op	     752 B/op	      19 allocs/op
>>> method: SimpleHeuristic
BenchmarkAll/SimpleHeuristic-favorable-4        	20000000	        67.0 ns/op	      32 B/op	       1 allocs/op
BenchmarkAll/SimpleHeuristic-unfavorable-4      	  300000	      3681 ns/op	     588 B/op	      17 allocs/op
BenchmarkAll/SimpleHeuristic-mixed-4            	 1000000	      1815 ns/op	     297 B/op	       8 allocs/op
>>> method: BytesUnmarshalFirst
BenchmarkAll/BytesUnmarshalFirst-favorable-4    	  500000	      3784 ns/op	     590 B/op	      17 allocs/op
BenchmarkAll/BytesUnmarshalFirst-unfavorable-4  	  300000	      4209 ns/op	     925 B/op	      24 allocs/op
BenchmarkAll/BytesUnmarshalFirst-mixed-4        	  500000	      3842 ns/op	     736 B/op	      20 allocs/op
>>> method: BytesSurroundFirst
BenchmarkAll/BytesSurroundFirst-favorable-4     	  300000	      5055 ns/op	     900 B/op	      22 allocs/op
BenchmarkAll/BytesSurroundFirst-unfavorable-4   	  300000	      4345 ns/op	     574 B/op	      17 allocs/op
BenchmarkAll/BytesSurroundFirst-mixed-4         	  300000	      4100 ns/op	     717 B/op	      19 allocs/op
PASS
ok  	_/Users/ian/tmp/unmarshal	22.955s
```
