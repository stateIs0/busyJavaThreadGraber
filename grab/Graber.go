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

func dumpTopThreadStack(subThreadList []*SubThread, pid string) {
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
		for _, subThread := range subThreadList {
			if strings.Contains(line, subThread.pid16) {
				output.WriteString(line + ",pid16= " + subThread.pid16 + ", CPUPercent= " + strconv.Itoa(int(subThread.CPUPercent)) + ", " +
					"parentCPUPercent = " + strconv.Itoa(int(subThread.parentCPUPercent)) + "\r\n")
				for i := idx + 1 ; i < idx + 20; i++ {
					if i >= len(split) {
						break
					}
					space := strings.TrimSpace(split[i])
					if  len(space) <= 0{
						output.WriteString("==========>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>" + "\n")
						break
					}
					output.WriteString(split[i] + "\n")
				}
			}
		}
	}

	log.Println("dump 成功, file ", newFile)
	output.Close()

}
