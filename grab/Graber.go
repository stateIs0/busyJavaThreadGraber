package grab

import (
	"awesomeProject1/cpu"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Layout = "2006-01-02 15:04:05"

type Police struct {
	Pid       string
	tick      int64
	threshold float64
}

func NewPolice(Pid string, tick int64, threshold float64) *Police {
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
			// 监听 CPU, CPU 触发指标时, 就根据 top 获取 thread, 并抓取堆栈
			_, f, _ := cpu.Get(p.Pid)
			// 触发
			if f > p.threshold {
				p.parseThreadContentAndDump()
			}
		}
	}()
}

func (p *Police) parseThreadContentAndDump() {
	thread := getTopJavaThread(p.Pid)
	dumpTopThreadStack(thread, p.Pid)
}

func getTopJavaThread(pid string) []string {
	cmd := " top -Hp " + pid
	c := exec.Command("bash", "-c", cmd)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))
	return nil
}

func dumpTopThreadStack(treads []string, pid string) {
	fileName := pid + ".txt"
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

	output, _ := os.Open(pid + time.Now().String() + ".txt")

	for _, line := range split {
		for _, threadNum := range treadList {
			if strings.Contains(line, threadNum) {
				output.WriteString(line + "\r\n")
			}
		}
	}

}
