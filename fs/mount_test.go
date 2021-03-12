package fs

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/git"
	"gitreefs/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

type mountTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
}

func TestMountTestSuite(t *testing.T) {
	logger.InitLoggers("logs/mount_test-%v.log", "INFO", "")
	suite.Run(t, new(mountTestSuite))
}

func (mntSuite *mountTestSuite) SetupTest() {
	var err error
	mntSuite.clonesPath, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	mntSuite.mountPoint, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	remotes := []string{
		"https://github.com/apiirolab/dc-heacth.git",
		"https://github.com/apiirolab/EVO-Exchange-BE-2019.git",
	}
	var gitExecutable string
	gitExecutable, err = exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	for _, remote := range remotes {
		logger.Info("Cloning %v", remote)
		cloneCmd := &exec.Cmd{
			Path: gitExecutable,
			Args: []string{"git", "clone", "--no-checkout", remote},
			Dir:  mntSuite.clonesPath,
		}
		err = cloneCmd.Start()
		if err != nil {
			panic(err)
		}
		err = cloneCmd.Wait()
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

func (mntSuite *mountTestSuite) assertWalk(
	commitishPath string,
	expectedDirsCount int,
	expectedFilesCount int,
	expectedMinFileSize int,
	expectedMaxFileSize int,
) {
	dirsCount, filesCount := 0, 0
	minFileSize, maxFileSize := int(git.MaxFileSizeBytes), 0
	filepath.Walk(commitishPath, func(path string, info os.FileInfo, err error) error {
		mntSuite.NotNil(info, "missing info for %v", path)
		if info == nil {
			return nil
		}
		if info.IsDir() {
			dirsCount++
			mntSuite.EqualValues(0, info.Size())
		} else {
			filesCount++
			mntSuite.True(info.Size() >= 0, "file at %v has invalid size", path)
			mntSuite.True(info.Mode().IsRegular())
			content, err := ioutil.ReadFile(path)
			mntSuite.Nil(err, "failed to read file at %v", path)
			fileSizeFromRead := len(content)
			mntSuite.EqualValues(info.Size(), fileSizeFromRead, "read different file size for %v", path)
			if fileSizeFromRead > maxFileSize {
				maxFileSize = fileSizeFromRead
			}
			if fileSizeFromRead < minFileSize {
				minFileSize = fileSizeFromRead
			}
		}
		return nil
	})
	mntSuite.EqualValues(expectedDirsCount, dirsCount, "unexpected dirs count")
	mntSuite.EqualValues(expectedFilesCount, filesCount, "unexpected files count")
	mntSuite.EqualValues(expectedMinFileSize, minFileSize, "unexpected min file size")
	mntSuite.EqualValues(expectedMaxFileSize, maxFileSize, "unexpected max file size")
}

func (mntSuite *mountTestSuite) TestWalkVirtualFileSystem() {
	files, err := ioutil.ReadDir(mntSuite.mountPoint)
	mntSuite.Nil(err)
	mntSuite.Empty(files)

	files, err = ioutil.ReadDir(path.Join(mntSuite.mountPoint, "dc-heacth"))
	mntSuite.Nil(err)
	mntSuite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(mntSuite.mountPoint, "dc-heacth", "wat"))
	mntSuite.NotNil(err)
	mntSuite.True(os.IsNotExist(err))

	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "dc-heacth", "2ca742044ba451d00c6854a465fdd4280d9ad1f5"),
		28, 181,
		215, 47804,
		)
	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "dc-heacth", "2ca7420"),
		28, 181,
		215, 47804,
		)
	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "dc-heacth", "master"),
		28, 183,
		215, 47814,
		)

	files, err = ioutil.ReadDir(path.Join(mntSuite.mountPoint, "EVO-Exchange-BE-2019"))
	mntSuite.Nil(err)
	mntSuite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(mntSuite.mountPoint, "EVO-Exchange-BE-2019", "wat"))
	mntSuite.NotNil(err)
	mntSuite.True(os.IsNotExist(err))

	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f49b65a559dee369dba2360cc92cb01cf5"),
		31, 79,
		33, 15969,
	)
	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f"),
		31, 79,
		33, 15969,
	)
	mntSuite.assertWalk(
		path.Join(mntSuite.mountPoint, "EVO-Exchange-BE-2019", "master"),
		75, 256,
		0, 677127,
	)
}
