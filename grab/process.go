package grab

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"strconv"
)

func GetThreads(pid int32, a float64) []string {
	var rootProcess *process.Process
	processes, _ := process.Processes()
	for _, p := range processes {
		if p.Pid == pid {
			rootProcess = p
			break
		}
	}

	percent, err := rootProcess.CPUPercent()

	log.Println("rootProcess CPUPercent = ", percent)

	if err != nil {
		log.Println(err)
		return nil
	}

	if percent < a {
		return nil
	}

	threads := []string{}

	if percent > a {
		fmt.Println("children:")
		children, _ := rootProcess.Children()
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
