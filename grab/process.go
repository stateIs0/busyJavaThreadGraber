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

func getParentThreadState(pid1 int32, channle chan float64) {
	newProcess, err := process.NewProcess(pid1)
	if err != nil {
		return
	}

	parentCPUPercent, err := newProcess.Percent(3 * time.Second)
	if err != nil {
		log.Println(err)
		return
	}

	if parentCPUPercent == 0 {
		return
	}
	channle <- parentCPUPercent
}

func GetThreads(pid int32, threshold float64) []*SubThread {

	getParentThreadStateResult := make(chan float64)
	// 获取进程的状态
	go func() { getParentThreadState(pid, getParentThreadStateResult) }()
	// 获取所有的子进程
	thread := getSubThread(pid)
	// 获取子进程详情
	detailSubThread := getThreadDetail(thread)

	var parentCPUPercent = 0.0

	select {
	// 等待结果
	case data := <-getParentThreadStateResult:
		parentCPUPercent = data
	// 超时 5s
	case <-time.After(5 * time.Second):
		break
	}

	log.Println("pid ", pid, " rootProcess parentCPUPercent = ", parentCPUPercent, ", threshold=", threshold)

	if parentCPUPercent < threshold {
		return nil
	}

	for _, subThread := range detailSubThread {
		subThread.parentCPUPercent = parentCPUPercent
	}

	log.Println("threads len --->>", detailSubThread)
	return detailSubThread
}

func getThreadDetail(threads []int) []*SubThread {
	subThreads := []*SubThread{}
	stop := make(chan string)
	wg := sync.WaitGroup{}
	chann := make(chan *SubThread, len(threads))

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
			s := &SubThread{
				pid:        pdi,
				CPUPercent: percent,
				pid16:      fmt.Sprintf("%x", pdi),
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
	if len(subThreads) >= 0 {
		return subThreads[0:10]
	}else {
		return subThreads[0:]
	}
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
