package main

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/zenpaw-labs/skypaw/cmd"
)

func main() {
	// go path_utils.AddToPath()
	cpuFile, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	cmd.Execute()

	memFile, _ := os.Create("mem.prof")
	defer memFile.Close()
	runtime.GC()
	pprof.WriteHeapProfile(memFile)
}
