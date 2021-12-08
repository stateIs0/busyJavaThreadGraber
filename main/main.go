package main

import (
	"awesomeProject1/grab"
	"flag"
	"fmt"
	"sync"
	"time"
)

var Pid string
var tick int64
var threshold float64

func main() {
	// parse
	parseArgs()
	// start
	grab.NewPolice(Pid, tick, threshold).Start()

	// wait
	group := sync.WaitGroup{}
	group.Add(1)
	fmt.Println(time.Now().Format(grab.Layout) + " grab start and wait...")
	group.Wait()
}

// -pid 23751 -tick 1 -threshold 800.0
func parseArgs() {
	flag.StringVar(&Pid, "pid", "", "java pid")
	flag.Int64Var(&tick, "tick", 1, "check cpu time tick")
	flag.Float64Var(&threshold, "threshold", 1.0, "grab cpu threshold")
	flag.Parse()
}
