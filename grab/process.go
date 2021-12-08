package grab

import (
	"fmt"
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

	parentCPUPercent, err := newProcess.Percent(3 * time.Second)
	if err != nil {
		log.Println(err)
		return nil
	}
	name, err := newProcess.Name()
	if err != nil {
		return nil
	}
	log.Println("pid ", name, " rootProcess parentCPUPercent = ", parentCPUPercent, ", threshold=", threshold)
	if parentCPUPercent == 0 {
		return nil
	}

	thread := getSubThread(pid)
	detailSubThread := getThreadDetail(thread, parentCPUPercent)

	if parentCPUPercent < threshold {
		return nil
	}
	log.Println("threads len --->>", detailSubThread)
	return detailSubThread
}

func getThreadDetail(threads []int, parentCPUPercent float64) []SubThread{
	subThreads := []SubThread{}
	stop := make(chan string)
	wg := sync.WaitGroup{}
	chann := make(chan SubThread, len(threads))

	go func() {
		for true {
			select {
			case data, ok := <-chann:
				wg.Done()
				if ok {
					if data.CPUPercent >= 0 {
						subThreads = append(subThreads, data)
					}
				}
			case <-stop:
				return
			}
		}
	}()

	for _, tt := range threads {
		pdi := tt
		go func() {
			wg.Add(1)
			subPro, _ := process.NewProcess(int32(pdi))
			percent, _ := subPro.Percent(2 * time.Second)
			s := SubThread{
				pid:              pdi,
				CPUPercent:       percent,
				pid16:            fmt.Sprintf("%x", pdi),
				parentCPUPercent: parentCPUPercent,
			}
			chann <- s
		}()
	}

	wg.Wait()
	stop <- ""
	sort.SliceStable(subThreads, func(i, j int) bool {
		if subThreads[i].CPUPercent > subThreads[j].CPUPercent {
			return true
		}
		return false
	})

	return subThreads[0:10]
}

func getSubThread(pid int32) []int {

	result := []int{}

	cmd := "ps -T -p " + strconv.Itoa(int(pid))
	c := exec.Command("bash", "-c", cmd)

	output, err := c.CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	if len(string(output)) > 0 {
		str := string(output)
		split := strings.Split(str, "\n")
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
			result = append(result, atoi)
		}
	}

	return result

}

type SubThread struct {
	pid              int
	pid16            string
	parentCPUPercent float64
	CPUPercent       float64
}
