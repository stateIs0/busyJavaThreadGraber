package grab

import (
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetThreads(pid int32, threshold float64) []SubThread {

	newProcess, err := process.NewProcess(pid)
	if err != nil {
		return nil
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
	return Handler(output)

}

func Handler(output []byte) []SubThread {
	threads := []SubThread{}
	wg := sync.WaitGroup{}
	stop := make(chan string)
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
			go func() {
				wg.Add(1)
				subPro, _ := process.NewProcess(int32(atoi))
				percent, _ := subPro.Percent(3 * time.Second)
				s := SubThread{
					pid:        atoi,
					CPUPercent: percent,
				}
				chann <- s
			}()
		}

		go func() {
			for true {
				select {
				case data, ok := <-chann:
					wg.Done()
					if ok {
						if data.CPUPercent >= 50 {
							threads = append(threads, data)
						}
					}
				case <-stop:
					return
				}
			}
		}()

	}
	wg.Wait()
	stop <- ""
	sort.SliceStable(threads, func(i, j int) bool {
		if threads[i].CPUPercent > threads[j].CPUPercent {
			return true
		}
		return false
	})
	log.Println("threads len --->>",threads[0:5])
	return threads[0:5]
}

type SubThread struct {
	pid        int
	CPUPercent float64
}
