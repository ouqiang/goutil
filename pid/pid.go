/**
 * Created by GoLand.
 * User: ouqiang<ouqiang@4399inc.com>
 * Date: 2019/7/30
 * Time: 9:25
 */

package pid

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

// WriteToFile 写入pid到文件中
func WriteToFile(filename string) error {
	pid := os.Getpid()
	if pid <= 0 {
		return errors.New("failed to get pid")
	}
	err := ioutil.WriteFile(filename, []byte(strconv.Itoa(pid)), 0644)

	return err
}

// ReadFromFile 从文件中读取pid
func ReadFromFile(filename string) (int64, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	pid, _ := strconv.Atoi(string(buf))
	if pid <= 0 {
		return 0, errors.New("failed to parse pid")
	}

	return int64(pid), nil
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
