package main

import (
	"fmt"
	"flag"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"io/ioutil"
	"strconv"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}

}

func run(args []string) {
	fmt.Printf("Running %v \n", args)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = & syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child(args []string) {
	fmt.Printf("Running from proc in namespace %v \n", args)
	
	err := cgroup()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = syscall.Sethostname([]byte("container"))
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

	if err = syscall.Unmount("proc", 0); err != nil {
		panic(err)
	}
}

func cgroup() (error){
	cgroup := "/sys/fs/cgroup"
	pids := filepath.Join(cgroup, "pids")
	os.Mkdir(filepath.Join(pids, "container"), 0755)

	err := ioutil.WriteFile(filepath.Join(pids, "container/pids.max"), []byte("20"), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(pids, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		return err
	}

	mem := filepath.Join(cgroup, "memory")
	os.Mkdir(filepath.Join(mem, "container"), 0755)

	err = ioutil.WriteFile(filepath.Join(mem, "container/memory.limit_in_bytes"), []byte("1000000"), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(mem, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		return nil
	}

	cpu := filepath.Join(cgroup, "cpu")
	os.Mkdir(filepath.Join(cpu, "container"), 0755)

	err = ioutil.WriteFile(filepath.Join(cpu, "container/cpu.shares"), []byte("512"), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(cpu, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(pids, "container/notify_on_release"), []byte("1"), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(mem, "container/notify_on_release"), []byte("1"), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(cpu, "container/notify_on_release"), []byte("1"), 0700)
	if err != nil {
		return err
	}

	return nil
}
