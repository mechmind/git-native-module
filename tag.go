// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import "github.com/mechmind/git-go/rawgit"

// Tag represents a Git tag.
type Tag struct {
	Name    string
	ID      sha1
	repo    *Repository
	Object  sha1 // The id of this commit object
	Type    string
	Tagger  *Signature
	Message string

	rawTag *rawgit.Tag
}

func (tag *Tag) Commit() (*Commit, error) {
	return tag.repo.getCommit(tag.Object)
}

func raw2tag(repo *Repository, raw *rawgit.Tag) *Tag {
	return &Tag{
		Name:    raw.Name,
		ID:      sha1(*raw.GetOID()),
		repo:    repo,
		Object:  sha1(raw.TargetOID),
		Tagger:  raw2signature(raw.Tagger),
		Message: raw.Message,

		rawTag: raw,
	}
}
