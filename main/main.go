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
var threshold int
var sleep int

var Layout = "2006-01-02 15:04:05"

func main() {
	// parse
	parseArgs()
	// start
	grab.NewPolice(Pid, tick, threshold, sleep).Start()

	// wait
	group := sync.WaitGroup{}
	group.Add(1)
	fmt.Println(time.Now().Format(Layout) + " grab start and wait...")
	group.Wait()
}

// ./main -pid 23751 -tick 1 -threshold 800.0 -sleep 2
func parseArgs() {
	flag.StringVar(&Pid, "pid", "", "java pid")
	flag.Int64Var(&tick, "tick", 1, "check cpu time tick")
	flag.IntVar(&threshold, "threshold", 1, "grab cpu threshold")
	flag.IntVar(&sleep, "sleep", 3, "garb cpu sleep")
	flag.Parse()
}
