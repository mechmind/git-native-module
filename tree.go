package git

import (
	"bytes"
	"path"
	"sort"
	"strings"

	"github.com/mechmind/git-go/rawgit"
)

type Tree struct {
	ID   sha1
	repo *Repository

	// parent tree
	ptree *Tree

	entries       Entries
	entriesParsed bool

	raw *rawgit.Tree
}

type Entries []*TreeEntry

type TreeEntry struct {
	ID   sha1
	Type ObjectType

	mode EntryMode
	name string

	ptree *Tree

	commited bool

	size  int64
	sized bool
}

type EntryMode int

// There are only a few file modes in Git. They look like unix file modes, but they can only be
// one of these.
const (
	ENTRY_MODE_BLOB    EntryMode = 0100644
	ENTRY_MODE_EXEC    EntryMode = 0100755
	ENTRY_MODE_SYMLINK EntryMode = 0120000
	ENTRY_MODE_COMMIT  EntryMode = 0160000
	ENTRY_MODE_TREE    EntryMode = 0040000
)

func NewTree(repo *Repository, id sha1) *Tree {
	return &Tree{
		ID:   id,
		repo: repo,
	}
}

var escapeChar = []byte("\\")

// UnescapeChars reverses escaped characters.
func UnescapeChars(in []byte) []byte {
	if bytes.Index(in, escapeChar) == -1 {
		return in
	}

	endIdx := len(in) - 1
	isEscape := false
	out := make([]byte, 0, endIdx+1)
	for i := range in {
		if in[i] == '\\' && !isEscape {
			isEscape = true
			continue
		}
		isEscape = false
		out = append(out, in[i])
	}
	return out
}

func (t *Tree) SubTree(rpath string) (*Tree, error) {
	// FIXME: reimplement
	if len(rpath) == 0 {
		return t, nil
	}

	paths := strings.Split(rpath, "/")
	var (
		err error
		g   = t
		p   = t
		te  *TreeEntry
	)
	for _, name := range paths {
		te, err = p.GetTreeEntryByPath(name)
		if err != nil {
			return nil, err
		}

		g, err = t.repo.getTree(te.ID)
		if err != nil {
			return nil, err
		}
		g.ptree = p
		p = g
	}
	return g, nil
}

func raw2treeEntries(ptree *Tree, items []rawgit.TreeItem) (Entries, error) {
	entries := make(Entries, len(items))
	for idx, item := range items {
		entry := &TreeEntry{
			ID:   sha1(*item.GetOID()),
			Type: ObjectType(item.GetOType().String()),

			mode: EntryMode(item.Mode),
			name: item.Name,

			ptree: ptree,
		}

		entries[idx] = entry
	}

	return entries, nil
}

func (t *Tree) ListEntries() (Entries, error) {
	if t.entriesParsed {
		return t.entries, nil
	}
	t.entriesParsed = true

	entries, err := raw2treeEntries(t, t.raw.Items)
	if err != nil {
		return nil, err
	}

	t.entries = entries
	return t.entries, nil
}

func (te *TreeEntry) Name() string {
	return te.name
}

func (te *TreeEntry) Size() int64 {
	if te.IsDir() {
		return 0
	} else if te.sized {
		return te.size
	}

	info, _, err := te.ptree.repo.repo.StatObject(sha2oidp(te.ID))
	if err != nil {
		return 0
	}

	te.sized = true
	te.size = int64(info.Size)

	return te.size
}

func (te *TreeEntry) IsSubModule() bool {
	return te.mode == ENTRY_MODE_COMMIT
}

func (te *TreeEntry) IsDir() bool {
	return te.mode == ENTRY_MODE_TREE
}

func (te *TreeEntry) Blob() *Blob {
	return &Blob{
		repo:      te.ptree.repo,
		TreeEntry: te,
	}
}

var sorter = []func(t1, t2 *TreeEntry) bool{
	func(t1, t2 *TreeEntry) bool {
		return (t1.IsDir() || t1.IsSubModule()) && !t2.IsDir() && !t2.IsSubModule()
	},
	func(t1, t2 *TreeEntry) bool {
		return t1.name < t2.name
	},
}

func (tes Entries) Len() int      { return len(tes) }
func (tes Entries) Swap(i, j int) { tes[i], tes[j] = tes[j], tes[i] }
func (tes Entries) Less(i, j int) bool {
	t1, t2 := tes[i], tes[j]
	var k int
	for k = 0; k < len(sorter)-1; k++ {
		sort := sorter[k]
		switch {
		case sort(t1, t2):
			return true
		case sort(t2, t1):
			return false
		}
	}
	return sorter[k](t1, t2)
}

func (tes Entries) Sort() {
	sort.Sort(tes)
}

func (t *Tree) GetTreeEntryByPath(relpath string) (*TreeEntry, error) {
	if len(relpath) == 0 {
		return &TreeEntry{
			ID:    t.ID,
			Type:  OBJECT_TREE,
			mode:  ENTRY_MODE_TREE,
			ptree: t.ptree,
		}, nil
	}

	relpath = path.Clean(relpath)
	parts := strings.Split(relpath, "/")
	var err error
	tree := t
	for i, name := range parts {
		if i == len(parts)-1 {
			entries, err := tree.ListEntries()
			if err != nil {
				return nil, err
			}
			for _, v := range entries {
				if v.name == name {
					return v, nil
				}
			}
		} else {
			tree, err = tree.SubTree(name)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, ErrNotExist{"", relpath}
}

func (t *Tree) GetBlobByPath(relpath string) (*Blob, error) {
	entry, err := t.GetTreeEntryByPath(relpath)
	if err != nil {
		return nil, err
	}

	if !entry.IsDir() {
		return entry.Blob(), nil
	}

	return nil, ErrNotExist{"", relpath}
}
