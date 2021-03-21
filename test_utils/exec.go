package testutils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func startAndWait(cmd *exec.Cmd) {
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func ExecCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	startAndWait(cmd)
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
	startAndWait(cmd)
}

func NfsMount(mountPoint string) {
	switch runtime.GOOS {
	case "linux":
		ExecCommand("mount", "-o", "port=2049,mountport=2049", "-t", "nfs", "localhost:/", mountPoint)
	case "darwin":
		ExecCommand("mount", "-o", "port=2049,mountport=2049,nfsvers=3,noacl,tcp", "-t", "nfs", "localhost:/", mountPoint)
	default:
		panic(fmt.Sprintf("unsupported os: %v", runtime.GOOS))
	}
}
