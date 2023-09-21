package main

import (
	"sync"
	"testing"
)

/*
~/Desktop/test-cache-line-perf ❯ go test -bench=. -count 5 -run=^#                         1m 10s 01:32:29 PM
goos: darwin
goarch: arm64
pkg: test-cache-line-perf
Benchmark_no_dirty_cache_line-12        1000000000               0.2427 ns/op
Benchmark_no_dirty_cache_line-12        1000000000               0.2239 ns/op
Benchmark_no_dirty_cache_line-12        1000000000               0.2326 ns/op
Benchmark_no_dirty_cache_line-12        1000000000               0.2358 ns/op
Benchmark_no_dirty_cache_line-12        1000000000               0.2245 ns/op
Benchmark_dirty_cache_line-12           1000000000               0.4307 ns/op
Benchmark_dirty_cache_line-12           1000000000               0.4458 ns/op
Benchmark_dirty_cache_line-12           1000000000               0.4251 ns/op
Benchmark_dirty_cache_line-12           1000000000               0.4559 ns/op
Benchmark_dirty_cache_line-12           1000000000               0.4442 ns/op
PASS
ok      test-cache-line-perf    68.388s
*/

/*
Explanation:
Each core(thread in our case) has its own L1, and L2 cache
Each cache obtains data from memory by copying it into a cache line (128 bytes on my machine)
So when a thread wants to read a single byte, it will still have to read in a whole cach line
The issue is what happens when multiple threads have the same cache line on their own local caches?

Lets say thread 1 has cache line [X,Y,Z] on its local cache (L1 or L2)
Lets say thread 2 has the same memory location in its cache line
If thread 1 modifies this cache line, the cache on thread 2 has to be declared as dirty
This means that thread 2 will actually have to refetch the data from memory, and not rely on cache

The tests are as follows:
Run 10 threads
Test1:
Each thread increamenets a segment of the array that is at a different cache line
So thread 1 increments memory segment A
Thread 2 increments memory segment B
And so on
This means that the cache is never dirty from other threads writing to the cache line
This results in faster execution


Test 2
All threads now actually access the same cache line
They read the same value, modify it, and then re-write it
This means that whenever a thread writes to a cache line, this line has to be marked
as dirty in all other threads caches, and is thus almost twice as slower.
This is the False Sharing Pattern.

The results may be even more pronounced if the threads were OS threads, but go opts
to mulitplex multiple goroutines on fewer OS threads.

Resources:
https://en.wikipedia.org/wiki/False_sharing
https://www.youtube.com/watch?v=WDIkqP4JbkE
*/

var arr []int = make([]int, CacheLineSizeBytes*4*CPUs+1)

func Benchmark_no_dirty_cache_line(t *testing.B) {
	// start := time.Now()

	wg := sync.WaitGroup{}
	for i := 0; i < CPUs; i++ {
		wg.Add(1)
		go func(threadNum int) {
			for ctr := 0; ctr < 100000000; ctr++ {
				arr[(CacheLineSizeBytes*4)*threadNum]++
			}
			wg.Done()
		}(i)
	}
	wg.Wait() //wait for all threads to finish execution

	// elapsed := time.Since(start)
	// fmt.Printf("Microseconds elapsed: %v\n", elapsed.Microseconds())
}

func Benchmark_dirty_cache_line(t *testing.B) {
	// start := time.Now()

	wg := sync.WaitGroup{}
	for i := 0; i < CPUs; i++ {
		wg.Add(1)
		go func(threadNum int) {
			for ctr := 0; ctr < 100000000; ctr++ {
				arr[threadNum]++
			}
			wg.Done()
		}(i)
	}
	wg.Wait() //wait for all threads to finish execution

	// elapsed := time.Since(start)
	// fmt.Printf("Microseconds elapsed: %v\n", elapsed.Microseconds())
}
