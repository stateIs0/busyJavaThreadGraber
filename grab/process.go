package grab

import (
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetThreads(pid int32, threshold float64) []string {

	newProcess, err := process.NewProcess(pid)
	if err != nil {
		return nil
	}

	percent, err := newProcess.Percent(3 * time.Second)
	if err != nil {
		log.Println(err)
		return nil
	}
	name, err := newProcess.Name()
	if err != nil {
		return nil
	}
	log.Println("pid ", name, " rootProcess percent = ", percent, ", threshold=", threshold)
	if percent == 0 {
		return nil
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	if percent < threshold {
		return nil
	}

	cmd := "ps -T -p " + strconv.Itoa(int(pid))
	c := exec.Command("bash", "-c", cmd)

	output, err := c.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil
	}
	threads := []string{}
	log.Println("len --->>" + strconv.Itoa(len(output)))
	if len(string(output)) > 0 {

		str := string(output)
		split := strings.Split(str, "\r\n")
		for i, line := range split {
			if i == 0 {
				continue
			}
			lineArr := strings.Split(line, " ")
			subThread := lineArr[0]
			threads = append(threads, subThread)
		}
	}

	return threads

}
