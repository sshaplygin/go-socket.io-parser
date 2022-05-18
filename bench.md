# Benchmarks

## Marshal

```bash
GOMAXPROCS=1 go test -bench=BenchmarkMarshal -benchmem -benchtime=10s
GOMAXPROCS=4 go test -bench=BenchmarkMarshal -benchmem -benchtime=10s 
GOMAXPROCS=10 go test -bench=BenchmarkMarshal -benchmem -benchtime=10s 
```

results:
```
BenchmarkMarshal        11542710              1026 ns/op            4936 B/op         11 allocs/op
BenchmarkMarshal-4      10486740              1079 ns/op            4936 B/op         11 allocs/op
BenchmarkMarshal-10     10489590              1124 ns/op            4938 B/op         11 allocs/op
```

## Unmarshal

```bash
GOMAXPROCS=1 go test -bench=BenchmarkUnmarshal -benchmem -benchtime=10s
GOMAXPROCS=4 go test -bench=BenchmarkUnmarshal -benchmem -benchtime=10s 
GOMAXPROCS=10 go test -bench=BenchmarkUnmarshal -benchmem -benchtime=10s 
```

```
BenchmarkUnmarshal      11866780               984.1 ns/op           640 B/op         16 allocs/op
BenchmarkUnmarshal-4    12247012               980.7 ns/op           640 B/op         16 allocs/op
BenchmarkUnmarshal-10           11397696               993.0 ns/op           640 B/op         16 allocs/op
```

### Compare changes
```bash
go test -run=NONE -bench=. ./... > old.txt
# make changes
go test -run=NONE -bench=. ./... > new.txt

benchcmp old.txt new.txt
```

### Bench with profiles

```bash
go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out ./...
```
