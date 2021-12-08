package grab

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"strconv"
)

func GetThreads(pid int32, threshold float64) []string {

	processes, err := process.Processes()
	if err != nil {
		return nil
	}
	var root *process.Process
	for _, p := range processes {
		if p.Pid == pid {
			root = p
			break
		}
	}

	CPUPercent, _ := root.CPUPercent()

	log.Println("pid ", pid, " rootProcess CPUPercent = ",CPUPercent , ", threshold=", threshold)

	if err != nil {
		log.Println(err)
		return nil
	}

	if CPUPercent < threshold {
		return nil
	}
	cmd := "ps -T -p" + strconv.Itoa(int(pid))
	c := exec.Command("bash ", "-c", cmd)
	output, _ := c.CombinedOutput()
	log.Println("--->>" + strconv.Itoa(int(pid)) + string(output))

	threads := []string{}

	if CPUPercent > threshold {
		children, _ := root.Children()
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
