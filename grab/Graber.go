package grab

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Police struct {
	Pid string
}

func NewPolice(Pid string) *Police {
	return &Police{
		Pid: Pid,
	}
}

func (p *Police) Start() {
	go func() {
		for {
			// 监听 CPU, CPU 触发指标时, 就根据 top 获取 thread, 并抓取堆栈
			p.process()
		}
	}()
}

func (p *Police) process() {
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
	exec.Command("bash", "-c", cmd)

	// 10 进制转成 16 进制
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

	for _, line := range split {
		for _, threadNum := range treadList {
			if strings.Contains(line, threadNum) {
				fmt.Println(line)
			}
		}
	}

}
