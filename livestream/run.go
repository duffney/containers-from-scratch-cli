package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type limits struct {
	Limits []limit
}

type limit struct {
	name  string
	path  string
	param string
	value []byte
}

func run(arg []string) {
	fmt.Printf("running container with arg: %v\n", arg)

	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cgroup()
	if err != nil {
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child(arg []string) {
	fmt.Printf("running from proc in namespace: %v\n", arg)

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := syscall.Sethostname([]byte("container"))
	if err != nil {
		panic(err)
	}

	if err = syscall.Chroot("ubuntu-rootfs/"); err != nil {
		panic(err)
	}

	if err = syscall.Chdir("/"); err != nil {
		panic(err)
	}

	if err = syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func cgroup() error {
	cgrouplimits := limits{
		[]limit{
			{
				"pids",
				"/sys/fs/cgroup/pids/container",
				"pids.max",
				[]byte("20"),
			},
			{
				"memory",
				"/sys/fs/cgroup/memory/container",
				"memory.limit_in_bytes",
				[]byte("10000000"),
			},
			{
				"cpu",
				"/sys/fs/cgroup/cpu/container",
				"cpu.shares",
				[]byte("512"),
			},
		},
	}

	for _, l := range cgrouplimits.Limits {
		fmt.Printf("setting %s to %s\n", l.path, l.value)

		os.Mkdir(l.path, 0755)

		err := ioutil.WriteFile(filepath.Join(l.path, l.param), l.value, 0700)
		if err != nil {
			fmt.Printf("error setting %s to %s\n", l.path, l.value)
			return err
		}

		err = ioutil.WriteFile(filepath.Join(l.path, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
		if err != nil {
			fmt.Printf("error setting %s to %s\n", l.path, l.value)
			return err
		}

	}

	return nil
}
