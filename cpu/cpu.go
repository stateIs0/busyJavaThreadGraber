package cpu

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Layout = "2006-01-02 15:04:05"

//在proc文件系统中，可以通过/proc/[pid]/stat获得进程消耗的时间片，
//输出的第14、15、16、17列分别对应进程用户态CPU消耗、内核态的消耗、
//用户态等待子进程的消耗、内核态等待子进程的消耗(man proc)。
//所以进程的CPU消耗可以使用如下命令：
//cat /proc/9583/stat|awk '{print "cpu_process_total_slice " $14+$15+$16+$17}'
//cpu_process_total_slice 1068099
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

func Get2(pid int) float64 {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		processes = append(processes, &Process{pid, cpu})
	}
	for _, p := range processes {
		if p.pid == pid {
			log.Println("Process ", p.pid, " takes ", p.cpu, " % of the CPU")
			return p.cpu
		}
	}
	return 0
}

type Process struct {
	pid int
	cpu float64
}

func getCPUSample(pid string) (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/" + pid + "/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), " ")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[14] == "cpu" {
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
