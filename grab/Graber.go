package grab

import (
	"fmt"
	"io/ioutil"
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
	arr := []string{
		thread[0].pid,
		thread[1].pid,
		thread[2].pid,
		thread[3].pid,
		thread[4].pid,
	}
	dumpTopThreadStack(arr, strconv.Itoa(int(p.Pid)))
}

func dumpTopThreadStack(treads []string, pid string) {
	fileName := pid + "_" + time.Now().Format(Layout) + ".txt"
	os.Create(fileName)
	cmd := "jstack -l " + pid + " > " + fileName
	command := exec.Command("bash", "-c", cmd)

	// 可能没权限.
	combinedOutput, err := command.CombinedOutput()
	fmt.Println(combinedOutput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 10 进制转成 16 进制 printf '%x\n' $threadId
	treadList := []string{}

	for _, tread := range treads {
		atto, err := strconv.Atoi(tread)
		if err != nil {
			return
		}
		formatInt := strconv.FormatInt(int64(atto), 16)
		treadList = append(treadList, formatInt)
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)

	split := strings.Split(string(content), "\r\n")

	newFile := pid + time.Now().Format(Layout) + ".dump"
	output, _ := os.Create(newFile)

	for _, line := range split {
		for _, threadNum := range treadList {
			if strings.Contains(line, threadNum) {
				output.WriteString(line + "\r\n")
			}
		}
	}

	log.Println("dump 成功, file ", newFile)

}
