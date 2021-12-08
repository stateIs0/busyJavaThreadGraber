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

func NewPolice(Pid int32, tick int64, threshold int, sleep int) *Police {
	return &Police{
		Pid:       Pid,
		tick:      tick,
		threshold: threshold,
		sleep:     sleep,
	}
}

func (p *Police) Start() {
	go func() {
		for {
			if p.tick == 0 {
				p.tick = 1
			}
			time.Sleep(time.Duration(p.tick) * time.Second)
			p.parseThreadContentAndDump()
		}
	}()
}

func (p *Police) parseThreadContentAndDump() {

	thread := GetThreads(p.Pid, float64(p.threshold))

	if thread == nil || len(thread) == 0 {
		return
	}
	dumpTopThreadStack(thread, strconv.Itoa(int(p.Pid)))
}

func dumpTopThreadStack(subThreadList []SubThread, pid string) {
	cmd := "jstack -l " + pid
	command := exec.Command("bash", "-c", cmd)

	// 可能没权限.
	jstackContent, err := command.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}

	split := strings.Split(string(jstackContent), "\n")

	newFile := pid + time.Now().Format(Layout) + ".dump"
	jstackFile := pid + time.Now().Format(Layout) + ".jstack"
	output, _ := os.Create(newFile)
	jstackFileoutput, _ := os.Create(jstackFile)

	jstackFileoutput.Write(jstackContent)
	jstackFileoutput.Close()

	for idx, line := range split {
		for _, threadNum := range subThreadList {
			if strings.Contains(line, threadNum.pid16) {
				output.WriteString(split[idx-1] + ", threadNum = " + threadNum.pid16 + "\n")
				output.WriteString(line + ", CPUPercent= " + strconv.Itoa(int(threadNum.CPUPercent)) + ", " +
					"parentCPUPercent = " + strconv.Itoa(int(threadNum.parentCPUPercent)) + "\r\n")
				for i := idx; i < idx + 10; i++ {
					output.WriteString(split[i] + "\n")
				}
			}
		}
	}

	log.Println("dump 成功, file ", newFile)
	output.Close()

}
