// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"strings"

	"github.com/mechmind/git-go/rawgit"
)

const TAG_PREFIX = "refs/tags/"

// IsTagExist returns true if given tag exists in the repository.
func IsTagExist(repoPath, name string) bool {
	return IsReferenceExist(repoPath, TAG_PREFIX+name)
}

func (repo *Repository) IsTagExist(name string) bool {
	oid, _ := repo.repo.ReadRef(TAG_PREFIX + name)
	return oid != ""
}

func (repo *Repository) CreateTag(name, revision string) error {
	return repo.repo.WriteRef(TAG_PREFIX+name, revision)
}

// GetTag returns a Git tag by given name.
func (repo *Repository) GetTag(name string) (*Tag, error) {
	oid, err := repo.repo.ResolveRef(TAG_PREFIX + name)
	if err != nil {
		return nil, err
	}

	info, _, err := repo.repo.StatObject(oid)
	if err != nil {
		return nil, err
	}

	if info.GetOType() == rawgit.OTypeCommit {
		return &Tag{
			ID:     sha1(*oid),
			Object: sha1(*oid),
			Type:   string(OBJECT_COMMIT),
			Name:   name,
			repo:   repo,
		}, nil
	}

	if info.GetOType() == rawgit.OTypeTag {
		obj, err := repo.repo.OpenTag(oid)
		if err != nil {
			return nil, err
		}

		tag := raw2tag(repo, obj)
		tag.Name = name
		return tag, nil
	}

	return nil, fmt.Errorf("invalid tag target: %s", info.GetOType().String())
}

// GetTags returns all tags of the repository.
func (repo *Repository) GetTags() ([]string, error) {
	rawRefs, err := repo.repo.ListRefs(TAG_PREFIX)
	if err != nil {
		return nil, err
	}

	refs := []string{}
	for _, rawRef := range rawRefs {
		refs = append(refs, strings.TrimPrefix(rawRef, TAG_PREFIX))
	}

	return refs, nil
}
