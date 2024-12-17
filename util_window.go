// +build windows

package goutil

import (
        "github.com/ouqiang/goutil/slice"
        "os"
        "os/exec"
        "syscall"
        "io"
)

// Daemon  守护进程
func Daemon(w io.Writer) {
        execFile, err := os.Executable()
        if err != nil {
                panic(err)
        }
        args := slice.Remove(os.Args[1:], "-d")
        cmd := exec.Command(execFile, args...)
        cmd.Stdin = nil
        if w != nil {
                cmd.Stdout = w
                cmd.Stderr = w
        } else {
                cmd.Stdout = nil
                cmd.Stderr = nil
        }
        cmd.SysProcAttr = &syscall.SysProcAttr{
                HideWindow: true,
        }
        err = cmd.Start()
        if err != nil {
                panic(err)
        }
        os.Exit(0)
}
