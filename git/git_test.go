package git

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type gitTestSuite struct {
	suite.Suite
	remote    string
	clonePath string
	provider  *RepositoryProvider
}

func TestGitTestSuite(t *testing.T) {
	suite.Run(t, new(gitTestSuite))
}

func cloneLocal(remote string) (clonePath string, err error) {
	clonePath, err = ioutil.TempDir("", "")
	if err != nil {
		return
	}
	proc := exec.Command("git", "clone", "--no-checkout", remote, clonePath)
	err = proc.Start()
	if err != nil {
		return
	}
	err = proc.Wait()
	return
}

func (gitSuite *gitTestSuite) SetupTest() {
	gitSuite.remote = "https://github.com/apiirolab/dc-heacth.git"
	clonePath, err := cloneLocal(gitSuite.remote)
	if err != nil {
		panic(err)
		return
	}
	gitSuite.clonePath = clonePath
	gitSuite.provider, err = NewRepository(clonePath)
}

func (gitSuite *gitTestSuite) TearDownTest() {
	os.RemoveAll(gitSuite.clonePath)
}

func (gitSuite *gitTestSuite) TestListTreeForRegularCommit() {
	tree, err := gitSuite.provider.ListTree("2ca742044ba451d00c6854a465fdd4280d9ad1f5")
	gitSuite.Nil(err, "git.ListTree: %w", err)
	gitSuite.EqualValues(209, len(tree), "tree size not as expected")
	gitSuite.Contains(tree, "", "no root entry")

	gitSuite.Contains(tree, "src", "no src dir")
	dirEntry := tree["src"]
	gitSuite.True(dirEntry.IsDir)
	gitSuite.EqualValues(0, dirEntry.Size)
	gitSuite.Equal("src", dirEntry.Name)
	gitSuite.Equal("", dirEntry.ParentPath)

	gitSuite.Contains(tree, "src/main/java/com/dchealth/service/common", "no common dir")
	dirEntry = tree["src/main/java/com/dchealth/service/common"]
	gitSuite.True(dirEntry.IsDir)
	gitSuite.EqualValues(0, dirEntry.Size)
	gitSuite.Equal("common", dirEntry.Name)
	gitSuite.Equal("src/main/java/com/dchealth/service", dirEntry.ParentPath)

	gitSuite.Contains(tree, "src/main/java/com/dchealth/service/common/YunUserService.java", "no java file")
	fileEntry := tree["src/main/java/com/dchealth/service/common/YunUserService.java"]
	gitSuite.False(fileEntry.IsDir)
	gitSuite.EqualValues(28092, fileEntry.Size)
	gitSuite.Equal("YunUserService.java", fileEntry.Name)
	gitSuite.Equal("src/main/java/com/dchealth/service/common", fileEntry.ParentPath)
}

func (gitSuite *gitTestSuite) TestListTreeForNonExisting() {
	_, err := gitSuite.provider.ListTree("wat")
	gitSuite.NotNil(err)
	_, err = gitSuite.provider.ListTree("23ce6f6bd72532aa410afeb8939ed6911c526f60f1411c1a40952928f90e15ad")
	gitSuite.NotNil(err)
}

func (gitSuite *gitTestSuite) TestListTreeForShortSha() {
	tree, err := gitSuite.provider.ListTree("2ca7420")
	gitSuite.Nil(err, "git.ListTree: %w", err)
	gitSuite.EqualValues(209, len(tree), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestListTreeForMainBranchName() {
	tree, err := gitSuite.provider.ListTree("master")
	gitSuite.Nil(err, "git.ListTree: %w", err)
	gitSuite.EqualValues(211, len(tree), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestListTreeForBranchName() {
	tree, err := gitSuite.provider.ListTree("remotes/origin/lfx")
	gitSuite.Nil(err, "git.ListTree: %w", err)
	gitSuite.EqualValues(209, len(tree), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestFileContents() {
	contents, err := gitSuite.provider.FileContents("2ca742044ba451d00c6854a465fdd4280d9ad1f5", "src/main/java/com/dchealth/service/common/YunUserService.java")
	gitSuite.Nil(err, "git.ListTree: %w", err)
	gitSuite.EqualValues(28092, len(contents), "file contents size not as expected")
}

func (gitSuite *gitTestSuite) TestFileContentsForNonExisting() {
	contents, err := gitSuite.provider.FileContents("2ca742044ba451d00c6854a465fdd4280d9ad1f5", "src/YunUserService.java")
	gitSuite.NotNil(err)
	gitSuite.EqualValues(0, len(contents), "file contents size not as expected")
	contents, err = gitSuite.provider.FileContents("wat", "src/main/java/com/dchealth/service/common/YunUserService.java")
	gitSuite.NotNil(err)
	gitSuite.EqualValues(0, len(contents), "file contents size not as expected")
}
