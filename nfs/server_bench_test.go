// +build bench

package main

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/logger"
	testutils "gitreefs/test_utils"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

type nfsBenchmarkTestSuite struct {
	suite.Suite
	clonesPath string
	clonePath  string
	dataPath   string
}

func TestNfsBenchmarkTestSuite(t *testing.T) {
	logger.InitLoggers("logs/nfs_bench_test-%v-%v.log", "INFO", "-")
	suite.Run(t, new(nfsBenchmarkTestSuite))
}

func (nfsSuite *nfsBenchmarkTestSuite) SetupTest() {

	var err error

	nfsSuite.clonesPath, nfsSuite.clonePath = testutils.SetupClone()

	nfsSuite.dataPath, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(nfsSuite.dataPath, 0777)
	if err != nil {
		panic(err)
	}

	logger.Info("Serving")
	go func() {
		err = Serve(nfsSuite.clonesPath, "localhost", "2049", "data")
		panic(err)
	}()
}

func (nfsSuite *nfsBenchmarkTestSuite) mount() (mountPoint string) {
	var err error
	mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("Mounting at %v", mountPoint)
	execCommand("mount", "-o", "port=2049,mountport=2049", "-t", "nfs", "localhost:/", mountPoint)
	return mountPoint
}

func (mntSuite *nfsBenchmarkTestSuite) unmount(mountPoint string) {
	defer os.RemoveAll(mountPoint)
	execCommand("umount", mountPoint)
	if testutils.PRINT_MEMORY {
		runtime.GC()
		testutils.PrintMemoryUsage()
	}
}

func (nfsSuite *nfsBenchmarkTestSuite) TearDownTest() {
	logger.Info("Unmounting")
	os.RemoveAll(nfsSuite.clonesPath)
	os.RemoveAll(nfsSuite.dataPath)
}

func (nfsSuite *nfsBenchmarkTestSuite) TestBenchmarkWalkVirtualFileSystem() {
	testutils.BenchmarkWalkVirtualFileSystem(
		nfsSuite.clonesPath,
		nfsSuite.clonePath,
		nfsSuite.mount,
		nfsSuite.unmount,
	)
}
