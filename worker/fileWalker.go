package worker

import (
	"github.com/qfarm/qfarm"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"log"
)

type FilesMap struct {
    FilesMap map[string]*qfarm.Node
	Root string
}

func BuildTree(repoDir string) (*FilesMap, error) {
	tree := &FilesMap{FilesMap: make(map[string]*qfarm.Node)}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".git") || strings.Contains(path, ".idea") || strings.Contains(path, "vendor") || strings.Contains(path, "Godeps") {
			return nil
		}

		tree.FilesMap[path] = &qfarm.Node{
			Path:  path,
			Dir:   info.IsDir(),
			Nodes: make([]qfarm.Node, 0),
		}

		// if file is not dir - read content as bytes array.
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			tree.FilesMap[path].Content = bytes
		}

		return nil
	}

	if err := filepath.Walk(repoDir, walkFunc); err != nil {
		return nil, err
	}

	for path, node := range tree.FilesMap {
		parentPath := filepath.Dir(path)
		parent, exists := tree.FilesMap[parentPath]
		if exists {
			node.ParentPath = parent.Path
			tmp := *node
			tmp.Content = nil
			tmp.Nodes = nil
			parent.Nodes = append(parent.Nodes, tmp)
		} else {
			// If a parent doesn't exist, this is the root.
			tree.Root = path
		}
	}

	return tree, nil
}

func (t *FilesMap) ApplyIssue(i *qfarm.Issue) error {
	toApply := []string{t.Root, i.Path}

	// find path with all dirs to update
	subPath := strings.Replace(i.Path, t.Root, "", -1)
	for i, val := range subPath {
		if i != 0 && val == '/' {
			toApply = append(toApply, t.Root + subPath[0:i])
		}
	}

	for _, key := range toApply {
		val, ok := t.FilesMap[key]

		if ok {
			if val.Issues == nil {
				val.Issues = make([]*qfarm.Issue, 0)
			}
			val.Issues = append(val.Issues, i)
			val.IssuesNo++
			if i.Severity == qfarm.Error {
				val.ErrorsNo++
			}
			if i.Severity == qfarm.Warning {
				val.WarningsNo++
			}
			t.FilesMap[key] = val
		} else {
			log.Printf("WARNING: Can't find %s key in FilesMap", val)
		}
	}

	return nil
}