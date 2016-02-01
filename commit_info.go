package git

import (
	"fmt"
	"path"
	"path/filepath"
)

// FIXME: investigate later
type commitInfo struct {
	entryName string
	infos     []interface{}
	err       error
}

// GetCommitsInfo takes advantages of concurrey to speed up getting information
// of all commits that are corresponding to these entries.
// TODO: limit max goroutines at same time
func (tes Entries) GetCommitsInfo(commit *Commit, treePath string) ([][]interface{}, error) {
	if len(tes) == 0 {
		return nil, nil
	}

	revChan := make(chan commitInfo, 10)

	infoMap := make(map[string][]interface{}, len(tes))
	for i := range tes {
		if tes[i].Type != OBJECT_COMMIT {
			go func(i int) {
				cinfo := commitInfo{entryName: tes[i].Name()}
				c, err := commit.GetCommitByPath(filepath.Join(treePath, tes[i].Name()))
				if err != nil {
					cinfo.err = fmt.Errorf("GetCommitByPath (%s/%s): %v", treePath, tes[i].Name(), err)
				} else {
					cinfo.infos = []interface{}{tes[i], c}
				}
				revChan <- cinfo
			}(i)
			continue
		}

		// Handle submodule
		go func(i int) {
			cinfo := commitInfo{entryName: tes[i].Name()}
			sm, err := commit.GetSubModule(path.Join(treePath, tes[i].Name()))
			if err != nil && !IsErrNotExist(err) {
				cinfo.err = fmt.Errorf("GetSubModule (%s/%s): %v", treePath, tes[i].Name(), err)
				revChan <- cinfo
				return
			}

			smUrl := ""
			if sm != nil {
				smUrl = sm.Url
			}

			c, err := commit.GetCommitByPath(filepath.Join(treePath, tes[i].Name()))
			if err != nil {
				cinfo.err = fmt.Errorf("GetCommitByPath (%s/%s): %v", treePath, tes[i].Name(), err)
			} else {
				cinfo.infos = []interface{}{tes[i], NewSubModuleFile(c, smUrl, tes[i].ID.String())}
			}
			revChan <- cinfo
		}(i)
	}

	i := 0
	for info := range revChan {
		if info.err != nil {
			return nil, info.err
		}

		infoMap[info.entryName] = info.infos
		i++
		if i == len(tes) {
			break
		}
	}

	commitsInfo := make([][]interface{}, len(tes))
	for i := 0; i < len(tes); i++ {
		commitsInfo[i] = infoMap[tes[i].Name()]
	}
	return commitsInfo, nil
}
