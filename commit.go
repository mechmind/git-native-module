package git

import (
	"bufio"
	"container/list"
	"net/http"
	"strings"
	"time"

	"github.com/mechmind/git-go/history"
	"github.com/mechmind/git-go/rawgit"
)

type Commit struct {
	Tree
	ID            sha1 // The ID of this commit object
	Author        *Signature
	Committer     *Signature
	CommitMessage string

	//parents []sha1 // SHA1 strings
	submoduleCache *objectCache

	rawCommit *rawgit.Commit
}

type Signature struct {
	Email string
	Name  string
	When  time.Time
}

func raw2commit(repo *Repository, raw *rawgit.Commit) (*Commit, error) {
	commit := &Commit{
		ID:            sha1(*raw.GetOID()),
		Author:        raw2signature(raw.Author),
		Committer:     raw2signature(raw.Committer),
		CommitMessage: raw.Message,

		rawCommit: raw,
	}

	return commit, nil
}

func raw2signature(raw rawgit.UserTime) *Signature {
	return &Signature{
		Name:  raw.Name,
		Email: raw.Email,
		When:  raw.Time,
	}
}

// Message returns the commit message. Same as retrieving CommitMessage directly.
func (c *Commit) Message() string {
	return c.CommitMessage
}

// Summary returns first line of commit message.
func (c *Commit) Summary() string {
	return strings.Split(c.CommitMessage, "\n")[0]
}

// ParentID returns oid of n-th parent (0-based index).
// It returns nil if no such parent exists.
func (c *Commit) ParentID(n int) (sha1, error) {
	if n >= len(c.rawCommit.ParentOIDs) {
		return sha1{}, ErrNotExist{"", ""}
	}
	return sha1(*c.rawCommit.ParentOIDs[n]), nil
}

// Parent returns n-th parent (0-based index) of the commit.
func (c *Commit) Parent(n int) (*Commit, error) {
	id, err := c.ParentID(n)
	if err != nil {
		return nil, err
	}

	parent, err := c.repo.getCommit(id)
	if err != nil {
		return nil, err
	}
	return parent, nil
}

// ParentCount returns number of parents of the commit.
// 0 if this is the root commit,  otherwise 1,2, etc.
func (c *Commit) ParentCount() int {
	return len(c.rawCommit.ParentOIDs)
}

func isImageFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "image/") != -1 {
		return contentType, true
	}
	return contentType, false
}

func (c *Commit) IsImageFile(name string) bool {
	blob, err := c.GetBlobByPath(name)
	if err != nil {
		return false
	}

	dataRc, err := blob.Data()
	if err != nil {
		return false
	}
	buf := make([]byte, 1024)
	n, _ := dataRc.Read(buf)
	if n > 0 {
		buf = buf[:n]
	}
	_, isImage := isImageFile(buf)
	return isImage
}

// GetCommitByPath return the commit of relative path object.
func (c *Commit) GetCommitByPath(relpath string) (*Commit, error) {
	pathCb := history.MakePathChecker(c.repo.repo, relpath)
	pager := history.MakePager(c.repo.repo, pathCb, 0, 1)
	cmp := history.MakePathComparator(c.repo.repo, relpath)

	hist := history.New(c.repo.repo)
	result, err := hist.WalkFilteredHistory(sha2oidp(c.ID), pager, cmp)
	if err != nil {
		return nil, err
	}

	if result.Len() == 0 {
		// no entries
		return nil, nil
	}

	return raw2commit(c.repo, result.Front().Value.(*rawgit.Commit))
}

// AddAllChanges marks local changes to be ready for commit.
func AddChanges(repoPath string, all bool, files ...string) error {
	/*
		cmd := NewCommand("add")
		if all {
			cmd.AddArguments("--all")
		}
		_, err := cmd.AddArguments(files...).RunInDir(repoPath)
	*/
	// FIXME
	panic("not implemented")
	return nil
}

