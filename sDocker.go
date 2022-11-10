// Let's first understand what namespaces are.

// network namespaces
// So basically, when you install Linux, by default the entire OS share the same routing table and the same IP address.
// The namespace forms a cluster of all global system resources which can only be used by the processes
// within the namespace, providing resource isolation.
//
// Docker containers use this technology to form their own cluster of resources which would be used only by that
// namespace, i.e. that container. so every container has its own IP address and work in isolation without facing
// resource sharing conflicts with other containers running on the same system.
//
// When the container is created using the -p flag, it maps the internal port 8080 to a higher external port 8080.
// So now the port 8080 of the host is mapped to the containers port 8080 and hence they are connected.

// PID namespaces
// where PID namespaces isolate the process ID number space,
// meaning that processes in different PID namespaces can have the same PID.
//
// PID namespaces allow containers to provide functionality such as suspending/resuming the set of processes in
// the container and migrating the container to a new host while the processes inside the container maintain the
// same PIDs.

// more details present here : https://man7.org/linux/man-pages/man2/unshare.2.html

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker             run image <cmd> <params>
// go run sDocker.go  run       <cmd> <params>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "init":
		child()
	default:
		panic("no such command")
	}
}

func run() {
	fmt.Printf("running %v\n", "init")
	cmd := exec.Command("/proc/self/exe", append([]string{"init"}, os.Args[2:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	fmt.Printf("running %v %d\n", os.Args[2:], os.Getpid())
	// syscall.SetHostname([]byte("Container"))
	syscall.Chroot("path/to/another/or copy of/a linux/filesys")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdin
	cmd.Run()
	syscall.Unmount("/proc", 0)
}

func Cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "sDocker"), 0755)
	if err != nil && os.IsExist(err) {
		panic(err)
	}
	f := fileErr(os.Create(filepath.Join(pids, "sDocker/pids.max")))
	must(f.Write([]byte("20")))
	// Removes the new cgroup in place after the container exits
	f = fileErr(os.Create(filepath.Join(pids, "sDocker/notify_on_release")))
	must(f.Write([]byte("1")))
	f = fileErr(os.Create(filepath.Join(pids, "sDocker/cgroup.procs")))
	must(f.Write([]byte(strconv.Itoa(os.Getpid()))))
}
func must(n int, err error) {
	if err != nil {
		panic(err)
	}
}
func fileErr(f *os.File, err error) *os.File {
	if err != nil {
		panic(err)
	}
	return f
}
