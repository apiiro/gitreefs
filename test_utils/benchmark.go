package testutils

import (
	"gitreefs/core/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

const (
	REMOTE       = "https://github.com/apiirolab/elasticsearch.git"
	REPO_NAME    = "elasticsearch"
	COMMIT       = "b926bf0"
	PRINT_MEMORY = true
)

//const (
//	REMOTE       = "https://github.com/apiirolab/EVO-Exchange-BE-2019"
//	REPO_NAME    = "EVO-Exchange-BE-2019"
//	COMMIT       = "c47980a"
//	PRINT_MEMORY = true
//)

func execCommand(workingDirectory string, args []string) {
	executablePath, err := exec.LookPath(args[0])
	if err != nil {
		panic(err)
	}
	cmd := &exec.Cmd{
		Path: executablePath,
		Args: args,
		Dir:  workingDirectory,
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func SetupClone() (clonesPath string, clonePath string) {
	clonesPath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("Cloning...")
	execCommand(clonesPath, []string{"git", "clone", REMOTE})
	return clonesPath, path.Join(clonesPath, REPO_NAME)
}

type Op func()

func timed(op Op) (elapsedSeconds float64) {
	start := time.Now()
	op()
	return time.Since(start).Seconds()
}

type op struct {
	description string
	seconds     float64
}

func PrintMemoryUsage() {
	if !PRINT_MEMORY {
		return
	}
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	logger.Info("\tHeap = %v GB ", float64(stats.Alloc)/(1024*1024*1024))
}

func walk(atPath string, readFile bool) {
	filepath.Walk(atPath, func(path string, info os.FileInfo, err error) error {
		if readFile && info != nil && !info.IsDir() {
			ioutil.ReadFile(path)
		}
		return nil
	})
}

func BenchmarkWalkVirtualFileSystem(
	clonesPath string,
	clonePath string,
	mount func() string,
	unmount func(string),
) {

	var times []op

	logger.Info("Checking out %v", COMMIT)

	times = append(times, op{
		description: "git checkout",
		seconds: timed(func() {
			execCommand(clonePath, []string{"git", "checkout", COMMIT})
		}),
	})

	logger.Info("Walking %v", clonePath)

	times = append(times, op{
		description: "walk clone (no files)",
		seconds: timed(func() {
			walk(clonePath, false)
		}),
	})

	times = append(times, op{
		description: "walk clone (with files)",
		seconds: timed(func() {
			walk(clonePath, true)
		}),
	})

	logger.Info("Archiving %v", COMMIT)

	contentDirPath := path.Join(clonesPath, "archive")
	os.MkdirAll(contentDirPath, 0777)
	times = append(times, op{
		description: "archive and decompress",
		seconds: timed(func() {
			tarFilePath := path.Join(clonesPath, "archive.tar")
			execCommand(clonePath, []string{"git", "archive", COMMIT, "--format=tar", "--output", tarFilePath})
			execCommand(contentDirPath, []string{"tar", "-xf", tarFilePath})
		}),
	})

	times = append(times, op{
		description: "walk archive (no files)",
		seconds: timed(func() {
			walk(clonePath, false)
		}),
	})

	times = append(times, op{
		description: "walk archive (with files)",
		seconds: timed(func() {
			walk(clonePath, true)
		}),
	})

	mountPoint := mount()
	defer unmount(mountPoint)

	PrintMemoryUsage()

	virtualPath := path.Join(mountPoint, REPO_NAME, COMMIT)
	times = append(times, op{
		description: "walk virtual (no files)",
		seconds: timed(func() {
			walk(virtualPath, false)
		}),
	})

	PrintMemoryUsage()

	times = append(times, op{
		description: "walk virtual #2 (no files)",
		seconds: timed(func() {
			walk(virtualPath, false)
		}),
	})

	PrintMemoryUsage()

	mountPoint = mount()
	defer unmount(mountPoint)

	PrintMemoryUsage()

	virtualPath = path.Join(mountPoint, REPO_NAME, COMMIT)
	times = append(times, op{
		description: "walk virtual (with files)",
		seconds: timed(func() {
			walk(virtualPath, true)
		}),
	})

	PrintMemoryUsage()

	times = append(times, op{
		description: "walk virtual #2 (with files)",
		seconds: timed(func() {
			walk(virtualPath, true)
		}),
	})

	PrintMemoryUsage()

	for _, timedOp := range times {
		logger.Info("%v - %v sec", timedOp.description, timedOp.seconds)
	}
}
