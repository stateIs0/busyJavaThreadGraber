package cpu

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var Layout = "2006-01-02 15:04:05"

func Get(pid string, sleep int) (float64, float64, float64) {
	idle0, total0 := getCPUSample(pid)
	time.Sleep(time.Duration(sleep) * time.Second)
	idle1, total1 := getCPUSample(pid)

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
	// usage busy total
	fmt.Printf(time.Now().Format(Layout)+": CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
	return cpuUsage, totalTicks - idleTicks, totalTicks
}

func getCPUSample(pid string) (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/" + pid + "/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
