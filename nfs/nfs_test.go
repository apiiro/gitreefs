package main

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/git"
	"gitreefs/core/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

type nfsTestSuite struct {
	suite.Suite
	clonesPath string
	mountPoint string
}

func TestMountTestSuite(t *testing.T) {
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
	nfsSuite.clonesPath, err = ioutil.TempDir("", "")
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
			Dir:  nfsSuite.clonesPath,
		}
		err = cloneCmd.Start()
		if err != nil {
			panic(err)
		}
		err = cloneCmd.Wait()
	}

	logger.Info("Serving")
	go func() {
		err = Serve(nfsSuite.clonesPath, "localhost","2049")
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

func (nfsSuite *nfsTestSuite) assertWalk(
	commitishPath string,
	expectedDirsCount int,
	expectedFilesCount int,
	expectedMinFileSize int,
	expectedMaxFileSize int,
) {
	dirsCount, filesCount := 0, 0
	minFileSize, maxFileSize := int(git.MaxFileSizeBytes), 0
	filepath.Walk(commitishPath, func(path string, info os.FileInfo, err error) error {
		nfsSuite.NotNil(info, "missing info for %v", path)
		if info == nil {
			return nil
		}
		if info.IsDir() {
			dirsCount++
			nfsSuite.EqualValues(0, info.Size())
		} else {
			filesCount++
			nfsSuite.True(info.Size() >= 0, "file at %v has invalid size", path)
			nfsSuite.True(info.Mode().IsRegular())
			content, err := ioutil.ReadFile(path)
			nfsSuite.Nil(err, "failed to read file at %v", path)
			fileSizeFromRead := len(content)
			nfsSuite.EqualValues(info.Size(), fileSizeFromRead, "read different file size for %v", path)
			if fileSizeFromRead > maxFileSize {
				maxFileSize = fileSizeFromRead
			}
			if fileSizeFromRead < minFileSize {
				minFileSize = fileSizeFromRead
			}
		}
		return nil
	})
	nfsSuite.EqualValues(expectedDirsCount, dirsCount, "unexpected dirs count")
	nfsSuite.EqualValues(expectedFilesCount, filesCount, "unexpected files count")
	nfsSuite.EqualValues(expectedMinFileSize, minFileSize, "unexpected min file size")
	nfsSuite.EqualValues(expectedMaxFileSize, maxFileSize, "unexpected max file size")
}

func (nfsSuite *nfsTestSuite) TestWalkVirtualFileSystem() {
	files, err := ioutil.ReadDir(nfsSuite.mountPoint)
	nfsSuite.Nil(err)
	nfsSuite.Empty(files)

	files, err = ioutil.ReadDir(path.Join(nfsSuite.mountPoint, "dc-heacth"))
	nfsSuite.Nil(err)
	nfsSuite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(nfsSuite.mountPoint, "dc-heacth", "wat"))
	nfsSuite.NotNil(err)
	nfsSuite.True(os.IsNotExist(err))

	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "dc-heacth", "2ca742044ba451d00c6854a465fdd4280d9ad1f5"),
		28, 181,
		215, 47804,
	)
	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "dc-heacth", "2ca7420"),
		28, 181,
		215, 47804,
	)
	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "dc-heacth", "master"),
		28, 183,
		215, 47814,
	)

	files, err = ioutil.ReadDir(path.Join(nfsSuite.mountPoint, "EVO-Exchange-BE-2019"))
	nfsSuite.Nil(err)
	nfsSuite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(nfsSuite.mountPoint, "EVO-Exchange-BE-2019", "wat"))
	nfsSuite.NotNil(err)
	nfsSuite.True(os.IsNotExist(err))

	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f49b65a559dee369dba2360cc92cb01cf5"),
		31, 79,
		33, 15969,
	)
	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f"),
		31, 79,
		33, 15969,
	)
	nfsSuite.assertWalk(
		path.Join(nfsSuite.mountPoint, "EVO-Exchange-BE-2019", "master"),
		75, 256,
		0, 677127,
	)
}
