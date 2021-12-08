package grab

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"strconv"
	"time"
)

func GetThreads(pid int32, threshold float64) []string {

	newProcess, err := process.NewProcess(pid)
	if err != nil {
		return nil
	}

	percent, err := newProcess.Percent(3 * time.Second)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("pid ", pid, " rootProcess CPUPercent = ", percent, ", threshold=", threshold)

	if err != nil {
		log.Println(err)
		return nil
	}

	if percent < threshold {
		return nil
	}

	threads := []string{}

	if percent > threshold {
		children, _ := newProcess.Children()
		fmt.Println("children:", children)
		for _, p := range children {
			cpuPercent, err := p.CPUPercent()
			if err != nil {
				log.Println(err)
				return nil
			}
			log.Println(strconv.Itoa(int(p.Pid))+", cpuPercent =", cpuPercent)
			if int(cpuPercent) > 10 {
				threads = append(threads, strconv.Itoa(int(p.Pid)))
			}
		}
	}

	return threads

}