func CommitChanges(repoPath, message string, author *Signature) error {
	/*
		cmd := NewCommand("commit", "-m", message)
		if author != nil {
			cmd.AddArguments(fmt.Sprintf("--author='%s <%s>'", author.Name, author.Email))
		}
		_, err := cmd.RunInDir(repoPath)

		// No stderr but exit status 1 means nothing to commit.
		if err != nil && err.Error() == "exit status 1" {
			return nil
		}
	*/
	// FIXME
	panic("not implemented")
	return nil
}

// CommitsCount returns number of total commits of until given revision.
func CommitsCount(repoPath, revision string) (int64, error) {
	repo, err := OpenRepository(repoPath)
	if err != nil {
		return 0, err
	}

	commit, err := repo.GetCommit(revision)
	if err != nil {
		return 0, err
	}

	return commit.CommitsCount()
}

func (c *Commit) CommitsCount() (int64, error) {
	counter, result := history.MakeCounter(nil)

	hist := history.New(c.repo.repo)
	_, err := hist.WalkHistory(sha2oidp(c.ID), counter)
	if err != nil {
		return 0, err
	}

	return int64(result()), nil
}

var CommitsRangeSize = 50

func (c *Commit) CommitsByRange(page int) (*list.List, error) {
	pager := history.MakePager(c.repo.repo, nil, (page-1)*CommitsRangeSize, CommitsRangeSize)
	hist := history.New(c.repo.repo)
	result, err := hist.WalkHistory(sha2oidp(c.ID), pager)
	if err != nil {
		return nil, err
	}

	return c.repo.raw2commitList(result)
}

func (c *Commit) CommitsBefore() (*list.List, error) {
	hist := history.New(c.repo.repo)
	result, err := hist.WalkHistory(sha2oidp(c.ID), nil)
	if err != nil {
		return nil, err
	}

	return c.repo.raw2commitList(result)
}

func (c *Commit) CommitsBeforeLimit(num int) (*list.List, error) {
	pager := history.MakePager(c.repo.repo, nil, 0, num)
	hist := history.New(c.repo.repo)
	result, err := hist.WalkHistory(sha2oidp(c.ID), pager)
	if err != nil {
		return nil, err
	}

	return c.repo.raw2commitList(result)
}

func (c *Commit) CommitsBeforeUntil(commitID string) (*list.List, error) {
	/*
		endCommit, err := c.repo.GetCommit(commitID)
		if err != nil {
			return nil, err
		}
		return c.repo.CommitsBetween(c, endCommit)
	*/
	// FIXME
	panic("not implemented!")
	return nil, nil
}

func (c *Commit) SearchCommits(keyword string) (*list.List, error) {
	searcher, err := history.MakeHistorySearcher(keyword)
	if err != nil {
		return nil, err
	}

	hist := history.New(c.repo.repo)
	result, err := hist.WalkHistory(sha2oidp(c.ID), searcher)
	if err != nil {
		return nil, err
	}

	return c.repo.raw2commitList(result)
}

func (c *Commit) GetSubModules() (*objectCache, error) {
	if c.submoduleCache != nil {
		return c.submoduleCache, nil
	}

	entry, err := c.GetTreeEntryByPath(".gitmodules")
	if err != nil {
		return nil, err
	}
	rd, err := entry.Blob().Data()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(rd)
	c.submoduleCache = newObjectCache()
	var ismodule bool
	var path string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "[submodule") {
			ismodule = true
			continue
		}
		if ismodule {
			fields := strings.Split(scanner.Text(), "=")
			k := strings.TrimSpace(fields[0])
			if k == "path" {
				path = strings.TrimSpace(fields[1])
			} else if k == "url" {
				c.submoduleCache.Set(path, &SubModule{path, strings.TrimSpace(fields[1])})
				ismodule = false
			}
		}
	}

	return c.submoduleCache, nil
}

func (c *Commit) GetSubModule(entryname string) (*SubModule, error) {
	modules, err := c.GetSubModules()
	if err != nil {
		return nil, err
	}

	module, has := modules.Get(entryname)
	if has {
		return module.(*SubModule), nil
	}
	return nil, nil
}
