package grab

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	s := "6873 (a.out) R 6723 6873 6723 34819 6873 8388608 77 0 0 0 41958 31 0 0 25 0 3 0 5882654 1409024 56 4294967295 134512640 134513720 3215579040 0 2097798 0 0 0 0 0 0 0 17 0 0 0"

	split := strings.Split(s, " ")
	// 该任务在用户态运行的时间，单位为jiffies
	fmt.Println(split[13])
	// 该任务在核心态运行的时间，单位为jiffies
	fmt.Println(split[14])
}
