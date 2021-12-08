package grab

import (
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var newProcess *process.Process

func GetThreads(pid int32, threshold float64) []string {

	if newProcess == nil {
		log.Println("pid ", pid)
		ss, err := process.NewProcess(pid)
		if err != nil {
			return nil
		}
		newProcess = ss
	}

	percent, err := newProcess.Percent(3 * time.Second)
	if err != nil {
		log.Println(err)
		return nil
	}
	name, err := newProcess.Name()
	if err != nil {
		return nil
	}
	log.Println("pid ", name, " rootProcess percent = ", percent, ", threshold=", threshold)
	if percent == 0 {
		return nil
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	if percent < threshold {
		return nil
	}

	cmd := "ps -T -p " + strconv.Itoa(int(pid))
	c := exec.Command("bash", "-c", cmd)

	output, err := c.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil
	}
	threads := []string{}
	var count
	if len(string(output)) > 0 {

		str := string(output)
		split := strings.Split(str, "\n")

		chann := make(chan SubThread, len(split)-1)

		for i, line := range split {
			if i == 0 {
				continue
			}
			lineArr := strings.Split(line, " ")
			if len(lineArr) < 2 {
				continue
			}
			subThread := lineArr[1]
			atoi, err := strconv.Atoi(subThread)
			if err != nil {
				continue
			}
			count += 1
			go func() {
				subPro, _ := process.NewProcess(int32(atoi))
				percent, _ := subPro.Percent(3 * time.Second)
				s := SubThread{
					pid:        subThread,
					CPUPercent: percent,
				}
				chann <- s
			}()

		}

		for true {
			select {
			case data, ok := <-chann:
				if ok {
					if data.CPUPercent > 10 {
						threads = append(threads, data.pid)
					}
				}
			default:
				if len(threads) >= count {
					log.Println("threads len --->> 2 " + strconv.Itoa(len(threads)))
					return threads
				}
			}
		}
	}
	log.Println("threads len --->>" + strconv.Itoa(len(threads)))
	return threads

}

type SubThread struct {
	pid        string
	CPUPercent float64
}
