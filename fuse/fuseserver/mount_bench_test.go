// +build bench

package fuseserver

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/logger"
	testutils "gitreefs/test_utils"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

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

	mntSuite.clonesPath, mntSuite.clonePath = testutils.SetupClone()
}

func (mntSuite *mountBenchmarkTestSuite) TearDownTest() {
	os.RemoveAll(mntSuite.clonesPath)
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
	if testutils.PRINT_MEMORY {
		runtime.GC()
		testutils.PrintMemoryUsage()
	}
}

func (mntSuite *mountBenchmarkTestSuite) TestBenchmarkWalkVirtualFileSystem() {
	testutils.BenchmarkWalkVirtualFileSystem(
		mntSuite.clonesPath,
		mntSuite.clonePath,
		mntSuite.mount,
		mntSuite.unmount,
	)
}
