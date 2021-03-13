package fs

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"
)

type mountBenchmarkTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
}

func TestMountBenchmarkTestSuite(t *testing.T) {
	logger.InitLoggers("logs/mount_test-%v.log", "INFO", "")
	suite.Run(t, new(mountBenchmarkTestSuite))
}

func (mntSuite *mountBenchmarkTestSuite) SetupTest() {
	var err error
	mntSuite.clonesPath, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	mntSuite.mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	var gitExecutable string
	gitExecutable, err = exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	logger.Info("Cloning...")
	cloneCmd := &exec.Cmd{
		Path: gitExecutable,
		Args: []string{"git", "clone", "https://github.com/apiirolab/elasticsearch.git"},
		Dir:  mntSuite.clonesPath,
	}
	err = cloneCmd.Start()
	if err != nil {
		panic(err)
	}
	err = cloneCmd.Wait()
	if err != nil {
		panic(err)
	}

	logger.Info("Checking out b926bf0")
	cloneCmd = &exec.Cmd{
		Path: gitExecutable,
		Args: []string{"git", "checkout", "b926bf0"},
		Dir:  path.Join(mntSuite.clonesPath, "elasticsearch"),
	}
	err = cloneCmd.Start()
	if err != nil {
		panic(err)
	}
	err = cloneCmd.Wait()
	if err != nil {
		panic(err)
	}

	logger.Info("Mounting")
	_, err = Mount(mntSuite.clonesPath, mntSuite.mountPoint, false)
	if err != nil {
		panic(err)
	}
}

func (mntSuite *mountBenchmarkTestSuite) TearDownTest() {
	logger.Info("Unmounting")
	unmountErr := Unmount(mntSuite.mountPoint)
	os.RemoveAll(mntSuite.mountPoint)
	os.RemoveAll(mntSuite.clonesPath)
	if unmountErr != nil {
		panic(unmountErr)
	}
}

func walk(atPath string) {
	filepath.Walk(atPath, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			ioutil.ReadFile(path)
		}
		return nil
	})
}

func (mntSuite *mountBenchmarkTestSuite) TestBenchmarkWalkVirtualFileSystem() {
	start := time.Now()
	walk(path.Join(mntSuite.clonesPath, "elasticsearch"))
	elapsedPhysicalSec := time.Since(start).Seconds()
	logger.Info("Walking physical fs: %v sec", elapsedPhysicalSec)

	start = time.Now()
	walk(path.Join(mntSuite.mountPoint, "elasticsearch", "b926bf0"))
	elapsedVirtualSec := time.Since(start).Seconds()
	logger.Info("Walking virtual fs: %v sec", elapsedVirtualSec)

	mntSuite.Less(elapsedVirtualSec, elapsedPhysicalSec)
}
