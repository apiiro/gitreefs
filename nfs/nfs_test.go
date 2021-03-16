package main

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/logger"
	testutils "gitreefs/test_utils"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type nfsTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
}

func TestNfsTestSuite(t *testing.T) {
	logger.InitLoggers("logs/nfs_test-%v.log", "INFO", "-")
	suite.Run(t, new(nfsTestSuite))
}

func execCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func (nfsSuite *nfsTestSuite) SetupTest() {

	var err error
	nfsSuite.clonesPath = testutils.SetupClones()

	logger.Info("Serving")
	go func() {
		err = Serve(nfsSuite.clonesPath, "localhost","2049", "data")
		panic(err)
	}()

	nfsSuite.mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	execCommand("mount", "-o", "port=2049,mountport=2049", "-t", "nfs", "localhost:/", nfsSuite.mountPoint)
}

func (nfsSuite *nfsTestSuite) TearDownTest() {
	logger.Info("Unmounting")
	os.RemoveAll(nfsSuite.clonesPath)
	defer os.RemoveAll(nfsSuite.mountPoint)
	execCommand("umount", nfsSuite.mountPoint)
}

func (nfsSuite *nfsTestSuite) TestWalkFileSystem() {
	testutils.WalkFileSystem(&nfsSuite.Suite, nfsSuite.mountPoint)
}
