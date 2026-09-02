package main

import (
	"os"
	"os/exec"
	"syscall"
)

const rootfs = "/home/hadiali/Projects/myDocker/rootfs"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
			panic(err)
		}
		if err := syscall.Mkdir(rootfs+"/.oldroot", 0755); err != nil && !os.IsExist(err) {
			panic(err)
		}
		if err := syscall.Mount(rootfs, rootfs, "", syscall.MS_BIND, ""); err != nil {
			panic(err)
		}
		if err := syscall.PivotRoot(rootfs, rootfs+"/.oldroot"); err != nil {
			panic(err)
		}
		if err := syscall.Chdir("/"); err != nil {
			panic(err)
		}
		if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
			panic(err)
		}
		if err := syscall.Unmount("/.oldroot", syscall.MNT_DETACH); err != nil {
			panic(err)
		}
		if err := syscall.Exec("/bin/sh", []string{"/bin/sh"}, os.Environ()); err != nil {
			panic(err)
		}
	} else {
		cmd := exec.Command("/proc/self/exe", "child")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
			Setsid:     true,
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}