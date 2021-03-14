package git

import (
	"github.com/stretchr/testify/suite"
	"gitreefs/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type gitTestSuite struct {
	suite.Suite
	remote    string
	clonePath string
	provider  *RepositoryProvider
}

func TestGitTestSuite(t *testing.T) {
	logger.InitLoggers("logs/git_test-%v.log", "ERROR", "-")
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
	gitSuite.provider, err = NewRepositoryProvider(clonePath)
}

func (gitSuite *gitTestSuite) TearDownTest() {
	os.RemoveAll(gitSuite.clonePath)
}

func countTreeNodes(node *Entry) uint {
	var count uint = 1
	if node.IsDir {
		for _, child := range node.EntriesByName {
			count += countTreeNodes(child)
		}
	}
	return count
}

func lookupNode(node *Entry, path string) *Entry {
	return lookupNodeRecursive(node, strings.Split(path, "/"))
}

func lookupNodeRecursive(node *Entry, pathParts []string) (target *Entry) {
	if len(pathParts) == 0 {
		return node
	}
	if !node.IsDir {
		return nil
	}
	currentPart := pathParts[0]
	child, found := node.EntriesByName[currentPart]
	if !found {
		return nil
	}
	remainderParts := pathParts[1:]
	return lookupNodeRecursive(child, remainderParts)
}

func (gitSuite *gitTestSuite) TestListTreeForRegularCommit() {
	tree, err := gitSuite.provider.ListTree("2ca742044ba451d00c6854a465fdd4280d9ad1f5")
	gitSuite.Nil(err, "git.ListTree: %v", err)
	gitSuite.EqualValues(209, len(tree.EntriesByPath), "tree size not as expected")
	gitSuite.EqualValues(209, countTreeNodes(&tree.Entry), "tree size not as expected")
	gitSuite.Contains(tree.EntriesByPath, "", "no root entry")
	dirEntry := tree.EntriesByPath[""]
	gitSuite.NotNil(dirEntry)
	gitSuite.True(dirEntry.IsDir)
	gitSuite.EqualValues(0, dirEntry.Size)
	gitSuite.Equal(4, len(dirEntry.EntriesByName))

	gitSuite.Contains(tree.EntriesByPath, "src", "no src dir")
	dirEntry = lookupNode(&tree.Entry, "src")
	gitSuite.NotNil(dirEntry)
	gitSuite.True(dirEntry.IsDir)
	gitSuite.EqualValues(0, dirEntry.Size)
	gitSuite.Equal(1, len(dirEntry.EntriesByName))

	gitSuite.Contains(tree.EntriesByPath, "src/main/java/com/dchealth/service/common", "no common dir")
	dirEntry = lookupNode(&tree.Entry, "src/main/java/com/dchealth/service/common")
	gitSuite.NotNil(dirEntry)
	gitSuite.True(dirEntry.IsDir)
	gitSuite.EqualValues(0, dirEntry.Size)
	gitSuite.Equal(7, len(dirEntry.EntriesByName))

	gitSuite.Contains(tree.EntriesByPath, "src/main/java/com/dchealth/service/common/YunUserService.java", "no java file")
	fileEntry := lookupNode(&tree.Entry, "src/main/java/com/dchealth/service/common/YunUserService.java")
	gitSuite.NotNil(dirEntry)
	gitSuite.False(fileEntry.IsDir)
	gitSuite.EqualValues(28092, fileEntry.Size)
	gitSuite.Nil(fileEntry.EntriesByName)

	gitSuite.NotContains(tree.EntriesByPath, "foo", "found fake dir")
	gitSuite.NotContains(tree.EntriesByPath, "foo/bar", "found fake dir")
	gitSuite.Nil(lookupNode(&tree.Entry, "foo"))
	gitSuite.Nil(lookupNode(&tree.Entry, "foo/bar"))
}

func (gitSuite *gitTestSuite) TestListTreeForNonExisting() {
	_, err := gitSuite.provider.ListTree("wat")
	gitSuite.NotNil(err)
	_, err = gitSuite.provider.ListTree("23ce6f6bd72532aa410afeb8939ed6911c526f60f1411c1a40952928f90e15ad")
	gitSuite.NotNil(err)
}

func (gitSuite *gitTestSuite) TestListTreeForShortSha() {
	tree, err := gitSuite.provider.ListTree("2ca7420")
	gitSuite.Nil(err, "git.ListTree: %v", err)
	gitSuite.EqualValues(209, countTreeNodes(&tree.Entry), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestListTreeForMainBranchName() {
	tree, err := gitSuite.provider.ListTree("master")
	gitSuite.Nil(err, "git.ListTree: %v", err)
	gitSuite.EqualValues(211, len(tree.EntriesByPath), "tree size not as expected")
	gitSuite.EqualValues(211, countTreeNodes(&tree.Entry), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestListTreeForBranchName() {
	tree, err := gitSuite.provider.ListTree("remotes/origin/lfx")
	gitSuite.Nil(err, "git.ListTree: %v", err)
	gitSuite.EqualValues(209, len(tree.EntriesByPath), "tree size not as expected")
	gitSuite.EqualValues(209, countTreeNodes(&tree.Entry), "tree size not as expected")
}

func (gitSuite *gitTestSuite) TestFileContents() {
	contents, err := gitSuite.provider.FileContents("2ca742044ba451d00c6854a465fdd4280d9ad1f5", "src/main/java/com/dchealth/service/common/YunUserService.java")
	gitSuite.Nil(err, "git.ListTree: %v", err)
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
