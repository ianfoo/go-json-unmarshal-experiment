# Run combined benchmark
all:
	go test -v -bench=BenchmarkAll -benchmem -run=XXX .

.PHONY: all

