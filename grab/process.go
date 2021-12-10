package grab

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getParentThreadState(pid1 int32, channle chan float64, tick int) {

	newProcess, err := process.NewProcess(pid1)
	if err != nil {
		return
	}

	parentCPUPercent, err := newProcess.Percent(time.Duration(tick) * time.Second)
	if err != nil {
		log.Println(err)
		channle <- parentCPUPercent
		return
	}

	if parentCPUPercent == 0 {
		channle <- parentCPUPercent
		return
	}
	channle <- parentCPUPercent
}

func GrabBusyThreads(pid int32, threshold float64, tick int, threadNum int, user string) []*SubThread {

	getParentThreadStateResult := make(chan float64)
	// 获取进程的状态
	go func() { getParentThreadState(pid, getParentThreadStateResult, tick) }()

	var parentCPUPercent = 0.0

	select {
	// 等待进程的统计结果
	case data := <-getParentThreadStateResult:
		parentCPUPercent = data
	// 超时 5s
	case <-time.After(5 * time.Second):
		break
	}

	log.Println("pid ", pid, " Java 进程 cpu 率 = ", parentCPUPercent, ", 触发 dump 阈值 = ", threshold)

	if parentCPUPercent < threshold {
		return nil
	}
	// 获取线程详情
	m := map[string]*SubThread{}
	// 多执行几遍, top 可能会漏掉.
	for i := 0; i < 5; i++ {
		detailSubThread := getThreadDetail(strconv.Itoa(int(pid)), user, threadNum)
		for _, thread := range detailSubThread {
			old, ok := m[thread.pid16]
			if ok {
				// 如果老数据, 少于新数据, 就用新的.
				if old.CPUPercent < thread.CPUPercent {
					m[thread.pid16] = thread
				}
			} else {
				m[thread.pid16] = thread
			}
		}
	}
	detailSubThread := []*SubThread{}
	for _, thread := range m {
		detailSubThread = append(detailSubThread, thread)
	}

	if len(detailSubThread) <= 0 {
		log.Println("子线程数量为0, 难道不是 Java 进程?")
		time.Sleep(time.Duration(tick) * time.Second)
		return nil
	}

	for _, subThread := range detailSubThread {
		subThread.parentCPUPercent = parentCPUPercent
	}

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

指定进程, 指定 user, 指定 top n 线程数.

*/
func getThreadDetail(goPid string, user string, threadNum int) []*SubThread {
	shell := fmt.Sprintf("(top  -bn 1 -Hp %s | grep %s | head -%s | sed 's/\\x1b\\x28\\x42\\x1b\\[m//' | sed 's/\\x1b\\[1m//' | sed s/[[:space:]]/\\ /g)",
		goPid, user, strconv.Itoa(threadNum))

	command := exec.Command("bash", "-c", shell)

	// 可能没权限.
	row, err := command.CombinedOutput()
	if err != nil {
		log.Println("top fail ", err, ", shell =", shell)
	}

	var subThreads []*SubThread
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
			pid16:      fmt.Sprintf("%x", atoi),
			CPUPercent: float,
		}
		subThreads = append(subThreads, sub)
	}
	return subThreads

}

type SubThread struct {
	pid              int
	pid16            string
	parentCPUPercent float64
	CPUPercent       float64
}

func (s *SubThread) String() string {
	return fmt.Sprintf("pid %d, pid16 %x, parentCPUPercent %f, CPUPercent %f \n",
		s.pid, s.pid16, s.parentCPUPercent, s.CPUPercent)
}
