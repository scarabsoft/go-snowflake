go-snowflake
====
[![GoDoc](https://godoc.org/github.com/scarabsoft/go-snowflake?status.svg)](https://godoc.org/github.com/scarabsoft/go-snowflake) 
[![Go Report Card](https://goreportcard.com/badge/github.com/scarabsoft/go-snowflake)](https://goreportcard.com/report/github.com/scarabsoft/go-snowflake)
[![Code Coverage](https://codecov.io/gh/scarabsoft/go-snowflake/branch/main/graph/badge.svg)](https://codecov.io/gh/scarabsoft/go-snowflake)

snowflake is a [Go](https://golang.org/) package that provides
* A very simple Twitter like snowflake generator.
* Monotonic Clock calculations protect from clock drift.

## Status
This package should be considered stable and completed.  Any additions in the 
future will strongly avoid API changes to existing functions. 
  
### ID Format
The format is different from original Twitter snowflake format.
* The ID as a whole is a 64 bit integer stored in an uint64
* 42 bits are used to store a timestamp with millisecond precision, using a custom epoch.
*  8 bits are used to store a node id - a range from 0 through 255.
* 14 bits are used to store a sequence number - a range from 1 through 16383.

```
+----------------------------------------------------------+
|42 Bit Timestamp |  8 Bit NodeID  |   14 Bit Sequence ID |
+---------------------------------------------------------+
```

This allows for **16383** unique IDs to be generated every millisecond, per Node ID.

### Custom Format
You can alter the number of bits used for the node id and step number (sequence)
by setting the snowflake.NodeBits and snowflake.StepBits values.  Remember that
There is a maximum of 22 bits available that can be shared between these two 
values. You do not have to use all 22 bits.

### Custom Clock
By default this package uses the Unix Epoch of 0 or January 1, 1970 12:00:00 AM.
You can set your own epoch value by setting to a time in milliseconds
to use as the epoch.

```go
gen, err := snowflake.New(
    snowflake.WithClock(snowflake.NewUnixClockWithEpoch(now)),
)
```

### Custom Epoch
By default the generator uses the Unix Epoch of 0 or January 1, 1970 12:00:00 AM.
You can set your own epoch value by setting to a time in milliseconds
to use as the epoch.

```go
gen, err := snowflake.New(
    snowflake.WithClock(snowflake.NewUnixClockWithEpoch(now)),
)
```

### Custom Node Id
By default the generator uses 1 as the default nodeID. You can set it like:
```go
gen, err := snowflake.New(
    snowflake.WithNodeId(1),
)
```

### How it Works.
Each time you generate an ID, it works, like this.
* A timestamp with millisecond precision is stored using 42 bits of the ID.
* Then the nodeID is added in subsequent bits.
* Then the Sequence iteration is added, starting at 1 and incrementing for each ID generated in the same millisecond.
 If you generate more IDs in a millisecond than max available, the generator will stop producing Ids until the next millisecond 
 

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

```sh
go get github.com/scarabsoft/go-snowflake
```


### Usage

Import the package into your project then construct a new snowflake generator. Use the Next() method to generate
a new id. It returns a chanel where you will receive a Result. Which either contains the ID or an Error.

```go
gen, err := snowflake.New()

for i := 0; i < 10; i++ {
    r := <-gen.Next()
    fmt.Printf("%064b\n", r.ID)
}
```
 
Keep in mind that each node you create must have a unique node number, even 
across multiple servers.  If you do not keep node numbers unique the generator 
cannot guarantee unique IDs across all nodes.

Use only a clock implementation which increases monotonic. If you use a clock which does not make any progress, the generator
will block once the the sequences are exhausted for max 1 ms.

### Performance

```bash
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkTestBenchmark_Single
BenchmarkTestBenchmark_Single-8     	15689761	        72.59 ns/op
BenchmarkTestBenchmark_Parallel
BenchmarkTestBenchmark_Parallel-8   	10395352	       117.8 ns/op
```

With default settings, this snowflake generator should be sufficiently fast 
enough on most systems to generate 4096+ unique ids per millisecond. 
The maximum snowflake ID this format supports is 16383. 
