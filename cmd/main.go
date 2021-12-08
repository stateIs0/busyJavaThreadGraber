package main

import (
	"awesomeProject1/grab"
	"flag"
)

var Pid string

func main() {
	parseArgs()

	grab.NewPolice(Pid).Start()
}

func parseArgs() {
	flag.StringVar(&Pid, "Pid", "", "")
	flag.Parse()
}
