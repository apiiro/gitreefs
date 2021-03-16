package testutils

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/core/git"
	"gitreefs/core/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func SetupClones() string {
	clonesPath, err := ioutil.TempDir("", "")
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
			Dir:  clonesPath,
		}
		err = cloneCmd.Start()
		if err != nil {
			panic(err)
		}
		err = cloneCmd.Wait()
	}

	return clonesPath
}

func assertWalk(
	suite *suite.Suite,
	commitishPath string,
	expectedDirsCount int,
	expectedFilesCount int,
	expectedMinFileSize int,
	expectedMaxFileSize int,
) {
	dirsCount, filesCount := 0, 0
	minFileSize, maxFileSize := int(git.MaxFileSizeBytes), 0
	filepath.Walk(commitishPath, func(path string, info os.FileInfo, err error) error {
		suite.NotNil(info, "missing info for %v", path)
		if info == nil {
			return nil
		}
		if info.IsDir() {
			dirsCount++
			suite.EqualValues(0, info.Size())
		} else {
			filesCount++
			suite.True(info.Size() >= 0, "file at %v has invalid size", path)
			suite.True(info.Mode().IsRegular())
			content, err := ioutil.ReadFile(path)
			suite.Nil(err, "failed to read file at %v", path)
			fileSizeFromRead := len(content)
			suite.EqualValues(info.Size(), fileSizeFromRead, "read different file size for %v", path)
			if fileSizeFromRead > maxFileSize {
				maxFileSize = fileSizeFromRead
			}
			if fileSizeFromRead < minFileSize {
				minFileSize = fileSizeFromRead
			}
		}
		return nil
	})
	suite.EqualValues(expectedDirsCount, dirsCount, "unexpected dirs count")
	suite.EqualValues(expectedFilesCount, filesCount, "unexpected files count")
	suite.EqualValues(expectedMinFileSize, minFileSize, "unexpected min file size")
	suite.EqualValues(expectedMaxFileSize, maxFileSize, "unexpected max file size")
}

func WalkFileSystem(
	suite *suite.Suite,
	mountPoint string,
) {
	files, err := ioutil.ReadDir(mountPoint)
	suite.Nil(err)
	suite.Empty(files)

	files, err = ioutil.ReadDir(path.Join(mountPoint, "dc-heacth"))
	suite.Nil(err)
	suite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(mountPoint, "dc-heacth", "wat"))
	suite.NotNil(err)
	suite.True(os.IsNotExist(err))

	assertWalk(
		suite,
		path.Join(mountPoint, "dc-heacth", "2ca742044ba451d00c6854a465fdd4280d9ad1f5"),
		28, 181,
		215, 47804,
	)

	assertWalk(
		suite, path.Join(mountPoint, "dc-heacth", "2ca7420"),
		28, 181,
		215, 47804,
	)

	assertWalk(
		suite, path.Join(mountPoint, "dc-heacth", "master"),
		28, 183,
		215, 47814,
	)

	files, err = ioutil.ReadDir(path.Join(mountPoint, "EVO-Exchange-BE-2019"))
	suite.Nil(err)
	suite.Empty(files)

	_, err = ioutil.ReadDir(path.Join(mountPoint, "EVO-Exchange-BE-2019", "wat"))
	suite.NotNil(err)
	suite.True(os.IsNotExist(err))

	assertWalk(
		suite, path.Join(mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f49b65a559dee369dba2360cc92cb01cf5"),
		31, 79,
		33, 15969,
	)

	assertWalk(
		suite, path.Join(mountPoint, "EVO-Exchange-BE-2019", "1ca7c1f"),
		31, 79,
		33, 15969,
	)

	assertWalk(
		suite, path.Join(mountPoint, "EVO-Exchange-BE-2019", "master"),
		75, 256,
		0, 677127,
	)
}
