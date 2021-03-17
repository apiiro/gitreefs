package fuseserver

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/logger"
	testutils "gitreefs/test_utils"
	"io/ioutil"
	"os"
	"testing"
)

type mountTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
}

func TestMountTestSuite(t *testing.T) {
	logger.InitLoggers("logs/mount_test-%v-%v.log", "INFO", "-")
	suite.Run(t, new(mountTestSuite))
}

func (mntSuite *mountTestSuite) SetupTest() {

	mntSuite.clonesPath = testutils.SetupClones()

	var err error
	mntSuite.mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("Mounting")
	_, err = Mount(mntSuite.clonesPath, mntSuite.mountPoint, false)
	if err != nil {
		panic(err)
	}
}

func (mntSuite *mountTestSuite) TearDownTest() {
	logger.Info("Unmounting")
	unmountErr := Unmount(mntSuite.mountPoint)
	os.RemoveAll(mntSuite.mountPoint)
	os.RemoveAll(mntSuite.clonesPath)
	if unmountErr != nil {
		panic(unmountErr)
	}
}

func (mntSuite *mountTestSuite) TestWalkFileSystem() {
	testutils.WalkFileSystem(&mntSuite.Suite, mntSuite.mountPoint, true)
}
