package main

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/logger"
	testutils "gitreefs/test_utils"
	"io/ioutil"
	"os"
	"testing"
)

type nfsTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
	dataPath   string
}

func TestNfsTestSuite(t *testing.T) {
	logger.InitLoggers("logs/nfs_test-%v-%v.log", "INFO", "-")
	suite.Run(t, new(nfsTestSuite))
}

func (nfsSuite *nfsTestSuite) SetupTest() {

	var err error

	nfsSuite.clonesPath = testutils.SetupClones()

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

	nfsSuite.mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	testutils.NfsMount(nfsSuite.mountPoint)
}

func (nfsSuite *nfsTestSuite) TearDownTest() {
	logger.Info("Unmounting")
	os.RemoveAll(nfsSuite.clonesPath)
	defer os.RemoveAll(nfsSuite.mountPoint)
	testutils.ExecCommand("umount", nfsSuite.mountPoint)
	os.RemoveAll(nfsSuite.dataPath)
}

func (nfsSuite *nfsTestSuite) TestWalkFileSystem() {
	testutils.WalkFileSystem(&nfsSuite.Suite, nfsSuite.mountPoint, false)
}
