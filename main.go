package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/cpu"
)

var CacheLineSizeBytes int = int(unsafe.Sizeof(cpu.CacheLinePad{}))

var CPUs int = runtime.NumCPU()

func main() {
	fmt.Printf("Cache line size bytes: %v -- CPUs num: %v\n", CacheLineSizeBytes, CPUs)

}
