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

func run(args []string) {
	fmt.Printf("Running %v \n", args)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
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
				[]byte("1000000"),
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
		fmt.Println(filepath.Join(l.path, l.param))
		//os.Mkdir
		os.Mkdir(l.path, 0755)
		//Create cgroup limit
		err := ioutil.WriteFile(filepath.Join(l.path, l.param), l.value, 0700)
		if err != nil {
			return err
		}
		//ioutil.WriteFile (add proc to cgroup)
		err = ioutil.WriteFile(filepath.Join(l.path, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
		if err != nil {
			return err
		}
		//ioutil.WriteFile notify_on-release
		err = ioutil.WriteFile(filepath.Join(l.path, "notify_on_release"), []byte("1"), 0700)
		if err != nil {
			return err
		}
	}

	//pids := filepath.Join(cgroup, "pids")
	//os.Mkdir(filepath.Join(pids, "container"), 0755)

	//err := ioutil.WriteFile(filepath.Join(pids, "container/pids.max"), []byte("20"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(pids, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return err
	//}

	//mem := filepath.Join(cgroup, "memory")
	//os.Mkdir(filepath.Join(mem, "container"), 0755)

	//err = ioutil.WriteFile(filepath.Join(mem, "container/memory.limit_in_bytes"), []byte("1000000"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(mem, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return nil
	//}

	//cpu := filepath.Join(cgroup, "cpu")
	//os.Mkdir(filepath.Join(cpu, "container"), 0755)

	//err = ioutil.WriteFile(filepath.Join(cpu, "container/cpu.shares"), []byte("512"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(cpu, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(pids, "container/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(mem, "container/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(cpu, "container/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	return nil
}
