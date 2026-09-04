package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const rootfs = "/home/hadiali/Projects/myDocker/rootfs"
func generateID() string {
    b := make([]byte, 4)
    rand.Read(b)
    return hex.EncodeToString(b)
}
func main() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
			panic(err)
		}
if err := syscall.Mkdir(rootfs+"/proc", 0755); err != nil && !os.IsExist(err) {
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
		containerID:= generateID()
		if err := os.MkdirAll("/sys/fs/cgroup/mydocker/"+containerID, 0755); err != nil {
    panic(err)
}

if err := os.WriteFile("/sys/fs/cgroup/mydocker/"+containerID+"/memory.max", []byte("100000000"), 0644); err != nil {
    panic(err)
}
if err := os.WriteFile("/sys/fs/cgroup/mydocker/"+containerID+"/cpu.max", []byte("50000 100000"), 0644); err != nil {
    panic(err)
}
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Start(); err != nil {
    panic(err)
}

pid := cmd.Process.Pid
if err := os.WriteFile("/sys/fs/cgroup/mydocker/"+containerID+"/cgroup.procs", []byte(strconv.Itoa(pid)), 0644); err != nil {
    panic(err)
}

fmt.Println("Container ID:", containerID)
cmd.Wait()
os.RemoveAll("/sys/fs/cgroup/mydocker/" + containerID)
}
	
}