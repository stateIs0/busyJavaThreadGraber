package main

import (
	"awesomeProject1/grab"
	"flag"
	"fmt"
	"sync"
	"time"
)
var layout = "2006-01-02 15:04:05"

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
	fmt.Println(time.Now().Format(layout) + " grab start and wait...")
	group.Wait()
}

func parseArgs() {
	flag.StringVar(&Pid, "pid", "", "java pid")
	flag.Int64Var(&tick, "tick", 1, "check cpu time tick")
	flag.Float64Var(&threshold, "threshold", 1.0, "grab cpu threshold")
	flag.Parse()
}

