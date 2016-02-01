// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"strconv"
	"time"
)

// Tag represents a Git tag.
type Tag struct {
	Name    string
	ID      sha1
	repo    *Repository
	Object  sha1 // The id of this commit object
	Type    string
	Tagger  *Signature
	Message string
}

func (tag *Tag) Commit() (*Commit, error) {
	return tag.repo.getCommit(tag.Object)
}

// FIXME: parsing must be implemented in rawgit package

// Parse commit information from the (uncompressed) raw
// data from the commit object.
// \n\n separate headers from message
func parseTagData(data []byte) (*Tag, error) {
	tag := new(Tag)
	// we now have the contents of the commit object. Let's investigate...
	nextline := 0
l:
	for {
		eol := bytes.IndexByte(data[nextline:], '\n')
		switch {
		case eol > 0:
			line := data[nextline : nextline+eol]
			spacepos := bytes.IndexByte(line, ' ')
			reftype := line[:spacepos]
			switch string(reftype) {
			case "object":
				id, err := NewIDFromString(string(line[spacepos+1:]))
				if err != nil {
					return nil, err
				}
				tag.Object = id
			case "type":
				// A commit can have one or more parents
				tag.Type = string(line[spacepos+1:])
			case "tagger":
				sig, err := newSignatureFromCommitline(line[spacepos+1:])
				if err != nil {
					return nil, err
				}
				tag.Tagger = sig
			}
			nextline += eol + 1
		case eol == 0:
			tag.Message = string(data[nextline+1:])
			break l
		default:
			break l
		}
	}
	return tag, nil
}

// FIXME: include timezone for timestamp!
// FIXME: move parsing to rawgit module
func newSignatureFromCommitline(line []byte) (_ *Signature, err error) {
	sig := new(Signature)
	emailStart := bytes.IndexByte(line, '<')
	sig.Name = string(line[:emailStart-1])
	emailEnd := bytes.IndexByte(line, '>')
	sig.Email = string(line[emailStart+1 : emailEnd])

	// Check date format.
	firstChar := line[emailEnd+2]
	if firstChar >= 48 && firstChar <= 57 {
		timestop := bytes.IndexByte(line[emailEnd+2:], ' ')
		timestring := string(line[emailEnd+2 : emailEnd+2+timestop])
		seconds, _ := strconv.ParseInt(timestring, 10, 64)
		sig.When = time.Unix(seconds, 0)
	} else {
		sig.When, err = time.Parse("Mon Jan _2 15:04:05 2006 -0700", string(line[emailEnd+2:]))
		if err != nil {
			return nil, err
		}
	}
	return sig, nil
}
