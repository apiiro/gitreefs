package testutils

import (
	"gitreefs/core/logger"
	"os/exec"
)

func runCommand(cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("command [%v] failed to execute: %v\nS%v", cmd.Args, err, string(output))
		panic(err)
	}
}

func ExecCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	runCommand(cmd)
}

func ExecCommandWithDir(dir string, name string, arg ...string) {
	executable, err := exec.LookPath(name)
	if err != nil {
		panic(err)
	}
	cmd := &exec.Cmd{
		Path: executable,
		Args: append([]string{name}, arg...),
		Dir:  dir,
	}
	runCommand(cmd)
}

func NfsMount(mountPoint string) {
	ExecCommand("mount", "-o", "port=2049,mountport=2049,nfsvers=3,noacl,tcp", "-t", "nfs", "localhost:/", mountPoint)
}
