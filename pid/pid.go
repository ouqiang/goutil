package pid

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// WriteToFile 写入pid到文件中
func WriteToFile(filename string, pids ...int) error {
	var pid int
	if len(pids) > 0 {
		pid = pids[0]
	} else {
		pid = os.Getpid()
	}

	if pid <= 0 {
		return errors.New("failed to get pid")
	}
	err := ioutil.WriteFile(filename, []byte(strconv.Itoa(pid)), 0644)

	return err
}

// ReadFromFile 从文件中读取pid
func ReadFromFile(filename string) (int, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(buf)))
	if pid <= 0 {
		return 0, errors.New("failed to parse pid")
	}

	return pid, nil
}

// IsRunning 通过pid判断进程是否运行中
func IsRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))

	return err == nil
}
