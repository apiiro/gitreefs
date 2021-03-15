// +build bench

package fs

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"
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
//	PRINT_MEMORY = false
//)

type mountBenchmarkTestSuite struct {
	suite.Suite
	clonesPath string
	clonePath  string
}

func TestMountBenchmarkTestSuite(t *testing.T) {
	suite.Run(t, new(mountBenchmarkTestSuite))
}

func (mntSuite *mountBenchmarkTestSuite) SetupTest() {
	logger.InitLoggers("logs/mount_bench_test-%v.log", "INFO", "-")

	var err error
	mntSuite.clonesPath, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("Cloning...")
	execCommand(mntSuite.clonesPath, []string{"git", "clone", REMOTE})
	mntSuite.clonePath = path.Join(mntSuite.clonesPath, REPO_NAME)
}

func (mntSuite *mountBenchmarkTestSuite) TearDownTest() {
	os.RemoveAll(mntSuite.clonesPath)
}

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

func (mntSuite *mountBenchmarkTestSuite) mount() (mountPoint string) {
	mountPoint, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("Mounting at %v", mountPoint)
	_, err = Mount(mntSuite.clonesPath, mountPoint, false)
	if err != nil {
		panic(err)
	}
	return mountPoint
}

func (mntSuite *mountBenchmarkTestSuite) unmount(mountPoint string) {
	logger.Info("Unmounting at %v", mountPoint)
	unmountErr := Unmount(mountPoint)
	os.RemoveAll(mountPoint)
	if unmountErr != nil {
		panic(unmountErr)
	}
	if PRINT_MEMORY {
		runtime.GC()
		printMemoryUsage()
	}
}

func walk(atPath string, readFile bool) {
	filepath.Walk(atPath, func(path string, info os.FileInfo, err error) error {
		if readFile && info != nil && !info.IsDir() {
			ioutil.ReadFile(path)
		}
		return nil
	})
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

func printMemoryUsage() {
	if !PRINT_MEMORY {
		return
	}
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	logger.Info("\tHeap = %v GB ", float64(stats.Alloc)/(1024*1024*1024))
}

func (mntSuite *mountBenchmarkTestSuite) TestBenchmarkWalkVirtualFileSystem() {

	var times []op

	logger.Info("Checking out %v", COMMIT)

	times = append(times, op{
		description: "git checkout",
		seconds: timed(func() {
			execCommand(mntSuite.clonePath, []string{"git", "checkout", COMMIT})
		}),
	})

	logger.Info("Walking %v", mntSuite.clonePath)

	times = append(times, op{
		description: "walk clone (no files)",
		seconds: timed(func() {
			walk(mntSuite.clonePath, false)
		}),
	})

	times = append(times, op{
		description: "walk clone (with files)",
		seconds: timed(func() {
			walk(mntSuite.clonePath, true)
		}),
	})

	logger.Info("Archiving %v", COMMIT)

	contentDirPath := path.Join(mntSuite.clonesPath, "archive")
	os.MkdirAll(contentDirPath, 0777)
	times = append(times, op{
		description: "archive and decompress",
		seconds: timed(func() {
			tarFilePath := path.Join(mntSuite.clonesPath, "archive.tar")
			execCommand(mntSuite.clonePath, []string{"git", "archive", COMMIT, "--format=tar", "--output", tarFilePath})
			execCommand(contentDirPath, []string{"tar", "-xf", tarFilePath})
		}),
	})

	times = append(times, op{
		description: "walk archive (no files)",
		seconds: timed(func() {
			walk(mntSuite.clonePath, false)
		}),
	})

	times = append(times, op{
		description: "walk archive (with files)",
		seconds: timed(func() {
			walk(mntSuite.clonePath, true)
		}),
	})

	mountPoint := mntSuite.mount()
	defer mntSuite.unmount(mountPoint)

	printMemoryUsage()

	virtualPath := path.Join(mountPoint, REPO_NAME, COMMIT)
	times = append(times, op{
		description: "walk virtual (no files)",
		seconds: timed(func() {
			walk(virtualPath, false)
		}),
	})

	printMemoryUsage()

	times = append(times, op{
		description: "walk virtual #2 (no files)",
		seconds: timed(func() {
			walk(virtualPath, false)
		}),
	})

	printMemoryUsage()

	mountPoint = mntSuite.mount()
	defer mntSuite.unmount(mountPoint)

	printMemoryUsage()

	virtualPath = path.Join(mountPoint, REPO_NAME, COMMIT)
	times = append(times, op{
		description: "walk virtual (with files)",
		seconds: timed(func() {
			walk(virtualPath, true)
		}),
	})

	printMemoryUsage()

	times = append(times, op{
		description: "walk virtual #2 (with files)",
		seconds: timed(func() {
			walk(virtualPath, true)
		}),
	})

	printMemoryUsage()

	for _, timedOp := range times {
		logger.Info("%v - %v sec", timedOp.description, timedOp.seconds)
	}
}
