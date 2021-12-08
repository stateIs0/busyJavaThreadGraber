package main

import (
	"awesomeProject1/grab"
	"flag"
	"fmt"
	"sync"
	"time"
)

var Pid int
var tick int64
var threshold int

var Layout = "2006-01-02 15:04:05"

func main() {
	// parse
	parseArgs()
	// start
	grab.NewPolice(int32(Pid), tick, threshold).Start()

	// wait
	group := sync.WaitGroup{}
	group.Add(1)
	fmt.Println(time.Now().Format(Layout) + " grab start and wait...")
	group.Wait()
}

// ./main -pid 24310 -tick 1 -threshold 8 -sleep 2
func parseArgs() {
	flag.IntVar(&Pid, "pid", 23751, "java pid")
	flag.Int64Var(&tick, "tick", 1, "check cpu time tick")
	flag.IntVar(&threshold, "threshold", 1, "grab cpu threshold")
	flag.Parse()
}
