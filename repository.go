package git

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/mechmind/git-go/git"
	"github.com/mechmind/git-go/rawgit"
)

// repo.go ports

const (
	BRANCH_PREFIX      = "refs/heads/"
	_PRETTY_LOG_FORMAT = `--pretty=format:%H`
)

type Repository struct {
	Path string

	repo *git.Repository
}

func InitRepository(path string, bare bool) error {
	// FIXME: implement
	return nil
}

func OpenRepository(path string) (*Repository, error) {
	// FIXME: implement
	return nil, nil
}

type CloneRepoOptions struct {
	Mirror  bool
	Bare    bool
	Quiet   bool
	Timeout time.Duration
}

type Branch struct {
	Name string
	Path string
}

func Clone(from, to string, opts CloneRepoOptions) (err error) {
	// FIXME: implement
	return nil
}

func Pull(repoPath string, all bool) error {
	// FIXME: implement
	return nil
}

func Push(repoPath, remote, branch string) error {
	// FIXME: implement
	return nil
}

func ResetHEAD(repoPath string, hard bool, revision string) error {
	// FIXME: implement
	return nil
}

// repo_object.go ports

type ObjectType string

const (
	OBJECT_COMMIT ObjectType = "commit"
	OBJECT_TREE   ObjectType = "tree"
	OBJECT_BLOB   ObjectType = "blob"
	OBJECT_TAG    ObjectType = "tag"
)

// repo_branch.go ports

func IsReferenceExist(repoPath, name string) bool {
	repo, err := OpenRepository(repoPath)
	if err != nil {
		return false
	}

	oid, err := repo.repo.ResolveRef(name)
	return oid != nil
}

func IsBranchExist(repoPath, name string) bool {
	repo, err := OpenRepository(repoPath)
	if err != nil {
		return false
	}

	return repo.IsBranchExist(name)
}

func (repo *Repository) IsBranchExist(name string) bool {
	oid, _ := repo.repo.ResolveBranch(name)
	return oid != nil
}

func (repo *Repository) GetHEADBranch() (*Branch, error) {
	value, err := repo.repo.ReadRef("HEAD")
	if err != nil {
		return nil, err
	}

	branchPrefix := rawgit.RefPrefix + rawgit.RefBranchNS
	if !strings.HasPrefix(value, branchPrefix) {
		return nil, fmt.Errorf("invalid HEAD: %v", value)
	}

	branch := &Branch{
		Name: value[len(branchPrefix):],
		Path: value[len(rawgit.RefPrefix):],
	}

	return branch, nil
}

func (repo *Repository) SetDefaultBranch(name string) error {
	ref := rawgit.RefPrefix + rawgit.RefBranchNS + name
	return repo.repo.WriteRef("HEAD", ref)
}

func (repo *Repository) GetBranches() ([]string, error) {
	heads, err := repo.repo.ListRefs(rawgit.RefBranchNS)
	if err != nil {
		return nil, err
	}

	branches := []string{}
	for _, head := range heads {
		branches = append(branches, head[len(rawgit.RefBranchNS):])
	}
	return branches, nil
}

func (repo *Repository) AddRemote(name, url string, fetch bool) error {
	// FIXME: implement
	return nil
}

func (repo *Repository) RemoveRemote(name string) error {
	// FIXME: implement
	return nil
}

// repo_commit.go ports

func (repo *Repository) GetBranchCommitID(name string) (string, error) {
	oid, err := repo.repo.ResolveRef(rawgit.RefBranchNS + name)
	if err != nil {
		return "", err
	}

	return oid.String(), nil
}

func (repo *Repository) getCommit(id sha1) (*Commit, error) {
	commit, err := repo.repo.OpenCommit(sha2oidp(id))
	if err != nil {
		return nil, err
	}

	return raw2commit(repo, commit)
}

func (repo *Repository) GetTagCommitID(name string) (string, error) {
	// FIXME: implement
	return "", nil
}

func (repo *Repository) GetCommit(commitID string) (*Commit, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) GetBranchCommit(name string) (*Commit, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) GetTagCommit(name string) (*Commit, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) GetCommitByPath(relpath string) (*Commit, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) FileCommitsCount(revision, file string) (int64, error) {
	// FIXME: implement
	return -1, nil
}
func (repo *Repository) CommitsByFileAndRange(revision, file string, page int) (*list.List, error) {

	// FIXME: implement
	return nil, nil
}
func (repo *Repository) FilesCountBetween(startCommitID, endCommitID string) (int, error) {

	// FIXME: implement
	return -1, nil
}
func (repo *Repository) CommitsBetween(last *Commit, before *Commit) (*list.List, error) {

	// FIXME: implement
	return nil, nil
}
func (repo *Repository) CommitsBetweenIDs(last, before string) (*list.List, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) CommitsCountBetween(start, end string) (int64, error) {
	// FIXME: implement
	return -1, nil
}

func (repo *Repository) getRefCommitID(name string) (string, error) {
	// FIXME: implement
	return "", nil
}
func parseCommitData(data []byte) (*Commit, error) {
	// FIXME: implement
	return nil, nil
}

func (repo *Repository) getCommitByPathWithID(id sha1, relpath string) (*Commit, error) {

	// FIXME: implement
	return nil, nil
}
func (repo *Repository) commitsByRange(id sha1, page int) (*list.List, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) searchCommits(id sha1, keyword string) (*list.List, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) commitsBefore(l *list.List, parent *list.Element, id sha1, current, limit int) error {

	// FIXME: implement
	return nil
}
func (repo *Repository) getCommitsBefore(id sha1) (*list.List, error) {
	// FIXME: implement
	return nil, nil
}
func (repo *Repository) getCommitsBeforeLimit(id sha1, num int) (*list.List, error) {
	// FIXME: implement
	return nil, nil
}

// repo_tree.go ports

func (repo *Repository) getTree(id sha1) (*Tree, error) {
	_, _, err := repo.repo.StatObject(sha2oidp(id))
	if err != nil {
		return nil, ErrNotExist{id.String(), ""}
	}

	return NewTree(repo, id), nil
}

// Find the tree object in the repository.
func (repo *Repository) GetTree(idStr string) (*Tree, error) {
	id, err := NewIDFromString(idStr)
	if err != nil {
		return nil, err
	}
	return repo.getTree(id)
}

// own helpers

func (repo *Repository) raw2commitList(src *list.List) (*list.List, error) {
	var err error
	for cur := src.Front(); cur != nil; cur = cur.Next() {
		cur.Value, err = raw2commit(repo, cur.Value.(*rawgit.Commit))
		if err != nil {
			return nil, err
		}
	}

	return src, nil
}