package grab

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Layout = "2006-01-02_15:04:05"

type Police struct {
	Pid       int32
	tick      int64
	threshold int
	sleep     int
}

func NewPolice(Pid int32, tick int64, threshold int) *Police {
	return &Police{
		Pid:       Pid,
		tick:      tick,
		threshold: threshold,
	}
}

func (p *Police) Start() {
	go func() {
		for {
			if p.tick == 0 {
				p.tick = 1
			}
			time.Sleep(time.Duration(p.tick) * time.Second)
			p.process()
		}
	}()
}

func (p *Police) process() {

	busyThread := GrabBusyThreads(p.Pid, float64(p.threshold))

	if busyThread == nil || len(busyThread) == 0 {
		return
	}
	dumpThreadStack2File(busyThread, strconv.Itoa(int(p.Pid)))
}

func dumpThreadStack2File(subThreadList []*SubThread, pid string) {
	cmd := "jstack -l " + pid
	command := exec.Command("bash", "-c", cmd)

	// 可能没权限.
	jstackContent, err := command.CombinedOutput()
	if err != nil {
		log.Println("jstack error >>>>>>", err)
		return
	}

	split := strings.Split(string(jstackContent), "\n")

	newFile := pid + "_" + time.Now().Format(Layout) + ".txt"
	jstackFile := pid + "_" + time.Now().Format(Layout) + ".jstack"
	output, err := os.Create(newFile)
	if err != nil {
		log.Println(err)
	}
	jstackFileoutput, err := os.Create(jstackFile)
	if err != nil {
		log.Println(err)
	}
	jstackFileoutput.Write(jstackContent)
	jstackFileoutput.Close()

	for idx, line := range split {
		for _, subThread := range subThreadList {
			if !strings.Contains(line, subThread.pid16) {
				continue
			}
			output.WriteString(line + ",16 进制 ID = " + subThread.pid16 +
				", 该线程 CPU 使用率 = " + strconv.Itoa(int(subThread.CPUPercent)) + ", " +
				"Java 进程 CPU 使用率 = " + strconv.Itoa(int(subThread.parentCPUPercent)) + "\r\n")

			for i := idx + 1; i < idx+30; i++ {
				if i >= len(split) {
					break
				}
				// 有空行了,大概就是这个堆栈结束了.
				space := strings.TrimSpace(split[i])
				if len(space) <= 0 {
					output.WriteString("==========>>>>>>>>>>>>> 分隔符 >>>>>>>>>>>>>>>>>>>>" + "\n")
					break
				}
				output.WriteString(split[i] + "\n")
			}
		}
	}

	log.Println("dump 成功, file ", newFile)
	output.Close()

}
