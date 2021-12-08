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
	//thread := getSubThread(pid)
	// 获取子进程详情
	detailSubThread := getThreadDetail2(strconv.Itoa(int(pid)), "vale")
	if len(detailSubThread) <= 0 {
		log.Println("子线程数量为0, 不需要这个工具了.")
		return nil
	}

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

/*
29630 vale      20   0 6697940  91448  12452 S 13.3  0.6   1:52.74 Thread-1
29632 vale      20   0 6697940  91448  12452 S 13.3  0.6   1:55.23 Thread-3
29635 vale      20   0 6697940  91448  12452 S 13.3  0.6   1:55.16 Thread-6
29629 vale      20   0 6697940  91448  12452 S  6.7  0.6   1:49.73 Thread-0
29631 vale      20   0 6697940  91448  12452 R  6.7  0.6   1:53.50 Thread-2
29633 vale      20   0 6697940  91448  12452 R  6.7  0.6   1:50.87 Thread-4
29636 vale      20   0 6697940  91448  12452 S  6.7  0.6   1:49.94 Thread-7
29637 vale      20   0 6697940  91448  12452 S  6.7  0.6   1:56.87 Thread-8
29617 vale      20   0 6697940  91448  12452 S  0.0  0.6   0:00.00 java
29618 vale      20   0 6697940  91448  12452 S  0.0  0.6   0:00.06 java
*/
func getThreadDetail2(goPid string, user string) []*SubThread {
	shell := fmt.Sprintf("(top  -bn 1 -Hp %s | grep %s | head -10 | sed 's/\\x1b\\x28\\x42\\x1b\\[m//' | sed 's/\\x1b\\[1m//' | sed s/[[:space:]]/\\ /g)", goPid, user)
	command := exec.Command("bash", "-c", shell)

	// 可能没权限.
	row, err := command.CombinedOutput()
	if err != nil {
		log.Println("top fail ", err, ", shell =", shell)
	}

	subThreads := []*SubThread{}
	lines := strings.Split(string(row), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		atoi, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Println("atoi ", atoi)
			return nil
		}
		float, err := strconv.ParseFloat(fields[8], 32)
		if err != nil {
			log.Println("float ", float)
			return nil
		}
		sub := &SubThread{
			pid:        atoi,
			pid16:      fmt.Sprintf("%x", fields[0]),
			CPUPercent: float,
		}
		subThreads = append(subThreads, sub)
	}
	return subThreads

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
		go func(goPid int) {
			log.Println("--------> 执行 Percent ", goPid)
			wg.Add(1)
			subPro, _ := process.NewProcess(int32(goPid))
			percent, _ := subPro.Percent(3 * time.Second)
			log.Println("--------> 执行 Percent ", goPid, ", result ", percent)
			s := &SubThread{
				pid:        goPid,
				CPUPercent: percent,
				pid16:      fmt.Sprintf("%x", goPid),
			}
			chann <- s
		}(tt)
	}

	wg.Wait()
	stop <- ""
	sort.SliceStable(subThreads, func(i, j int) bool {
		if subThreads[i].CPUPercent > subThreads[j].CPUPercent {
			return true
		}
		return false
	})
	if len(subThreads) >= 10 {
		return subThreads[0:10]
	} else {
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
		log.Println("子线程数量:", len(split)-1)
		for i, line := range split {
			if i == 0 {
				continue
			}
			lineArr := strings.Split(line, " ")
			if len(lineArr) < 2 {
				continue
			}
			// 线程号
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

func (s *SubThread) String() string  {
	return fmt.Sprintf("pid %d, pid16 %s, parentCPUPercent %f, CPUPercent %f ",
		s.pid, s.pid16, s.parentCPUPercent, s.CPUPercent)
}