// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package git

import (
	"fmt"
	"time"
)

const _VERSION = "0.2.4"

func Version() string {
	return _VERSION
}

var (
	// Debug enables verbose logging on everything.
	// This should be false in case Gogs starts in SSH mode.
	Debug  = false
	Prefix = "[git-module] "
)

func log(format string, args ...interface{}) {
	if !Debug {
		return
	}

	fmt.Print(Prefix)
	if len(args) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// compatibility target
var gitVersion string = "1.8.0"

// Version returns current Git version from shell.
func BinVersion() (string, error) {
	return gitVersion, nil
}

// Fsck verifies the connectivity and validity of the objects in the database
func Fsck(repoPath string, timeout time.Duration, args ...string) error {
	/*
		// Make sure timeout makes sense.
		if timeout <= 0 {
			timeout = -1
		}
		_, err := NewCommand("fsck").AddArguments(args...).RunInDirTimeout(timeout, repoPath)
	*/
	panic("not implemented!")
	return nil
}
