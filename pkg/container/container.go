package container

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

var CmdReexec *exec.Cmd
var CmdRun *exec.Cmd

func initFunction() {
	reexec.Register("nsInitialisation", nsInitialisation)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	fmt.Printf("\n>> namespace setup code begin...<<\n\n")

	//container file path
	setMount("/home/function/functionA")
	// host file path
	//setMount("/home/william/serverless/ubuntu-python/ubuntu-python")

	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		os.Exit(1)
	}

	set_cgroups()

	nsRun()
}

func initFunctionB() {
	reexec.Register("nsInitialisationB", nsInitialisationB)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisationB() {
	fmt.Printf("\n>> namespace setup code begin...<<\n\n")

	//container file path
	setMount("/home/function/functionB")
	// host file path
	//setMount("/home/william/serverless/ubuntu-python/ubuntu-python")

	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		os.Exit(1)
	}

	set_cgroups()

	nsRun()
}

func set_cgroups() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.MkdirAll(filepath.Join(pids, "ourContainer"), 0777)
	ioutil.WriteFile(filepath.Join(pids, "ourContainer/pids.max"), []byte("10"), 0777)
	//up here we limit the number of child processes to 10

	ioutil.WriteFile(filepath.Join(pids, "ourContainer/notify_on_release"), []byte("1"), 0777)

	ioutil.WriteFile(filepath.Join(pids, "ourContainer/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0777)
	// up here we write container PIDs to cgroup.procs
}

func setMount(root string) error {
	if err := syscall.Chroot(root); err != nil {
		return err
	}
	// 设置容器里面的当前工作目录
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return err
	}
	return nil
}

func nsRun() {
	CmdRun = exec.Command(os.Getenv("function"))

	CmdRun.Stdin = os.Stdin
	CmdRun.Stdout = os.Stdout
	CmdRun.Stderr = os.Stderr

	CmdRun.Env = []string{"PS1=-[container]- # "}

	if err := CmdRun.Run(); err != nil {
		fmt.Printf("Error running the %s command - %s\n", os.Getenv("function"), err)
		os.Exit(1)
	}

	syscall.Unmount("/proc", 0)

}

func StartFunction(functionPath string) {
	os.Setenv("function", functionPath)
	//os.Setenv("function", "/bin/bash")
	initFunction()
	CmdReexec = reexec.Command(append([]string{"nsInitialisation"},
		os.Getenv("function"))...)
	run(CmdReexec)
}

func StartFunctionB(functionPath string) {
	os.Setenv("function", functionPath)
	//os.Setenv("function", "/bin/bash")
	initFunctionB()
	cmd := reexec.Command(append([]string{"nsInitialisationB"},
		os.Getenv("function"))...)
	run(cmd)
}

func run(cmd *exec.Cmd) {

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			//syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", err)
		os.Exit(1)
	}
}
