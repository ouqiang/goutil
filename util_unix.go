// +build !windows

package goutil


import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/ouqiang/goutil/slice"
)

// Daemon  守护进程
func Daemon(w io.Writer) {
	execFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	args := slice.Remove(os.Args[1:], "-d")
	cmd := exec.Command(execFile, args...)
	cmd.Stdin = nil
	if w != nil {
		if _, ok := w.(*os.File); !ok {
			panic("writer requires a file descriptor")
		}

		cmd.Stdout = w
		cmd.Stderr = w
	} else {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
