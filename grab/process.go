package grab

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"strconv"
)

func GetThreads(pid int32, threshold float64) []string {
	var rootProcess *process.Process
	processes, _ := process.Processes()
	for _, p := range processes {
		if p.Pid == pid {
			rootProcess = p
			break
		}
	}

	percent, err := rootProcess.CPUPercent()

	log.Println(pid, " rootProcess CPUPercent = ", percent*100, ", threshold=", threshold)

	if err != nil {
		log.Println(err)
		return nil
	}

	if percent*100 < threshold {
		return nil
	}

	threads := []string{}

	if percent*100 > threshold {
		children, _ := rootProcess.Children()
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
